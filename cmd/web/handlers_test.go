package main

import (
	"database/sql"
	"net/http"
	"net/url"
	"testing"

	"snippet.darieldejesus.com/internal/assert"
	"snippet.darieldejesus.com/internal/models/mocks"
)

func TestHome(t *testing.T) {
	t.Run("Unexpected error", func(t *testing.T) {
		app := newTestApplication(t)
		// Overwrite model to throw an unexpected error
		app.snippets = &mocks.SnippetModel{Err: sql.ErrConnDone}
		ts := newTestServer(t, app.routes())
		defer ts.Close()

		statusCode, _, _ := ts.get(t, "/")
		assert.Equal(t, statusCode, http.StatusInternalServerError)
	})

	t.Run("Success", func(t *testing.T) {
		app := newTestApplication(t)
		ts := newTestServer(t, app.routes())
		defer ts.Close()

		statusCode, _, body := ts.get(t, "/")

		assert.Equal(t, statusCode, http.StatusOK)
		assert.Contains(t, string(body), "Latest Snippets")
		assert.Contains(t, string(body), "Lorem ipsum")
		assert.Contains(t, string(body), "07 Dec 2022 at 11:12")
		assert.Contains(t, string(body), "#7")
	})
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name       string
		urlPath    string
		expectCode int
		expectBody string
	}{
		{
			name:       "Valid ID",
			urlPath:    "/snippet/view/7",
			expectCode: http.StatusOK,
			expectBody: "Lorem ipsum dolor sit amet...",
		},
		{
			name:       "Non-existent ID",
			urlPath:    "/snippet/view/9",
			expectCode: http.StatusNotFound,
		},
		{
			name:       "Negative ID",
			urlPath:    "/snippet/view/-1",
			expectCode: http.StatusNotFound,
		},
		{
			name:       "Decimal ID",
			urlPath:    "/snippet/view/1.23",
			expectCode: http.StatusNotFound,
		},
		{
			name:       "String ID",
			urlPath:    "/snippet/view/foo",
			expectCode: http.StatusNotFound,
		},
		{
			name:       "Empty ID",
			urlPath:    "/snippet/view/",
			expectCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, _, body := ts.get(t, test.urlPath)
			assert.Equal(t, code, test.expectCode)
			if test.expectBody != "" {
				assert.Contains(t, body, test.expectBody)
			}
		})
	}
}

func TestSnippetViewServerError(t *testing.T) {
	app := newTestApplication(t)
	// Overwrite model to throw an unexpected error
	app.snippets = &mocks.SnippetModel{Err: sql.ErrConnDone}
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, _ := ts.get(t, "/snippet/view/7")
	assert.Equal(t, statusCode, http.StatusInternalServerError)
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	authenticateUser(t, ts)

	code, _, body := ts.get(t, "/snippet/create")
	assert.Equal(t, code, http.StatusOK)
	validCSRFToken := extractCSRFToken(t, body)

	tests := []struct {
		name           string
		title          string
		content        string
		expires        string
		expectCode     int
		expectContain  string
		expectLocation string
	}{
		{
			name:          "Missing title",
			content:       "At vero eos et accusamus et iusto odio dignissimos",
			expires:       "365",
			expectCode:    http.StatusUnprocessableEntity,
			expectContain: "This field cannot be blank",
		},
		{
			name:          "Long title",
			title:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			content:       "At vero eos et accusamus et iusto odio dignissimos",
			expires:       "365",
			expectCode:    http.StatusUnprocessableEntity,
			expectContain: "This field cannot be more than 100 characters long",
		},
		{
			name:          "Missing content",
			title:         "Lorem ipsum",
			expires:       "365",
			expectCode:    http.StatusUnprocessableEntity,
			expectContain: "This field cannot be blank",
		},
		{
			name:          "Missing expiration",
			title:         "Lorem ipsum",
			content:       "At vero eos et accusamus et iusto odio dignissimos",
			expectCode:    http.StatusUnprocessableEntity,
			expectContain: "This field must equal 1, 7 or 365",
		},
		{
			name:          "Invalid expiration",
			title:         "Lorem ipsum",
			content:       "At vero eos et accusamus et iusto odio dignissimos",
			expires:       "1000",
			expectCode:    http.StatusUnprocessableEntity,
			expectContain: "This field must equal 1, 7 or 365",
		},
		{
			name:           "Valid submission",
			title:          "Lorem ipsum",
			content:        "At vero eos et accusamus et iusto odio dignissimos",
			expires:        "365",
			expectCode:     http.StatusSeeOther,
			expectLocation: "/snippet/view/8",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", test.title)
			form.Add("content", test.content)
			form.Add("expires", test.expires)
			form.Add("csrf_token", validCSRFToken)
			code, headers, body := ts.postForm(t, "/snippet/create", form)
			assert.Equal(t, code, test.expectCode)
			assert.Contains(t, body, test.expectContain)

			if test.expectLocation != "" {
				assert.Equal(t, headers.Get("location"), test.expectLocation)
			}
		})
	}
}

func TestSnippetCreateServerError(t *testing.T) {
	app := newTestApplication(t)
	// Overwrite model to throw an unexpected error
	app.snippets = &mocks.SnippetModel{Err: sql.ErrConnDone}
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	authenticateUser(t, ts)

	code, _, body := ts.get(t, "/snippet/create")
	assert.Equal(t, code, http.StatusOK)
	validCSRFToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("title", "Lorem ipsum")
	form.Add("content", "At vero eos et accusamus et iusto odio dignissimos")
	form.Add("expires", "365")
	form.Add("csrf_token", validCSRFToken)

	statusCode, _, _ := ts.postForm(t, "/snippet/create", form)
	assert.Equal(t, statusCode, http.StatusInternalServerError)
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		validName     = "John"
		validPassword = "Pa$$word123"
		validEmail    = "john@example.com"
		formTag       = `<form action="/user/signup" method="POST" novalidate>`
	)

	tests := []struct {
		name          string
		userName      string
		userEmail     string
		userPassword  string
		csrfToken     string
		expectCode    int
		expectFormTag string
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			expectCode:   http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			expectCode:   http.StatusBadRequest,
		},
		{
			name:          "Empty name",
			userName:      "",
			userEmail:     validEmail,
			userPassword:  validPassword,
			csrfToken:     validCSRFToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Empty email",
			userName:      validName,
			userEmail:     "",
			userPassword:  validPassword,
			csrfToken:     validCSRFToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Empty password",
			userName:      validName,
			userEmail:     validEmail,
			userPassword:  "",
			csrfToken:     validCSRFToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Invalid email",
			userName:      validName,
			userEmail:     "john@example.",
			userPassword:  validPassword,
			csrfToken:     validCSRFToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Short password",
			userName:      validName,
			userEmail:     validEmail,
			userPassword:  "pa$$",
			csrfToken:     validCSRFToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
		{
			name:          "Duplicate email",
			userName:      validName,
			userEmail:     "user@darieldejesus.com",
			userPassword:  validPassword,
			csrfToken:     validCSRFToken,
			expectCode:    http.StatusUnprocessableEntity,
			expectFormTag: formTag,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", test.userName)
			form.Add("email", test.userEmail)
			form.Add("password", test.userPassword)
			form.Add("csrf_token", test.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)
			assert.Equal(t, code, test.expectCode)

			if test.expectFormTag != "" {
				assert.Contains(t, body, test.expectFormTag)
			}
		})
	}
}

func TestUserSignupServerError(t *testing.T) {
	app := newTestApplication(t)
	// Overwrite model to throw an unexpected error
	app.users = &mocks.UserModel{Err: sql.ErrConnDone}
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/user/signup")
	assert.Equal(t, code, http.StatusOK)

	form := url.Values{}
	form.Add("name", "John")
	form.Add("email", "john@example.com")
	form.Add("password", "Pa$$word123")
	form.Add("csrf_token", extractCSRFToken(t, body))

	statusCode, _, _ := ts.postForm(t, "/user/signup", form)
	assert.Equal(t, statusCode, http.StatusInternalServerError)
}

func TestUserLogout(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	authenticateUser(t, ts)

	_, _, body := ts.get(t, "/")
	form := url.Values{}
	form.Add("csrf_token", extractCSRFToken(t, body))
	code, headers, _ := ts.postForm(t, "/user/logout", form)

	assert.Equal(t, code, http.StatusSeeOther)
	assert.Equal(t, headers.Get("location"), "/")
}

func TestUserLogin(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		validPassword = "Pa$$word123"
		validEmail    = "user@darieldejesus.com"
		invalidEmail  = "hacker@darieldejesus.com"
	)

	tests := []struct {
		name          string
		email         string
		password      string
		expectCode    int
		expectContain string
	}{
		{
			name:       "Valid submission",
			email:      validEmail,
			password:   validPassword,
			expectCode: http.StatusSeeOther,
		},
		{
			name:          "Missing email",
			password:      validPassword,
			expectCode:    http.StatusUnprocessableEntity,
			expectContain: "This field cannot be blank",
		},
		{
			name:          "Invalid email address",
			email:         "Platano@Power.",
			password:      validPassword,
			expectCode:    http.StatusUnprocessableEntity,
			expectContain: "This field must be a valid email address",
		},
		{
			name:          "Missin password",
			email:         validEmail,
			expectCode:    http.StatusUnprocessableEntity,
			expectContain: "This field cannot be blank",
		},
		{
			name:          "Unregistered Email",
			email:         invalidEmail,
			password:      validPassword,
			expectCode:    http.StatusUnprocessableEntity,
			expectContain: "Email or password is incorrect",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", test.email)
			form.Add("password", test.password)
			form.Add("csrf_token", validCSRFToken)

			code, _, body := ts.postForm(t, "/user/login", form)
			assert.Equal(t, code, test.expectCode)
			if test.expectContain != "" {
				assert.Contains(t, body, test.expectContain)
			}
		})
	}
}
