package main

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"snippetbox.bimasenaputra/internal/assert"
)

func TestHome(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, _ := ts.get(t, "/")
	assert.Equal(t, code, http.StatusOK)
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name string
		path string
		wantCode int
		wantBody string
	} {
		{
			name: "Valid ID",
			path: "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name: "Non-existent ID",
			path: "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Negative ID",
			path: "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Decimal ID",
			path: "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name: "String ID",
			path: "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Empty ID",
			path: "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, _, body := ts.get(t, test.path)
			assert.Equal(t, code, test.wantCode)

			if test.wantBody != "" {
				assert.StringContains(t, string(body), test.wantBody)
			}
		})
	}
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, _ := ts.get(t, "/snippet/create")
	assert.Equal(t, code, http.StatusOK)
}

func TestSnippetCreatePost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	param := url.Values{}

    param.Set("title", "title")
	param.Set("content", "content")
	param.Set("expires", "7")
    payload1 := bytes.NewBufferString(param.Encode())

	param.Set("title", "")
	payload2 := bytes.NewBufferString(param.Encode())
	param.Set("title", "title")

	param.Set("title", strings.Repeat("a", 101))
	payload3 := bytes.NewBufferString(param.Encode())
	param.Set("title", "title")

	param.Set("content", "")
	payload4 := bytes.NewBufferString(param.Encode())
	param.Set("content", "content")

	param.Set("expires", "0")
	payload5 := bytes.NewBufferString(param.Encode())
	param.Set("expires", "7")

	param.Set("expires", "expires")
	payload6 := bytes.NewBufferString(param.Encode())

	tests := []struct {
		name string
		payload *bytes.Buffer
		expected int
	} {
		{
			name: "Valid Request",
			payload: payload1,
			expected: http.StatusOK,
		},
		{
			name: "Empty Title",
			payload: payload2,
			expected: http.StatusUnprocessableEntity,
		},
		{
			name: "Title Has More Than 100 Characters",
			payload: payload3,
			expected: http.StatusUnprocessableEntity,
		},
		{
			name: "Empty Content",
			payload: payload4,
			expected: http.StatusUnprocessableEntity,
		},
		{
			name: "Invalid Expires",
			payload: payload5,
			expected: http.StatusUnprocessableEntity,
		},
		{
			name: "String Expires",
			payload: payload6,
			expected: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, _, _ := ts.post(t, "/snippet/create", test.payload)
			assert.Equal(t, code, test.expected)
		})
	}
}