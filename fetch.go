package main

import (
	"fmt"
	"github.com/SlyMarbo/rss"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type Feed struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Url         string        `json:"url"`
	Title       string        `json:"title"`
	Unread      int           `json:"unread"`
	Total       int           `json:"total"`
	DateCreated time.Time     `json:"dateCreated"`
}

type Item struct {
	// id ????
	Id      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Url     string        `json:"url"`
	Title   string        `json:"title"`
	Content string        `json:"content"`
	Date    time.Time     `json:"date"`
	Created time.Time     `json:"created"`
	Updated time.Time     `json:"updated"`
}

type FeedItem struct {
	FeedId      string
	ItemId      string
	DateCreated time.Time
	IsRead      bool
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("goreader")
	cItems := db.C("items")
	cFeeds := db.C("feeds")

	updatedItem := Item{}
	feeds := []Feed{}

	err = cFeeds.Find(nil).All(&feeds)
	if err != nil {
		panic("Unable to fetch feeds from db")
	}

	for _, feed := range feeds {
		fmt.Printf("fetching url %s\n", feed.Url)

		rssFeed, err := rss.Fetch(feed.Url)
		if err != nil {
			fmt.Printf("error fetching %s\n", feed.Url)
			continue
		}

		for _, item := range rssFeed.Items {
			// fmt.Printf("title %v \n", item.Title)
			//fmt.Printf("link %v \n", item.Link)
			//fmt.Printf("content %v \n", item.Content)
			//fmt.Printf("date %+v \n", item.Date)

			change := mgo.Change{
				Update: bson.M{"$set": bson.M{
					"title":   item.Title,
					"content": item.Content,
					"updated": time.Now(),
				}},
				ReturnNew: true,
			}

			// gettin' messy
			_, err := cItems.Find(bson.M{"url": item.Link}).Apply(change, &updatedItem)
			if err == nil {
				fmt.Printf("update %v\n", item.Link)
				// updatedItem
				// save feed_item relation
			} else {
				// err if attempt at updating fails, so insert it!
				fmt.Printf("insert %v\n", item.Link)
				err = cItems.Insert(bson.M{
					"url":     item.Link,
					"title":   item.Title,
					"content": item.Content,
					"date":    item.Date,
					"created": time.Now(),
				})
				if err == nil {
					// save feed_item relation
				} else {
					fmt.Printf("unable to insert record %+v\nerr %v\n", item, err)
				}
			}
		}

		// fmt.Printf("%v\n", time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)) if date is this, use now
	}

	/*
		result := Item{}
		cItems.Find(bson.M{"title": "moo"}).One(&result)
		if err != nil {
			panic(err)
		}
		fmt.Printf("mgo result %+v\n", result)
	*/
}
