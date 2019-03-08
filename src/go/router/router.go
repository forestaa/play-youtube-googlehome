package router

import (
	"log"
	"net/http"

	"github.com/forestaa/play-youtube-googlehome/src/go/youtubecast"
	"github.com/gorilla/websocket"
)

type Router struct {
	urlChan     chan string
	apidataChan chan youtubecast.APIData
}

type Message struct {
	Url string `json:"url"`
}

func NewRouter(urlChan chan string, apidataChan chan youtubecast.APIData) *Router {
	return &Router{
		urlChan:     urlChan,
		apidataChan: apidataChan,
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[Error] Router.ServeHTTP: socket server configuration error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go func() {
		for data := range rt.apidataChan {
			if err := ws.WriteJSON(data); err != nil {
				log.Printf("[Error] Router.ServeHTTP: failed to write json to websocket: %v", data)
			}
			log.Printf("[Info] Router.ServeHTTP: write to websocket: %v", data)
		}
	}()

	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			log.Printf("[Error] Router.ServeHTTP: failed to read json: %v", err)
		}
		log.Printf("[Info] Router.ServeHTTP: receive from websocket: %v", msg)

		rt.urlChan <- msg.Url
	}
}
