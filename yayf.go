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
	Len  int      // total len of Cids + Pids
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

func (s *Subscriptions) GetSubs() error {
	data, err := ioutil.ReadFile(SubscriptionsPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	s.Len = len(s.Cids) + len(s.Pids)
	return nil
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

func (s *Subscriptions) GetFeeds(ch chan Feeds) error {
	for _, id := range s.Cids {
		WG.Add(1)
		go GetFeed(ch, id, "channel_id")
	}
	for _, id := range s.Pids {
		WG.Add(1)
		go GetFeed(ch, id, "playlist_id")
	}
	WG.Wait()
	close(ch)
	return nil
}

func GetFeed(ch chan Feeds, id string, mode string) {
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
	var (
		subs    = &Subscriptions{}
		feeds   = make(map[string]map[string]string)
		toFetch []string
	)
	// read subs
	subs.GetSubs()
	// get feeds from each subs
	ch := make(chan Feeds, subs.Len)
	subs.GetFeeds(ch)
	// read records
	records := GetRecords()
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
