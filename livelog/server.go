package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/olivere/elastic.v2"
)

var client *elastic.Client

func init() {
	var err error
	// Create a client
	client, err = elastic.NewClient()
	if err != nil {
		fmt.Println(err)
	}

	// Create an index
	_, err = client.CreateIndex("alls").Do()
	if err != nil {
		fmt.Println(err)
	}
}

// build/lauch the server and prepare to write logs
func main() {
	server := http.Server{
		Addr:    ":5000",
		Handler: myHandler(),
	}

	server.ListenAndServe()
}

//build and return the server's handler
func myHandler() *mux.Router {
	mx := mux.NewRouter()
	mx.HandleFunc("/", POST).Methods("POST")
	//mx.HandleFunc("/{id}", GET).Methods("GET")
	return mx
}

func POST(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	var allFields AllFields
	json.Unmarshal(bytes, &allFields)
	logType := allFields.LogType
	if logType == "impress" {
		//impression(allFields)
	} else if logType == "click" {
		click(allFields)
	} else if logType == "install" {
		//install(allFields)
	} else if logType == "event" {
		//event(allFields)
	} else if logType == "open" {
		//open(allFields)
	}

}

//if click is contributed to an install, no need to care about event or open
//but if the click is not finding the install yet, consider about the events or opens
//do in opposite way, if install is already related to the clicks no need to consider aboutothers
func click(allFields AllFields) {
	//set up unqiue fields first
	var clickUni Click
	clickUni.Id = allFields.Id
	clickUni.Created = allFields.Created
	clickUni.DeviceIp = allFields.DeviceIp
	clickUni.CountryCode = allFields.CountryCode
	clickUni.RegionCode = allFields.RegionCode
	clickUni.PostalCode = allFields.PostalCode
	clickUni.Location = strconv.FormatFloat(allFields.Latitude, 'E', -1, 64) + "," + strconv.FormatFloat(allFields.Longitude, 'E', -1, 64)
	clickUni.WurflBrandName = allFields.WurflBrandName
	clickUni.WurflModelName = allFields.WurflModelName
	clickUni.WurflDeviceOs = allFields.WurflDeviceOs

	//search for realationship with install first
	termQuery := elastic.NewTermQuery("StatInstallId", clickUni.Id)
	searchResult, _ := client.Search().Index("alls").Type("Common").Query(&termQuery).From(0).Size(2).Pretty(true).Do()
	if searchResult.TotalHits() != 0 {
		hits := searchResult.Hits.Hits
		firstHit := hits[0]
		parentId := firstHit.Id

		//update the parent doc's clickInstallId -----the click id which realted to that install
		var clickID UpdateClickId
		clickID.ClickInstallId = clickUni.Id
		client.Update().Index("alls").Type("Common").Id(parentId).Doc(clickID).Do()

		//index the click doc as child doc to this parent id
		client.Index().Index("alls").Type("Click").Parent(parentId).BodyJson(clickUni).Do()

	} else { // if find no install, find the event /open which may related to this click

		termQuery = elastic.NewTermQuery("StatClickId", clickUni.Id)
		searchEvent, _ := client.Search().Index("alls").Type("Event").Query(&termQuery).Pretty(true).Do()
		for _, hit := range searchEvent.Hits.Hits {
			//update the events to have reengaement  click struct

			var updateclick UpdateClick
			updateclick.Click = clickUni
			docId := hit.Id
			client.Update().Index("alls").Type("Event").Id(docId).Doc(updateclick).Do()

		}

		//update the opens to have this reengaement click
		searchOpen, _ := client.Search().Index("alls").Type("Open").Query(&termQuery).Pretty(true).Do()
		for _, hit := range searchOpen.Hits.Hits {
			var updateclick UpdateClick
			updateclick.Click = clickUni
			docId := hit.Id
			client.Update().Index("alls").Type("Open").Id(docId).Doc(updateclick).Do()
		}

		//if there is no install/event/open related to this click, post it lonely
		if searchEvent.TotalHits() == 0 && searchOpen.TotalHits() == 0 {
			//post the parent first
			var parent Common
			parent.GoogleAid = allFields.GoogleAid
			parent.WindowsAid = allFields.WindowsAid
			parent.IosIfa = allFields.IosIfa
			parent.Language = allFields.Language
			parent.CurrencyCode = allFields.CurrencyCode
			parent.SiteId = allFields.SiteId
			parent.AdvertiserId = allFields.AdvertiserId
			parent.PackageName = allFields.PackageName
			parent.PublisherId = allFields.PublisherId
			parent.AdNetworkId = allFields.AdNetworkId
			parent.AgencyId = allFields.AgencyId
			parent.CampaignId = allFields.CampaignId
			parent.PublisherUserId = allFields.PublisherUserId
			parent.StatInstallId = allFields.StatInstallId
			//the click id which future install  may contributed to
			parent.ClickInstallId = clickUni.Id

			client.Index().Index("alls").Type("Common").BodyJson(parent).Do()

			//find this parent id and then post the rest click to be that children
			termQuery := elastic.NewTermQuery("StatInstallId", clickUni.Id)
			searchResult, _ := client.Search().Index("alls").Type("Common").Query(&termQuery).From(0).Size(2).Pretty(true).Do()
			parentId := searchResult.Hits.Hits[0].Id

			client.Index().Index("alls").Type("Click").Parent(parentId).BodyJson(clickUni).Do()
		}
	}

}

//for update the click inside the events/opens
type UpdateClick struct {
	Click Click
}

//for update the impression inside the events/opens
type UpdateImpression struct {
	Impression Impression
}

//for update the install id inside the common parent fileds
type UpdateInstallId struct {
	StatInstallId string
}

//for update the click_id inside the common parent fileds
type UpdateClickId struct {
	ClickInstallId string
}

type Common struct {
	GoogleAid       string
	WindowsAid      string
	IosIfa          string
	Language        string
	CurrencyCode    string
	SiteId          int64
	AdvertiserId    int64
	PackageName     string
	PublisherId     int64
	AdNetworkId     int64
	AgencyId        int64
	CampaignId      int64
	PublisherUserId string
	StatInstallId   string // the install related to all events /open/clicks
	ClickInstallId  string //The click id related to that install
}

type Impression struct {
	Id             string
	Created        time.Time
	DeviceIp       string
	CountryCode    string
	RegionCode     string
	PostalCode     int32
	Location       string // string(float64 lat) , string(float64 log)
	WurflBrandName string
	WurflModelName string
	WurflDeviceOs  string
}

type Click struct {
	Id             string
	Created        time.Time
	DeviceIp       string
	CountryCode    string
	RegionCode     string
	PostalCode     int32
	Location       string
	WurflBrandName string
	WurflModelName string
	WurflDeviceOs  string
}

type Install struct {
	Id               string
	Created          time.Time
	DeviceIp         string
	StatImpressionId string
	StatClickId      string
	CountryCode      string
	RegionCode       string
	PostalCode       int32
	Location         string
	WurflBrandName   string
	WurflModelName   string
	WurflDeviceOs    string
}

type Open struct {
	Id               string
	Created          time.Time
	DeviceIp         string
	StatImpressionId string
	StatClickId      string
	StatInstallId    string
	CountryCode      string
	RegionCode       string
	PostalCode       int32
	Location         string
	WurflBrandName   string
	WurflModelName   string
	WurflDeviceOs    string
	Click            Click
	Impression       Impression
}

type Event struct {
	Id               string
	Created          time.Time
	DeviceIp         string
	StatImpressionId string
	StatClickId      string
	StatInstallId    string
	StatOpenId       string
	CountryCode      string
	RegionCode       string
	PostalCode       int32
	Location         string
	WurflBrandName   string
	WurflModelName   string
	WurflDeviceOs    string
	Click            Click
	Impression       Impression
}

type AllFields struct {
	LogType          string
	Id               string
	Created          time.Time
	DeviceIp         string
	GoogleAid        string
	WindowsAid       string
	IosIfa           string
	Language         string
	StatInstallId    string
	StatOpenId       string
	StatClickId      string
	StatImpressionId string
	CurrencyCode     string
	SiteId           int64
	AdvertiserId     int64
	PackageName      string
	PublisherId      int64
	AdNetworkId      int64
	AgencyId         int64
	CampaignId       int64
	CountryCode      string
	RegionCode       string
	PostalCode       int32
	WurflBrandName   string
	WurflModelName   string
	WurflDeviceOs    string
	PublisherUserId  string
	Latitude         float64
	Longitude        float64
}
