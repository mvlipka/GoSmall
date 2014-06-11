package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	//"io"
	//"io/ioutil"
	"./frontend"
	"math/rand"
	"net/http"
	//"strings"
	"time"
)

var tokenValues = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	conn, err := redis.Dial("tcp", "localhost:6379")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	htmlServe := frontend.HtmlServe{Directory: "views"}

	htmlServe.CacheHtml()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Requesting: " + r.URL.Path)
		if r.Method == "GET" {
			if r.URL.Path == "/" {
				htmlServe.Serve(&w, "index.html")
			} else {
				n, err := conn.Do("GET", r.URL.Path)
				if err != nil {
					fmt.Println(err)
				}
				url, _ := redis.String(n, err)
				http.Redirect(w, r, url, 302)
				fmt.Println(n)
			}
		}
	})

	http.HandleFunc("/_urlhandler", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			url := r.FormValue("url")
			if validateUrl(url) == false {
				return
			}
			generatedToken := generateToken()
			fmt.Println(url + " generated url: " + generatedToken)
			_, err := conn.Do("SET", "/"+generatedToken, url)
			if err != nil {
				fmt.Println(err)
			}

		}
	})

	http.ListenAndServe(":80", nil)
}

func generateToken() string {
	rand.Seed(int64(time.Now().Second()))
	temp := ""
	for i := 0; i < 6; i++ {
		temp += string(tokenValues[rand.Intn(len(tokenValues))])
	}
	return temp
}

func validateUrl(url string) bool {
	result := true
	if url[:7] != "http://" {
		result = false
	}
	if url[:8] == "https://" {
		result = true
	}

	return result
}
