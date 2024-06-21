# EASYDDNS

## Compose Example

### 01
dns: aliyun
ip: device_interface

``` yaml
version: "3"
services:
  easyddns:
    image: ghcr.io/tasse00/easyddns:main
    container_name: easyddns
    environment:
      - DOMAIN=YOUR_DOMAIN
      - RR=YOUR_SUBDOMAIN
      - ENABLE_V4=false
      - ENABLE_V6=true
      - REFRESH_INTERVAL=12h
      - DNS_MANAGER=aliyun
      - ALIYUN_ACCESS_KEY_ID=YOUR_KEY_ID
      - ALIYUN_ACCESS_KEY_SECRET=YOUR_KEY_SECRET
      - ALIYUN_REGION_ID=YOUR_REGION
      - IP_DETECTOR=device_interface
```


### 02
dns: cloudflare
ip: device_interface

``` yaml
version: "3"
services:
  easyddns:
    image: ghcr.io/tasse00/easyddns:main
    container_name: easyddns
    environment:
      - DOMAIN=YOUR_DOMAIN
      - RR=YOUR_SUBDOMAIN
      - ENABLE_V4=false
      - ENABLE_V6=true
      - REFRESH_INTERVAL=12h
      - DNS_MANAGER=cloudflare
      - CLOUDFLARE_API_TOKEN=YOUR_API_TOKEN
      - IP_DETECTOR=device_interface
```