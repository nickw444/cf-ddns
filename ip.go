package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type IPService interface {
	GetExternalIP() (net.IP, error)
}

type FakeIPService struct {
	fakeIp net.IP
}

func (f *FakeIPService) GetExternalIP() (net.IP, error) {
	if f.fakeIp == nil {
		return nil, fmt.Errorf("FakeIPService: No IP specified")
	}
	return f.fakeIp, nil
}

type IpifyIPService struct {
	httpClient *http.Client
}

func NewIpifyIPService() *IpifyIPService {
	return &IpifyIPService{
		httpClient: &http.Client{},
	}
}

type IpifyAPIResponse struct {
	IP string
}

func (i *IpifyIPService) GetExternalIP() (net.IP, error) {
	r, err := i.httpClient.Get("https://api.ipify.org?format=json")
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	var resp IpifyAPIResponse
	json.NewDecoder(r.Body).Decode(&resp)
	return net.ParseIP(resp.IP), nil
}
