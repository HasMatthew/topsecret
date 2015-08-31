//for update the click inside the events/opens
import "time"

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

type Acticity struct {
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
	ParentId         string
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
