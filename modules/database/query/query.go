package query

import (
	"github.com/yusuf/go-app/modules/config"
	"github.com/yusuf/go-app/modules/database"
	"go.mongodb.org/mongo-driver/mongo"
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

func (g *GoAppDB) InsertUser() {
	return
}
