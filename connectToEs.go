package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	url string = "http://dp-kewei01-dev.sea1.office.priv:9200/find/types"
)

func main() {

	csvfile, err := os.Open("/tmp/stat_installs_1681.csv")
	//csvfile, err := os.Open("tempFileholder.csv")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = 28 // see the Reader struct information below
	reader.LazyQuotes = true

	//_, err = reader.Read()
	fieldKey, err := reader.Read()

	count := 1

	for {

		eachRecord, err := reader.Read()
		if err != io.EOF && err == nil {
			changetobeios8601(4, eachRecord)
			changetobeios8601(22, eachRecord)
			changetobeios8601(23, eachRecord)
			changetobeios8601(27, eachRecord)
			var jsonstring string
			jsonstring += "{ "
			for i := 0; i < len(eachRecord)-1; i++ {
				if i == 0 || i == 1 || i == 2 || i == 3 || i == 11 || i == 12 || i == 13 || i == 14 || i == 15 || i == 16 || i == 17 || i == 20 || i == 21 || i == 26 {
					jsonstring += "\"" + fieldKey[i] + "\"" + " : " + "\"" + eachRecord[i] + "\"" + " , "
				} else if eachRecord[i] == "NULL" {
					jsonstring += "\"" + fieldKey[i] + "\"" + " : " + strconv.Itoa(0) + ","
				} else {
					jsonstring += "\"" + fieldKey[i] + "\"" + " : " + eachRecord[i] + " , "
				}
			}
			jsonstring += "\"" + fieldKey[len(eachRecord)-1] + "\"" + " : " + eachRecord[len(eachRecord)-1] + " }"

			//post them to elasticsearch
			url := "http://dp-kewei01-dev.sea1.office.priv:9200/allintstalls/types"
			req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonstring)))
			if err != nil {
				fmt.Println(err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}

			resp.Body.Close()

		} else if err != nil {
			fmt.Println(err)
		} else {
			break
		}
		fmt.Println(count)

		count++

	}

}

func changetobeios8601(n int, eachRecord []string) {
	eachRecord[n] = strings.Replace(eachRecord[n], `"`, "", 2)
	convert, _ := time.Parse("2006-01-02 15:04:05", eachRecord[n])
	b, _ := convert.MarshalJSON()
	eachRecord[n] = string(b)
}
