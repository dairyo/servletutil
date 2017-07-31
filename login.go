package droolswbutil

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Session struct {
	ID  string
	Key string
}

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
	resp, err := client.PostForm(endpoint, v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	cs := jar.Cookies(u)
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
