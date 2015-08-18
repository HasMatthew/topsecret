import "time"

type Document struct {
	Common     commonInfo
	Click      click
	Impression impression
	Install    install
	Events     []event
	Opens      []open
}

type commonInfo struct {
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

type impression struct {
	Id             string
	Created        time.Time
	DeviceIp       string
	CountryCode    string
	RegionCode     string
	PostalCode     int32
	Latitude       float64
	Longitude      float64
	WurflBrandName string
	WurflModelName string
	WurflDeviceOs  string
}

type click struct {
	Id               string
	Created          time.Time
	DeviceIp         string
	StatImpressionId string
	CountryCode      string
	RegionCode       string
	PostalCode       int32
	Latitude         float64
	Longitude        float64
	WurflBrandName   string
	WurflModelName   string
	WurflDeviceOs    string
}

type install struct {
	Id             string
	Created        time.Time
	DeviceIp       string
	StatClickId    string
	CountryCode    string
	RegionCode     string
	PostalCode     int32
	Latitude       float64
	Longitude      float64
	WurflBrandName string
	WurflModelName string
	WurflDeviceOs  string
}

type open struct {
	Id             string
	Created        time.Time
	DeviceIp       string
	StatInstallId  string
	CountryCode    string
	RegionCode     string
	PostalCode     int32
	Latitude       float64
	Longitude      float64
	WurflBrandName string
	WurflModelName string
	WurflDeviceOs  string
	Click          click
	Impression     impression
}

type event struct {
	Id             string
	Created        time.Time
	DeviceIp       string
	StatOpenId     string
	CountryCode    string
	RegionCode     string
	PostalCode     int32
	Latitude       float64
	Longitude      float64
	WurflBrandName string
	WurflModelName string
	WurflDeviceOs  string
	Click          click
	Impression     impression
}

type allFields struct {
	Id              string
	Created         time.Time
	DeviceIp        string
	GoogleAid       string
	WindowsAid      string
	IosIfa          string
	Language        string
	StatEventId     string
	StatInstallId   string
	StatOpenId      string
	StatClickId     string
	CurrencyCode    string
	SiteId          int64
	AdvertiserId    int64
	PackageName     string
	PublisherId     int64
	AdNetworkId     int64
	AgencyId        int64
	CampaignId      int64
	CountryCode     string
	RegionCode      string
	PostalCode      int32
	WurflBrandName  string
	WurflModelName  string
	WurflDeviceOs   string
	PublisherUserId string
	BundleSiteId    int64
	IsBundle        bool
	Latitude        float64
	Longitude       float64
}

// {
// 	"Click": {},
// 	"Impression": {},
// 	"Install": {},
// 	"Events": [
// 		{
// 			"ID": "",
// 			"GoogleAID": "",
// 			"Click": {
// 				"ID": "",
// 			},
// 		}
// 	],
// 	"Opens": [
// 		{
// 			"ID": "",
// 			"GoogleAID": "",
// 			"Click": {
// 				"ID": "",
// 			},
// 		}
// 	]
// }