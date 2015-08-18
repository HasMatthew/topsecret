import "time"

type Document struct {
	Click      click
	Impression impression
	Install    install
	Events     []event
	Opens      []open
}

type commonInfo struct {
	google_aid        string
	windows_aid       string
	ios_ifa           string
	language          string
	currency_code     string
	site_id           int64
	advertiser_id     int64
	package_name      string
	publisher_id      int64
	ad_network_id     int64
	agency_id         int64
	campaign_id       int64
	publisher_user_id string
}

type impression struct {
	id               string
	created          time.Time
	device_ip        string
	country_code     string
	region_code      string
	postal_code      int32
	latitude         float64
	longitude        float64
	wurfl_brand_name string
	wurfl_model_name string
	wurfl_device_os  string
}

type click struct {
	id                 string
	created            time.Time
	device_ip          string
	stat_impression_id string
	country_code       string
	region_code        string
	postal_code        int32
	latitude           float64
	longitude          float64
	wurfl_brand_name   string
	wurfl_model_name   string
	wurfl_device_os    string
}

type install struct {
	id               string
	created          time.Time
	device_ip        string
	stat_click_id    string
	country_code     string
	region_code      string
	postal_code      int32
	latitude         float64
	longitude        float64
	wurfl_brand_name string
	wurfl_model_name string
	wurfl_device_os  string
}

type open struct {
	id               string
	created          time.Time
	device_ip        string
	stat_install_id  string
	country_code     string
	region_code      string
	postal_code      int32
	latitude         float64
	longitude        float64
	wurfl_brand_name string
	wurfl_model_name string
	wurfl_device_os  string
	Click            click
	Impression       impression
}

type event struct {
	id               string
	created          time.Time
	device_ip        string
	stat_open_id     string
	country_code     string
	region_code      string
	postal_code      int32
	latitude         float64
	longitude        float64
	wurfl_brand_name string
	wurfl_model_name string
	wurfl_device_os  string
	Click            click
	Impression       impression
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