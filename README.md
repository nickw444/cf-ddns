# Cloudflare DDNS Updater

Use Cloudflare a DDNS provider with this tool on crontab.

```
$> ./cf-ddns --help
usage: cf-ddns --cf-email=CF-EMAIL --cf-api-key=CF-API-KEY --cf-zone-id=CF-ZONE-ID [<flags>] <hostnames>...

Cloudflare DynDNS Updater

Flags:
  --help                   Show context-sensitive help (also try --help-long and
                           --help-man).
  --dummy-ip               Use a dummy IP service
  --dummy-ip-address=DUMMY-IP-ADDRESS
                           Dummy address to set
  --cf-email=CF-EMAIL      Cloudflare Email
  --cf-api-key=CF-API-KEY  Cloudflare API key
  --cf-zone-id=CF-ZONE-ID  Cloudflare Zone ID

Args:
  <hostnames>  Hostnames to update
```

## Compiling for MIPS (Ubnt Edgerouter)
```
GOOS=linux GOARCH=mips64 go build
```
