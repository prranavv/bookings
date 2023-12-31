package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/prranavv/bookings/internal/config"
	"github.com/prranavv/bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	//change this to true when in production
	testApp.InProduction = false
	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infolog
	errorlog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorlog
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session
	app = &testApp
	os.Exit(m.Run())
}

type mywriter struct{}

func (mw mywriter) Header() http.Header {
	var h http.Header
	return h
}

func (mw mywriter) WriteHeader(i int) {

}

func (mw mywriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
