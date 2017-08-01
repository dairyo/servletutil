package servletutil

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	// Create test server.
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request.
			if r.Method != "POST" {
				t.Errorf("expected method %s; got %s.", "POST", r.Method)
				return
			}
			if r.RequestURI != "/drools-wb/j_security_check" {
				t.Errorf("expected request URI %s; got %s.",
					"drools-wb/j_security_check", r.RequestURI)
				return
			}
			err := r.ParseForm()
			if err != nil {
				t.Errorf("err must be nil: %s", err)
				return
			}
			want := url.Values{
				"j_username": []string{"dummy_user"},
				"j_password": []string{"dummy_password"},
			}
			if !reflect.DeepEqual(r.PostForm, want) {
				t.Errorf("expected post form %#v; got %#v", r.PostForm, want)
				return
			}
			// Create response.
			w.Header().Set(
				"Set-Cookie",
				"JSESSIONID=jJ4kllb1J0vwdZvSL4Bg4pIb0YDDMZFbOz3__ku2.drools-wildfly; path=/drools-wb")
		}))
	defer ts.Close()
	// Execute login.
	session, err := Login(
		ts.URL+"/drools-wb/j_security_check",
		"dummy_user",
		"dummy_password")
	if err != nil {
		t.Fatalf("err must be nil: %s", err)
	}
	want := Session{
		ID:  "jJ4kllb1J0vwdZvSL4Bg4pIb0YDDMZFbOz3__ku2.drools-wildfly",
		Key: "JSESSIONID",
	}
	if !reflect.DeepEqual(*session, want) {
		t.Fatalf("expected session %#v; got %#v", *session, want)
	}
}
