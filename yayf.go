/* yayf.go */

// Why don't we make the record simpler?
// Just a list of links instead.
// This way, we can quickly check the current subscription list.

// Adding support for playlists as well.
// It follows the same format as the channels!

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
	RecordsPath       = "records.tpl" //"yayf.record"
	YTFeedUrl         = "https://www.youtube.com/feeds/videos.xml?%s=%s"
)

var (
	WG sync.WaitGroup
)

type Subscriptions struct {
	Cids []string `json:"Subscriptions"`
	Pids []string `json:"Playlists"`
	Len  int
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

func GetSubs() *Subscriptions {
	var subs Subscriptions
	data, err := ioutil.ReadFile(SubscriptionsPath)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Subscrions file is missing: yayf.subs")
	}
	err = json.Unmarshal(data, &subs)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Wrong JSON format")
	}
	subs.Len = len(subs.Cids) + len(subs.Pids)
	return &subs
}

func Exist(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetRecords() map[string]string {
	var records map[string]string
	data, err := ioutil.ReadFile(RecordsPath)
	if err != nil {
		// error due to non-exisitent file?
		if Exist(RecordsPath) {
			// record does not exist. create one.
			ioutil.WriteFile(RecordsPath, []byte("{}"), 0644)
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

func GetFeeds(ch chan Feeds, id string, mode string) {
	defer WG.Done()
	var (
		feeds Feeds
	)
	resp, err := http.Get(fmt.Sprintf(YTFeedUrl, mode, id))
	if err != nil {
		log.Fatal(err)
	}
	blobs, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	xml.Unmarshal(blobs, &feeds)
	feeds.ChannelID = id
	ch <- feeds
}

func main() {
	subs := GetSubs()
	var (
		ch      = make(chan Feeds, subs.Len)
		feeds   = make(map[string]map[string]string)
		toFetch []string
	)
	records := GetRecords()
	for _, id := range subs.Cids {
		WG.Add(1)
		go GetFeeds(ch, id, "channel_id")
	}
	for _, id := range subs.Pids {
		WG.Add(1)
		go GetFeeds(ch, id, "playlist_id")
	}
	WG.Wait()
	close(ch)
	for feed := range ch {
		feeds[feed.ChannelID] = map[string]string{}
		for _, e := range feed.Entries {
			if len(e.Title) > 20 {
				e.Title = e.Title[:20] + "..."
			}
			feeds[feed.ChannelID][e.Link] = e.Title
		}
	}
	for k, v := range feeds {
		fmt.Println("((( " + k + " )))")
		for kk, _ := range v {
			fmt.Println(kk)
		}
	}
	for _, feed := range feeds {
		for link, _ := range feed {
			if record, ok := records[link]; ok {
				fmt.Println("link exists in record", link, "|", record)
			} else {
				//don't use it now!
				//records[link] = ""
				toFetch = append(toFetch, link)
			}
		}
	}
	fmt.Println(toFetch)
}
