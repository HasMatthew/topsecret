package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// type Install struct {
// 	id               string
// 	tracking_id      string
// 	stat_click_id    string
// 	session_ip       string
// 	session_datetime time.Time
// 	publisher_id     int64
// 	ad_network_id    int64
// 	advertiser_id    int64
// 	site_id          int64
// 	campaign_id      int64
// 	site_event_id    int64
// 	publisher_ref_id string
// 	device_ip        string
// 	sdk              string
// 	device_carrier   string
// 	language         string
// 	package_name     string
// 	app_name         string
// 	country_id       int64
// 	region_id        int64
// 	user_agent       string
// 	request_url      string
// 	created          time.Time
// 	modified         time.Time
// 	latitude         float64
// 	longitude        float64
// 	match_type       string
// 	install_date     time.Time
// }

// temp.id = fields[0]
// temp.tracking_id = fields[1]
// temp.stat_click_id = fields[2]
// temp.session_ip = fields[3]
// temp.session_datetime, _ = time.Parse("2006-01-02 15:04:05", fields[4])
// temp.publisher_id, _ = strconv.ParseInt(fields[5], 10, 64)
// temp.ad_network_id, _ = strconv.ParseInt(fields[6], 10, 64)
// temp.advertiser_id, _ = strconv.ParseInt(fields[7], 10, 64)
// temp.site_id, _ = strconv.ParseInt(fields[8], 10, 64)
// temp.campaign_id, _ = strconv.ParseInt(fields[9], 10, 64)
// temp.site_event_id, _ = strconv.ParseInt(fields[10], 10, 64)
// temp.publisher_ref_id = fields[11]
// temp.device_ip = fields[12]
// temp.sdk = fields[13]
// temp.device_carrier = fields[14]
// temp.language = fields[15]
// temp.package_name = fields[16]
// temp.app_name = fields[17]
// temp.country_id, _ = strconv.ParseInt(fields[18], 10, 64)
// temp.region_id, _ = strconv.ParseInt(fields[19], 10, 64)
// temp.user_agent = fields[20]
// temp.request_url = fields[21]
// temp.created, _ = time.Parse("2006-01-02 15:04:05", fields[22])
// temp.modified, _ = time.Parse("2006-01-02 15:04:05", fields[23])
// temp.latitude, _ = strconv.ParseFloat(fields[24], 64)
// temp.longitude, _ = strconv.ParseFloat(fields[25], 64)
// temp.match_type = fields[26]
// temp.install_date, _ = time.Parse("2006-01-02 15:04:05", fields[27])

func main() {
	readLines("/tmp/stat_installs_1681.csv")
}

func IOS8601(timefield string) string {
	timefield = strings.Replace(timefield, `"`, "", 2)
	convert, _ := time.Parse("2006-01-02 15:04:05", timefield)
	b, _ := convert.MarshalJSON()
	timefield = string(b)
	return timefield
}

func readLines(path string) {

	var jsonString []string

	//	url := "http://dp-joshp01-dev.sea1.office.priv:9200/database2/mydata"
	dataFields := `"id","tracking_id","stat_click_id","session_ip","session_datetime","publisher_id","ad_network_id","advertiser_id","site_id","campaign_id","site_event_id","publisher_ref_id","device_ip","sdk","device_carrier","language","package_name","app_name","country_id","region_id","user_agent","request_url","created","modified","latitude","longitude","match_type","install_date"`
	dataFieldsSlice := strings.Split(dataFields, ",")
	lengthOfDataFieldSlice := len(dataFieldsSlice)

	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	count := 0

	for scanner.Scan() {

		if count == 0 {
			count++
			continue
		}
		// if count > 10000 {
		// 	break
		// }

		// clear the jsonString
		jsonString = jsonString[:0]

		jsonString = append(jsonString, "{")

		line := scanner.Text()

		//fmt.Println(line)

		fields := strings.Split(line, ",")
		lengthOfFields := len(fields)

		if lengthOfFields > lengthOfDataFieldSlice {
			for i := 0; i < lengthOfDataFieldSlice; i++ {
				if strings.HasPrefix(fields[i], `"`) && !(strings.HasSuffix(fields[i], `"`)) {
					fields[i] = fields[i] + "," + fields[i+1]
					for j := i + 1; j < lengthOfDataFieldSlice; j++ {
						fields[j] = fields[j+1]
					}
					lengthOfFields--
				}
			}
		}

		for i := 0; i < lengthOfDataFieldSlice; i++ {
			if fields[i] != "NULL" {
				continue
			} else {
				if i == 5 || i == 6 || i == 7 || i == 8 || i == 9 || i == 10 || i == 18 || i == 19 || i == 24 || i == 25 {
					fields[i] = "0"
				} else {
					fields[i] = "\"0\""
				}
			}
		}

		fields[4] = IOS8601(fields[4])
		fields[22] = IOS8601(fields[22])
		fields[23] = IOS8601(fields[23])
		fields[27] = IOS8601(fields[27])

		// for i := 0; i < lengthOfDataFieldSlice-1; i++ {
		// 	jsonString = append(jsonString, dataFieldsSlice[i])
		// 	jsonString = append(jsonString, ":")
		// 	jsonString = append(jsonString, fields[i])
		// 	jsonString = append(jsonString, ",")
		// }

		// jsonString = append(jsonString, dataFieldsSlice[lengthOfDataFieldSlice-1])
		// jsonString = append(jsonString, ":")
		// jsonString = append(jsonString, fields[lengthOfDataFieldSlice-1])
		// jsonString = append(jsonString, "}")

		// finalJsonString := strings.Join(jsonString, "")

		//	fmt.Println(finalJsonString)
		//	fmt.Println("\n")

		// req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(finalJsonString)))
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// client := &http.Client{}
		// resp, err := client.Do(req)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// resp.Body.Close()

		fmt.Println(count)
		// fmt.Println(finalJsonString)

		count++
	}

}
