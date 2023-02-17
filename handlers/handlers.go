package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yusuf/go-app/modules/config"
	"github.com/yusuf/go-app/modules/database"
	"github.com/yusuf/go-app/modules/database/query"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type GoApp struct {
	App *config.GoAppTools
	DB  database.DBRepo
}

func NewGoApp(app *config.GoAppTools, db *mongo.Client) *GoApp {
	return &GoApp{
		App: app,
		DB:  query.NewGoAppDB(app, db),
	}
}

func (ga *GoApp) Home() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"resp": "Welcome to Go App home page"})
	}
}
