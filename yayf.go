/* yayf.go */

package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	SubscriptionsPath = "yayf.subs"
	RecordsPath       = "yayf.record"
	YTFeedUrl         = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"
)

var (
	WG sync.WaitGroup
)

type ChannelIDs struct {
	Cids []string `json:"Subscriptions"`
}

type Feeds struct {
	ChannelID string
	Title     string  `xml:"title"`
	Entries   []Entry `xml:"entry"`
}

type Entry struct {
	Link  string `xml:"videoId"`
	Title string `xml:"title"`
}

func GetSubs() []string {
	var cids ChannelIDs
	data, err := ioutil.ReadFile(SubscriptionsPath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &cids)
	if err != nil {
		log.Fatal(err)
	}
	return cids.Cids
}

func GetRecords() map[string]map[string]string {
	var records map[string]map[string]string
	data, err := ioutil.ReadFile(RecordsPath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &records)
	if err != nil {
		log.Fatal(err)
	}
	return records
}

func GetFeeds(ch chan Feeds, cid string) {
	defer WG.Done()
	var (
		feeds Feeds
	)
	resp, err := http.Get(fmt.Sprintf(YTFeedUrl, cid))
	if err != nil {
		log.Fatal(err)
	}
	blobs, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	xml.Unmarshal(blobs, &feeds)
	feeds.ChannelID = cid
	ch <- feeds
}

func main() {
	cids := GetSubs()
	var (
		ch    = make(chan Feeds, len(cids))
		feeds = make(map[string]map[string]string)
	)
	records := GetRecords()
	for _, id := range cids {
		WG.Add(1)
		go GetFeeds(ch, id)
	}
	WG.Wait()
	close(ch)
	for feed := range ch {
		feeds[feed.ChannelID] = map[string]string{}
		for _, e := range feed.Entries {
			feeds[feed.ChannelID][e.Link] = e.Title[:30] + "..."
		}
	}
	for k, v := range feeds {
		fmt.Println("((( " + k + " )))")
		for kk, vv := range v {
			fmt.Println(kk, vv)
		}
	}
	fmt.Println(records)
}
