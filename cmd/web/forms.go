package main

type createSnippetForm struct {
	Title string
	Content string
	Expires int
	FieldErrors map[string]string
}