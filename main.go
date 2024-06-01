package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/signal"
	"project-survey-proceeder/internal/configuration"
	"project-survey-proceeder/internal/request"
	"project-survey-proceeder/internal/services"
	"syscall"
)

func main() {
	//messageProducer, err := kafka.InitProducer(messageProducerUrl)
	//if err != nil {
	//	fmt.Printf("Error creating message producer: %v\n", err)
	//	os.Exit(1)
	//}
	//defer messageProducer.CloseConnection()

	parser := configuration.NewParser()
	config, err := parser.Parse("appsettings.json")
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	serviceProvider := services.NewProvider(config)
	requestHandler := request.NewHandler(serviceProvider.GetDbRepo(), serviceProvider.GetUnitContextFiller(),
		serviceProvider.GetEventContextFiller(),
		serviceProvider.GetTargetingService(), serviceProvider.GetSurveyMarkupService())

	go func() {
		if err := fasthttp.ListenAndServe(config.Host, requestHandler.Handle); err != nil {
			fmt.Printf("Error starting HTTP server: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Printf("HTTP server started on %v\n", config.Host)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
