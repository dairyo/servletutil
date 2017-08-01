package servletutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	// Create test server.
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request.
			if r.Method != "POST" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintf(w, "expected method %s; got %s", "POST", r.Method)
				return
			}
			if r.RequestURI != "/drools-wb/j_security_check" {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "expected path %s; got %s",
					"/drools-wb/j_security_check", r.RequestURI)
				return
			}
			u := r.PostFormValue("j_username")
			if u != "dummy_user" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "expected user %s; got %s", "dummy_user", u)
				return
			}
			p := r.PostFormValue("j_password")
			if p != "dummy_password" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "expected password %s; got %s",
					"dummy_password", p)
				return
			}
			// Create response.
			w.Header().Set(
				"Set-Cookie",
				"JSESSIONID=jJ4kllb1J0vwdZvSL4Bg4pIb0YDDMZFbOz3__ku2.drools-wildfly; path=/drools-wb")
			w.WriteHeader(http.StatusOK)
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
		t.Fatalf("expected session %#v; got %#v", want, *session)
	}
}
