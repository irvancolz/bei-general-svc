package main

import (
	"fmt"
	"log"
	"os"

	"be-idx-tsg/internal/app/httprest/router"

	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}
	port := ":" + os.Getenv("ACTIVE_PORT")
	if err := router.Routes().Run(port); err != nil {
		log.Fatalln(err)
	}
}
