package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yusuf/go-app/modules/config"
	"github.com/yusuf/go-app/modules/database"
	"github.com/yusuf/go-app/modules/database/query"
	"github.com/yusuf/go-app/modules/encrypt"
	"github.com/yusuf/go-app/modules/model"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
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

func (ga *GoApp) SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user *model.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.Password, _ = encrypt.Hash(user.Password)

		if err := ga.App.Validate.Struct(&user); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				ga.App.InfoLogger.Println(err)
				return
			}
		}

		ok, status, err := ga.DB.InsertUser(user)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, errors.New("error while adding new user"))
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if !ok {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		switch status {
		case 1:
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Registered Successfully",
			})
			return

		case 2:
			ctx.JSON(http.StatusFound, gin.H{
				"message": "Existing Account, Go to the Login page",
			})
			return

		}
	}
}
