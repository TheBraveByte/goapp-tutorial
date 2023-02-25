package main

import (
	"context"
	"encoding/gob"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/yusuf/go-app/driver"
	"github.com/yusuf/go-app/handlers"
	"github.com/yusuf/go-app/modules/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var app config.GoAppTools
var validate *validator.Validate

func main() {
	gob.Register(map[string]interface{}{})
	gob.Register(primitive.NewObjectID())
	
	InfoLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)
	ErrorLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)

	validate = validator.New()
	app.InfoLogger = InfoLogger
	app.ErrorLogger = ErrorLogger
	app.Validate = validate

	err := godotenv.Load()
	if err != nil {
		app.ErrorLogger.Fatal("No .env file available")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		app.ErrorLogger.Fatalln("mongodb uri string not found : ")
	}
	// connecting to the database
	client := driver.Connection(uri)
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			app.ErrorLogger.Fatal(err)
			return
		}
	}()
	appRouter := gin.New()

	goApp := handlers.NewGoApp(&app, client)
	Routes(appRouter, goApp)

	err = appRouter.Run()
	if err != nil {
		app.ErrorLogger.Fatal(err)
	}
}
