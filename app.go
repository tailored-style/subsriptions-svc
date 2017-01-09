package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/codegangsta/negroni"

	_ "github.com/joho/godotenv/autoload"
	"log"
)

func main() {
	r := buildRoutes()

	n := negroni.New()
	n.Use(negroni.HandlerFunc(catchPanicsMiddleware))
	n.UseHandler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	n.Run(":" + port)
}

func catchPanicsMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	type errorOutput struct {
		Message string `json:"message"`
		StackTrace string `json:"stackTrace"`
	}

	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()

			log.Printf("Unexpected panic: %v\n%s", r, stack)

			out := &errorOutput{
				Message: fmt.Sprintf("%v", r),
				StackTrace: string(stack),
			}

			js, err := json.Marshal(out)
			if err != nil {
				panic(fmt.Sprintf("Failed to marshal error content for panic: %v", r))
			}

			http.Error(w, string(js), http.StatusInternalServerError)
		}
	}()

	next(w, r)
}
