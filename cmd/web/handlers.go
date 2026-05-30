package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.ServerError(w, err)
		return
	}
	data := app.newTemplateData()

	data.Snippets = snippets

	app.Render(w, 200, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)

	data := app.newTemplateData()

	data.Snippet = snippet

	app.Render(w, 200, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "test title"
	content := "test insert from the code..."
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
