package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.bimasenaputra/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	templateData := &templateData {
		Snippets: snippets,
	}

	app.render(w, "home.html", http.StatusOK, templateData)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFoundError(w)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFoundError(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	templateData := &templateData {
		Snippet: snippet,
	}

	app.render(w, "view.html", http.StatusOK, templateData)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	templateData := &templateData {
		Form: &createSnippetForm { Expires: 365, },
	}
	app.render(w, "create.html", http.StatusOK, templateData)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &createSnippetForm {
		Title: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}
	
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		templateData := &templateData {
			Form: form,
		}
		app.render(w, "create.html", http.StatusUnprocessableEntity, templateData)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, expires)

	if err != nil {
		app.serverError(w, err)
		return 
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}