package servletutil

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// Session contains session id and key. They are cookie values.
type Session struct {
	ID  string
	Key string
}

// Drools-wb user infomation.
type User struct {
	Username string
	Password string
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

// Logout drools-wb.
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

// Create drools-wb user.
func CreateUser(endpoint string, newUser, restUser User) (err error) {
	// Create client.
	client := http.Client{}
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	// Create user.
	u2, err := u.Parse("drools-wb/rest/user")
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, u2.String(),
		strings.NewReader(newUser.Username))
	if err != nil {
		return err
	}
	req.SetBasicAuth(restUser.Username, restUser.Password)
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	// Add a role.
	u3, err := u.Parse("drools-wb/rest/user/roles/" + newUser.Username)
	if err != nil {
		return err
	}
	req, err = http.NewRequest(http.MethodPut, u3.String(),
		strings.NewReader("[\"analyst\"]"))
	if err != nil {
		return err
	}
	req.SetBasicAuth(restUser.Username, restUser.Password)
	req.Header.Set("content-type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	// Set password
	u4, err := u.Parse("drools-wb/rest/user/password/" +
		newUser.Username + "?password=" + newUser.Password)
	if err != nil {
		return err
	}
	req, err = http.NewRequest(http.MethodPut, u4.String(),
		strings.NewReader(""))
	if err != nil {
		return err
	}
	req.SetBasicAuth(restUser.Username, restUser.Password)
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}
