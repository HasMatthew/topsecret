//NOTE:
// assume event happen after the install, if before the install, leave it in ES with no relationship

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MobileAppTracking/measurement/lib/structured"
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

	//log message to rsyslog file
	structured.AddHookToSyslog("tcp", "localhost:10514", syslog.LOG_EMERG, "live===log")
	//log message to os standard
	structured.SetOutput(os.Stderr)
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
	//read the raw post data nad marhsal it s
	bytes, _ := ioutil.ReadAll(r.Body)
	var allFields AllFields
	json.Unmarshal(bytes, &allFields)

	//get the time sturct work right
	if allFields.TempTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", allFields.TempTime)
		if err != nil {
			structured.Warn("", "", "can't parse the time", 0, nil)
		} else {
			allFields.Created = t
		}
	}

	//find the log type
	logType := allFields.LogType

	if logType == "impress" {
		impression(allFields)
	} else if logType == "click" {
		click(allFields)
	} else if logType == "install" {
		install(allFields)
	} else if logType == "event" {
		event(allFields)
	} else if logType == "open" {
		open(allFields)
	}

}

//install----X----event/open--------alone

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

	//search for realationship with install first ******click id aleays exist
	termQuery := elastic.NewTermQuery("ClickInstallId", clickUni.Id)
	searchResult, _ := client.Search().Index("alls").Type("Common").Query(&termQuery).Pretty(true).Do()

	//******************** find the only install
	if searchResult.TotalHits() == 1 {
		hits := searchResult.Hits.Hits

		firstHit := hits[0]
		parentId := firstHit.Id

		//update the parent doc's clickInstallId -----the click id which realted to that install
		var clickID UpdateClickId
		clickID.ClickInstallId = clickUni.Id
		client.Update().Index("alls").Type("Common").Id(parentId).Doc(clickID).Do()

		//index the click doc as child doc to this parent id
		client.Index().Index("alls").Type("Click").Parent(parentId).BodyJson(clickUni).Do()

	} else if searchResult.TotalHits() == 0 {
		//************find no install--- find the event and open first *****click id always exists

		termQuery = elastic.NewTermQuery("StatClickId", clickUni.Id)

		//update the events to have reengaement  click struct
		searchEvent, _ := client.Search().Index("alls").Type("Event").Query(&termQuery).Pretty(true).Do()
		for _, hit := range searchEvent.Hits.Hits {
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

		//************if there is no install and no matter whether it found a event/open, post the click

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
		//the click id which future install  may contributed to
		parent.ClickInstallId = clickUni.Id

		//index the parent common field first
		indexParent, _ := client.Index().Index("alls").Type("Common").BodyJson(parent).Do()
		parentId := indexParent.Id
		//index the click as the children
		client.Index().Index("alls").Type("Click").Parent(parentId).BodyJson(clickUni).Do()

	} else {
		//*************find multiple install related to this click
		structured.Warn(clickUni.Id, "Click", "the click has multiple install contributed to it ", int(allFields.SiteId), nil)
	}

}

func impression(allFields AllFields) {

}

//click/impression----X---alone
//assume all events happen after install
func install(allFields AllFields) {
	//************build unique install struct first
	var installUni Install
	installUni.Id = allFields.Id
	installUni.Created = allFields.Created
	installUni.DeviceIp = allFields.DeviceIp
	installUni.StatImpressionId = allFields.StatImpressionId
	installUni.StatClickId = allFields.StatClickId
	installUni.CountryCode = allFields.CountryCode
	installUni.RegionCode = allFields.RegionCode
	installUni.PostalCode = allFields.PostalCode
	installUni.Location = strconv.FormatFloat(allFields.Latitude, 'E', -1, 64) + "," + strconv.FormatFloat(allFields.Longitude, 'E', -1, 64)
	installUni.WurflBrandName = allFields.WurflBrandName
	installUni.WurflModelName = allFields.WurflModelName
	installUni.WurflDeviceOs = allFields.WurflDeviceOs

	noHit := true
	//*************************search for realtionship with click  by using Common field clickinsatllid if provide statclickid
	if installUni.StatClickId != "" {

		termQuery := elastic.NewTermQuery("ClickInstallId", installUni.StatClickId)
		searchResult, _ := client.Search().Index("alls").Type("Common").Query(&termQuery).Pretty(true).Do()
		if searchResult.TotalHits() != 0 {

			if searchResult.TotalHits() != 1 {
				structured.Warn(installUni.Id, "install", "the click is not only one when related to this install", int(allFields.SiteId), nil)
			} else {
				noHit = false
				hit := searchResult.Hits.Hits[0]
				parentId := hit.Id

				//update the statinstallid in common field
				var update UpdateInstallId
				update.StatInstallId = installUni.Id
				client.Update().Index("alls").Type("Common").Id(parentId).Doc(update).Do()

				//post install data as child of this parent doc
				client.Index().Index("alls").Type("Install").Parent(parentId).BodyJson(installUni).Do()
			}
		}
	} else if installUni.StatImpressionId != "" {

		//**********************search for relationship with impression bY using common field impressioninstallid if provide statimpression id
		termQuery := elastic.NewTermQuery("ImpressionInstallId", installUni.StatImpressionId)
		searchResult, _ := client.Search().Index("alls").Type("Common").Query(&termQuery).Pretty(true).Do()

		if searchResult.TotalHits() != 0 {
			if searchResult.TotalHits() != 1 {
				structured.Warn(installUni.Id, "installs", "this install has not onlt one impression related to it", int(allFields.SiteId), nil)
			} else {
				noHit = false
				hit := searchResult.Hits.Hits[0]
				parentId := hit.Id

				//update the statInstallId in common field
				var update UpdateInstallId
				update.StatInstallId = installUni.Id
				client.Update().Index("alls").Type("Common").Id(parentId).Doc(update).Do()

				//post install data as child of this parent doc
				client.Index().Index("alls").Type("Install").Parent(parentId).BodyJson(installUni).Do()
			}
		}
	}

	//***********post the install data alone when it find no realtionship with existing ***outside click or impression
	// if find no realtionship with click/ impression, 1. this install may not have click/imrpession id 2. the click/impression id is not posted yet

	if noHit {

		//post the common field first
		var common Common
		common.GoogleAid = allFields.GoogleAid
		common.WindowsAid = allFields.WindowsAid
		common.IosIfa = allFields.IosIfa
		common.Language = allFields.Language
		common.CurrencyCode = allFields.CurrencyCode
		common.SiteId = allFields.SiteId
		common.AdvertiserId = allFields.AdvertiserId
		common.PackageName = allFields.PackageName
		common.PublisherId = allFields.PublisherId
		common.AdNetworkId = allFields.AdNetworkId
		common.AgencyId = allFields.AgencyId
		common.CampaignId = allFields.CampaignId
		common.PublisherUserId = allFields.PublisherUserId
		//stat install id is the install's id
		common.StatInstallId = installUni.Id

		/********this is for avoiding the speical ---- event happen before install, still have click id even the click is in inner struct of opens/events
		********the trade off is not storing the related click outside but it would only happen when the first event is logging before install and click and also click
		********is logging after install****************/
		common.ClickInstallId = installUni.StatClickId
		common.ImpressionInstallId = installUni.StatImpressionId

		//post the parent data frist and then post the install as a child
		indexParent, _ := client.Index().Index("alls").Type("Common").BodyJson(common).Refresh(true).Do()
		parentId := indexParent.Id
		client.Index().Index("alls").Type("Install").Parent(parentId).BodyJson(installUni).Pretty(true).Do()
	}

}

func event(allFields AllFields) {
	var event Event
	event.Id = allFields.Id
	event.Created = allFields.Created
	event.DeviceIp = allFields.DeviceIp
	event.StatImpressionId = allFields.StatImpressionId
	event.StatClickId = allFields.StatClickId
	event.StatInstallId = allFields.StatInstallId
	event.StatOpenId = allFields.StatOpenId
	event.CountryCode = allFields.CountryCode
	event.RegionCode = allFields.RegionCode
	event.PostalCode = allFields.PostalCode
	event.Location = strconv.FormatFloat(allFields.Latitude, 'E', -1, 64) + "," + strconv.FormatFloat(allFields.Longitude, 'E', -1, 64)
	event.WurflBrandName = allFields.WurflBrandName
	event.WurflDeviceOs = allFields.WurflDeviceOs
	event.WurflModelName = allFields.WurflModelName

	//************find whether it related to the install first
	noHit := true

	if event.StatInstallId != "" {
		termQuery := elastic.NewTermQuery("StatInstallId", event.StatInstallId)
		searchResult, _ := client.Search().Index("alls").Type("Common").Query(&termQuery).Pretty(true).Do()

		if searchResult.TotalHits() == 1 {
			//already find that only install related
			noHit = false
			hit := searchResult.Hits.Hits[0]
			parentId := hit.Id

			//unmarhsal the parent
			var common Common
			json.Unmarshal(*hit.Source, &common)
			nohitClickImp := true

			//*************************if event has stat click id and different with parent
			if event.StatClickId != "" && common.ClickInstallId != event.StatClickId {
				termQuery = elastic.NewTermQuery("Id", event.StatClickId)
				searchResult, _ = client.Search().Index("alls").Type("Click").Query(&termQuery).Pretty(true).Do()
				if searchResult.TotalHits() > 0 {
					if searchResult.TotalHits() == 1 {
						nohitClickImp = false
						//put the click struct inside the event sturct
						var found Click
						hits := searchResult.Hits.Hits[0]
						json.Unmarshal(*hits.Source, &found)
						event.Click = found

						//put the event as a child
						client.Index().Index("alls").Type("Event").Parent(parentId).BodyJson(event).Do()

					} else {
						structured.Warn(event.Id, "Event", "this event has multiple clicks related to it ", int(allFields.SiteId), nil)
					}
				}

			} else if event.StatImpressionId != "" && common.ImpressionInstallId != event.StatImpressionId {
				//*************************if event has stat impression id and different with parent

				termQuery = elastic.NewTermQuery("Id", event.StatImpressionId)
				searchResult, _ = client.Search().Index("alls").Type("Impression").Query(&termQuery).Pretty(true).Do()
				if searchResult.TotalHits() > 0 {
					if searchResult.TotalHits() == 1 {
						nohitClickImp = false
						//put the impression struct inside the event sturct
						var found Impression
						hits := searchResult.Hits.Hits[0]
						json.Unmarshal(*hits.Source, &found)
						event.Impression = found

						//put the event as a child
						client.Index().Index("alls").Type("Event").Parent(parentId).BodyJson(event).Do()

					} else {
						structured.Warn(event.Id, "Event", "this event has multiple impression related to it ", int(allFields.SiteId), nil)
					}
				}

			}

			if nohitClickImp {
				//if the event no related with any click or impresssion
				client.Index().Index("alls").Type("Event").Parent(parentId).BodyJson(event).Do()
			}

		} else if searchResult.TotalHits() > 1 {
			structured.Warn(event.Id, "event", "has multiple event contributed to that", int(allFields.SiteId), nil)
		}

		//else is no hit with current install leave nohit to be true

	}
	//**********no install related to it, post all parent and children doc have ttl and then ignore it in ES
	if noHit {
		var common Common
		common.GoogleAid = allFields.GoogleAid
		common.WindowsAid = allFields.WindowsAid
		common.IosIfa = allFields.IosIfa
		common.Language = allFields.Language
		common.CurrencyCode = allFields.CurrencyCode
		common.SiteId = allFields.SiteId
		common.AdvertiserId = allFields.AdvertiserId
		common.PackageName = allFields.PackageName
		common.PublisherId = allFields.PublisherId
		common.AdNetworkId = allFields.AdNetworkId
		common.AgencyId = allFields.AgencyId
		common.CampaignId = allFields.CampaignId
		common.PublisherUserId = allFields.PublisherUserId

		indexParent, _ := client.Index().Index("alls").Type("Common").BodyJson(common).Do()
		parentId := indexParent.Id

		client.Index().Index("alls").Type("Event").Parent(parentId).BodyJson(event).Do()
	}

}

func open(allFileds AllFields) {

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

//for update the impression_id inside the common parent fileds
type UpdateImpressionId struct {
	ImpressionInstallId string
}

type Common struct {
	GoogleAid           string
	WindowsAid          string
	IosIfa              string
	Language            string
	CurrencyCode        string
	SiteId              int64
	AdvertiserId        int64
	PackageName         string
	PublisherId         int64
	AdNetworkId         int64
	AgencyId            int64
	CampaignId          int64
	PublisherUserId     string
	StatInstallId       string // the install related to all events /open/clicks
	ClickInstallId      string //The click id related to that install
	ImpressionInstallId string //The impression id related to that install
	TempEventID         string
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
	LogType          string `json:"log_type"`
	Id               string `json:"id"`
	Created          time.Time
	TempTime         string  `json:"created"`
	DeviceIp         string  `json:"device_ip"`
	GoogleAid        string  `json:"google_aid"`
	WindowsAid       string  `json:"windows_aid"`
	IosIfa           string  `json:"ios_ifa"`
	Language         string  `json:"language"`
	StatInstallId    string  `json:"stat_install_id"`
	StatOpenId       string  `json:"stat_open_id`
	StatClickId      string  `json:"stat_click_id"`
	StatImpressionId string  `json:"stat_impression_id"`
	CurrencyCode     string  `json:"currency_code"`
	SiteId           int64   `json:"site_id"`
	AdvertiserId     int64   `json:"advertiser_id"`
	PackageName      string  `json:"package_name"`
	PublisherId      int64   `json:"publisher_id"`
	AdNetworkId      int64   `json:"ad_network_id"`
	AgencyId         int64   `json:"agency_id"`
	CampaignId       int64   `json:"campaign_id"`
	CountryCode      string  `json:"country_code"`
	RegionCode       string  `json:"region_code"`
	PostalCode       int32   `json:"postal_code"`
	WurflBrandName   string  `json:"wurfl_brand_name"`
	WurflModelName   string  `json:"wurfl_model_name"`
	WurflDeviceOs    string  `json:"wurfl_device_os"`
	PublisherUserId  string  `json:"publisher_user_id"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
}
