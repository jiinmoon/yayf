yayf : Yet, Another Youtube Feed-aggregator
===========================================

![yayf-do-I-really-needs-it?](resources/yayf-frontpage.jpg "Here's an attempt
at drawing in xkcd style. Drawing with mouse sucks.")

This is a small project with go that iterates over the registered youtube
subscription feeds, and notifies you of any new uploads. And returns the full
list of youtube-urls that which can be used to download further.

There are already feed aggregators and another softwares that can do this but
I am doing this for more of an exploration of golang.

UPDATES
-------

- First proto yayf.proto works.
- Now, need to find efficient strcuts to store various data.
- Once feed is received from YT, for each CID, its entry should be checked for
  whether the entry is new or not.
- Feed data struct should match the records.

TODO
----

- [x] Figure out xml schema of youtube RSS feeds.
- [x] Fetch a single channel feeds.
- [x] Use goroutine to apply over multiple channels at same time.
- [x] First proto yayf.go finished; proof of concept is valid.
- [x] Implement GetFeeds().
- [ ] Resolve record v feeds.

Contacts
--------

You may contact me at <mailto:jmoon@jiinmoon.com>.

