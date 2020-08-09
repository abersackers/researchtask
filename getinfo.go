package main

//all packages still under expirementation
import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
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
		startTime, endTime, responseStatus := httpData(websites[i])
		data := []string{strconv.Itoa(int(startTime)), strconv.Itoa(int(endTime)), responseStatus}
		writer.Write(data)
	}

}

//analysis code of each website
//Currently returns the start time, end time, and https response code and prints
//the resonse body (currently in testing stage of how to get bytes and such)
//Goals: figure out how to implement request ID, ensure that URL queried can be added,
//and check with sudheesh as to if response bytes is only about the body or includes the headers as well)
func httpData(website string) (int64, int64, string) {

	startTime := time.Now().UnixNano()
	resp, err := http.Get(website)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	endTime := time.Now().UnixNano()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)

	fmt.Println(responseString)

	return startTime, endTime, resp.Status

}
