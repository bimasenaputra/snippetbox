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