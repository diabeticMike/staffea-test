package main

import (
	"fmt"
	"github.com/staffea-test/repository"
	"log"
)

func main() {
	db, err := repository.NewConn()
	if err != nil {
		log.Fatalln(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	defer sqlDB.Close()

	fmt.Println(db)
}
