package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/mapstructure"
	"github.com/mmcdole/gofeed"
)

func setRoutes() (router *httprouter.Router) {
	// Load and parse templates (from binary or disk)
	templateBox = rice.MustFindBox(templatePath)
	templateBox.Walk("", newTemplate)

	// mux handler
	router = httprouter.New()

	// Index route
	router.GET("/", index)

	// Accounts route
	router.GET("/accounts", index)

	// Email route
	router.GET("/email", email)

	// Exchange route
	router.GET("/exchange", exchange)

	// Comments router
	router.GET("/comments/:id/", comments)

	// Feeds : ID
	//router.GET("/feed/0", feed)
	router.GET("/feed/:feedid/:page", feed)
	router.GET("/feedlist", feedList)
	router.GET("/feedtest", feedTest)
	router.POST("/feed/post", addFeed)

	// Forums : ID
	router.GET("/forum/:feedid/:page", forum)
	router.GET("/forumlist", forumList)

	// Links route
	router.GET("/links", links)

	// Links route
	router.GET("/streams", streams)

	// Search route
	router.GET("/search", search)

	// Settings route
	router.GET("/settings", settings)

	// Videos : ID
	router.GET("/video/:feedid/:page", video)
	router.GET("/videolist", videoList)

	// Proxy route
	router.GET("/proxy", proxy)

	router.GET("/test", test)

	// Serve static assets
	fs := rice.MustFindBox(staticPath).HTTPBox()
	router.ServeFiles("/static/*filepath", fs)

	return router
}

func accounts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := r.FormValue("url")
	feedType := r.FormValue("type")
	setFeed(url, feedType)

	w.Write([]byte("<html><head><meta http-equiv='refresh' content='0;url=" + r.Referer() + "' /></head></html>"))
}

// Feeds
func email(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//model := Model{Name: ps.ByName("name")}

	//feedhash := ps.ByName("feedhash")
	//page, _ := strconv.Atoi(ps.ByName("page"))

	model := Model{
		"Email",
		"email",
	}

	renderTemplate(w, templatePath+"\\email.html", model)

}

func exchange(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	days, err := strconv.Atoi(r.FormValue("days"))
	if days < 0 || err != nil {
		days = 1
	}

	exchange := Exchange{
		"Exchange",
		getOpenExchange(days),
		getCommoditiesExchange(days),
	}
	renderTemplate(w, templatePath+"\\exchange.html", exchange)
}

func addFeed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := r.FormValue("feedurl")
	feedType := r.FormValue("type")

	if feedType == "Bitchute" {
		feedType = "Video"
	} else if feedType == "Dailymotion" {
		feedType = "Video"
	} else if feedType == "Liveleak" {
		feedType = "Video"
	} else if feedType == "Vimeo" {
		feedType = "Video"
	} else if feedType == "Youtube" {
		feedType = "Video"
	}

	setFeed(url, feedType)

	w.Write([]byte("<html><head><meta http-equiv='refresh' content='0;url=" + r.Referer() + "' /></head></html>"))
}

// Feeds
func feed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//model := Model{Name: ps.ByName("name")}

	//feedhash := ps.ByName("feedhash")
	feedid, _ := strconv.Atoi(ps.ByName("feedid"))
	page, _ := strconv.Atoi(ps.ByName("page"))
	thisPage := page
	pageLimit := 16

	var (
		prev string
		next string
	)

	if page <= 1 {
		prev = strconv.Itoa(page + 1)
		next = strconv.Itoa(page)
	} else {
		prev = strconv.Itoa(page + 1)
		next = strconv.Itoa(page - 1)
		page = page * pageLimit
	}

	feedItems := FeedItems{
		"Feeds",
		thisPage,
		"/feed/" + strconv.Itoa(feedid) + "/" + prev,
		"/feed/" + strconv.Itoa(feedid) + "/" + next,
		getFeedItems(feedid, "Feed", pageLimit, page),
	}

	renderTemplate(w, templatePath+"\\feed.html", feedItems)
}

// List of feeds
func feedList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	feeds := Feeds{
		"Feeds",
		getFeed("Feed"),
	}
	renderTemplate(w, templatePath+"\\feedList.html", feeds)
}

// Feedtest
func feedTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderTemplate(w, templatePath+"\\feedTest.html", nil)
}

// Forums
func forum(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//model := Model{Name: ps.ByName("name")}

	//feedhash := ps.ByName("feedhash")
	feedid, _ := strconv.Atoi(ps.ByName("feedid"))
	page, _ := strconv.Atoi(ps.ByName("page"))
	thisPage := page
	pageLimit := 16

	var (
		prev string
		next string
	)

	if page <= 1 {
		prev = strconv.Itoa(page + 1)
		next = strconv.Itoa(page)
	} else {
		prev = strconv.Itoa(page + 1)
		next = strconv.Itoa(page - 1)
		page = page * pageLimit
	}

	feedItems := FeedItems{
		"Forums",
		thisPage,
		"/forum/" + strconv.Itoa(feedid) + "/" + prev,
		"/forum/" + strconv.Itoa(feedid) + "/" + next,
		getFeedItems(feedid, "Forum", pageLimit, page),
	}

	renderTemplate(w, templatePath+"\\forum.html", feedItems)

}

// List of forums
func forumList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	feeds := Feeds{
		"Forums",
		getFeed("Forum"),
	}
	renderTemplate(w, templatePath+"\\forumList.html", feeds)
}

// Videos
func video(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//model := Model{Name: ps.ByName("name")}

	//feedhash := ps.ByName("feedhash")
	feedid, _ := strconv.Atoi(ps.ByName("feedid"))
	page, _ := strconv.Atoi(ps.ByName("page"))
	thisPage := page
	pageLimit := 50

	var (
		prev string
		next string
	)

	if page <= 1 {
		prev = strconv.Itoa(page + 1)
		next = strconv.Itoa(page)
	} else {
		prev = strconv.Itoa(page + 1)
		next = strconv.Itoa(page - 1)
		page = page * pageLimit
	}

	feedItems := FeedItems{
		"Videos",
		thisPage,
		"/video/0/" + prev,
		"/video/0/" + next,
		getFeedItems(feedid, "Video", pageLimit, page),
	}

	renderTemplate(w, templatePath+"\\video.html", feedItems)

}

// List of videos
func videoList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	feeds := Feeds{
		"Videos",
		getFeed("Video"),
	}
	renderTemplate(w, templatePath+"\\videoList.html", feeds)
}

func links(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	links := Links{
		"Links",
		getLinks(),
	}
	renderTemplate(w, templatePath+"\\links.html", links)
}

func streams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	streams := Streams{
		"Streams",
		getStreams(),
	}
	renderTemplate(w, templatePath+"\\streams.html", streams)
}

func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderTemplate(w, templatePath+"\\index.html", nil)
}

func comments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID, _ := strconv.Atoi(ps.ByName("ID"))
	URL, _ := url.QueryUnescape(r.URL.Query().Get("url"))

	commentModal := CommentModal{
		URL + " - Comments",
		ID,
		URL,
	}

	renderTemplate(w, templatePath+"\\comments.html", commentModal)
}

func proxy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	URL, _ := url.QueryUnescape(r.URL.Query().Get("url"))
	//fmt.Println(URL)
	proxyURL, err := url.Parse("http://127.0.0.1:6901")
	checkError(err)
	thisClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	response, err := thisClient.Get(URL)

	if err != nil {
		checkError(err)
	} else {
		defer response.Body.Close()
		_, err = io.Copy(w, response.Body)
		if err != nil {
			checkError(err)
		}
	}
}

func search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q := r.URL.Query().Get("q")
	terms := strings.Split(q, " ")

	search := SearchResults{
		"Search",
		getSearch(terms),
	}
	renderTemplate(w, templatePath+"\\searchList.html", search)
}

func settings(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	settings := Settings{
		"Settings",
		//getSettings(1),
		[]Setting{},
	}
	renderTemplate(w, templatePath+"\\settings.html", settings)
}

func test(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Get feed information
	fp := gofeed.NewParser()
	//fp.RSSTranslator = NewMyCustomTranslator()
	feed, err := fp.ParseURL("http://www.dailymotion.com/rss/user/Styxhexenhammer666")
	checkError(err)

	//ytChan := YoutubeChannel{}
	//err3 := mapstructure.Decode(feed.Extensions, &ytChan)
	//checkError(err3)
	//PrettyPrint(this.Extensions)
	//PrettyPrint(feed.Title)

	for _, this := range feed.Items {

		test := Dailymotion{}
		err2 := mapstructure.Decode(this.Extensions, &test)
		checkError(err2)
		//PrettyPrint(this.Extensions)
		//PrettyPrint(vimeo)

		/*for _, this := range test.Media.Content[0].Children {
			if this[0].Name == "thumbnail" {
				r, _ := regexp.Compile(`\/.*\/(.*)_thumb`)
				match := r.FindAllStringSubmatch(this[0].Attrs["url"], 1)
				url := "https://www.liveleak.com/ll_embed?f=" + match[0][1]
				PrettyPrint(url)
			}
			if this[0].Name == "thumbnail" {
				PrettyPrint(this[0].Attrs["url"])
			}
			if this[0].Name == "description" {
				PrettyPrint(html.EscapeString(strings.Replace(this[0].Value, "'", "\\'", -1)))
			}
		}*/
		//PrettyPrint(test.Media)

		for _, this := range test.Media.Player {
			PrettyPrint(this.Attrs["url"])
		}

		for _, this := range test.Media.Thumbnail {
			PrettyPrint(this.Attrs["url"])
		}
	}

	renderTemplate(w, templatePath+"\\test.html", feed)
}

// PrettyPrint some shit
func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}
