package main

import (
	"html"
	"html/template"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/torden/go-strutil"
)

var (
	templateMap = template.FuncMap{
		"Multiply": func(a float64, b float64) float64 {
			return a * b
		},
		"Divide": func(a float64, b float64) float64 {
			return a / b
		},
		"DivideRound": func(a float64, b float64) float64 {
			return float64(int(a/b*10000)) / 10000
		},
		"ToGramsRound": func(a float64) float64 {
			return float64(int(a/31.10348*10000)) / 10000
		},

		"Escape": func(s string) template.HTML {
			return template.HTML(html.EscapeString(s))
		},

		"Last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
		"Lower": func(s string) string {
			return strings.ToLower(s)
		},
		"Sanitize": func(s string) string {
			return sanitize.HTML(sanitize.Accents(s))
		},
		"Strip": func(s template.HTML) string {
			strproc := strutils.NewStringProc()
			ret, err := strproc.StripTags(string(s))
			checkError(err)
			return ret
		},

		"TrimSpace": func(s string) string {
			return strings.Join(strings.Fields(s), "	")
		},

		"ToString": func(i int) string {
			return strconv.Itoa(i)
		},
		"Upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"Unescape": func(s string) template.HTML {
			s = html.UnescapeString(s)
			replace := map[string]string{"\\\\": "\\", `\'`: "'", "\\\\0": "\\0", "\\n": "\n", "\\r": "\r", `\"`: `"`, "\\Z": "\x1a"}

			for b, a := range replace {
				s = strings.Replace(s, b, a, -1)
			}

			return template.HTML(s)
		},
		"UnixTimeReadable": func(i int) string {
			inputTime := time.Unix(int64(i), 0)
			hour := inputTime.Hour()
			minute := inputTime.Minute()
			second := inputTime.Second()

			//layout := "Mon 2 15:04:05"
			//output, err := time.Parse(layout, inputTime.String())
			//checkError(err)

			return strconv.Itoa(hour) + " : " + strconv.Itoa(minute) + " : " + strconv.Itoa(second)
			//return output.String()
		},
	}
)

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
