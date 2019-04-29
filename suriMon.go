package main

import (
	"fmt"
	"log"
)

type suricataMonitor struct {
	topic string
}

func (s suricataMonitor) GetTopic() string {
	return s.topic
}

func (s suricataMonitor) Run(es eventStream) {

	for {
		select {
		case msg := <-es.Reader():
			go func(val []byte) {
				js := NewRevJson(val)
				if js == nil {
					log.Println("Error getting json object")
					return
				}
				res, _ := js.getValStr("event_type")
				if res != "alert" {
					return
				}
				/*
					result, _ := js.compareValStr("alert.category", "Generic Protocol Command Decode")
					if result == 0 {
						return
					}
				*/
				desc, err := js.getValStr("alert.alert.signature")
				if err != nil {
					log.Println("Error decoding alert.alert.signature")
					log.Println(err)
					return
				}

				src, err := js.getValStr("alert.alert.src_ip")
				if err != nil {
					log.Println("Error decoding alert.alert.src_ip")
					log.Println(err)
					return
				}

				dst, err := js.getValStr("alert.alert.dst_ip")
				if err != nil {
					log.Println("Error decoding alert.alert.dst_ip")
					log.Println(err)
					return
				}
				text := fmt.Sprintf("%s: %s -> %s", desc, src, dst)
				es.PostMsg(text)

			}(msg)
		}
	}

}
