package main

import (
	"encoding/gob"
	"flag"
	"fmt"
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
	gob.Register(map[string]int{})

	//read flags
	inProduction := flag.Bool("production", true, "Application is in Productions")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbname := flag.String("dbname", "", "Database name")
	dbhost := flag.String("dbhost", "localhost", "Database Host")

	dbuser := flag.String("dbuser", "", "Database User")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable,prefer,require)")
	flag.Parse()

	if *dbname == "" || *dbuser == "" || *dbPass == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}
	mailchan := make(chan models.MailData)
	app.MailChan = mailchan
	//change this to true when in production
	app.InProduction = *inProduction
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
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbhost, *dbPort, *dbname, *dbuser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = *useCache
	render.NewRenderer(&app)
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	return db, err
}
