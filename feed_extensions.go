package main

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
	ext "github.com/mmcdole/gofeed/extensions"
)

type Dailymotion struct {
	Media DailymotionMedia
}

type DailymotionMedia struct {
	Player    []ext.Extension
	Thumbnail []ext.Extension
}

func decodeDailymotion(item ext.Extensions, dailymotion Dailymotion) (string, string) {
	var (
		content string
		image   string
	)

	err := mapstructure.Decode(item, &dailymotion)
	checkError(err)

	for _, this := range dailymotion.Media.Player {
		r, _ := regexp.Compile(`\/.*\/(.*)`)
		match := r.FindAllStringSubmatch(this.Attrs["url"], 1)
		content = "http://www.dailymotion.com/embed/video/" + match[0][1]
	}

	for _, this := range dailymotion.Media.Thumbnail {
		image = this.Attrs["url"]
	}

	return content, image
}

type Liveleak struct {
	Media LiveleakMedia
}

type LiveleakMedia struct {
	Content []ext.Extension
}

func decodeLiveleak(item ext.Extensions, liveleak Liveleak) (string, string) {
	var (
		content string
		image   string
	)

	err := mapstructure.Decode(item, &liveleak)
	checkError(err)

	for _, this := range liveleak.Media.Content[0].Children {
		if this[0].Name == "thumbnail" {
			r, _ := regexp.Compile(`\/.*\/(.*)_thumb_`)
			match := r.FindAllStringSubmatch(this[0].Attrs["url"], 1)

			if match != nil {
				embed := "https://www.liveleak.com/ll_embed?f=" + match[0][1]
				//content = embed

				// Grab the source location from remote player
				response, err := http.Get(embed)
				if err != nil {
					log.Fatal(err)
				} else {
					defer response.Body.Close()
					body, err := ioutil.ReadAll(response.Body)
					if err != nil {
						log.Fatal(err)
					} else {
						//fmt.Println(string(body))
						r, _ := regexp.Compile(`(?U)source src="(.*)"`)
						match := r.FindAllStringSubmatch(string(body), 1)
						if match != nil {
							content = embed
						}
					}
				}
			}

			image = this[0].Attrs["url"]
		}
	}

	return content, image
}

type Vimeo struct {
	Media VimeoMedia
}

type VimeoMedia struct {
	Content []ext.Extension
}

func decodeVimeo(item ext.Extensions, vimeo Vimeo) (string, string) {
	var (
		content string
		image   string
	)

	err := mapstructure.Decode(item, &vimeo)
	checkError(err)

	for _, this := range vimeo.Media.Content[0].Children {
		if this[0].Name == "player" {
			content = this[0].Attrs["url"]
		}
		if this[0].Name == "thumbnail" {
			image = this[0].Attrs["url"]
		}
	}

	return content, image
}

type Youtube struct {
	YT    YoutubeVideo
	Media YoutubeMedia
}

type YoutubeChannel struct {
	ChannelID []ext.Extension
}

type YoutubeVideo struct {
	ChannelID []ext.Extension
	VideoID   []ext.Extension
}

type YoutubeMedia struct {
	Group []ext.Extension
}

func decodeYoutube(item ext.Extensions, youtube Youtube) (string, string, string) {
	var (
		content     string
		image       string
		description string
	)

	err := mapstructure.Decode(item, &youtube)
	checkError(err)

	for _, this := range youtube.Media.Group[0].Children {
		if this[0].Name == "content" {
			content = this[0].Attrs["url"]
		}
		if this[0].Name == "thumbnail" {
			image = this[0].Attrs["url"]
		}
		if this[0].Name == "description" {
			description = html.EscapeString(strings.Replace(this[0].Value, "'", "\\'", -1))
		}
	}

	return content, image, description
}
