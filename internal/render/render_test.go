package render

import (
	"net/http"
	"testing"

	"github.com/prranavv/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getsession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "124")
	result := AddDefaultData(&td, r)
	if result.Flash != "124" {
		t.Error("flash value of 124 did not work")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathtotemplate = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc
	r, err := getsession()
	if err != nil {
		t.Error(err)
	}
	var ww mywriter
	err = RenderTemplate(ww, r, "home.page.html", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}
	err = RenderTemplate(ww, r, "non-existent.page.html", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that did not exist")
	}
}

func getsession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathtotemplate = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
