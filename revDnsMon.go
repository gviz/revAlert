package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type revEvent struct {
	Name    string `json:"name"`
	Type    string `json:"event_type"`
	Host    string `json:"host"`
	Msg     string `json:"message"`
	encoded []byte
	err     error
}
type revMonitor struct {
	topic string
}

func (r revMonitor) GetTopic() string {
	return r.topic
}

func (r revMonitor) Run(es eventStream) {

	for {
		select {
		case msg := <-es.Reader():
			go func(val []byte) {
				var js revEvent
				err := json.Unmarshal(val, &js)
				if err != nil {
					log.Println("Error getting json object")
					log.Println(err)
					return
				}

				text := fmt.Sprintf("%s: %s -> %s", js.Name, js.Host, js.Msg)
				es.PostMsg(text)

			}(msg)
		}
	}

}
