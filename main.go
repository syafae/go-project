package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/syafae/femProject/internal/app"
	"github.com/syafae/femProject/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "server backend port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	app.Logger.Printf("We are running on port %d", port)
	r := routes.SetUpRoutes(app)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      r,
	}
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}

}
