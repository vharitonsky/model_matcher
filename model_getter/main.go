package main

import (
	"bufio"
	"fmt"
	"github.com/fzzy/radix/redis"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer c.Close()
	c.Cmd("select", 0)
	log.Print("Initializing models")
	models_count, categories_count := 0, 0
	start := time.Now()
	cat_file, err := os.Open("../data/cats.txt")
	defer cat_file.Close()
	if err != nil {
		log.Fatal(err)
	}
	cats := make([]string, 0)
	model_file_reader := bufio.NewReader(cat_file)
	for {
		cat_line, err := model_file_reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		cat_id := string(cat_line[:len(cat_line)-1])
		models := make([]string, 0)
		categories_count += 1

		file, err := os.Open("../data/models/m_" + cat_id + ".txt")
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		reader := bufio.NewReader(file)
		for {
			model_line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatal(err)
				}
			}
			models = append(models, string(model_line[:len(model_line)-1]))
		}
		models_count += len(models)
		c.Append("del", "_model_matcher_cat_"+cat_id)
		c.Append("lpush", "_model_matcher_cat_"+cat_id, models)
		c.GetReply()
		c.GetReply()
		cats = append(cats, cat_id)
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
	log.Print(fmt.Sprintf("Matcher initialized with %d models from %d categories in %s", models_count, categories_count, elapsed))
}
