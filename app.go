package main

import (
	"os"

	"github.com/codegangsta/negroni"
	"github.com/joho/godotenv"
	"log"
)

func init()  {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err.Error())
	}
}

func main() {
	r := buildRoutes()

	n := negroni.New()
	n.UseHandler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	n.Run(":" + port)
}
