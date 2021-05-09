package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/html"
)

func services() {
	exchangeServiceHandle := scheduler.Every(1).Hours()
	exchangeServiceHandle.Run(exchangeService)

	cleanServiceHandle := scheduler.Every(1).Day()
	cleanServiceHandle.Run(cleanService)

	feedServiceHandle := scheduler.Every(10).Minutes()
	feedServiceHandle.Run(feedService)
}

// Exchange service (updates and adds exchange rates)
var exchangeService = func() {
	setOpenExchange()
	setCommoditiesExchange(2)
}

// Clean service (removes old feed entries)
var cleanService = func() {
	cleanFeed(30) // 30 days
}

// Feed service (updates and adds new feed items)
var feedService = func() {
	var (
		wg         sync.WaitGroup
		allDomains []string
	)

	for _, feed := range getFeed("0") {
		// Don't repeatedly hit the same domains
		delay := 1 // (seconds)
		thisURL, _ := url.Parse(feed.FeedLink)
		match := 0

		for _, domain := range allDomains {
			if thisURL.Host == domain {
				if match == 0 {
					match = delay
				} else {
					match = match + delay
				}
			}
		}

		allDomains = append(allDomains, thisURL.Host)
		sleep := time.Duration(match) * time.Second

		wg.Add(1)
		go func(feed Feed, delay time.Duration) {
			time.Sleep(sleep)

			setFeed(feed.FeedLink, feed.FeedType)
			fmt.Println("Updated : " + feed.FeedType + " : " + feed.FeedLink)
			//setLog("Update", "Feed Service: "+feed.FeedType+" : "+feed.FeedLink)

			wg.Done()
		}(feed, sleep)
	}

	wg.Wait()

	//warmCache()
}

// CacheProxy bla
/*func CacheProxy() {
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = "127.0.0.1"
		},
	}

	handler := httpcache.NewHandler(httpcache.NewMemoryCache(), proxy)
	handler.Shared = true

	log.Printf("proxy listening on http://%s", "0.0.0.0:6980")
	log.Fatal(http.ListenAndServe("0.0.0.0:6980", handler))
}*/

// NewProxy sets up proxy instance
func NewProxy() *goproxy.ProxyHttpServer {
	// Proxy server
	proxy := goproxy.NewProxyHttpServer()
	//proxy.Verbose = true
	proxy.OnResponse(goproxy_html.IsHtml).Do(goproxy_html.HandleString(
		func(s string, ctx *goproxy.ProxyCtx) string {
			//s = strings.Replace(s, "http://", "http://127.0.0.1:6900/proxy?url=http://", -1)
			//s = strings.Replace(s, "https://", "http://127.0.0.1:6900/proxy?url=https://", -1)

			return s
		}))
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.Header.Del("X-Frame-Options")
			r.Header.Set("Access-Control-Allow-Origin", "*")

			return r, nil
		})

	return proxy
}

func warmCache() {
	// Prewarm all multi-page feeds with requests.
	// Total to cache for multi-page
	pages := []string{"/forum/0/", "/feed/0/", "/video/0/"}
	n := 5
	for _, page := range pages {

		for i := 1; i <= n; i++ {
			go func(page string, i int) {
				get("http://127.0.0.1:6900" + page + strconv.Itoa(i))
			}(page, i)
		}

	}
}

func get(url string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		_, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Warmed : " + url)
		}
	}
}
