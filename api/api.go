package api

import (
	"log"

	"github.com/MarkSmersh/nil-chat/utils"
)

func Init() {
	dburl := utils.GetDBUrl()

	s := NewServer()

	if err := s.ConnectDB(dburl); err != nil {
		log.Fatal(err.Error())
	}

	s.InitRoutes()

	if err := s.Start("", 1488); err != nil {
		log.Fatal(err.Error())
	}
}
