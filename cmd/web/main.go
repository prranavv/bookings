package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/prranavv/bookings/pkg/config"
	"github.com/prranavv/bookings/pkg/handlers"
	"github.com/prranavv/bookings/pkg/render"
)

var app config.AppConfig
var session *scs.SessionManager

const portNumber = ":8080"

func main() {

	//change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	render.NewTemplates(&app)
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	//log.Fatal(http.ListenAndServe(portNumber, nil))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	log.Println("Server is running on port 8080")
	err = srv.ListenAndServe()
	log.Fatal(err)
}
