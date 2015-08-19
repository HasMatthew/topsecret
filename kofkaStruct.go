import "time"

type GeoLocation struct {
	Lat float64
	Lon float64
}

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
	location       GeoLocation
	WurflBrandName string
	WurflModelName string
	WurflDeviceOs  string
}

type click struct {
	Id             string
	Created        time.Time
	DeviceIp       string
	CountryCode    string
	RegionCode     string
	PostalCode     int32
	location       GeoLocation
	WurflBrandName string
	WurflModelName string
	WurflDeviceOs  string
}

type install struct {
	Id             string
	Created        time.Time
	DeviceIp       string
	StatClickId    string
	CountryCode    string
	RegionCode     string
	PostalCode     int32
	location       GeoLocation
	WurflBrandName string
	WurflModelName string
	WurflDeviceOs  string
}

type open struct {
	Id             string
	Created        time.Time
	DeviceIp       string
	StatClickId    string
	StatInstallId  string
	CountryCode    string
	RegionCode     string
	PostalCode     int32
	location       GeoLocation
	WurflBrandName string
	WurflModelName string
	WurflDeviceOs  string
}

type event struct {
	Id             string
	Created        time.Time
	DeviceIp       string
	StatClickId    string
	StatInstallId  string
	StatOpenId     string
	CountryCode    string
	RegionCode     string
	PostalCode     int32
	location       GeoLocation
	WurflBrandName string
	WurflModelName string
	WurflDeviceOs  string
}

type allFields struct {
	LogType         string
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
	location        GeoLocation
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