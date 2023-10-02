package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/prranavv/bookings/internal/config"
	"github.com/prranavv/bookings/internal/handlers"
	"github.com/prranavv/bookings/internal/helpers"
	"github.com/prranavv/bookings/internal/models"
	"github.com/prranavv/bookings/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var infolog *log.Logger
var errorlog *log.Logger

const portNumber = ":8080"

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
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

func run() error {
	//what am i goin to put in the session
	gob.Register(models.Reservation{})
	//change this to true when in production
	app.InProduction = false
	infolog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infolog
	errorlog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorlog
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false
	render.NewTemplates(&app)
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	return nil
}
