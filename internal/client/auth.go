package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (v *VTUClient) Login(email, password string) error {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", v.Base+"/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	// Explicitly store new cookies to ensure fresh token is used
	baseURL, _ := url.Parse(v.Base)
	v.Client.Jar.SetCookies(baseURL, resp.Cookies())

	return nil
}