package main

import (
	"os"

	"github.com/codegangsta/negroni"

	_ "github.com/joho/godotenv/autoload"
)

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
