package main

import (
	"log"
	"os"
	"time"

	"tasse00.com/easyddns/dns"
	"tasse00.com/easyddns/ip"
)

func refreshIpV4(ipdetector ip.DetectIp, dnsmanager dns.ManageDNS, domain string) error {

	ipV4, err := ipdetector.DetectIpV4()
	if err != nil {
		return err
	}
	log.Println("current ipv4: ", ipV4)

	rcdVal, err := dnsmanager.GetIpV4(domain)
	if err != nil {
		return err
	}
	log.Println(domain, "record value", rcdVal)

	if rcdVal != ipV4 {
		log.Println(domain, "will update ipv4 ...")
		err = dnsmanager.UpdateIpV4(domain, ipV4)
		if err != nil {
			return err
		}
		log.Println(domain, "updated ipv4.")
	} else {
		log.Println(domain, "ipv4 is the same.")
	}
	return nil
}

func refreshIpV6(ipdetector ip.DetectIp, dnsmanager dns.ManageDNS, domain string) error {
	ipV6, err := ipdetector.DetectIpV6()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("current ipv6: ", ipV6)

	rcdVal, err := dnsmanager.GetIpV6(domain)
	if err != nil {
		return err
	}
	log.Println(domain, "record value", rcdVal)
	if rcdVal != ipV6 {
		log.Println(domain, "will update ipv6 ...")
		err = dnsmanager.UpdateIpV6(domain, ipV6)
		if err != nil {
			return err
		}
		log.Println(domain, "updated ipv6.")
	} else {
		log.Println(domain, "ipv6 is the same.")
	}
	return nil
}
func main() {

	ipDetectorType := os.Getenv("IP_DETECTOR")
	dnsManagerType := os.Getenv("DNS_MANAGER")
	domain := os.Getenv("DOMAIN")
	enableV4 := os.Getenv("ENABLE_V4") == "true"
	enableV6 := os.Getenv("ENABLE_V6") == "true"
	refreshInterval := os.Getenv("REFRESH_INTERVAL")

	println("Ip Detector:", ipDetectorType)
	println("Dns Manager:", dnsManagerType)
	println("Domain:", domain)
	println("Enable V4:", enableV4)
	println("Enable V6:", enableV6)
	println("Refresh Interval:", refreshInterval)

	ipDetector, err := ip.GetIpDetector(ipDetectorType)
	if err != nil {
		log.Fatalln(err)
	}

	dnsManager, err := dns.GetDnsManager(dnsManagerType)
	if err != nil {
		log.Fatalln(err)
	}

	dur, err := time.ParseDuration(refreshInterval)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		if enableV4 {
			err = refreshIpV4(ipDetector, dnsManager, domain)
			if err != nil {
				log.Println(err)
			}
		}

		if enableV6 {
			err = refreshIpV6(ipDetector, dnsManager, domain)
			if err != nil {
				log.Println(err)
			}
		}

		// string to int

		if err != nil {
			log.Fatalln(err)
		}

		log.Println("sleeping...", dur)
		time.Sleep(dur)
	}

}
