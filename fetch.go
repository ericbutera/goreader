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
			fmt.Printf("\n\n")
			// fmt.Printf("%s\n", item)
			err = c.Insert(bson.M{
				"url":     item.Link,
				"title":   item.Title,
				"content": item.Content,
				"date":    item.Date,
				"created": time.Now(),
				"updated": time.Now()})
		}

		fmt.Printf("%v\n", time.Now())
		fmt.Printf("%v\n", time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC))
	}

	result := Item{}
	c.Find(bson.M{"title": "moo"}).One(&result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("mgo result %+v\n", result)
}
