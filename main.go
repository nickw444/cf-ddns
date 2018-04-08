package main

import (
	"net"
	"os"

	"github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"crypto/tls"
)

var log = logrus.New()

func main() {
	var (
		app = kingpin.New("cf-ddns", "Cloudflare DynDNS Updater")

		ipAddress = app.Flag("ip-address", "Skip resolving external IP and use provided IP").String()
		noVerify  = app.Flag("no-verify", "Don't verify ssl certificates").Bool()

		cfEmail  = app.Flag("cf-email", "Cloudflare Email").Required().String()
		cfApiKey = app.Flag("cf-api-key", "Cloudflare API key").Required().String()
		cfZoneId = app.Flag("cf-zone-id", "Cloudflare Zone ID").Required().String()

		hostnames = app.Arg("hostnames", "Hostnames to update").Required().Strings()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	var ip IPService
	var dns *CFDNSUpdater
	var err error

	if *ipAddress != "" {
		ip = &FakeIPService{
			fakeIp: net.ParseIP(*ipAddress),
		}
	} else {
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: *noVerify},
			},
		}
		ip = &IpifyIPService{HttpClient: httpClient}
	}

	if dns, err = NewCFDNSUpdater(*cfZoneId, *cfApiKey, *cfEmail, log.WithField("component", "cf-dns-updater")); err != nil {
		log.Panic(err)
	}

	res, err := ip.GetExternalIP()
	if err != nil {
		log.Panic(err)
	}

	for _, hostname := range *hostnames {
		err := dns.UpdateRecordA(hostname, res)
		if err != nil {
			log.Panic(err)
		}
	}
}
