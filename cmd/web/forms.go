package main

import (
	"snippetbox.bimasenaputra/internal/validator"
)

type createSnippetForm struct {
	Title string
	Content string
	Expires int
	validator.Validator
}

type searchForm struct {
	Query string
	validator.Validator
}