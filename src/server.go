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
	defer func() {
		for _, device := range devices {
			device.Close()
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("views")))
	http.HandleFunc("/music", func(w http.ResponseWriter, r *http.Request) {
		log.Print("[INFO] access to /music")

		r.ParseForm()

		out, err := exec.Command("youtube-dl", "-x", "-g", r.PostForm["url"][0]).Output()
		if err != nil {
			log.Printf("[ERROR] Failed to exec youtube-dl: %v", err)
			return
		}

		url, err := url.Parse(string(out))
		if err != nil {
			log.Printf("[ERROR] Failed to parse url: %v", err)
			return
		}

		for _, device := range devices {
			if err := device.Play(ctx, url); err != nil {
				log.Printf("[ERROR] Failed to play: %v", err)
			}
		}

	})

	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "query: %s\n", r.URL.RawQuery)

	r.ParseForm()
	form := r.PostForm
	fmt.Fprintf(w, "form: \n%v\n", form)

	params := r.Form
	fmt.Fprintf(w, "parameter: \n%v\n", params)
}
