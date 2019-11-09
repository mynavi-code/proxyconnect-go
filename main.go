package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage %s host port\n", os.Args[0])
		os.Exit(1)
	}

	u, err := url.Parse(os.Getenv("http_proxy"))
	if err != nil || u.Hostname() == "" || u.Port() == "" {
		fmt.Fprintf(os.Stderr, "Error: Environment variable \"http_proxy\" is not set\n")
		os.Exit(1)
	}

	proxyHost := u.Hostname() + ":" + u.Port()
	destHost := os.Args[1] + ":" + os.Args[2]

	fmt.Fprintf(os.Stderr, "ProxyConnect %s -> %s\n", proxyHost, destHost)

	url := url.URL{
		Scheme: u.Scheme,
		Host:   proxyHost,
		Opaque: destHost,
	}

	pipeRead, pipeWrite := io.Pipe()

	request := &http.Request{
		Method: "CONNECT",
		URL:    &url,
		Header: map[string][]string{},
		Body:   pipeRead,
	}

	request.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(u.User.String())))

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP Error: %s\n", err)
		os.Exit(1)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "HTTP Error: %s %s\n", response.StatusCode, response.Status)
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
