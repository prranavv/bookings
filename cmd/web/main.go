package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/prranavv/bookings/internal/config"
	"github.com/prranavv/bookings/internal/driver"
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
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)
	listenForMail()

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	log.Println("Server is running on port 8080")
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	//what am i goin to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailchan := make(chan models.MailData)
	app.MailChan = mailchan
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
	//connect to database
	log.Println("connecting to database.....")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go user=tsawlergo password=1234")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false
	render.NewRenderer(&app)
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	return db, err
}
