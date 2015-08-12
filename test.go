package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	readClicks("/tmp/stat_clicks_1681.csv")
	readInstalls("/tmp/stat_installs_1681.csv")
}

func IOS8601(timefield string) string {
	timefield = strings.Replace(timefield, `"`, "", 2)
	convert, _ := time.Parse("2006-01-02 15:04:05", timefield)
	b, _ := convert.MarshalJSON()
	timefield = string(b)
	return timefield
}

// type Install struct {
// 	id               string			0
// 	tracking_id      string			1
// 	stat_click_id    string			2
// 	session_ip       string			3
// 	session_datetime time.Time		4
// 	publisher_id     int64			5
// 	ad_network_id    int64			6
// 	advertiser_id    int64			7
// 	site_id          int64			8
// 	campaign_id      int64			9
// 	site_event_id    int64			10
// 	publisher_ref_id string			11
// 	device_ip        string			12
// 	sdk              string			13
// 	device_carrier   string			14
// 	language         string			15
// 	package_name     string			16
// 	app_name         string			17
// 	country_id       int64			18
// 	region_id        int64			19
// 	user_agent       string			20
// 	request_url      string			21
// 	created          time.Time		22
// 	modified         time.Time		23
// 	latitude         float64		24
// 	longitude        float64		25
// 	match_type       string			26
// 	install_date     time.Time		27
// }

// type Click struct {
// id					string		0
// tracking_id			string		1
// publisher_id			int			2
// ad_network_id		int			3
// advertiser_id		int			4
// site_id				int			5
// campaign_id			int			6
// publisher_ref_id		string		7
// device_ip			string		8
// sdk					string		9
// device_carrier		string		10
// language				string		11
// package_name			string		12
// app_name				string		13
// country_id			int			14
// region_id			int			15
// user_agent			string		16
// request_url			string		17
// created				time.Time	18
// modified				time.Time	19
// latitude				float64		20
// longitude			float64		21
// }

func readClicks(path string) {

	var jsonString []string

	url := "http://dp-joshp01-dev.sea1.office.priv:9200/clicks/clickdata"
	dataFieldsClicks := `"id","tracking_id","publisher_id","ad_network_id","advertiser_id","site_id","campaign_id","publisher_ref_id","device_ip","sdk","device_carrier","language","package_name","app_name","country_id","region_id","user_agent","request_url","created","modified","latitude","longitude"`
	dataFieldsSlice := strings.Split(dataFieldsClicks, ",")
	lengthOfDataFieldSlice := len(dataFieldsSlice)

	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	count := 0
	count1 := 0

	for scanner.Scan() {

		if count == 0 {
			count++
			continue
		}

		// clear the jsonString
		jsonString = jsonString[:0]

		jsonString = append(jsonString, "{")

		line := scanner.Text()

		fields := strings.Split(line, ",")
		lengthOfFields := len(fields)

		if lengthOfFields < lengthOfDataFieldSlice {
			count++
			count1++
			continue
		}

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
				if i == 0 || i == 1 || i == 7 || i == 8 || i == 9 || i == 10 || i == 11 || i == 12 || i == 13 || i == 16 || i == 17 {
					fields[i] = "\"0\""
				} else {
					fields[i] = "0"
				}
			}
		}

		fields[18] = IOS8601(fields[18])
		fields[19] = IOS8601(fields[19])

		for i := 0; i < lengthOfDataFieldSlice-1; i++ {
			jsonString = append(jsonString, dataFieldsSlice[i])
			jsonString = append(jsonString, ":")
			jsonString = append(jsonString, fields[i])
			jsonString = append(jsonString, ",")
		}

		jsonString = append(jsonString, dataFieldsSlice[lengthOfDataFieldSlice-1])
		jsonString = append(jsonString, ":")
		jsonString = append(jsonString, fields[lengthOfDataFieldSlice-1])
		jsonString = append(jsonString, "}")

		finalJsonString := strings.Join(jsonString, "")

		// fmt.Println(finalJsonString)
		// fmt.Println("\n")

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(finalJsonString)))
		if err != nil {
			fmt.Println(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		resp.Body.Close()

		fmt.Println(count, ",", count1)
		// fmt.Println(finalJsonString)

		count++
	}

}

func readInstalls(path string) {

	var jsonString []string

	url := "http://dp-joshp01-dev.sea1.office.priv:9200/installs/installdata"
	dataFieldsInstalls := `"id","tracking_id","stat_click_id","session_ip","session_datetime","publisher_id","ad_network_id","advertiser_id","site_id","campaign_id","site_event_id","publisher_ref_id","device_ip","sdk","device_carrier","language","package_name","app_name","country_id","region_id","user_agent","request_url","created","modified","latitude","longitude","match_type","install_date"`
	dataFieldsSlice := strings.Split(dataFieldsInstalls, ",")
	lengthOfDataFieldSlice := len(dataFieldsSlice)

	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	count := 0
	count1 := 0

	for scanner.Scan() {

		if count == 0 {
			count++
			continue
		}

		// clear the jsonString
		jsonString = jsonString[:0]

		jsonString = append(jsonString, "{")

		line := scanner.Text()

		fields := strings.Split(line, ",")
		lengthOfFields := len(fields)

		if lengthOfFields < lengthOfDataFieldSlice {
			count++
			count1++
			continue
		}

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

		for i := 0; i < lengthOfDataFieldSlice-1; i++ {
			jsonString = append(jsonString, dataFieldsSlice[i])
			jsonString = append(jsonString, ":")
			jsonString = append(jsonString, fields[i])
			jsonString = append(jsonString, ",")
		}

		jsonString = append(jsonString, dataFieldsSlice[lengthOfDataFieldSlice-1])
		jsonString = append(jsonString, ":")
		jsonString = append(jsonString, fields[lengthOfDataFieldSlice-1])
		jsonString = append(jsonString, "}")

		finalJsonString := strings.Join(jsonString, "")

		// fmt.Println(finalJsonString)
		// fmt.Println("\n")

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(finalJsonString)))
		if err != nil {
			fmt.Println(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		resp.Body.Close()

		fmt.Println(count, ",", count1)
		// fmt.Println(finalJsonString)

		count++
	}

}
