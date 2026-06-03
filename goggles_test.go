package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRun_SendsGooglebotUA_AndAppliesSelector(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		fmt.Fprint(w, `<html><body><h1>Hello</h1><p>ignored</p></body></html>`)
	}))
	defer srv.Close()

	var buf bytes.Buffer
	if err := run([]string{"-selector", "h1", srv.URL}, &buf); err != nil {
		t.Fatalf("run: %v", err)
	}
	if !strings.Contains(buf.String(), "Hello") {
		t.Errorf("output missing 'Hello': %q", buf.String())
	}
	if strings.Contains(buf.String(), "ignored") {
		t.Errorf("selector leaked unrelated content: %q", buf.String())
	}
	if !strings.Contains(gotUA, "Googlebot") {
		t.Errorf("UA missing 'Googlebot': %q", gotUA)
	}
}

func TestRun_RejectsNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	err := run([]string{srv.URL}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "404") {
		t.Errorf("want 404 error, got %v", err)
	}
}

func TestRun_CustomUserAgent(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		fmt.Fprint(w, `<body>ok</body>`)
	}))
	defer srv.Close()

	if err := run([]string{"-ua", "custom/1.0", srv.URL}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if gotUA != "custom/1.0" {
		t.Errorf("UA = %q, want custom/1.0", gotUA)
	}
}
