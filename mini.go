package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	wrapper "github.com/MobileAppTracking/measurement/logrusWrapper"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	logWriter *syslog.Writer
	db        *sql.DB
)

func init() {

	wrapper.AddHookToSyslog("tcp", "localhost:10514", syslog.LOG_EMERG, "Exact Details: ")

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	var err error
	//set up the logwriter
	logWriter, err = syslog.Dial("tcp", "localhost:10514", syslog.LOG_EMERG, "mini---project")
	if err != nil {
		log.Fatal(err)
	}

	//set up the database
	db, err = sql.Open("mysql", "root@tcp(localhost:3306)/logs?parseTime=true")
	if err != nil {
		logWriter.Err("can't open databases")
		return
	}

	// seed the random generator to generate IDs
	rand.Seed(time.Now().UTC().UnixNano())

}

// build/lauch the server and prepare to write logs
func main() {
	server := http.Server{
		Addr:    ":5000",
		Handler: myHandler(),
	}

	server.ListenAndServe()

	//close the database and log writer after the server stop running
	logWriter.Close()
	db.Close()
}

//build and return the server's handler
func myHandler() *mux.Router {
	mx := mux.NewRouter()
	mx.HandleFunc("/", Poster).Methods("POST")
	mx.HandleFunc("/{id}", GET).Methods("GET")
	return mx
}

//build the struct which holds the temp data

type Click struct {
	ID           string
	Type         string
	AdvertiserID int
	SiteID       int
	IP           string
	IosIfa       string
	GoogleAid    string
	WindowsAid   string
	Date_time    time.Time // store the time.Time struct in it and use to make the timestamp for elasticsearch later
}

type PostResponses struct {
	ErrMessage string
	Id         string
	HttpStatus string
}

//function that handles the GET method
//retrieve the json data form from database/server and output to the browser
func GET(writer http.ResponseWriter, reader *http.Request) {

	// time the method
	var starttime = time.Now()

	///get the id from the hashmap
	id := mux.Vars(reader)["id"]

	//select data from sql databases according to the id
	row := db.QueryRow("SELECT id, advertiser_id, site_id, ip, ios_ifa, google_aid, windows_aid, date_time FROM clicks WHERE id=?", id)

	//store the data from sql database in a temp struct
	var c Click

	err := row.Scan(&c.ID, &c.AdvertiserID, &c.SiteID, &c.IP, &c.IosIfa, &c.GoogleAid, &c.WindowsAid, &c.Date_time)

	//check for errors in scan  (404 and 500)
	if err == sql.ErrNoRows {

		// fmt.Println(err)

		writer.WriteHeader(http.StatusNotFound)

		// io.WriteString(writer, `{"message" : "Error 404"}`)
		// io.WriteString(writer, `{"httpstatus" : "404"}`)
		// log.Print(`{"message" : "Error 404"}, `)
		// log.Println(`{"httpstatus" : "404"}`)

		wrapper.Error(c.ID, c.Type, "Error 404: no rows found", c.SiteID)

		return
	} else if err != nil {

		// fmt.Println(err)

		writer.WriteHeader(http.StatusInternalServerError)

		// io.WriteString(writer, "Error 500")
		// log.Print(`{"message" : "Error 500"}, "`)
		// log.Println(`{"httpstatus" : "500"}`)

		wrapper.Error(c.ID, c.Type, "Error 500: Internal server error", c.SiteID)

		return
	}

	//marshal the data from the temp struct to json
	bytes, err := json.Marshal(&c)
	if err != nil {
		fmt.Println(err)
		return
	}

	// //output the raw bytes to the browser
	// writer.WriteHeader(http.StatusOK)
	// io.WriteString(writer, string(bytes))

	// log the event
	wrapper.info(c.ID, c.Type, "GET successful", c.SiteID, time.Now().Sub(starttime))

	//log the time
	// log.Println("GET: took", time.Now().Sub(starttime), "to execute")

}

//the function which handle the post method
//post the json data from broswer to the server and sql databases
func Poster(w http.ResponseWriter, r *http.Request) {
	//time when do request
	RequestStart := time.Now()

	//get the raw bytes of input data
	bytes, errs := ioutil.ReadAll(r.Body)
	if errs != nil {
		errString := fmt.Sprintf("buffer overflow %s", errs)
		response(w, errString, "", http.StatusBadRequest)
		return
	}

	//store the raw bytes to a temporary struct and log the Json invalid format
	var point Click
	errs = json.Unmarshal(bytes, &point)
	if errs != nil {
		errString := fmt.Sprintf("invalid Json format: %s", errs)
		response(w, errString, "", http.StatusBadRequest)
		logWriter.Err(errString)
		return
	}

	//validate the input and log error input message
	if point.AdvertiserID == 0 || point.SiteID == 0 {
		errString := "your advertiserID or site ID may equals to 0"
		response(w, errString, "", http.StatusBadRequest)
		logWriter.Err(errString)
		return
	}

	//generate a ramdom id for the post data and also get the ip address
	id := Id(point.AdvertiserID)
	ip := r.RemoteAddr

	//store the data from the struct to the sql databases and log the error or latency time
	QueryStart := time.Now()

	_, errs = db.Exec("INSERT INTO clicks(id, advertiser_id, site_id, ip, ios_ifa, google_aid, windows_aid, date_time) VALUES(?, ?, ?, ?, ?, ?, ?,?)",
		id, point.AdvertiserID, point.SiteID, ip, point.IosIfa, point.GoogleAid, point.WindowsAid, RequestStart)

	if errs != nil {
		fmt.Println(errs) //show to the sever protecter inside of users
		errString := "sorry, there is an error"
		response(w, errString, "", http.StatusInternalServerError)
		errString = fmt.Sprintf("database connection error : %s", errs)
		logWriter.Err(errString)
		return
	}

	//logWriter.Info(point.ID, point.Type, "the time for inserting data to clicks table is ", time.Since(QueryStart))

	//sucess and log the request latency
	response(w, "", id, http.StatusOK)
	logWriter.Info(point.ID, point.Type, "Post successful!", point.SiteID, time.Since(RequestStart), time.Since(QueryStart))
}

//write the post reponse (faliure /success) to the client in Json format
func response(w http.ResponseWriter, errMessage string, id string, status int) {
	w.WriteHeader(status)

	validate := PostResponses{errMessage, id, strconv.Itoa(status)}
	bytes, errs := json.Marshal(&validate)
	if errs != nil {
		fmt.Println(errs) // this errors is only for execution no need to output to user
	} else {
		w.Write(bytes)
	}

}

//generate a random id to represent the unique id
func Id(adId int) string {
	t := time.Now()
	year, month, day := t.Date()
	var id = Hex(4) + "-" + strconv.Itoa(year) +
		strconv.Itoa(int(month)) + strconv.Itoa(day) + "-" + strconv.Itoa(adId)
	return id
}

//generate a random string encoded hex value with given byte
func Hex(chunks int) string {
	var buffer bytes.Buffer

	bytes := make([]byte, 4)
	for i := 0; i < chunks; i++ {
		binary.LittleEndian.PutUint32(bytes, rand.Uint32())
		buffer.WriteString(hex.EncodeToString(bytes))
	}

	return buffer.String()
}
