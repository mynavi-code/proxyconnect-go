package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage %s host port\n", os.Args[0])
		os.Exit(1)
	}

	proxy_url := os.Getenv("proxyconnect_url")
	if proxy_url == "" {
		proxy_url = os.Getenv("http_proxy")
	}
	if proxy_url == "" {
		fmt.Fprintf(os.Stderr, "Error: Environment variable \"proxyconnect_url\" and \"http_proxy\" are not set\n")
		os.Exit(1)
	}

	u, err := url.Parse(proxy_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Proxy setting can not parse\n")
		os.Exit(1)
	}

	proxyHost := u.Hostname() + ":" + u.Port()

	if u.Hostname() == "" || u.Port() == "" {
		fmt.Fprintf(os.Stderr, "Error: Proxy setting is '%s'\n", proxyHost)
		os.Exit(1)
	}

	destHost := os.Args[1] + ":" + os.Args[2]

	fmt.Fprintf(os.Stderr, "ProxyConnect %s -> %s\n", proxyHost, destHost)

	ur := url.URL{
		Scheme: u.Scheme,
		Host:   proxyHost,
		Opaque: destHost,
	}

	pipeRead, pipeWrite := io.Pipe()

	request := &http.Request{
		Method: "CONNECT",
		URL:    &ur,
		Header: map[string][]string{},
		Body:   pipeRead,
	}

	user, _ := url.PathUnescape(u.User.String())
	if user != "" {
		vs := strings.Split(user, ":")
		username := vs[0]
		password := ""
		if len(vs) >= 2 {
			password = strings.Repeat("*", len(vs[1]))
		}
		fmt.Fprintf(os.Stderr, "Proxy-Authorization: %s %s\n", username, password)
		request.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user)))
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Error: %s\n", err)
		os.Exit(1)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "HTTP Error: %s (%d)\n", response.Status, response.StatusCode)
		os.Exit(1)
	}

	go func(c io.ReadCloser) {
		defer c.Close()
		for true {
			var b [4096]byte
			n, _ := c.Read(b[:])
			os.Stdout.Write(b[:n])
		}
	}(response.Body)

	for true {
		var b [4096]byte
		n, _ := os.Stdin.Read(b[:])
		if n == 0 {
			break
		}
		pipeWrite.Write(b[:n])
	}
}
