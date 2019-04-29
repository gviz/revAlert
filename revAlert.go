package main

import "context"

/*
curl -X POST -H 'Content-type: application/json'
--data '{"text":"Allow me to reintroduce myself!"}'
 https://hooks.slack.com/services/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
*/

const (
	url           = "https://hooks.slack.com/services/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"
	suricataTopic = "suricata-raw"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	chn := slackChannel{url: url}
	s := stream{
		channel: chn,
		brokers: []string{"localhost:9092"},
		ctx:     ctx,
	}
	s.Init()
	s.NewMonitor(suricataMonitor{topic: suricataTopic})
	//chn.postMessage(slackMsg{msg: "Test 2"})

	select {
	case <-ctx.Done():
	}
}
