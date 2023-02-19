package database

import "github.com/yusuf/go-app/modules/model"

type DBRepo interface {
	InsertUser(user *model.User) (bool, int, error)
}
