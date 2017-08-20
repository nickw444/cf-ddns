package main

import (
	"net"
	"os"

	"github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var log = logrus.New()

func main() {
	var (
		app = kingpin.New("cf-ddns", "Cloudflare DynDNS Updater")

		dummyIp        = app.Flag("dummy-ip", "Use a dummy IP service").Bool()
		dummyIpAddress = app.Flag("dummy-ip-address", "Dummy address to set").String()

		cfEmail  = app.Flag("cf-email", "Cloudflare Email").Required().String()
		cfApiKey = app.Flag("cf-api-key", "Cloudflare API key").Required().String()
		cfZoneId = app.Flag("cf-zone-id", "Cloudflare Zone ID").Required().String()

		hostnames = app.Arg("hostnames", "Hostnames to update").Required().Strings()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	var ip IPService
	var dns *CFDNSUpdater
	var err error

	if *dummyIp {
		if *dummyIpAddress == "" {
			log.Panic("dummy-ip was specified without dummy-ip-address")
		}
		ip = &FakeIPService{
			fakeIp: net.ParseIP(*dummyIpAddress),
		}
	} else {
		ip = NewIpifyIPService()
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
