package main

import (
	"net"
	"os"

	"crypto/tls"
	"github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
)

var log = logrus.New()
var Version string

func main() {
	var (
		app = kingpin.New("cf-ddns", "Cloudflare DynDNS Updater").Version(Version)

		ipAddress = app.Flag("ip-address", "Skip resolving external IP and use provided IP").String()
		noVerify  = app.Flag("no-verify", "Don't verify ssl certificates").Bool()

		cfEmail   = app.Flag("cf-email", "Cloudflare Email").Required().String()
		cfApiKey  = app.Flag("cf-api-key", "Cloudflare API key").Required().String()
		cfZoneId  = app.Flag("cf-zone-id", "Cloudflare Zone ID").Required().String()
		supportV6 = app.Flag("support-v6", "Whether IPv6 records should be updated/created").Bool()

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
		ip = &IpifyIPService{HttpClient: httpClient, supportV6: *supportV6}
	}

	if dns, err = NewCFDNSUpdater(*cfZoneId, *cfApiKey, *cfEmail, log.WithField("component", "cf-dns-updater")); err != nil {
		log.Panic(err)
	}

	res, err := ip.GetExternalIP()
	if err != nil {
		log.Panic(err)
	}

	for _, hostname := range *hostnames {
		if res.v4 != nil {
			err := dns.UpdateRecordA(hostname, res.v4)
			if err != nil {
				log.Panic(err)
			}
		}
		if res.v6 != nil {
			err := dns.UpdateRecordAAAA(hostname, res.v6)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}
