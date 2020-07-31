/* yayf.go

TODO:

- [ ] handle file ios for channel lists.
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

func getSubs() []string {
	var CI ChannelIDs
	data, err := ioutil.ReadFile(SubscriptionsPath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &CI)
	if err != nil {
		log.Fatal(err)
	}
	return CI.Cids
}

func main() {
	// First, need to read from subscriptions file - grab all the channel
	// ids into ChannelIDs struct.
	Cids := getSubs()
	for _, c := range Cids {
		fmt.Println(c)
	}
}
