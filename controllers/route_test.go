package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func attemptLogin(t *testing.T, ctx *testContext, client *http.Client, username, password, optionalPath string) *http.Response {
	resp, err := http.Get(fmt.Sprintf("%s/login", ctx.adminServer.URL))
	if err != nil {
		t.Fatalf("error requesting the /login endpoint: %v", err)
	}
	got := resp.StatusCode
	expected := http.StatusOK
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	if client == nil {
		client = &http.Client{}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login%s", ctx.adminServer.URL, optionalPath), strings.NewReader(url.Values{
		"username": {username},
		"password": {password},
	}.Encode()))
	if err != nil {
		t.Fatalf("error creating new /login request: %v", err)
	}

	req.Header.Set("Cookie", resp.Header.Get("Set-Cookie"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the /login endpoint: %v", err)
	}
	return resp
}

func TestLoginCSRF(t *testing.T) {
	ctx := setupTest(t)
	defer tearDown(t, ctx)

	// A cross-site unsafe request (as flagged by the browser via the
	// Sec-Fetch-Site header) must be rejected by the cross-origin protection.
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", ctx.adminServer.URL),
		strings.NewReader(url.Values{
			"username": {"admin"},
			"password": {"gophish"},
		}.Encode()))
	if err != nil {
		t.Fatalf("error creating the /login request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Sec-Fetch-Site", "cross-site")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error requesting the /login endpoint: %v", err)
	}

	got := resp.StatusCode
	expected := http.StatusForbidden
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}
}

func TestInvalidCredentials(t *testing.T) {
	ctx := setupTest(t)
	defer tearDown(t, ctx)
	resp := attemptLogin(t, ctx, nil, "admin", "bogus", "")
	got := resp.StatusCode
	expected := http.StatusUnauthorized
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}
}

func TestSuccessfulLogin(t *testing.T) {
	ctx := setupTest(t)
	defer tearDown(t, ctx)
	resp := attemptLogin(t, ctx, nil, "admin", "gophish", "")
	got := resp.StatusCode
	expected := http.StatusOK
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}
}

func TestSuccessfulRedirect(t *testing.T) {
	ctx := setupTest(t)
	defer tearDown(t, ctx)
	next := "/campaigns"
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	resp := attemptLogin(t, ctx, client, "admin", "gophish", fmt.Sprintf("?next=%s", next))
	got := resp.StatusCode
	expected := http.StatusFound
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}
	url, err := resp.Location()
	if err != nil {
		t.Fatalf("error parsing response Location header: %v", err)
	}
	if url.Path != next {
		t.Fatalf("unexpected Location header received. expected %s got %s", next, url.Path)
	}
}

func TestAccountLocked(t *testing.T) {
	ctx := setupTest(t)
	defer tearDown(t, ctx)
	resp := attemptLogin(t, ctx, nil, "houdini", "gophish", "")
	got := resp.StatusCode
	expected := http.StatusUnauthorized
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}
}
