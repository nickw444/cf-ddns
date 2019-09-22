package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type GetExternalIpResponse struct {
	v4 net.IP
	v6 net.IP
}

type IPService interface {
	GetExternalIP() (*GetExternalIpResponse, error)
}

type FakeIPService struct {
	fakeIp net.IP
}

func (f *FakeIPService) GetExternalIP() (*GetExternalIpResponse, error) {
	if f.fakeIp == nil {
		return nil, fmt.Errorf("FakeIPService: No IP specified")
	}
	return &GetExternalIpResponse{v4: f.fakeIp}, nil
}

type IpifyIPService struct {
	HttpClient *http.Client
	supportV6  bool
}

type IpifyAPIResponse struct {
	IP string
}

func (i *IpifyIPService) GetExternalIP() (*GetExternalIpResponse, error) {
	v4, err := i.getIp("https://api.ipify.org?format=json")
	if err != nil {
		return nil, err
	}

	var v6 net.IP
	if i.supportV6 {
		v6, err = i.getIp("https://api6.ipify.org?format=json")
		if err != nil {
			return nil, err
		}
	}

	return &GetExternalIpResponse{v4: v4, v6: v6}, nil
}

func (i *IpifyIPService) getIp(endpoint string) (net.IP, error) {
	r, err := i.HttpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	var resp IpifyAPIResponse
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(resp.IP), nil
}
