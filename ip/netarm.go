package ip

import (
	"io/ioutil"
	"net/http"
)

type Netarm struct{}

func (n *Netarm) DetectIpV4() (string, error) {
	url := "https://ipv4.netarm.com/"
	v4, err := getUrlResponseAsString(url)
	if err != nil {
		return "", err
	}
	return v4, nil
}

func (n *Netarm) DetectIpV6() (string, error) {
	url := "https://ipv6.netarm.com/"
	v6, err := getUrlResponseAsString(url)
	if err != nil {
		return "", err
	}
	return v6, nil
}

func getUrlResponseAsString(url string) (retval string, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
