package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// select how may pieces of random data to create
var sizeOfDatabase int = 100000

type Event struct {
	Type         string
	ID           string
	AdvertiserID int
	PublisherID  int
	SiteID       int
	IP           string
	IosIfa       string
	GoogleAid    string
	WindowsAid   string
}

// hold all the characters that are used in random strings
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// generate a random string that is n characters
func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// generate an event that will have IosIfa OR GoogleAid OR WindowsAid, but never more than 1 type
func makeEvent() Event {

	var thing Event

	typeOfDevice := rand.Intn(3)

	thing.PublisherID = rand.Intn(90000) + 10000
	thing.AdvertiserID = rand.Intn(90000) + 10000
	thing.ID = generateID(thing.AdvertiserID)
	thing.SiteID = rand.Intn(90000) + 10000
	thing.IP = strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256))
	thing.Type = ""

	switch typeOfDevice {
	case 0:
		thing.IosIfa = randString(12)
	case 1:
		thing.GoogleAid = randString(12)
	case 2:
		thing.WindowsAid = randString(12)
	}

	return thing
}

//generate a random id to represent the unique id
func generateID(adId int) string {
	t := time.Now()
	year, month, day := t.Date()
	var id = Hex(4) + "-" + strconv.Itoa(year) +
		strconv.Itoa(int(month)) + strconv.Itoa(day) + "-" + strconv.Itoa(adId)
	return id
}

//generate a random string encoded hex value with given byte
func Hex(size int) string {
	var buffer bytes.Buffer

	bytes := make([]byte, 4)
	for i := 0; i < size; i++ {
		binary.LittleEndian.PutUint32(bytes, rand.Uint32())
		buffer.WriteString(hex.EncodeToString(bytes))
	}

	return buffer.String()
}

func main() {

	// define the url to post to in elasticSearch
	//url := "http://dp-joshp01-dev.sea1.office.priv:9200/indexs/types"

	// seed the random number generator to prevent repeated values each time the program runs
	rand.Seed(time.Now().UTC().UnixNano())

	// the variables that will be used
	var thing Event // this is the event that holds the information for each of the events in the database
	var test int    // this is used to test if the event is posted to the database via a random numver in the nested for loop below

	for i := 0; i < 10; i++ {

		// generate the event that may be added to the database
		thing = makeEvent()

		for j := 0; j < 4; j++ {

			test = rand.Intn(3)
			if test == 0 {

				// generate a unique ID for this instance of the event
				thing.ID = generateID(thing.AdvertiserID)

				switch j {
				case 0:
					thing.Type = "impression"
				case 1:
					thing.Type = "click"
				case 2:
					thing.Type = "install"
				case 3:
					thing.Type = "open"
				}

				eventJSONstring, err := json.Marshal(thing)
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(string(eventJSONstring))

				/*req, err := http.NewRequest("POST", url, bytes.NewBuffer(eventJSONstring))
				if err != nil {
					fmt.Println(err)
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
				}

				//temp, _ := ioutil.ReadAll(resp.Body)

				//	fmt.Println(string(temp))

				resp.Body.Close()*/

			} // close if statement
		} // close inner for loop (j)
	} // close outer for loop (i)

}
