package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/olivere/elastic.v2"
)

type ClickUni struct {
	Id               string    //0
	Tracking_id      string    //1
	Publisher_ref_id string    //7
	Device_ip        string    //8
	User_agent       string    //16
	Request_url      string    //17
	Created          time.Time //18
	Modified         time.Time //19
	// Latitude         float64   //20
	// Longitude        float64   //21
	Location string // 20 21
}

type InstallUnq struct {
	Id               string
	Tracking_id      string
	Stat_click_id    string
	Session_ip       string
	Session_datetime time.Time
	Site_event_id    int64
	Publisher_ref_id string
	Device_ip        string
	User_agent       string
	Request_url      string
	Created          time.Time
	Modified         time.Time
	// Latitude         float64
	// Longitude        float64
	Location     string
	Match_type   string
	Install_date time.Time
}

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

type Combine struct {
	Publisher_id   int64  //2
	Ad_network_id  int64  //3
	Advertiser_id  int64  //4
	Site_id        int64  //5
	Campaign_id    int64  //6
	Sdk            string //9
	Device_carrier string //10
	Language       string //11
	Package_name   string //12
	App_name       string //13
	Country_id     int64  //14
	Region_id      int64  //15
	Clicks         ClickUni
	Installs       InstallUnq
}

func main() {
	//make a csv reader for that csv file
	csvfile, err := os.Open("/tmp/stat_clicks_1681.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = 22
	reader.LazyQuotes = true

	//read through the first line and then get all field keys
	_, err = reader.Read()
	if err != nil {
		fmt.Println(err)
	}

	//post click data and combine the install unqiue data if install.stat_click_id = click.id
	PostClicks(reader)
}

func PostClicks(reader *csv.Reader) {
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
		eachRecord, _ := reader.Read()
		var temp Combine
		temp.Publisher_id, _ = strconv.ParseInt(eachRecord[2], 10, 64)
		temp.Ad_network_id, _ = strconv.ParseInt(eachRecord[3], 10, 64)
		temp.Advertiser_id, _ = strconv.ParseInt(eachRecord[4], 10, 64)
		temp.Site_id, _ = strconv.ParseInt(eachRecord[5], 10, 64)
		temp.Campaign_id, _ = strconv.ParseInt(eachRecord[6], 10, 64)
		temp.Country_id, _ = strconv.ParseInt(eachRecord[14], 10, 64)
		temp.Region_id, _ = strconv.ParseInt(eachRecord[15], 10, 64)

		temp.Sdk = eachRecord[9]
		temp.Device_carrier = eachRecord[10]
		temp.Language = eachRecord[11]
		temp.Package_name = eachRecord[12]
		temp.App_name = eachRecord[13]

		var clicks ClickUni
		clicks.Id = eachRecord[0]
		clicks.Tracking_id = eachRecord[1]
		clicks.Publisher_ref_id = eachRecord[7]
		clicks.Device_ip = eachRecord[8]
		clicks.User_agent = eachRecord[16]
		clicks.Request_url = eachRecord[17]

		clicks.Created, _ = time.Parse("2006-01-02 15:04:05", eachRecord[18])
		clicks.Modified, _ = time.Parse("2006-01-02 15:04:05", eachRecord[19])

		// clicks.Latitude, _ = strconv.ParseFloat(eachRecord[20], 64)
		// clicks.Longitude, _ = strconv.ParseFloat(eachRecord[21], 64)
		if eachRecord[20] != "NULL" && eachRecord[21] != "NULL" {

			clicks.Location = eachRecord[20] + "," + eachRecord[21]
		}

		temp.Clicks = clicks

		//search with a term query
		termQuery := elastic.NewTermQuery("Installs.Stat_click_id", clicks.Id)
		searchResult, err := client.Search().
			Index("alldata").
			Query(&termQuery).
			From(0).Size(10). // take documents 0-9
			Pretty(true).     // pretty print request and response JSON
			Do()              // execute
		if err != nil {
			fmt.Println(err)
		}

		//process if the hit is only one
		if searchResult.TotalHits() == 0 {
			_, err = client.Index().
				Index("alldata").
				Type("clinstall").
				BodyJson(temp).
				Do()
			if err != nil {
				fmt.Println(err)
			}
		} else if searchResult.TotalHits() != 1 {
			fmt.Println("something unexpected happend")
		} else {
			//loop through all hits but we only want one this time
			for _, hit := range searchResult.Hits.Hits {

				//marshal to be a install struct so we can get their unique datas and combine with the clicks data
				var install Install
				err := json.Unmarshal(*hit.Source, &install)
				if err != nil {
					fmt.Println(err)
				}

				temp.Installs = install.Installs // combine the install data to the clicks data if they are related

				//find out the id of that hit in the index
				docId := hit.Id
				fmt.Println(docId)

				//update to be the combination of click and install
				client.Update().Index("alldata").Type("clinstall").Id(docId).Doc(temp).Do()

				break
			}
		}

		//break
	}

}
