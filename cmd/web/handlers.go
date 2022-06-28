package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"snippetbox.bimasenaputra/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	minId, err := app.snippets.GetMinID()
	if err != nil {
		app.serverError(w, err)
		return
	}

	templateData := &templateData {
		Snippets: snippets,
		HasPrev: false,
		HasNext: snippets[len(snippets)-1].ID != minId,
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
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

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

func (app *application) snippetLatest(w http.ResponseWriter, r *http.Request) {
	opt := r.URL.Query().Get("direction")
	if strings.TrimSpace(opt) == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if opt == "next" {
		snippets, err := app.snippets.NextLatestPaging(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		minId, err := app.snippets.GetMinID()
		if err != nil {
			app.serverError(w, err)
			return
		}

		maxId, err := app.snippets.GetMaxID()
		if err != nil {
			app.serverError(w, err)
			return
		}

		templateData := &templateData{
			Snippets: snippets,
			HasNext: snippets[len(snippets)-1].ID != minId,
			HasPrev: snippets[0].ID != maxId,
		}

		app.render(w, "snippets.html", http.StatusOK, templateData)
	} else if opt == "prev" {
		snippets, err := app.snippets.PrevLatestPaging(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		minId, err := app.snippets.GetMinID()
		if err != nil {
			app.serverError(w, err)
			return
		}

		maxId, err := app.snippets.GetMaxID()
		if err != nil {
			app.serverError(w, err)
			return
		}

		templateData := &templateData{
			Snippets: snippets,
			HasNext: snippets[len(snippets)-1].ID != minId,
			HasPrev: snippets[0].ID != maxId,
		}

		app.render(w, "snippets.html", http.StatusOK, templateData)
	} else {
		app.clientError(w, http.StatusBadRequest)
	}
}

func (app *application) snippetSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if strings.TrimSpace(query) == "" {
		templateData := &templateData{
			Form: &searchForm {
				Query: "",
			},
		}
		app.render(w, "search.html", http.StatusOK, templateData)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if strings.TrimSpace(r.URL.Query().Get("id")) != "" && err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	direction := r.URL.Query().Get("direction")
	
	if id == 0 || strings.TrimSpace(direction) == "" {
		snippets, err := app.snippets.LatestContainsTitle(query)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if len(snippets) == 0 {
			templateData := &templateData{
				Snippets: snippets,
				Form: &searchForm {
					Query: query,
				},
			}
			app.render(w, "search.html", http.StatusOK, templateData)
			return
		}

		minId, err := app.snippets.GetMinIDByTitle(query)
		if err != nil {
			app.serverError(w, err)
			return
		}

		templateData := &templateData{
			Snippets: snippets,
			Form: &searchForm {
				Query: query,
			},
			HasPrev: false,
			HasNext: snippets[len(snippets)-1].ID != minId,
		}

		app.render(w, "search.html", http.StatusOK, templateData)
		return
	}

	if direction == "next" {
		minId, err := app.snippets.GetMinIDByTitle(query)
		if err != nil {
			app.serverError(w, err)
			return
		}

		maxId, err := app.snippets.GetMaxIDByTitle(query)
		if err != nil {
			app.serverError(w, err)
			return
		}

		snippets, err := app.snippets.NextLatestContainsTitle(id, query)
		if err != nil {
			app.serverError(w, err)
			return
		}

		templateData := &templateData{
			Snippets: snippets,
			Form: &searchForm {
				Query: query,
			},
			HasPrev: snippets[0].ID != maxId,
			HasNext: snippets[len(snippets)-1].ID != minId,
		}

		app.render(w, "snippets.html", http.StatusOK, templateData)
	} else if direction == "prev" {
		minId, err := app.snippets.GetMinIDByTitle(query)
		if err != nil {
			app.serverError(w, err)
			return
		}

		maxId, err := app.snippets.GetMaxIDByTitle(query)
		if err != nil {
			app.serverError(w, err)
			return
		}

		snippets, err := app.snippets.PrevLatestContainsTitle(id, query)
		if err != nil {
			app.serverError(w, err)
			return
		}

		templateData := &templateData{
			Snippets: snippets,
			Form: &searchForm {
				Query: query,
			},
			HasPrev: snippets[0].ID != maxId,
			HasNext: snippets[len(snippets)-1].ID != minId,
		}

		app.render(w, "snippets.html", http.StatusOK, templateData)
	} else {
		app.clientError(w, http.StatusBadRequest)
	}
}

func (app *application) snippetSearchPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &searchForm{
		Query: r.PostForm.Get("query"),
	}

	form.CheckField(validator.NotBlank(form.Query), "query", "This field cannot be blank")

	if !form.Valid() {
		templateData := &templateData {
			Form: form,
		}
		app.render(w, "search.html", http.StatusUnprocessableEntity, templateData)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets/search?q=%s", form.Query), http.StatusSeeOther)
}