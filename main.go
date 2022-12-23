package main

import (
	"github.com/google/uuid"
	"github.com/staffea-test/entity"
	"github.com/staffea-test/repository"
	"github.com/staffea-test/web"
	"log"
)

func main() {
	db, err := repository.NewConn()
	if err != nil {
		log.Fatalln(err)
	}

	// temporary migrations and seeds
	db.AutoMigrate(&entity.User{}, &entity.Invite{})
	db.Exec("truncate table users")
	db.Create(entity.Invite{ID: uuid.MustParse("524d832f-2624-4c4a-957a-4e48112d3df3")})
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	defer sqlDB.Close()

	userDB := repository.NewRepo(db)
	ctl := web.NewController(userDB)
	mw := web.NewMiddleware()

	web.ListenTo(ctl, mw)
}
