package hiplog

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type HipChatEvent struct {
	Event string           `json:"event"`
	Item  HipChatEventItem `json:"item"`
}

type HipChatEventItem struct {
	Message HipChatEventMessage `json:"message"`
	Room    HipChatRoom         `json:"room"`
}

type HipChatEventMessage struct {
	Date string
	//File          HipChatFile `json:"file,omitempty"`
	From          string
	Message       string
	Color         string
	Type          string `json:"type,omitempty"`
	Id            string
	MessageFormat string `json:"message_format"`
}

type HipChatFile struct {
	Name     string
	Size     int
	ThumbUrl string
	Url      string
}

type HipChatRoom struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func writeToFile(f *os.File, sourceRoom HipChatRoom, sourceMessage HipChatEventMessage) error {
	msg := fmt.Sprintf("[%s] %s\n", sourceRoom.Name, sourceMessage.Message)
	_, err := f.WriteString(msg)
	return err
}

func handler(w http.ResponseWriter, r *http.Request, outFile *os.File) {
	var notifyEvent HipChatEvent

	json.NewDecoder(r.Body).Decode(&notifyEvent)

	err := writeToFile(outFile, notifyEvent.Item.Room, notifyEvent.Item.Message)
	if err != nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}
}

func main() {
	filePath := "hip.log"
	out, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("can not open the log file: %s, err: %v", filePath, err)
	}
	defer out.Close()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, out)
	})
	http.ListenAndServe(":8081", nil)
}
