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

	url := "http://dp-joshp01-dev.sea1.office.priv:9200/database2/mydata"
	dataFields := `"id","tracking_id","stat_click_id","session_ip","session_datetime","publisher_id","ad_network_id","advertiser_id","site_id","campaign_id","site_event_id","publisher_ref_id","device_ip","sdk","device_carrier","language","package_name","app_name","country_id","region_id","user_agent","request_url","created","modified","latitude","longitude","match_type","install_date"`
	dataFieldsSlice := strings.Split(dataFields, ",")
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
		// if count > 10000 {
		// 	break
		// }

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
