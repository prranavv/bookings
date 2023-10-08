package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/prranavv/bookings/internal/config"
	"github.com/prranavv/bookings/internal/models"
	"github.com/prranavv/bookings/internal/render"
)

var pathtotemplate = "./../../templates"
var app config.AppConfig
var session *scs.SessionManager

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	//change this to true when in production
	app.InProduction = false
	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infolog
	errorlog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorlog
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = true
	render.NewRenderer(&app)
	repo := NewTestRepo(&app)
	NewHandlers(repo)
	os.Exit(m.Run())
}

func getroutes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.General)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availabity)
	mux.Post("/search-availability", Repo.PostAvailabity)
	mux.Post("/search-availability-json", Repo.AvailabityJSON)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	return mux
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	//mycache := make(map[string]*template.Template)
	mycache := map[string]*template.Template{}

	//get all of the files name
	pages, err := filepath.Glob(fmt.Sprintf("%s/*page.html", pathtotemplate))
	if err != nil {
		return mycache, err
	}
	//range through all files ending with *.page.html
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return mycache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathtotemplate))
		if err != nil {
			return mycache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathtotemplate))
			if err != nil {
				return mycache, err
			}
		}
		mycache[name] = ts
	}
	return mycache, nil
}
