package main

import (
	"EduCommentSync/internal/service"
	"log"
)

func main() {
	srv := service.New()
	err := srv.Run()
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
