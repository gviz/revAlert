package main

import (
	"log"

	"github.com/nlopes/slack"
)

type slackChannel struct {
	url string
}

type slackMsg struct {
	msg   string
	mType int
}

func (s *slackChannel) postMessage(text slackMsg) {
	attach := slack.Attachment{
		Text: text.msg,
	}

	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attach},
	}

	err := slack.PostWebhook(s.url, &msg)
	if err != nil {
		log.Println(err)
	}

}
