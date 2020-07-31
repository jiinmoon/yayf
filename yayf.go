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
	SUB_PATH    = "yayf.subs"
	RECORD_PATH = "yayf.record"
)

var (
	wg sync.WaitGroup
)

type Channel_Ids struct {
	Subscriptions []string
}

type Subscriptions struct {
	ChannelID string
	Title     string  `xml:"title"`
	Entries   []Entry `xml:"entry"`
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

func get_feed(ch_id string, c chan Subscriptions) {
	defer wg.Done()
	var (
		subs Subscriptions
	)
	resp, _ := http.Get(fmt.Sprintf(YT_FEED_URL, ch_id))
	blobs, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(blobs, &subs)
	subs.ChannelID = ch_id
	resp.Body.Close()
	c <- subs
}

func put_subs(subs []Subscriptions) {
	sub_file := make(map[string]Subscriptions)
	for _, sub := range subs {
		sub_file[sub.ChannelID] = sub
	}
	record, _ := json.MarshalIndent(sub_file, "", " ")
	_ = ioutil.WriteFile(RECORD_PATH, record, 0644)
}

func main() {
	Channel_IDs := get_subs(SUB_PATH)
	var (
		c    = make(chan Subscriptions, len(Channel_IDs))
		subs []Subscriptions
	)
	for _, ch_id := range Channel_IDs {
		wg.Add(1)
		go get_feed(ch_id, c)
	}
	wg.Wait()
	close(c)
	for sub := range c {
		for _, e := range sub.Entries {
			fmt.Println(e.Title[:10])
		}
		subs = append(subs, sub)
	}
	put_subs(subs)
}
