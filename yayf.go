package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

const (
	YT_FEED_URL = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"
)

var (
	CHANNEL_IDS = [3]string{
		"UCwi3BrUqM4xStpbCyxsb3TA", // mono okito
		"UCFaYLR_1aryjfB7hLrKGRaQ", // Michael Sugure
		"UC3I2GFN_F8WudD_2jUZbojA", // KEXP
	}
	wg sync.WaitGroup
)

type Channel_Ids struct {
	Subscriptions []string
}

type Subscriptions struct {
	Title   string  `xml:"title"`
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Id    string `xml:"videoId"`
	Title string `xml:"title"`
}

func get_subs(path string) []string {
	var CI Channel_Ids
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(data, &CI)
	if err != nil {
		fmt.Println(err)
	}
	return CI.Subscriptions
}

func get_feed(url string, c chan Subscriptions) {
	defer wg.Done()
	var (
		subs Subscriptions
	)
	resp, _ := http.Get(url)
	blobs, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(blobs, &subs)
	resp.Body.Close()
	c <- subs
}

func main() {
	CI := get_subs("./yayf.config")
	var (
		curr_url string
		c        = make(chan Subscriptions, len(CI))
	)
	for _, ch_id := range CI {
		curr_url = fmt.Sprintf(YT_FEED_URL, ch_id)
		wg.Add(1)
		go get_feed(curr_url, c)
	}
	wg.Wait()
	close(c)
	for sub := range c {
		fmt.Println("***" + sub.Title + "***")
		for _, e := range sub.Entries {
			fmt.Println(e.Id, e.Title)
		}
	}
}
