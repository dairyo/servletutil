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
			if r.RequestURI == "/drools-wb" {
				if r.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					fmt.Fprintf(w, "expected method %s; got %s", http.MethodGet, r.Method)
					return
				}
				http.SetCookie(w, &http.Cookie{
					Name:  "JSESSIONID",
					Value: "5q2CmRQRZLB81T9Gsbkt44iplQCNHvV5lmZkI0u9.drools-wildfly",
					Path:  "/drools-wb",
				})
				w.WriteHeader(http.StatusOK)
				return
			} else if r.RequestURI == "/drools-wb/j_security_check" {
				if r.Method != "POST" {
					w.WriteHeader(http.StatusMethodNotAllowed)
					fmt.Fprintf(w, "expected method %s; got %s", "POST", r.Method)
					return
				}
				c, err := r.Cookie("JSESSIONID")
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "no cookie for JSESSIONID: %s", err)
					return
				}
				if c.Value != "5q2CmRQRZLB81T9Gsbkt44iplQCNHvV5lmZkI0u9.drools-wildfly" {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "expected cookie %s; got %s", "5q2CmRQRZLB81T9Gsbkt44iplQCNHvV5lmZkI0u9.drools-wildfly", c.Value)
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
				http.SetCookie(w, &http.Cookie{
					Name:  "JSESSIONID",
					Value: "IKmkpzGWSKMO_Z3cjZIz8jH615ZMC95msfr-muRG.drools-wildfly",
					Path:  "/drools-wb",
				})
				w.WriteHeader(http.StatusOK)
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "invalid path")
				return
			}
		}))
	defer ts.Close()
	// Execute login.
	session, err := Login(ts.URL, "dummy_user", "dummy_password")
	if err != nil {
		t.Fatalf("err must be nil: %s", err)
	}
	want := Session{
		ID:  "IKmkpzGWSKMO_Z3cjZIz8jH615ZMC95msfr-muRG.drools-wildfly",
		Key: "JSESSIONID",
	}
	if !reflect.DeepEqual(*session, want) {
		t.Fatalf("expected session %#v; got %#v", want, *session)
	}
}
