package query

import (
	"context"
	"errors"
	"github.com/yusuf/go-app/modules/config"
	"github.com/yusuf/go-app/modules/database"
	"github.com/yusuf/go-app/modules/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"regexp"
	"time"
)

type GoAppDB struct {
	App *config.GoAppTools
	DB  *mongo.Client
}

func NewGoAppDB(app *config.GoAppTools, db *mongo.Client) database.DBRepo {
	return &GoAppDB{
		App: app,
		DB:  db,
	}
}

func (g *GoAppDB) InsertUser(user *model.User) (bool, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	regMail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	ok := regMail.MatchString(user.Email)
	if !ok {
		g.App.ErrorLogger.Println("invalid registered details")
		return false, 0, errors.New("invalid registered details")
	}

	filter := bson.D{{Key: "email", Value: user.Email}}

	var res bson.M
	err := User(g.DB, "user").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			user.ID = primitive.NewObjectID()
			_, insertErr := User(g.DB, "user").InsertOne(ctx, user)
			if insertErr != nil {
				g.App.ErrorLogger.Fatalf("cannot add user to the database : %v ", insertErr)
			}
			return true, 1, nil
		}
		g.App.ErrorLogger.Fatal(err)
	}
	return true, 2, nil
}
