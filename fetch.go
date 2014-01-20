package main

import (
	"fmt"
	"github.com/SlyMarbo/rss"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

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

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("goreader").C("items")

	// feed, err := rss.Fetch("http://localhost:8080/reddit.xml")
	feeds := []string{"http://localhost:8080/reddit.xml", "http://localhost:8080/hn.xml"}

	updatedItem := Item{}
	for _, url := range feeds {
		fmt.Printf("fetching url %s\n", url)

		feed, err := rss.Fetch(url)
		if err != nil {
			fmt.Printf("error fetching %s\n", url)
			continue
		}

		for _, item := range feed.Items {
			fmt.Printf("title %v \n", item.Title)
			fmt.Printf("content %v \n", item.Content)
			fmt.Printf("link %v \n", item.Link)
			fmt.Printf("date %+v \n", item.Date)
			fmt.Printf("id %v \n", item.ID)
			fmt.Printf("read %v \n", item.Read)
			fmt.Printf("\n\n") // fmt.Printf("%s\n", item)

			change := mgo.Change{
				Update: bson.M{"$set": bson.M{
					"title":   item.Title,
					"content": item.Content,
					"updated": time.Now(),
				}},
				ReturnNew: true,
			}

			// gettin' messy
			info, err := c.Find(bson.M{"url": item.Link}).Apply(change, &updatedItem)
			fmt.Printf("c.Find info %+v\nerr %+v\n", info, err)
			if err == nil {
				// updatedItem
				// save feed_item relation
			} else {
				// err if attempt at updating fails, so insert it!
				fmt.Printf("update failed, inserting %v\n", url)
				err = c.Insert(bson.M{
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
		c.Find(bson.M{"title": "moo"}).One(&result)
		if err != nil {
			panic(err)
		}
		fmt.Printf("mgo result %+v\n", result)
	*/
}
