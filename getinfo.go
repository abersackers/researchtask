package main

//all packages still under experimentation
import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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

	//input titles into csv for dataframes later
	writer.Write([]string{"Request ID", "URL Queried", "Query Reponse Status", "Start Time", "End Time", "Response Size"})

	//implement parallelism
	var waitGroup sync.WaitGroup

	// Give some time for listenForever to start
	time.Sleep(time.Nanosecond * 10)

	for i := 0; i < len(websites); i++ {
		//parallelize for loop by adding process and adding go func (i)
		waitGroup.Add(1)

		go func(i int) {
			// At the end of the goroutine, tell the WaitGroup
			//   that another thread has completed.
			defer waitGroup.Done()
			requestID, urlQueried, respStatus, startTime, endTime, respSize := httpData(websites[i])
			if respStatus != "false" {
				data := []string{requestID.String(), urlQueried, respStatus, strconv.Itoa(int(startTime)), strconv.Itoa(int(endTime)), strconv.Itoa(int(respSize))}
				writer.Write(data)
			}
		}(i)
	}

	//wait for this parallel process to end before exiting the function
	waitGroup.Wait()
}

//analysis code of each website
//parallelism started working with better wifi connection
func httpData(website string) (uuid.UUID, string, string, int64, int64, float64) {

	requestID := uuid.New() //randomly generate locally since it does not need to be sent in network request
	urlQueried := website   //storing to preserve order (can be deleted)

	//log times and perform Get request on website
	startTime := time.Now().UnixNano()

	resp, err := http.Get(website)
	if err != nil {
		//use these different formats to see if there is a possible solution
		if strings.HasPrefix(website, "https://") && !(strings.HasPrefix(website, "https://www.")) {
			//https:// failed try https://www.
			return httpData("https://www." + website[8:len(website)])
		} else if strings.HasPrefix(website, "https://www.") {
			//https:// failed try http://
			return httpData("http://" + website[12:len(website)])
		} else if strings.HasPrefix(website, "http://") && !strings.HasPrefix(website, "http://www.") {
			//http:// failed try http://www.
			return httpData("http://www." + website[7:len(website)])
		}
		// indicated failure and will not be added to csv
		return requestID, urlQueried, "false", startTime, time.Now().UnixNano(), -1
	}
	defer resp.Body.Close()

	endTime := time.Now().UnixNano()

	//calculate response statistics
	respStatus := resp.Status
	body, err := ioutil.ReadAll(resp.Body)

	//need to figure out better error handling here and why concurrency fails at this step
	if err != nil {
		fmt.Println("error: " + website)
		httpData(website)
		panic(err)
	}
	respSize := len(body) //interpreted as number of bytes from io

	//return all values in order that they were asked
	return requestID, urlQueried, respStatus, startTime, endTime, float64(respSize)

}
