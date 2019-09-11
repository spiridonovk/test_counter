package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)
import "log"
import "io/ioutil"

const workersSize = 5

func main() {
	counter := make(map[string]int)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	blockChannel := make(chan struct{}, workersSize)
	r := regexp.MustCompile(`\b(Go)\b`)
	urls := getInputData()

	wg.Add(len(urls))
	for _, url := range urls {
		blockChannel <- struct{}{}

		go func(str string) {
			defer func() {
				<-blockChannel
				wg.Done()
			}()
			res, err := doRequest(str)
			if err != nil {
				log.Println("Url: "+str+" err: ", err.Error())
			}
			size := countSubstring(res, r)
			mu.Lock()
			counter[str] = size
			mu.Unlock()
		}(url)
	}
	wg.Wait()
	printResult(counter)
}

func getInputData() []string {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Print(err.Error())
	}
	str := string(bytes)
	urls := strings.Split(str, "\n")
	return urls[:len(urls)-1]
}

func doRequest(url string) (string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err

	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil

}

func countSubstring(str string, r *regexp.Regexp) int {
	matches := r.FindAllString(str, -1)
	return len(matches)
}
func printResult(counter map[string]int) {
	i := 0
	for key, val := range counter {
		i = i + val
		fmt.Printf("Count for %s:  %d \n", key, val)
	}
	fmt.Println("Total: ", i)

}
