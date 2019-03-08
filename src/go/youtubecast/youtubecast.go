package youtubecast

import (
	"bufio"
	"context"
	"log"
	"net/url"
	"os/exec"

	"github.com/ikasamah/homecast"
)

type YoutubeDevice struct {
	castdevice  *homecast.CastDevice
	urlChan     chan string
	apidataChan chan APIData
}

type APIData struct {
	API   string `json:"api"`
	Title string `json:"title"`
}

type MusicInfo struct {
	URL          string
	AudioURL     string
	Title        string
	ThumbnailURL string
}

func NewYoutubeDevice(urlChan chan string, apidataChan chan APIData, castDevice *homecast.CastDevice) *YoutubeDevice {
	return &YoutubeDevice{
		castdevice:  castDevice,
		urlChan:     urlChan,
		apidataChan: apidataChan,
	}
}

func (ydevice *YoutubeDevice) Serve(ctx context.Context) {
	for urlstr := range ydevice.urlChan {
		log.Printf("[Info] YoutubeDevice.Serve: recieve url: %v", urlstr)
		ydevice.CastURL(ctx, urlstr)
	}
}

func (ydevice *YoutubeDevice) CastURL(ctx context.Context, urlstr string) {
	cmd := exec.Command("youtube-dl", "--get-title", "-x", "-g", "--get-thumbnail", urlstr)
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Failed to exec youtube-dl: %v", err)
	}
	scanner := bufio.NewScanner(outReader)
	cmd.Start()
	log.Printf("[Info] YoutubeDevice.CastURL: execute youtube-dl")

	nItem := 0
	nFields := 0
	info := MusicInfo{URL: urlstr}
	for scanner.Scan() {
		switch nFields % 3 {
		case 0:
			info.Title = scanner.Text()
		case 1:
			info.AudioURL = scanner.Text()
		default:
			info.ThumbnailURL = scanner.Text()
			log.Printf("[Info] nItem: %v, info: %v", nItem, info)

			youtubeURL, err := url.Parse(info.AudioURL)
			if err != nil {
				log.Printf("[ERROR] Failed to parse url: %v", err)
				return
			}

			media := homecast.MediaData{
				URL:   youtubeURL,
				Title: info.Title,
			}
			medias := []homecast.MediaData{media}

			switch nItem {
			case 0:
				if err := ydevice.castdevice.QueueLoad(ctx, medias); err != nil {
					log.Printf("[ERROR] Failed to play: %v", err)
				}
				ydevice.apidataChan <- APIData{
					API:   "QUEUE_LOAD",
					Title: info.Title,
				}
			default:
				if err := ydevice.castdevice.QueueInsert(ctx, medias); err != nil {
					log.Printf("[ERROR] Failed to play: %v", err)
				}
				ydevice.apidataChan <- APIData{
					API:   "QUEUE_INSERT",
					Title: info.Title,
				}
			}
			nItem++
			info = MusicInfo{URL: urlstr}
		}
		nFields++
	}
}
