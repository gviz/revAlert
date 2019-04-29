package main

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
)

//Input Stream Handling .
type stream struct {
	brokers    []string
	readTopic  string
	alertTopic string
	consumer   sarama.Consumer
	pConsumer  sarama.PartitionConsumer
	producer   sarama.AsyncProducer
	ctx        context.Context
	channel    slackChannel
	msgChan    chan slackMsg
	streams    []eventStream
}

type eventReader func() []byte

type eventStream struct {
	ctx        context.Context
	cancel     context.CancelFunc
	stream     chan []byte
	msgChannel chan slackMsg
}

func (ev *eventStream) Cancel() {
	ev.cancel()
}

func (ev *eventStream) PostAlert(msg slackMsg) {
	go func() {
		ev.msgChannel <- msg
	}()
}

func (ev *eventStream) PostMsg(msg string) {
	go func() {
		ev.msgChannel <- slackMsg{msg: msg}
	}()
}

func (ev *eventStream) Reader() chan []byte {
	return ev.stream
}

type eventMonitor interface {
	GetTopic() string       // Name of kafka topic
	Run(stream eventStream) //Actual logic
}

func (s *stream) Close() {
	for _, strm := range s.streams {
		strm.cancel()
	}
}

func (s *stream) Init() {
	s.msgChan = make(chan slackMsg)
	go func() {
		ctx, cancel := context.WithCancel(s.ctx)
		defer cancel()
		for {
			select {
			case <-ctx.Done():
			case msg := <-s.msgChan:
				s.channel.postMessage(msg)
			}
		}
	}()
}
func (s *stream) AddStream(topic string) eventStream {
	if len(s.brokers) == 0 {
		log.Panic("No brokers specified for kafka..")
	}
	log.Printf("Kafka Brokers: %v", s.brokers)

	consumer, err := sarama.NewConsumer(s.brokers, nil)
	if err != nil {
		log.Fatalln(err)
	}

	s.consumer = consumer
	pConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalln(err)
	}

	es := eventStream{
		ctx:        s.ctx,
		stream:     make(chan []byte),
		msgChannel: s.msgChan,
	}
	sCtx, cancel := context.WithCancel(s.ctx)
	es.cancel = cancel
	go func(es eventStream) {
		defer es.cancel()
		for {
			select {
			case msg := <-pConsumer.Messages():
				es.stream <- msg.Value
			case <-sCtx.Done():
				close(es.stream)
				pConsumer.Close()
			}
		}

	}(es)
	return es
}

func (s *stream) NewMonitor(monitor eventMonitor) {
	handler := s.AddStream(monitor.GetTopic())
	go monitor.Run(handler)
}
