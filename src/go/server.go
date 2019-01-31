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
		urlbody := r.PostForm["url"]
		if len(urlbody) <= 0 {
			log.Print("[Error] Failed to read url property: Body of Post doesn't contain URL of Music")
			return
		}

		infos, err := getInfo(urlbody[0])
		if err != nil {
			log.Printf("[ERROR] Failed to get music infos: %v", err)
			return
		}

		medias := make([]homecast.MediaData, len(infos))
		for i, info := range infos {
			url, err := url.Parse(info.AudioURL)
			if err != nil {
				log.Printf("[ERROR] Failed to parse url: %v", err)
				return
			}
			medias[i] = homecast.MediaData{
				URL:   url,
				Title: info.Title,
			}
		}

		if err := device.QueueLoad(ctx, medias); err != nil {
			log.Printf("[ERROR] Failed to play: %v", err)
		}

		t, err := template.ParseFiles("./views/template.html")
		if err != nil {
			log.Printf("[ERROR] Failed to load template file: %v", err)
			return
		}
		t.Execute(w, infos[0])
	})

	log.Print("[INFO] start http server")
	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("[FATAL] ListenAndServe: ", err)
	}
}

type MusicInfo struct {
	URL          string
	AudioURL     string
	Title        string
	ThumbnailURL string
}

func getInfo(url string) ([]MusicInfo, error) {
	out, err := exec.Command("youtube-dl", "--get-title", "-x", "-g", "--get-thumbnail", url).Output()
	if err != nil {
		return nil, fmt.Errorf("Failed to exec youtube-dl: %v", err)
	}

	lines := strings.Split(string(out), "\n")
	musics := make([]MusicInfo, len(lines)/3)
	for i := 0; i < len(lines)/3; i++ {
		musics[i] = MusicInfo{
			URL:          url,
			AudioURL:     lines[3*i+1],
			Title:        lines[3*i+0],
			ThumbnailURL: lines[3*i+2],
		}
	}

	return musics, nil
}
