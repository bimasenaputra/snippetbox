package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"

	"snippetbox.bimasenaputra/internal/mocks"
)

func newTestApplication(t *testing.T) *application {
	// Change path directory so it can fetch all required html files
	_, filename, _, _ := runtime.Caller(0)

    dir := path.Join(path.Dir(filename), "../..")
    err := os.Chdir(dir)
    if err != nil {
        t.Fatal(err)
    }  
	
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog: log.New(io.Discard, "", 0),
		snippets: &mocks.SnippetModel{},
		templateCache: templateCache,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)

	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) post(t *testing.T, urlPath string, payload *bytes.Buffer) (int, http.Header, string) {
	rs, err := ts.Client().Post(ts.URL + urlPath, "application/x-www-form-urlencoded", payload)

	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)	
}