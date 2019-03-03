package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"html/template"
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
		urlbody := r.PostForm["url"]
		if len(urlbody) <= 0 {
			log.Print("[Error] Failed to read url property: Body of Post doesn't contain URL of Music")
			return
		}

		cmd := exec.Command("youtube-dl", "--get-title", "-x", "-g", "--get-thumbnail", urlbody[0])
		outReader, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("Failed to exec youtube-dl: %v", err)
		}

		infoChan := make(chan MusicInfo)
		scanner := bufio.NewScanner(outReader)
		go func() {
			info := MusicInfo{URL: urlbody[0]}
			nItem := 0
			nFields := 0
			for scanner.Scan() {
				switch nFields % 3 {
				case 0:
					info.Title = scanner.Text()
				case 1:
					info.AudioURL = scanner.Text()
				default:
					info.ThumbnailURL = scanner.Text()
					infoChan <- info
					log.Printf("[Info] nItem: %v, info: %v", nItem, info)

					url, err := url.Parse(info.AudioURL)
					if err != nil {
						log.Printf("[ERROR] Failed to parse url: %v", err)
						return
					}

					media := homecast.MediaData{
						URL:   url,
						Title: info.Title,
					}
					medias := []homecast.MediaData{media}

					switch nItem {
					case 0:
						if err := device.QueueLoad(ctx, medias); err != nil {
							log.Printf("[ERROR] Failed to play: %v", err)
						}
					default:
						if err := device.QueueInsert(ctx, medias); err != nil {
							log.Printf("[ERROR] Failed to play: %v", err)
						}
					}
					nItem++
					info = MusicInfo{URL: urlbody[0]}
				}
				nFields++
			}
			close(infoChan)
		}()
		cmd.Start()

		infos := make([]MusicInfo, 25)
		i := 0
		for info := range infoChan {
			infos[i] = info
			i++
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
