package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	YT_FEED_URL  = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"
	SAMPLE_CH_ID = "UCwi3BrUqM4xStpbCyxsb3TA"
)

type EntryIndex struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Id    string `xml:"id"`
	Title string `xml:"title"`
}

func main() {
	curr_url := fmt.Sprintf(YT_FEED_URL, SAMPLE_CH_ID)
	resp, _ := http.Get(curr_url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	var eIndex EntryIndex
	xml.Unmarshal(bytes, &eIndex)
	for _, e := range eIndex.Entries {
		if len(e.Title) > 30 {
			e.Title = e.Title[:30] + "..."
		}
		fmt.Println(":: ", e.Id, e.Title)
	}
}
