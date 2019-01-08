package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"

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

		out, err := exec.Command("youtube-dl", "-x", "-g", urls[0]).Output()
		if err != nil {
			log.Printf("[ERROR] Failed to exec youtube-dl: %v", err)
			return
		}

		url, err := url.Parse(string(out))
		if err != nil {
			log.Printf("[ERROR] Failed to parse url: %v", err)
			return
		}

		if err := device.Play(ctx, url); err != nil {
			log.Printf("[ERROR] Failed to play: %v", err)
		}
	})

	log.Print("[INFO] start http server")
	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("[FATAL] ListenAndServe: ", err)
	}
}
