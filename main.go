package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/lox/httpcache"
	//"./lib/user"
)

//type Config struct {
//}

// Model of stuff to render a page
type Model struct {
	Title string
	Name  string
}

// CommentModal stuffs
type CommentModal struct {
	Title string
	ID    int
	URL   string
}

// Exchange sources
type Exchange struct {
	Title       string
	Exchanges   []ExchangeRates
	Commodities []ExchangeRatesCommodities
}

// Feeds sources
type Feeds struct {
	Title string
	Feeds []Feed
}

// Feed source
type Feed struct {
	ID          int
	FeedHash    string
	FeedType    string
	Title       sql.NullString
	Description sql.NullString
	Link        sql.NullString
	FeedLink    string
	Updated     string
	Published   string
	Author      sql.NullString
	Language    sql.NullString
	Image       sql.NullString
	Copyright   sql.NullString
}

// FeedItems sources
type FeedItems struct {
	Title string
	Page  int
	Prev  string
	Next  string
	Feeds []FeedItem
}

// FeedItem source
type FeedItem struct {
	ID           int
	FeedItemHash string
	FeedID       int
	FeedTitle    sql.NullString
	Title        sql.NullString
	Description  sql.NullString
	Categories   sql.NullString
	Content      sql.NullString
	FeedLink     string
	Link         string
	Updated      string
	Published    string
	Language     sql.NullString
	FeedAuthor   sql.NullString
	Author       sql.NullString
	FeedImage    sql.NullString
	Image        sql.NullString
	Copyright    sql.NullString
}

// Links sources
type Links struct {
	Title string
	Links []Link
}

// Link sources
type Link struct {
	ID    int
	Title string
	URL   string
}

// Streams sources
type Streams struct {
	Title   string
	Streams []Stream
}

// Stream sources
type Stream struct {
	ID       int
	Name     string
	Desc     string
	HostID   int
	Owner    string
	Type     string
	Head     string
	Proto    string
	Domain   string
	URL      string
	Params   string
	Interval int
}

// SearchResults All
type SearchResults struct {
	Title         string
	SearchResults []SearchResult
}

// SearchResult individual
type SearchResult struct {
	Title     string
	Link      string
	Text      string
	Type      string
	Published string
	Updated   string
}

// Settings ALL
type Settings struct {
	Title   string
	Setting []Setting
}

// Setting for feeder
type Setting struct {
	ID               int
	FeedInterval     int
	ExchangeInterval int
	EmailHostAddress string
}

// Templates with functions available to them
var staticPath = "assets\\static"
var templatePath = "assets\\templates"
var templates = template.New("").Funcs(templateMap)
var templateBox *rice.Box

func newTemplate(path string, _ os.FileInfo, _ error) error {
	if path == "" {
		return nil
	}
	templateString, err := templateBox.String(path)
	if err != nil {
		log.Panicf("Unable to parse: path=%s, err=%s", path, err)
	}
	templates.New(filepath.Join(templatePath, path)).Parse(templateString)
	return nil
}

// Render a template given a model
func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// The server itself
func main() {
	// Set DB defaults (alters connection settings)
	initDB()
	dbDefaults()
	//loadProcedures()
	//os.Setenv("ALL_PROXY", "192.168.1.146:9050")

	//cache := httpcache.Cache(router, 20*time.Second)

	//setFeed("https://www.bitchute.com/profile/YPqBtYr4C8zI/", "Video")
	//getSearch([]string{"this"})

	//
	// Services
	//
	go services()

	// Set all path routes, Serve interface
	go func() {
		log.Fatal(http.ListenAndServe(":6900", setRoutes()))
	}()

	/*go func() {
		CacheProxy()
	}()*/

	//warmCache()

	// Serve Proxy
	go log.Fatal(http.ListenAndServe(":6901", NewProxy()))
}
