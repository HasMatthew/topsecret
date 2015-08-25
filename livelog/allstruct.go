package main

import "time"

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
	Click            ClickStr
	Impression       ImpressionStr
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
	Click            ClickStr
	Impression       ImpressionStr
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

//code trash :

// 	//***********************special case: if event happen beofre both click/install, well actually it is rare***********************//
// 	//find related impression or click from inner struct of open/events, put the click back to outside,
// 	//I decide not to put all related events and opens back to this install because it is really rare and instead, I would put ttl
// 	//on both parent and children of events /opens if they are post alone

// 	QueryOne := elastic.NewTermQuery("Click.Id", installUni.StatClickId)
// 	QueryTwo := elastic.NewTermQuery("Impression.Id", installUni.StatImpressionId)
// 	searchEvents, _ := client.Search().Index("alls").Type("Event").Query(&QueryOne).From(0).Size(1).Do()
// 	if searchEvents.TotalHits() == 0 {
// 		searchEvents, _ = client.Search().Index("alls").Type("Event").Query(&QueryTwo).From(0).Size(1).Do()
// 		if searchEvents.TotalHits() == 0 {

// 		} else {
// 			hit := searchEvents.Hits.Hits[0]
// 			var event Event
// 			json.Unmarshal(*hit.Source, &event)
// 			impression:= event.Impression
// 			//update the click id inside of common one
// 			var update UpdateImpression
// 			update.Impression = impression.Id
// 			client.Update().Index("alls").Type("Common").Id(parentId).Doc(update).Do()
// 			//get the click struct out from event to install
// 			client.Index().Index("alls").Type("Click").Parent(parentId).BodyJson(click).Do()
// 		}
// 	} else {
// 		hit := searchEvents.Hits.Hits[0]
// 		var event Event
// 		json.Unmarshal(*hit.Source, &event)
// 		click := event.Click
// 		//update the click id inside of common one
// 		var update UpdateClickId
// 		update.ClickInstallId = click.Id
// 		client.Update().Index("alls").Type("Common").Id(parentId).Doc(update).Do()
// 		//get the click struct out from event to install
// 		client.Index().Index("alls").Type("Click").Parent(parentId).BodyJson(click).Do()

// 	}

// }
