package main

//all packages still under experimentation
import (
	"bufio"
	"encoding/csv"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//has been tested and properly written csv file for smallLinks.txt
func main() {
	//read in txt file and run the analysis
	websites, err := readLines(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	csvCreate(websites)
}

//reads the txt file inputted, takes the data and adds https:// in front of every entry
//and adds that all to one array for use by other methods
func readLines(file string) ([]string, error) {
	data, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	var websites []string
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		websites = append(websites, "https://"+scanner.Text())
	}
	return websites, scanner.Err()
}

//creates a csv file called requestData.csv and runs the analysis on every entry in the inputted array
//Currently the analysis is incomplete as some of it still needs to be implemented
func csvCreate(websites []string) {
	file, err := os.Create("requestData.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < len(websites); i++ {
		requestID, urlQueried, respStatus, startTime, endTime, respSize := httpData(websites[i])
		data := []string{requestID.String(), urlQueried, respStatus, strconv.Itoa(int(startTime)), strconv.Itoa(int(endTime)), strconv.Itoa(int(respSize))}
		writer.Write(data)
	}

}

//analysis code of each website
//goal: figure out response bytes
//timing out on largeLinks.txt, current iteration is on github
//seems to be a timeout issue which will require further investigation, however
//since it works well on smallLinks.txt will come back later
func httpData(website string) (uuid.UUID, string, string, int64, int64, int64) {

	requestID := uuid.New() //randomly generate locally since it does not need to be sent in network request
	urlQueried := website   //storing to preserve order (can be deleted)

	//log times and perform Get request on website
	startTime := time.Now().UnixNano()

	resp, err := http.Get(website)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	endTime := time.Now().UnixNano()

	//calculate response statistics
	respStatus := resp.Status
	respSize := resp.ContentLength //returns -1 a lot so have to check up (might be due to type of int64)

	//return all values in order that they were asked
	return requestID, urlQueried, respStatus, startTime, endTime, respSize

}
