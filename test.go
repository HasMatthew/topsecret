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
	fmt.Println("-----------------------------------")
	//	readInstalls("/tmp/stat_installs_1681.csv")
}

func IOS8601(timefield string) string {
	timefield = strings.Replace(timefield, `"`, "", 2)
	convert, _ := time.Parse("2006-01-02 15:04:05", timefield)
	b, _ := convert.MarshalJSON()
	timefield = string(b)
	return timefield
}

func addToJsonString(field string, value string, JSONstring []string) []string {

	JSONstring = append(JSONstring, ",")
	JSONstring = append(JSONstring, field)
	JSONstring = append(JSONstring, ":")
	JSONstring = append(JSONstring, value)

	return JSONstring
}
func readClicks(path string) {

	var jsonString []string
	var jsonStringClick []string
	var jsonStringInstall []string

	url := "http://dp-joshp01-dev.sea1.office.priv:9200/realSampleData/testData"
	dataFieldsClicks := `"id","tracking_id","publisher_id","ad_network_id","advertiser_id","site_id","campaign_id","publisher_ref_id","device_ip","sdk","device_carrier","language","package_name","app_name","country_id","region_id","user_agent","request_url","created","modified","latitude","longitude"`
	dataFieldsSlice := strings.Split(dataFieldsClicks, ",")
	lengthOfDataFieldSlice := len(dataFieldsSlice)

	dataFieldsInstalls := `"id","tracking_id","stat_click_id","session_ip","session_datetime","publisher_id","ad_network_id","advertiser_id","site_id","campaign_id","site_event_id","publisher_ref_id","device_ip","sdk","device_carrier","language","package_name","app_name","country_id","region_id","user_agent","request_url","created","modified","latitude","longitude","match_type","install_date"`
	dataFieldsSliceInstalls := strings.Split(dataFieldsInstalls, ",")
	// lengthOfDataFieldSliceInstalls := len(dataFieldsSliceInstalls)

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

		if count > 1 {
			count++
			break
		}

		// clear the jsonStrings
		jsonString = jsonString[:0]
		jsonStringClick = jsonStringClick[:0]
		jsonStringInstall = jsonStringInstall[:0]

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

		//------- same for clicks  ---------
		// publisher_id			int			2
		// ad_network_id		int			3
		// advertiser_id		int			4
		// site_id				int			5
		// campaign_id			int			6
		// sdk					string		9
		// device_carrier		string		10
		// language				string		11
		// package_name			string		12
		// app_name				string		13
		// country_id			int			14
		// region_id			int			15

		// ------  all click fields ---------
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

		jsonString = append(jsonString, "{")
		jsonStringClick = append(jsonStringClick, "{")
		jsonStringInstall = append(jsonStringInstall, "{")

		jsonStringClick = append(jsonStringClick, dataFieldsSlice[0])
		jsonStringClick = append(jsonStringClick, ":")
		jsonStringClick = append(jsonStringClick, fields[0])

		jsonStringClick = addToJsonString(dataFieldsSlice[1], fields[1], jsonStringClick)

		jsonString = append(jsonString, dataFieldsSlice[2])
		jsonString = append(jsonString, ":")
		jsonString = append(jsonString, fields[2])

		jsonString = addToJsonString(dataFieldsSlice[3], fields[3], jsonString)
		jsonString = addToJsonString(dataFieldsSlice[4], fields[4], jsonString)
		jsonString = addToJsonString(dataFieldsSlice[5], fields[5], jsonString)
		jsonString = addToJsonString(dataFieldsSlice[6], fields[6], jsonString)

		jsonStringClick = addToJsonString(dataFieldsSlice[7], fields[7], jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSlice[8], fields[8], jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSlice[9], fields[9], jsonStringClick)

		jsonString = addToJsonString(dataFieldsSlice[10], fields[10], jsonString)
		jsonString = addToJsonString(dataFieldsSlice[11], fields[11], jsonString)
		jsonString = addToJsonString(dataFieldsSlice[12], fields[12], jsonString)
		jsonString = addToJsonString(dataFieldsSlice[13], fields[13], jsonString)
		jsonString = addToJsonString(dataFieldsSlice[14], fields[14], jsonString)
		jsonString = addToJsonString(dataFieldsSlice[15], fields[15], jsonString)

		jsonStringClick = addToJsonString(dataFieldsSlice[16], fields[16], jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSlice[17], fields[17], jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSlice[18], fields[18], jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSlice[19], fields[19], jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSlice[20], fields[20], jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSlice[21], fields[21], jsonStringClick)

		jsonStringInstall = append(jsonStringInstall, dataFieldsSliceInstalls[0])
		jsonStringInstall = append(jsonStringInstall, ":")
		jsonStringInstall = append(jsonStringInstall, "null")

		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[1], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[2], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[3], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[4], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[10], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[11], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[12], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[20], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[21], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[22], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[23], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[24], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[25], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[26], "null", jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[27], "null", jsonStringInstall)

		jsonString = append(jsonString, `, "Click_data" : `)
		jsonString = append(jsonString, jsonStringClick...)
		jsonString = append(jsonString, "},")
		jsonString = append(jsonString, `"Install_data" : `)
		jsonString = append(jsonString, jsonStringInstall...)
		jsonString = append(jsonString, "}}")

		finalJsonString := strings.Join(jsonString, "")

		fmt.Println(finalJsonString)
		fmt.Println("\n")

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
		//	fmt.Println(finalJsonString)

		count++
	}

}

func readInstalls(path string) {

	var jsonString []string
	var jsonStringClick []string
	var jsonStringInstall []string

	url := "http://dp-joshp01-dev.sea1.office.priv:9200/realSampleData/testData"
	dataFieldsInstalls := `"id","tracking_id","stat_click_id","session_ip","session_datetime","publisher_id","ad_network_id","advertiser_id","site_id","campaign_id","site_event_id","publisher_ref_id","device_ip","sdk","device_carrier","language","package_name","app_name","country_id","region_id","user_agent","request_url","created","modified","latitude","longitude","match_type","install_date"`
	dataFieldsSliceInstalls := strings.Split(dataFieldsInstalls, ",")
	lengthOfDataFieldSliceInstalls := len(dataFieldsSliceInstalls)

	dataFieldsClicks := `"id","tracking_id","publisher_id","ad_network_id","advertiser_id","site_id","campaign_id","publisher_ref_id","device_ip","sdk","device_carrier","language","package_name","app_name","country_id","region_id","user_agent","request_url","created","modified","latitude","longitude"`
	dataFieldsSliceClicks := strings.Split(dataFieldsClicks, ",")
	//	lengthOfDataFieldClicksSlice := len(dataFieldsSliceClicks)

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
		if count > 100 {
			break
		}

		// clear the jsonString
		jsonString = jsonString[:0]
		jsonStringClick = jsonStringClick[:0]
		jsonStringInstall = jsonStringInstall[:0]

		line := scanner.Text()

		fields := strings.Split(line, ",")
		lengthOfFields := len(fields)

		if lengthOfFields < lengthOfDataFieldSliceInstalls {
			count++
			count1++
			continue
		}

		if lengthOfFields > lengthOfDataFieldSliceInstalls {
			for i := 0; i < lengthOfDataFieldSliceInstalls; i++ {
				if strings.HasPrefix(fields[i], `"`) && !(strings.HasSuffix(fields[i], `"`)) {
					fields[i] = fields[i] + "," + fields[i+1]
					for j := i + 1; j < lengthOfDataFieldSliceInstalls; j++ {
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

		//  --------- same for Installs  ---------
		//  publisher_id     int64			5	gr
		//  ad_network_id    int64			6	gr
		//  advertiser_id    int64			7	gr
		//  site_id          int64			8	gr
		//  campaign_id      int64			9	gr
		// 	sdk              string			13	gr
		// 	device_carrier   string			14	gr
		// 	language         string			15	gr
		// 	package_name     string			16	gr
		// 	app_name         string			17	gr
		// 	country_id       int64			18	gr
		// 	region_id        int64			19	gr

		// type Install struct {
		// 	id               string			0	u
		// 	tracking_id      string			1	u
		// 	stat_click_id    string			2	u
		// 	session_ip       string			3	gr
		// 	session_datetime time.Time		4	u
		// 	publisher_id     int64			5	gr
		// 	ad_network_id    int64			6	gr
		// 	advertiser_id    int64			7	gr
		// 	site_id          int64			8	gr
		// 	campaign_id      int64			9	gr
		// 	site_event_id    int64			10	gr
		// 	publisher_ref_id string			11	gr
		// 	device_ip        string			12	gr
		// 	sdk              string			13	gr
		// 	device_carrier   string			14	gr
		// 	language         string			15	gr
		// 	package_name     string			16	gr
		// 	app_name         string			17	gr
		// 	country_id       int64			18	gr
		// 	region_id        int64			19	gr
		// 	user_agent       string			20	too variable
		// 	request_url      string			21	u
		// 	created          time.Time		22	u
		// 	modified         time.Time		23	u
		// 	latitude         float64		24	u
		// 	longitude        float64		25	u
		// 	match_type       string			26	u
		// 	install_date     time.Time		27	u

		jsonString = append(jsonString, "{")
		jsonStringClick = append(jsonStringClick, "{")
		jsonStringInstall = append(jsonStringInstall, "{")

		jsonStringInstall = append(jsonStringInstall, dataFieldsSliceInstalls[0])
		jsonStringInstall = append(jsonStringInstall, ":")
		jsonStringInstall = append(jsonStringInstall, fields[0])

		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[1], fields[1], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[2], fields[2], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[3], fields[3], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[4], fields[4], jsonStringInstall)

		jsonString = append(jsonString, dataFieldsSliceInstalls[5])
		jsonString = append(jsonString, ":")
		jsonString = append(jsonString, fields[5])

		jsonString = addToJsonString(dataFieldsSliceInstalls[6], fields[6], jsonString)
		jsonString = addToJsonString(dataFieldsSliceInstalls[7], fields[7], jsonString)
		jsonString = addToJsonString(dataFieldsSliceInstalls[8], fields[8], jsonString)
		jsonString = addToJsonString(dataFieldsSliceInstalls[9], fields[9], jsonString)

		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[10], fields[10], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[11], fields[11], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[12], fields[12], jsonStringInstall)

		jsonString = addToJsonString(dataFieldsSliceInstalls[13], fields[13], jsonString)
		jsonString = addToJsonString(dataFieldsSliceInstalls[14], fields[14], jsonString)
		jsonString = addToJsonString(dataFieldsSliceInstalls[15], fields[15], jsonString)
		jsonString = addToJsonString(dataFieldsSliceInstalls[16], fields[16], jsonString)
		jsonString = addToJsonString(dataFieldsSliceInstalls[17], fields[17], jsonString)
		jsonString = addToJsonString(dataFieldsSliceInstalls[18], fields[18], jsonString)
		jsonString = addToJsonString(dataFieldsSliceInstalls[19], fields[19], jsonString)

		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[20], fields[20], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[21], fields[21], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[22], fields[22], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[23], fields[23], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[24], fields[24], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[25], fields[25], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[26], fields[26], jsonStringInstall)
		jsonStringInstall = addToJsonString(dataFieldsSliceInstalls[27], fields[27], jsonStringInstall)

		jsonStringClick = append(jsonStringClick, dataFieldsSliceClicks[0])
		jsonStringClick = append(jsonStringClick, ":")
		jsonStringClick = append(jsonStringClick, "null")

		jsonStringClick = addToJsonString(dataFieldsSliceClicks[7], "null", jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSliceClicks[8], "null", jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSliceClicks[9], "null", jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSliceClicks[16], "null", jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSliceClicks[17], "null", jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSliceClicks[18], "null", jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSliceClicks[19], "null", jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSliceClicks[20], "null", jsonStringClick)
		jsonStringClick = addToJsonString(dataFieldsSliceClicks[21], "null", jsonStringClick)

		jsonString = append(jsonString, `, "Click_data" : `)
		jsonString = append(jsonString, jsonStringClick...)
		jsonString = append(jsonString, "},")
		jsonString = append(jsonString, `"Install_data" : `)
		jsonString = append(jsonString, jsonStringInstall...)
		jsonString = append(jsonString, "}}")

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
