package main

import (
	"github.com/go-redis/redis"
	"encoding/json"
	"os"
	"bufio"
	"io"
	"path"
	"strings"
	"net/http"
	"net"
	"time"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"sync"
	"log"
)

var fileStr string
var wg sync.WaitGroup

var httpClient = http.Client{
	Transport: &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			deadline := time.Now().Add(5 * time.Second)
			c, err := net.DialTimeout(netw, addr, time.Second*5)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
		},
	},
}

var conn = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func push(k string, host string, title string) {
	jsonData := make(map[string]string)
	jsonData["host"] = host
	jsonData["title"] = title
	b, _ := json.Marshal(jsonData)
	err := conn.LPush(k, b).Err()
	if err != nil {
		panic(err)
	}
}

func run(k string, url string) {
	r, err := httpClient.Get("http://" + url)
	if err != nil {
		if r==nil{
			wg.Done()
			return
		}
		panic(err)
		wg.Done()
		return
	}
	if r != nil {
		if r.StatusCode == 200 {
			node, err := html.Parse(r.Body)
			if err == nil {
				r.Body.Close()
			}
			if err != nil {
				log.Fatal(err)
			}
			doc := goquery.NewDocumentFromNode(node)
			push(k, url, doc.Find("title").Text())
		}
	}

	defer wg.Done()
}

func main() {
	fileStr = "./test.com-host.list"
	fKey := path.Base(fileStr)
	f, err := os.Open(fileStr)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	rd := bufio.NewReader(f)
	for {
		s, err := rd.ReadString('\n')
		if (err != nil || io.EOF == err) {
			break
		}
		wg.Add(1)
		s = strings.Replace(s, " ", "", -1)
		s = strings.Replace(s, "\r", "", -1)
		s = strings.Replace(s, "\n", "", -1)
		s = strings.Replace(s, "\t", "", -1)
		go run(fKey, s)
	}
	wg.Wait()
}
