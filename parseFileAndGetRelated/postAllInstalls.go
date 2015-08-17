package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"gopkg.in/olivere/elastic.v2"
)

// var client *elastic.Client

// func init() {
// 	// Create a client
// 	client, err := elastic.NewClient()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Create an index
// 	_, err = client.CreateIndex("install").Do()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("hehe")
// }

type Install struct {
	Publisher_id   int64  //5
	Ad_network_id  int64  //6
	Advertiser_id  int64  //7
	Site_id        int64  //8
	Campaign_id    int64  //9
	Sdk            string //13
	Device_carrier string //14
	Language       string //15
	Package_name   string //16
	App_name       string //17
	Country_id     int64  //18
	Region_id      int64  //19
	Installs       InstallUnq
}

type InstallUnq struct {
	Id               string    //0
	Tracking_id      string    //1
	Stat_click_id    string    //2
	Session_ip       string    //3
	Session_datetime time.Time //4
	Site_event_id    int64     //10
	Publisher_ref_id string    //11
	Device_ip        string    //12
	User_agent       string    //20
	Request_url      string    //21
	Created          time.Time //22
	Modified         time.Time //23
	// Latitude         float64   //24
	// Longitude        float64   //25
	Location     string    // 24 25
	Match_type   string    //26
	Install_date time.Time //27
}

func main() {
	//make a csv reader for that csv file
	csvfile, err := os.Open("/tmp/stat_installs_1681.csv")
	//csvfile, err := os.Open("sample20data.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = 28
	reader.LazyQuotes = true

	//read through the first line and then get all field keys
	_, err = reader.Read()
	if err != nil {
		fmt.Println(err)
	}

	//post all installs data to ES first
	PostInstalls(reader)

}

func PostInstalls(reader *csv.Reader) {
	// Create a client
	client, err := elastic.NewClient()
	if err != nil {
		fmt.Println(err)
	}

	// Create an index
	_, err = client.CreateIndex("alldata").Do()
	if err != nil {
		fmt.Println(err)
	}

	//loop through each record
	for {
		eachRecord, err := reader.Read()
		if err == nil {
			//build the struct al special case can be handled by struct
			var post Install
			post.Publisher_id, _ = strconv.ParseInt(eachRecord[5], 10, 64)
			post.Ad_network_id, _ = strconv.ParseInt(eachRecord[6], 10, 64)
			post.Advertiser_id, _ = strconv.ParseInt(eachRecord[7], 10, 64)
			post.Site_id, _ = strconv.ParseInt(eachRecord[8], 10, 64)
			post.Campaign_id, _ = strconv.ParseInt(eachRecord[9], 10, 64)
			post.Country_id, _ = strconv.ParseInt(eachRecord[18], 10, 64)
			post.Region_id, _ = strconv.ParseInt(eachRecord[19], 10, 64)

			post.Sdk = eachRecord[13]
			post.Device_carrier = eachRecord[14]
			post.Language = eachRecord[15]
			post.Package_name = eachRecord[16]
			post.App_name = eachRecord[17]

			var uniques InstallUnq
			uniques.Id = eachRecord[0]
			uniques.Tracking_id = eachRecord[1]
			uniques.Stat_click_id = eachRecord[2]
			uniques.Session_ip = eachRecord[3]
			uniques.Publisher_ref_id = eachRecord[11]
			uniques.Device_ip = eachRecord[12]
			uniques.User_agent = eachRecord[20]
			uniques.Request_url = eachRecord[21]
			uniques.Match_type = eachRecord[26]

			uniques.Site_event_id, _ = strconv.ParseInt(eachRecord[10], 10, 64)
			// uniques.Latitude, _ = strconv.ParseFloat(eachRecord[24], 64)
			// uniques.Longitude, _ = strconv.ParseFloat(eachRecord[25], 64)
			if eachRecord[24] != "NULL" && eachRecord[25] != "NULL" {

				uniques.Location = eachRecord[24] + "," + eachRecord[25]
			}

			uniques.Session_datetime, _ = time.Parse("2006-01-02 15:04:05", eachRecord[4])
			uniques.Created, _ = time.Parse("2006-01-02 15:04:05", eachRecord[22])
			uniques.Modified, _ = time.Parse("2006-01-02 15:04:05", eachRecord[23])
			uniques.Install_date, _ = time.Parse("2006-01-02 15:04:05", eachRecord[27])

			post.Installs = uniques

			//index the document

			_, err = client.Index().
				Index("alldata").
				Type("clinstall").
				BodyJson(post).
				Do()
			if err != nil {
				fmt.Println(err)
			}

		} else if err != nil && err != io.EOF {
			fmt.Println(err)
		} else {
			break
		}
	}
}
