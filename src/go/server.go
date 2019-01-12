package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"text/template"

	"github.com/ikasamah/homecast"
)

func main() {
	port := flag.Int("port", 3000, "Listen port")
	flag.Parse()

	ctx := context.Background()
	devices := homecast.LookupAndConnect(ctx)
	if len(devices) <= 0 {
		log.Fatal("[Fatal] Failed to find GoogleHome in this LAN")
	}
	device := devices[0]

	defer func() {
		device.Close()
	}()

	http.Handle("/", http.FileServer(http.Dir("views")))
	http.HandleFunc("/music", func(w http.ResponseWriter, r *http.Request) {
		log.Print("[INFO] access to /music")

		r.ParseForm()
		urls := r.PostForm["url"]
		if len(urls) <= 0 {
			log.Print("[Error] Failed to read url property: Body of Post doesn't contain URL of Music")
			return
		}

		out, err := exec.Command("youtube-dl", "--get-title", "-x", "-g", "--get-thumbnail", urls[0]).Output()
		if err != nil {
			log.Printf("[ERROR] Failed to exec youtube-dl: %v", err)
			return
		}

		lines := strings.Split(string(out), "\n")
		page := page{
			URL:          urls[0],
			AudioURL:     lines[1],
			Title:        lines[0],
			ThumbnailURL: lines[2],
		}

		url, err := url.Parse(page.AudioURL)
		if err != nil {
			log.Printf("[ERROR] Failed to parse url: %v", err)
			return
		}

		if err := device.Play(ctx, url); err != nil {
			log.Printf("[ERROR] Failed to play: %v", err)
		}

		t, err := template.ParseFiles("./views/template.html")
		if err != nil {
			log.Printf("[ERROR] Failed to load template file: %v", err)
			return
		}
		t.Execute(w, page)
	})

	log.Print("[INFO] start http server")
	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("[FATAL] ListenAndServe: ", err)
	}
}

type page struct {
	URL          string
	AudioURL     string
	Title        string
	ThumbnailURL string
}
