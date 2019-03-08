package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/forestaa/play-youtube-googlehome/src/go/router"
	"github.com/forestaa/play-youtube-googlehome/src/go/youtubecast"

	"github.com/ikasamah/homecast"
)

func main() {
	port := flag.Int("port", 4000, "Listen port")
	flag.Parse()

	ctx := context.Background()
	urlChan := make(chan string)
	apidataChan := make(chan youtubecast.APIData)

	go func() {
		devices := homecast.LookupAndConnect(ctx)
		if len(devices) <= 0 {
			log.Fatal("[Fatal] Failed to find GoogleHome in this LAN")
		}
		device := devices[0]
		// var device *homecast.CastDevice

		defer func() {
			device.Close()
		}()

		ydevice := youtubecast.NewYoutubeDevice(urlChan, apidataChan, device)
		ydevice.Serve(ctx)
	}()

	router := router.NewRouter(urlChan, apidataChan)
	http.Handle("/", router)
	log.Print("[INFO] start http server")
	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("[FATAL] ListenAndServe: ", err)
	}
}
