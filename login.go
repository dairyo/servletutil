package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func Login(endpoint, username, password string) (sid string, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}
	client := http.Client{Jar: jar}

	values := url.Values{}
	values.Set("j_username", username)
	values.Set("j_password", password)

	req, err := http.NewRequest(
		"POST",
		endpoint,
		strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	set_cookie_url, err := url.Parse(endpoint)
	cookies := jar.Cookies(set_cookie_url)
	sid = cookies[0].Value

	return sid, nil
}

func main() {
	fmt.Println(Login(
		"http://192.168.50.51:8080/drools-wb/j_security_check",
		"admin",
		"admin123"))
}
