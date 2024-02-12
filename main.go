package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"os"
	"os/signal"
	"project-survey-proceeder/context"
	"project-survey-proceeder/kafka"
	"project-survey-proceeder/pools"
	"project-survey-proceeder/request"
	"syscall"
)

var messageProducerUrl = "localhost:9092"
var host = ":5126"

func main() {
	messageProducer, err := kafka.InitProducer(messageProducerUrl)
	if err != nil {
		fmt.Printf("Error creating message producer: %v\n", err)
		os.Exit(1)
	}
	defer messageProducer.CloseConnection()

	parserPool := &fastjson.ParserPool{}
	uaPool := &pools.UserAgentPool{}

	prCtx := &context.ProceederContext{MessageProducer: messageProducer, ParserPool: parserPool, UserAgentPool: uaPool}
	requestHandler := request.RequestHandler{ProceederContext: prCtx}

	go func() {
		if err = fasthttp.ListenAndServe(host, requestHandler.HandleRequest); err != nil {
			fmt.Printf("Error starting HTTP server: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Printf("HTTP server started on %v\n", host)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
