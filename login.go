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
	// Create http client.
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{Jar: jar}
	// Get cookie for login.
	u2, err := u.Parse("drools-wb")
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(u2.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	// Login to get session ID.
	u3, err := u.Parse("drools-wb/j_security_check")
	if err != nil {
		return nil, err
	}
	resp, err = client.PostForm(u3.String(), url.Values{
		"j_username": []string{username}, "j_password": []string{password}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	u4, err := u.Parse("drools-wb")
	if err != nil {
		return nil, err
	}
	cs := jar.Cookies(u4)
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

func Logout(endpoint string, session *Session) (err error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	u2, err := u.Parse("drools-wb")
	if err != nil {
		return err
	}
	// Create http client.
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	jar.SetCookies(u2, []*http.Cookie{&http.Cookie{
		Name:  session.Key,
		Value: session.ID,
		Path:  "drools-wb",
	}})
	client := http.Client{Jar: jar}
	// Logout.
	u3, err := u.Parse("drools-wb/logout.jsp")
	if err != nil {
		return err
	}
	resp, err := client.Get(u3.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	return nil
}
