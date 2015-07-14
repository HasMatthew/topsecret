package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"log/syslog"
	mrand "math/rand"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	logWriter *syslog.Writer
	db        *sql.DB
)

func init() {
	var err error
	logWriter, err = syslog.Dial("tcp", "localhost:10514", syslog.LOG_EMERG, "mini---porject")
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("mysql", "root@tcp(localhost:3306)/logs")
	if err != nil {
		logWriter.Err("can't open databases")
		return
	}
}

//build and lauch the server
func main() {
	server := http.Server{
		Addr:    ":5000",
		Handler: myHandler(),
	}

	server.ListenAndServe()

	logWriter.Close()
	db.Close()
}

//build and return the server's handler
func myHandler() *mux.Router {
	mx := mux.NewRouter()
	mx.HandleFunc("/", Poster).Methods("POST")
	mx.HandleFunc("/{id}", Geter).Methods("GET")
	return mx
}

//build the struct which holds the temp data

type Click struct {
	ID           string
	AdvertiserID int
	SiteID       int
	IP           string
	IosIfa       string
	GoogleAid    string
	WindowsAid   string
}

type PostResponses struct {
	ErrMessage string
	Id         string
	HttpStatus string
}

//the function which handle the get methods
//get the json data form from database and server to browser
func Geter(w http.ResponseWriter, r *http.Request) {
	///get the id from the name of maps
	id := mux.Vars(r)["id"]

	//select data from sql databases according to the id
	row := db.QueryRow("SELECT id, advertiser_id, site_id, ip, ios_ifa FROM clicks WHERE id=?", id)

	//store the data from sql database to the temp struct
	var c Click
	err := row.Scan(&c.ID, &c.AdvertiserID, &c.SiteID, &c.IP, &c.IosIfa)
	if err != nil {
		fmt.Println(err)
		return
	}

	//marshal the data from temp struct to the json
	bytes, errs := json.Marshal(&c)
	if errs != nil {
		fmt.Println(errs)
		return
	}

	//output the raw bytes to the browser
	io.WriteString(w, string(bytes))

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

	//store the raw bytes to a temporary struct
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

	_, errs = db.Exec("INSERT INTO clicks(id, advertiser_id, site_id, ip, ios_ifa) VALUES(?, ?, ?, ?, ?)",
		id, point.AdvertiserID, point.SiteID, ip, point.IosIfa)

	if errs != nil {
		errString := "sorry, there is an error"
		response(w, errString, "", http.StatusInternalServerError)
		errString = fmt.Sprintf("database connection error : %s", errs)
		logWriter.Err(errString)
		return
	}

	responseTime("the time for inserting data to clicks table is ", QueryStart)

	//sucess and log the request latency
	response(w, "", id, http.StatusOK)
	responseTime("the time for this Post request is ", RequestStart)
}

//report the query / request latency
func responseTime(message string, startTime time.Time) {
	responseDuration := time.Since(startTime)
	logWriter.Info(message + responseDuration.String())
}

//write the post reponse (faliure /success) to the client in  Json format
func response(w http.ResponseWriter, errMessage string, id string, status int) {
	w.WriteHeader(status)

	validate := PostResponses{errMessage, id, strconv.Itoa(status)}
	bytes, errs := json.Marshal(&validate)
	if errs != nil {
		fmt.Println(errs) // this errors is only for execution no need to output to user
	} else {
		io.WriteString(w, string(bytes))
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
		binary.LittleEndian.PutUint32(bytes, mrand.Uint32())
		buffer.WriteString(hex.EncodeToString(bytes))
	}

	return buffer.String()
}
