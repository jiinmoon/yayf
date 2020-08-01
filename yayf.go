/* yayf.go */

// Why don't we make the record simpler?
// Just a list of links instead.
// This way, we can quickly check the current subscription list.

// Playlists!

package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		fmt.Println("Subscrions file is missing: yayf.subs")
	}
	err = json.Unmarshal(data, &cids)
	if err != nil {
		log.Fatal(err)
	}
	return cids.Cids
}

func IsExist(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetRecords() map[string]map[string]string {
	var records map[string]map[string]string
	data, err := ioutil.ReadFile(RecordsPath)
	if err != nil {
		// error due to non-exisitent file?
		if IsExist(RecordsPath) {
			// record does not exist. create one.
			ioutil.WriteFile(RecordsPath, []byte{}, 0644)
		} else {
			// record exists, but other err.
			log.Fatal(err)
		}
		return records
	} else {
		err = json.Unmarshal(data, &records)
		if err != nil {
			log.Fatal(err)
		}
		return records
	}
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
