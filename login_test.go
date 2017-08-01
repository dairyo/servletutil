package servletutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	// Create test server.
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request.
			if r.Method != "POST" {
				t.Errorf("expected method %s; got %s.", "POST", r.Method)
				t.SkipNow()
			}
			if r.RequestURI != "/drools-wb/j_security_check" {
				t.Errorf("expected request URI %s; got %s.",
					"drools-wb/j_security_check", r.RequestURI)
				t.SkipNow()
			}
			if r.PostFormValue("j_username") != "dummy_user" {
				t.Errorf("expected username %s; got %s.",
					"dummy_user", r.PostFormValue("j_username"))
				t.SkipNow()
			}
			if r.PostFormValue("j_password") != "dummy_password" {
				t.Errorf("expected password %s got; %s.",
					"dummy_password", r.PostFormValue("j_password"))
				t.SkipNow()
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
		t.Errorf("err must be nil: %s", err)
		t.SkipNow()
	}
	if session == nil {
		t.Error("session must not be nil.")
		t.SkipNow()
	}
	if session.Key != "JSESSIONID" {
		t.Errorf("expected key %s; got %s", "JSESSIONID", session.Key)
		t.SkipNow()
	}
	if session.ID != "jJ4kllb1J0vwdZvSL4Bg4pIb0YDDMZFbOz3__ku2.drools-wildfly" {
		t.Errorf("expected id %s; got %s",
			"jJ4kllb1J0vwdZvSL4Bg4pIb0YDDMZFbOz3__ku2.drools-wildfly",
			session.ID)
		t.SkipNow()
	}
}
