/* yayf.go

TODO:

- [X] handle file ios for channel lists.
- [ ] handle file ios for previous records.


*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	SubscriptionsPath = "yayf.subs"
	RecordsPath       = "yayf.record"
)

type ChannelIDs struct {
	Cids []string `json:"Subscriptions"`
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
	//fmt.Println(records)
	//fmt.Println(records["cid1"])
	//fmt.Println(records["cid1"].(map[string]interface{})["link1"])
	return records
}

func main() {
	// First, need to read from subscriptions file - grab all the channel
	// ids into ChannelIDs struct.
	cids := GetSubs()
	for _, c := range cids {
		fmt.Println(c)
	}
	// Now, we need to read the previous records so that we know which
	// entries are new (new subscriptions added? new entries for existing
	// subscriptions?)
	records := GetRecords()
	for _, r := range records {
		fmt.Println(r)
	}
}
