package main

import (
	"log"
	"os"
	"time"

	"tasse00.com/easyddns/dns"
	"tasse00.com/easyddns/ip"
)

func refreshIpV4(ipdetector ip.DetectIp, dnsmanager dns.ManageDNS, domain, rr string) error {

	ipV4, err := ipdetector.DetectIpV4()
	if err != nil {
		return err
	}
	log.Println("current ipv4: ", ipV4)

	rcdVal, err := dnsmanager.GetIpV4(domain, rr)
	if err != nil {
		return err
	}
	log.Println(domain, "record value", rcdVal)

	if rcdVal != ipV4 {
		log.Println(domain, "will update ipv4 ...")
		err = dnsmanager.UpdateIpV4(domain, rr, ipV4)
		if err != nil {
			return err
		}
		log.Println(domain, "updated ipv4.")
	} else {
		log.Println(domain, "ipv4 is the same.")
	}
	return nil
}

func refreshIpV6(ipdetector ip.DetectIp, dnsmanager dns.ManageDNS, domain, rr string) error {
	ipV6, err := ipdetector.DetectIpV6()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("current ipv6: ", ipV6)

	rcdVal, err := dnsmanager.GetIpV6(domain, rr)
	if err != nil {
		return err
	}
	log.Println(domain, "record value", rcdVal)
	if rcdVal != ipV6 {
		log.Println(domain, "will update ipv6 ...")
		err = dnsmanager.UpdateIpV6(domain, rr, ipV6)
		if err != nil {
			return err
		}
		log.Println(domain, "updated ipv6.")
	} else {
		log.Println(domain, "ipv6 is the same.")
	}
	return nil
}

func getEnvOrDefault(key string, dft string) string {
	if v := os.Getenv(key); v != "" {
		return v
	} else {
		return dft
	}
}

func main() {

	ipDetectorTypeV4 := getEnvOrDefault("IP_DETECTOR_V4", "netarm")
	ipDetectorTypeV6 := getEnvOrDefault("IP_DETECTOR_V6", "netarm")

	ipDetectorType := getEnvOrDefault("IP_DETECTOR", "")

	if ipDetectorType != "" {
		ipDetectorTypeV4 = ipDetectorType
		ipDetectorTypeV6 = ipDetectorType
	}

	dnsManagerType := getEnvOrDefault("DNS_MANAGER", "mock")

	domain := getEnvOrDefault("DOMAIN", "example.com")
	rr := getEnvOrDefault("RR", "sub")
	enableV4 := getEnvOrDefault("ENABLE_V4", "true") == "true"
	enableV6 := getEnvOrDefault("ENABLE_V6", "true") == "true"
	refreshInterval := getEnvOrDefault("REFRESH_INTERVAL", "6h")

	println("Ip Detector V4:", ipDetectorTypeV4)
	println("Ip Detector V6:", ipDetectorTypeV6)
	println("Dns Manager:", dnsManagerType)
	println("Domain:", domain)
	println("RR:", rr)
	println("Enable V4:", enableV4)
	println("Enable V6:", enableV6)
	println("Refresh Interval:", refreshInterval)

	ipDetectorV4, err := ip.GetIpDetector(ipDetectorTypeV4)
	if err != nil {
		log.Fatalln(err)
	}
	ipDetectorV6, err := ip.GetIpDetector(ipDetectorTypeV6)
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
			err = refreshIpV4(ipDetectorV4, dnsManager, domain, rr)
			if err != nil {
				log.Println(err)
			}
		}

		if enableV6 {
			err = refreshIpV6(ipDetectorV6, dnsManager, domain, rr)
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
