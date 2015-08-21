package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/MobileAppTracking/measurement/lib/structured"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Document struct {
	Common     CommonInfoStr
	Click      ClickStr
	Impression ImpressionStr
	Install    InstallStr
	Events     EventStr
	Opens      OpenStr
}

type CommonInfoStr struct {
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
}

type ImpressionStr struct {
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

type ClickStr struct {
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

type InstallStr struct {
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

type OpenStr struct {
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
	ClickInOpen      ClickStr
}

type EventStr struct {
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
	ClickInEvent     ClickStr
}

type AllFieldsStr struct {
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

func init() {

	// initialize structured logging
	structured.AddHookToSyslog("tcp", "localhost:10514", syslog.LOG_EMERG, "mini---project")
	structured.AddHookToElasticsearch("localhost", "9200", "clients", "user", "")

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
	db.Close()
}

//build and return the server's handler
func myHandler() *mux.Router {
	mx := mux.NewRouter()
	mx.HandleFunc("/", POST).Methods("POST")
	mx.HandleFunc("/{id}", GET).Methods("GET")
	return mx
}

type PostResponses struct {
	ErrMessage string
	Id         string
	HttpStatus string
}

//the function which handle the post method
//post the json data from broswer to the server and sql databases
func POST(w http.ResponseWriter, r *http.Request) {
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
	var temp AllFieldsStr
	errs = json.Unmarshal(bytes, &temp)
	if errs != nil {
		errString := fmt.Sprintf("invalid Json format: %s", errs)
		response(w, errString, "", http.StatusBadRequest)
		structured.Error("", "", errString, 0, nil)
		return
	}

	//validate the input and log error input message
	if point.AdvertiserID == 0 || point.SiteID == 0 {
		errString := "your advertiserID or site ID may equals to 0"
		response(w, errString, "", http.StatusBadRequest)
		structured.Error("", "", errString, 0, nil)
		return
	}

	//sucess and log the request latency
	response(w, "", id, http.StatusOK)
	structured.Info(point.ID, point.Type, "Post successful!", point.SiteID,
		structured.ExtraFields{structured.RequestLatency: time.Since(RequestStart),
			structured.QueryLatency: time.Since(QueryStart)})
}

func UpdateCommonData(originDoc Document, clickData AllFieldsStr) Document {

	if originDoc.Common.AdNetworkId == 0 {
		originDoc.Common.AdNetworkId = clickData.AdNetworkId
	}
	if originDoc.Common.AdvertiserId == 0 {
		originDoc.Common.AdvertiserId = clickData.AdvertiserId
	}
	if originDoc.Common.AgencyId == 0 {
		originDoc.Common.AgencyId = clickData.AgencyId
	}
	if originDoc.Common.CampaignId == 0 {
		originDoc.Common.CampaignId = clickData.CampaignId
	}
	if originDoc.Common.CurrencyCode == "" {
		originDoc.Common.CurrencyCode = clickData.CurrencyCode
	}
	if originDoc.Common.GoogleAid == "" {
		originDoc.Common.GoogleAid = clickData.GoogleAid
	}
	if originDoc.Common.IosIfa == "" {
		originDoc.Common.IosIfa = clickData.IosIfa
	}
	if originDoc.Common.Language == "" {
		originDoc.Common.Language = clickData.Language
	}
	if originDoc.Common.PackageName == "" {
		originDoc.Common.PackageName = clickData.PackageName
	}
	if originDoc.Common.PublisherId == 0 {
		originDoc.Common.PublisherId = clickData.PublisherId
	}
	if originDoc.Common.PublisherUserId == "" {
		originDoc.Common.PublisherUserId = clickData.PublisherUserId
	}
	if originDoc.Common.SiteId == 0 {
		originDoc.Common.SiteId = clickData.SiteId
	}
	if originDoc.Common.WindowsAid == "" {
		originDoc.Common.WindowsAid = clickData.WindowsAid
	}

	return originDoc
}

func AddClickToOpen(originDoc Document, clickData AllFieldsStr) Document {

	originDoc.Opens.ClickInOpen.CountryCode = clickData.CountryCode
	originDoc.Opens.ClickInOpen.Created = clickData.Created
	originDoc.Opens.ClickInOpen.DeviceIp = clickData.DeviceIp
	originDoc.Opens.ClickInOpen.Id = clickData.Id
	originDoc.Opens.ClickInOpen.Location = location
	originDoc.Opens.ClickInOpen.PostalCode = clickData.PostalCode
	originDoc.Opens.ClickInOpen.RegionCode = clickData.RegionCode
	originDoc.Opens.ClickInOpen.WurflBrandName = clickData.WurflBrandName
	originDoc.Opens.ClickInOpen.WurflDeviceOs = clickData.WurflDeviceOs
	originDoc.Opens.ClickInOpen.WurflModelName = clickData.WurflModelName

	return originDoc
}

func AddClickToEvent(originDoc Document, clickData AllFieldsStr) Document {

	originDoc.Events.ClickInEvent.CountryCode = clickData.CountryCode
	originDoc.Events.ClickInEvent.Created = clickData.Created
	originDoc.Events.ClickInEvent.DeviceIp = clickData.DeviceIp
	originDoc.Events.ClickInEvent.Id = clickData.Id
	originDoc.Events.ClickInEvent.Location = location
	originDoc.Events.ClickInEvent.PostalCode = clickData.PostalCode
	originDoc.Events.ClickInEvent.RegionCode = clickData.RegionCode
	originDoc.Events.ClickInEvent.WurflBrandName = clickData.WurflBrandName
	originDoc.Events.ClickInEvent.WurflDeviceOs = clickData.WurflDeviceOs
	originDoc.Events.ClickInEvent.WurflModelName = clickData.WurflModelName

	return originDoc
}

func AddImpression(originDoc Document, impressionData AllFieldsStr) Document {

	originDoc.Impression.CountryCode = impressionData.CountryCode
	originDoc.Impression.Created = impressionData.Created
	originDoc.Impression.DeviceIp = impressionData.DeviceIp
	originDoc.Impression.Id = impressionData.Id
	originDoc.Impression.Location = location
	originDoc.Impression.PostalCode = impressionData.PostalCode
	originDoc.Impression.RegionCode = impressionData.RegionCode
	originDoc.Impression.WurflBrandName = impressionData.WurflBrandName
	originDoc.Impression.WurflDeviceOs = impressionData.WurflDeviceOs
	originDoc.Impression.WurflModelName = impressionData.WurflModelName

	return originDoc
}

func AddClick(originDoc Document, clickData AllFieldsStr) Document {

	originDoc.Click.CountryCode = clickData.CountryCode
	originDoc.Click.Created = clickData.Created
	originDoc.Click.DeviceIp = clickData.DeviceIp
	originDoc.Click.Id = clickData.Id
	originDoc.Click.Location = location
	originDoc.Click.PostalCode = clickData.PostalCode
	originDoc.Click.RegionCode = clickData.RegionCode
	originDoc.Click.WurflBrandName = clickData.WurflBrandName
	originDoc.Click.WurflDeviceOs = clickData.WurflDeviceOs
	originDoc.Click.WurflModelName = clickData.WurflModelName

	return originDoc
}

func AddInstall(originDoc Document, installData AllFieldsStr) {

	originDoc.Install.CountryCode = installData.CountryCode
	originDoc.Install.Created = installData.Created
	originDoc.Install.DeviceIp = installData.DeviceIp
	originDoc.Install.Id = installData.Id
	originDoc.Install.Location = location
	originDoc.Install.PostalCode = installData.PostalCode
	originDoc.Install.RegionCode = installData.RegionCode
	originDoc.Install.StatClickId = installData.StatClickId
	originDoc.Install.StatImpressionId = installData.StatImpressionId
	originDoc.Install.WurflBrandName = installData.WurflBrandName
	originDoc.Install.WurflDeviceOs = installData.WurflDeviceOs
	originDoc.Install.WurflModelName = installData.WurflModelName

}

func AddOpen(originDoc Document, clickData AllFieldsStr) {

	originDoc.Opens.CountryCode = clickData.CountryCode
	originDoc.Opens.Created = clickData.Created
	originDoc.Opens.DeviceIp = clickData.DeviceIp
	originDoc.Opens.Id = clickData.Id
	originDoc.Opens.Location = location
	originDoc.Opens.PostalCode = clickData.PostalCode
	originDoc.Opens.RegionCode = clickData.RegionCode
	originDoc.Opens.StatImpressionId = clickData.StatImpressionId
	originDoc.Opens.StatClickId = clickData.StatClickId
	originDoc.Opens.StatInstallId = clickData.StatInstallId
	originDoc.Opens.WurflBrandName = clickData.WurflBrandName
	originDoc.Opens.WurflDeviceOs = clickData.WurflDeviceOs
	originDoc.Opens.WurflModelName = clickData.WurflModelName

	return originDoc
}

func NewImpression(data AllFieldsStr) Document {

	var impression Document
	var location string

	location = "\"" + data.Latitude + "," + data.Longitude + "\""

	impression.Common.AdNetworkId = data.AdNetworkId
	impression.Common.AdvertiserId = data.AdvertiserId
	impression.Common.AgencyId = data.AgencyId
	impression.Common.CampaignId = data.CampaignId
	impression.Common.CurrencyCode = data.CurrencyCode
	impression.Common.GoogleAid = data.GoogleAid
	impression.Common.IosIfa = data.IosIfa
	impression.Common.Language = data.Language
	impression.Common.PackageName = data.PackageName
	impression.Common.PublisherId = data.PublisherId
	impression.Common.PublisherUserId = data.PublisherUserId
	impression.Common.SiteId = data.SiteId
	impression.Common.WindowsAid = data.WindowsAid

	impression.Impression.CountryCode = data.CountryCode
	impression.Impression.Created = data.Created
	impression.Impression.DeviceIp = data.DeviceIp
	impression.Impression.Id = data.Id
	impression.Impression.Location = location
	impression.Impression.PostalCode = data.PostalCode
	impression.Impression.RegionCode = data.RegionCode
	impression.Impression.WurflBrandName = data.WurflBrandName
	impression.Impression.WurflDeviceOs = data.WurflDeviceOs
	impression.Impression.WurflModelName = data.WurflModelName

	return impression

}

func NewClick(data AllFieldsStr) Document {

	var click Document
	var location string

	location = "\"" + data.Latitude + "," + data.Longitude + "\""

	click.Common.AdNetworkId = data.AdNetworkId
	click.Common.AdvertiserId = data.AdvertiserId
	click.Common.AgencyId = data.AgencyId
	click.Common.CampaignId = data.CampaignId
	click.Common.CurrencyCode = data.CurrencyCode
	click.Common.GoogleAid = data.GoogleAid
	click.Common.IosIfa = data.IosIfa
	click.Common.Language = data.Language
	click.Common.PackageName = data.PackageName
	click.Common.PublisherId = data.PublisherId
	click.Common.PublisherUserId = data.PublisherUserId
	click.Common.SiteId = data.SiteId
	click.Common.WindowsAid = data.WindowsAid

	click.Click.CountryCode = data.CountryCode
	click.Click.Created = data.Created
	click.Click.DeviceIp = data.DeviceIp
	click.Click.Id = data.Id
	click.Click.Location = location
	click.Click.PostalCode = data.PostalCode
	click.Click.RegionCode = data.RegionCode
	click.Click.WurflBrandName = data.WurflBrandName
	click.Click.WurflDeviceOs = data.WurflDeviceOs
	click.Click.WurflModelName = data.WurflModelName

	return click
}

func NewInstall(data AllFieldsStr) Document {

	var install Document
	var location string

	location = "\"" + data.Latitude + "," + data.Longitude + "\""

	install.Common.AdNetworkId = data.AdNetworkId
	install.Common.AdvertiserId = data.AdvertiserId
	install.Common.AgencyId = data.AgencyId
	install.Common.CampaignId = data.CampaignId
	install.Common.CurrencyCode = data.CurrencyCode
	install.Common.GoogleAid = data.GoogleAid
	install.Common.IosIfa = data.IosIfa
	install.Common.Language = data.Language
	install.Common.PackageName = data.PackageName
	install.Common.PublisherId = data.PublisherId
	install.Common.PublisherUserId = data.PublisherUserId
	install.Common.SiteId = data.SiteId
	install.Common.WindowsAid = data.WindowsAid

	install.Install.CountryCode = data.CountryCode
	install.Install.Created = data.Created
	install.Install.DeviceIp = data.DeviceIp
	install.Install.Id = data.Id
	install.Install.Location = location
	install.Install.PostalCode = data.PostalCode
	install.Install.RegionCode = data.RegionCode
	install.Install.StatClickId = data.StatClickId
	install.Install.StatImpressionId = data.StatImpressionId
	install.Install.WurflBrandName = data.WurflBrandName
	install.Install.WurflDeviceOs = data.WurflDeviceOs
	install.Install.WurflModelName = data.WurflModelName

	return install

}

func NewOpen(data AllFieldsStr) Document {

	var open Document
	var location string

	location = "\"" + data.Latitude + "," + data.Longitude + "\""

	open.Common.AdNetworkId = data.AdNetworkId
	open.Common.AdvertiserId = data.AdvertiserId
	open.Common.AgencyId = data.AgencyId
	open.Common.CampaignId = data.CampaignId
	open.Common.CurrencyCode = data.CurrencyCode
	open.Common.GoogleAid = data.GoogleAid
	open.Common.IosIfa = data.IosIfa
	open.Common.Language = data.Language
	open.Common.PackageName = data.PackageName
	open.Common.PublisherId = data.PublisherId
	open.Common.PublisherUserId = data.PublisherUserId
	open.Common.SiteId = data.SiteId
	open.Common.WindowsAid = data.WindowsAid

	open.Opens.CountryCode = data.CountryCode
	open.Opens.Created = data.Created
	open.Opens.DeviceIp = data.DeviceIp
	open.Opens.Id = data.Id
	open.Opens.Location = location
	open.Opens.PostalCode = data.PostalCode
	open.Opens.RegionCode = data.RegionCode
	open.Opens.StatImpressionId = data.StatImpressionId
	open.Opens.StatClickId = data.StatClickId
	open.Opens.StatInstallId = data.StatInstallId
	open.Opens.WurflBrandName = data.WurflBrandName
	open.Opens.WurflDeviceOs = data.WurflDeviceOs
	open.Opens.WurflModelName = data.WurflModelName

	return open

}

func NewEvent(data AllFieldsStr) Document {

	var event Document
	var location string

	location = "\"" + data.Latitude + "," + data.Longitude + "\""

	event.Common.AdNetworkId = data.AdNetworkId
	event.Common.AdvertiserId = data.AdvertiserId
	event.Common.AgencyId = data.AgencyId
	event.Common.CampaignId = data.CampaignId
	event.Common.CurrencyCode = data.CurrencyCode
	event.Common.GoogleAid = data.GoogleAid
	event.Common.IosIfa = data.IosIfa
	event.Common.Language = data.Language
	event.Common.PackageName = data.PackageName
	event.Common.PublisherId = data.PublisherId
	event.Common.PublisherUserId = data.PublisherUserId
	event.Common.SiteId = data.SiteId
	event.Common.WindowsAid = data.WindowsAid

	event.Events.CountryCode = data.CountryCode
	event.Events.Created = data.Created
	event.Events.DeviceIp = data.DeviceIp
	event.Events.Id = data.Id
	event.Events.Location = location
	event.Events.PostalCode = data.PostalCode
	event.Events.RegionCode = data.RegionCode
	event.Events.StatImpressionId = data.StatImpressionId
	event.Events.StatClickId = data.StatClickId
	event.Events.StatOpenId = data.StatOpenId
	event.Events.StatInstallId = data.StatInstallId
	event.Events.WurflBrandName = data.WurflBrandName
	event.Events.WurflDeviceOs = data.WurflDeviceOs
	event.Events.WurflModelName = data.WurflModelName

	return event

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
