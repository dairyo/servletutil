package droolswbutil

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Session struct {
	Id  string
	Key string
}

func Login(endpoint, username, password string) (session *Session, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	setCookieUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	client := http.Client{Jar: jar}
	values := url.Values{}
	values.Set("j_username", username)
	values.Set("j_password", password)

	resp, err := client.PostForm(endpoint, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	for _, cookie := range jar.Cookies(setCookieUrl) {
		if cookie.Name == "JSESSIONID" {
			return &Session{Id: cookie.Value, Key: cookie.Name}, nil
		}
	}

	return nil, errors.New("no cookies")
}
