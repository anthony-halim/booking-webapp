package main

import (
	"fmt"
	"github.com/anthony-halim/booking-webapp/internal/config"
	"github.com/anthony-halim/booking-webapp/internal/handlers"
	"github.com/anthony-halim/booking-webapp/internal/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber string = ":8080"
var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	// Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Get template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	// Populate app config
	app.TemplateCache = tc
	app.UseCache = false

	// Link with handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// Link with renders
	render.NewTemplates(&app)

		fmt.Println("Starting application on port", portNumber)

	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}