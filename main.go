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

var messageProducerUrl = "localhost:9092"
var host = ":5126"

var server = "project-survey-diploma.database.windows.net"
var port = 1433
var user = "diplomaAdmin"
var password = "ktXw#A84tY!$7ig"
var database = "ProjectSurveyDb"

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
	requestHandler := request.NewHandler(serviceProvider.GetDbRepo(), serviceProvider.GetContextFiller(),
		serviceProvider.GetTargetingService(), serviceProvider.GetSurveyMarkupService())

	go func() {
		if err := fasthttp.ListenAndServe(host, requestHandler.Handle); err != nil {
			fmt.Printf("Error starting HTTP server: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Printf("HTTP server started on %v\n", host)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	//connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
	//	server, user, password, port, database)
	//
	//newReader := reader.NewSqlReader(connString)
	//
	//cache := dbcache.NewRepo(newReader)
	//
	//cache.Reload()
}
