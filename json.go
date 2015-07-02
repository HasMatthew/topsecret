package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var jsonString string = `{"id": "1235", "ip": "1.2.3.4", "advertiser_id": 8}`

type click struct {
	Id           int
	Ip           string
	AdvertiserID int `json:"advertiser_id"`
}

type Click struct {
	ID           string
	AdvertiserID int
	SiteID       int
	IP           string
	IosIfa       string
	GoogleAid    string
	WindowsAid   string
}

func random() int {
	rand.Seed(time.Now().UTC().UnixNano())
	a := rand.Float64()*90000.0 + 10000.0
	return int(a)
}

func insertinto(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Inserting into clicks")
	reading, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(reading))

	var p Click

	err := json.Unmarshal(reading, &p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)

	conn, err := sql.Open("mysql", "root@tcp(localhost:3306)/logs")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.Exec("INSERT INTO clicks(id, advertiser_id, site_id, ip, ios_ifa) VALUES(?, ?, ?, ?, ?)", p.ID, p.AdvertiserID, p.SiteID, p.IP, p.IosIfa)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func deletefrom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "deleting from clicks")
	reading, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(reading))

	var p Click

	err := json.Unmarshal(reading, &p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)

	conn, err := sql.Open("mysql", "root@tcp(localhost:3306)/logs")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.Exec("DELETE FROM clicks(id, advertiser_id, site_id, ip, ios_ifa) WHERE(?, ?, ?, ?, ?)", p.ID, p.AdvertiserID, p.SiteID, p.IP, p.IosIfa)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func retrieve(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Query())

	var ID string
	ID = r.URL.String()

	ID = strings.TrimPrefix(ID, "/retrieve?id=")

	//fmt.Fprint(w, "\nretrieving from clicks")
	reading, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(reading))

	// var ID string
	// var AdvertiserID int
	// var SiteID int
	// var IP string
	// var IosIfa string

	// err := json.Unmarshal(reading, &ID)
	// if err != nil {
	// 	fmt.Printf("Error running Unmarshal: %s", err)
	// }
	// fmt.Println(ID)

	conn, err := sql.Open("mysql", "root@tcp(localhost:3306)/logs")
	if err != nil {
		fmt.Printf("Error running Open: %s", err)
		return
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT id, advertiser_id, site_id, ip, ios_ifa FROM clicks WHERE id=?", &ID)
	if err != nil {
		fmt.Printf("Error running Query: %s", err)
	}
	for rows.Next() {
		click := new(Click)
		if err = rows.Scan(&click.ID, &click.AdvertiserID, &click.SiteID, &click.IP, &click.IosIfa); err != nil {
			fmt.Printf("Error running scan: %s", err)
		}
		fmt.Fprint(w, "\n{\"ID\" : \"", click.ID, "\", \"AdvertiserID\" : \"", click.AdvertiserID, "\", \"SiteID\" : \"", click.SiteID, "\", \"IP\" : \"", click.IP, "\", \"IosIfa\" : \"", click.IosIfa, "\"}")
	}

}

var fromJson click

func main() {

	http.HandleFunc("/insert", insertinto)
	http.HandleFunc("/delete", deletefrom)
	http.HandleFunc("/retrieve", retrieve)
	http.ListenAndServe(":8000", nil)

	var byteslice []byte = []byte(jsonString)

	err := json.Unmarshal(byteslice, &fromJson)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", fromJson)
	fmt.Println(fromJson.Ip)

	returnValue, err := json.Marshal(&fromJson)

	fmt.Println(string(returnValue))

}

// new comment
