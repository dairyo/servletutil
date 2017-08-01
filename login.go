package servletutil

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// Session contains session id and key. They are cookie values.
type Session struct {
	ID  string
	Key string
}

// Login to drools-wb endpoint with username and password.
func Login(endpoint, username, password string) (session *Session, err error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{Jar: jar}
	v := url.Values{}
	v.Set("j_username", username)
	v.Set("j_password", password)
	u2, err := u.Parse("drools-wb/j_security_check")
	if err != nil {
		return nil, err
	}
	resp, err := client.PostForm(u2.String(), v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	u3, err := u.Parse("drools-wb")
	if err != nil {
		return nil, err
	}
	cs := jar.Cookies(u3)
	if len(cs) == 0 {
		return nil, errors.New("no cookies for endpoint")
	}
	for _, c := range cs {
		if c.Name == "JSESSIONID" {
			return &Session{ID: c.Value, Key: c.Name}, nil
		}
	}
	return nil, errors.New("session not found")
}
