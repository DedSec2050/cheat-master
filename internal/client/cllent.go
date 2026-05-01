package client

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

type VTUClient struct {
	Client *http.Client
	Base   string
}

func NewClient() *VTUClient {
	jar, _ := cookiejar.New(nil)

	return &VTUClient{
		Client: &http.Client{
			Jar:     jar,
			Timeout: 15 * time.Second,
		},
		Base: "https://online.vtu.ac.in/api/v1",
	}
}