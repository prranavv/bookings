package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/prranavv/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	// {"post-search-avail", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2021-01-01"},
	// 	{key: "end", value: "2021-01-02"},
	// }, http.StatusOK},
	// {"post-search-avail-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2021-01-01"},
	// 	{key: "end", value: "2021-01-02"},
	// }, http.StatusOK},
	// {"make reservation post", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "me@here.com"},
	// 	{key: "phone", value: "8923849284"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getroutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "generals quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getctx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong code: got %d,wanted %d", rr.Code, http.StatusOK)
	}
	//test case where reservation is not in session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getctx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong code: got %d,wanted %d", rr.Code, http.StatusOK)
	}
	//test with non existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getctx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	reservation.RoomID = 100
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong code: got %d,wanted %d", rr.Code, http.StatusOK)
	}
}

func getctx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func TestRepository_PostReservation(t *testing.T) {
	// reqBody := "start_date=2030-11-11"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-01")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=joen")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=smth")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@email.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	//another way to make req body
	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "john")
	postedData.Add("last_name", "doe")
	postedData.Add("email", "john@doe.com")
	postedData.Add("phone", "1234567890")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content_Type", "application/x-www-form-url-encoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailabity)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong code: got %d,wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	//first case rooms not available
	reqbody := "start=2050-02-3"
	reqbody = fmt.Sprintf("%s&%s", reqbody, "end=2050-09-01")
	reqbody = fmt.Sprintf("%s&%s", reqbody, "room_id=1")
	//create request
	req, _ := http.NewRequest("POST", "/search-availabilty", strings.NewReader(reqbody))
	//get context with sessionr
	ctx := getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	//make handler handlerfunc
	handler := http.HandlerFunc(Repo.AvailabityJSON)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
}
