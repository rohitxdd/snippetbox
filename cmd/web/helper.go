package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *application) Render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServerError(w, err)
		return
	}

	w.WriteHeader(status)

	buff := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buff, "base", data)

	if err != nil {
		app.ServerError(w, err)
		return
	}

	buff.WriteTo(w)
}

func (app *application) newTemplateData() *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
