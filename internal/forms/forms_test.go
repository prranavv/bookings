package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}
	posteddata := url.Values{}
	posteddata.Add("a", "a")
	posteddata.Add("b", "b")
	posteddata.Add("c", "c")
	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = posteddata
	form = New(r.PostForm)
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	has := form.Has("whatever")
	if has {
		t.Error("form shows has filed when there is no field")
	}

	posteddata := url.Values{}
	posteddata.Add("a", "a")
	form = New(posteddata)
	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it shud ")
	}
}

func TestForm_Minlength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("shows min lenght for nonexistent field")
	}
	isErr := form.Errors.Get("x")
	if isErr == "" {
		t.Error("should have an errror but din get one")
	}

	posteddata := url.Values{}
	posteddata.Add("some_field", "some_value")
	form = New(posteddata)
	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows minlenght of 100 met wehn data is shorter")
	}
	posteddata = url.Values{}
	posteddata.Add("another", "abcd123")
	form = New(posteddata)
	form.MinLength("another", 1)
	if !form.Valid() {
		t.Error("shows minlength of 1 is not met when it is")
	}
	isErr = form.Errors.Get("another")
	if isErr != "" {
		t.Error("should have an errror but din get one")
	}
}

func TestForm_isEmail(t *testing.T) {
	postedvalues := url.Values{}
	form := New(postedvalues)
	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shes valid email for non-existent field")
	}
	postedvalues = url.Values{}
	postedvalues.Add("email", "me@afsd.com")
	form = New(postedvalues)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when it should have")
	}
	postedvalues = url.Values{}
	postedvalues.Add("email", "me")
	form = New(postedvalues)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("got an valid email when it shouldnt have")
	}
}
