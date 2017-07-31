package droolswbutils

import (
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestLoginSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"http://dummy.com/drools-wb/j_security_check",
		func(req *http.Request) (*http.Response, error) {
			if req.FormValue("j_username") != "dummy_user" {
				t.Errorf("username must be dummy_user but %s",
					req.FormValue("j_username"))
			}
			if req.FormValue("j_password") != "dummy_pass" {
				t.Errorf("password must be dummy_pass but %s",
					req.FormValue("j_password"))
			}
			res := httpmock.NewStringResponse(200, "")
			res.Header.Add(
				"Set-Cookie",
				"JSESSIONID=jJ4kllb1J0vwdZvSL4Bg4pIb0YDDMZFbOz3__ku2.drools-wildfly; path=/drools-wb")
			return res, nil
		},
	)

	session, err := Login(
		"http://dummy.com/drools-wb/j_security_check",
		"dummy_user",
		"dummy_pass")
	if err != nil {
		t.Errorf("err must be nil: %s", err)
	}
	if session == nil {
		t.Error("session must not be nil.")
	}
	if session.Key != "JSESSIONID" {
		t.Errorf("session key must be JSESSIONID: %s", session.Key)
	}
	if session.Id != "jJ4kllb1J0vwdZvSL4Bg4pIb0YDDMZFbOz3__ku2.drools-wildfly" {
		t.Errorf("session id is wrong: %s", session.Id)
	}
}
