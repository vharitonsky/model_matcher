//This package downloads models from database and uploads them into redis
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fzzy/radix/redis"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

type Configuration struct {
	SqlUrl          string
	CatIdQuery      string
	ModelLinesQuery string
	RedisAddr       string
}

var (
	configuration = Configuration{}
)

func init() {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Configured with ", configuration.SqlUrl)
}

func main() {
	c, err := redis.DialTimeout("tcp", configuration.RedisAddr, time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer c.Close()
	c.Cmd("select", 0)
	log.Print("Initializing models")
	models_count := 0
	start := time.Now()
	db, err := sql.Open("postgres", configuration.SqlUrl)
	if err != nil {
		log.Fatal(err)
	}
	cats := make([]string, 0)
	rows, err := db.Query(configuration.CatIdQuery, 0)
	for rows.Next() {
		var cat_id string
		rows.Scan(&cat_id)
		cats = append(cats, cat_id)
	}
	defer rows.Close()
	for _, cat_id := range cats {
		models := make([]string, 0)
		last_id := 0
		limit := 1000
		for {
			processed := 0
			rows, err := db.Query(configuration.ModelLinesQuery, 0, cat_id, last_id, limit)
			defer rows.Close()
			if err != nil {
				log.Fatal(err)
			}
			for rows.Next() {
				var model_line string
				rows.Scan(&last_id, &model_line)
				models = append(models, model_line)
				processed++
			}
			log.Print(fmt.Sprintf("Queried %d models", processed))
			if processed < limit {
				break
			}

		}
		models_count += len(models)
		c.Append("del", "_model_matcher_cat_"+cat_id)
		c.Append("lpush", "_model_matcher_cat_"+cat_id, models)
		c.GetReply()
		c.GetReply()
	}
	c.Append("del", "_model_matcher_cats")
	c.Append("lpush", "_model_matcher_cats", cats)
	c.GetReply()
	c.GetReply()
	var next_version string
	current_version, _ := c.Cmd("get", "_model_matcher_version").Int()
	next_version = strconv.Itoa(current_version + 1)
	c.Cmd("set", "_model_matcher_version", next_version)
	log.Print("Current version is ", next_version)
	elapsed := time.Since(start)
	log.Print(fmt.Sprintf("Matcher initialized with %d models from %d categories in %s", models_count, len(cats), elapsed))
}
