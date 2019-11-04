package main

// https://github.com/golang/go/issues/17227

import (
	"fmt"
	"net/url"
	"os"
	"io"
	"net"
	"net/http"
	_ "net/http/httputil"
	"context"
	"time"
	"encoding/base64"
)

// https://github.com/golang/go/issues/17227

func main() {
	proxy := os.Getenv("http_proxy")
	u, err := url.Parse(proxy)
	if err == nil {
		_ = u
		fmt.Fprintf(os.Stderr, "%s %s %s \n", u.Hostname(), u.Opaque, u.Port())
		fmt.Fprintf(os.Stderr, "%s %s \n", u.User.String(), u.Port())
	}

	fmt.Fprintf(os.Stderr, "%s\n", base64.StdEncoding.EncodeToString([]byte(u.User.String())))
	fmt.Fprintf(os.Stderr, "%s\n", os.Args[1] + ":" + os.Args[2])
	//os.Exit(255)

	url := url.URL{
		Scheme: "http",
		Host:   u.Hostname() + ":" + u.Port(),
		Opaque: os.Args[1] + ":" + os.Args[2],
	}

	db, _ := os.Create("debug.log")
	defer db.Close()

	pr, pw := io.Pipe()

	req := &http.Request{
		Method: "CONNECT",
		URL: &url,
		Header: map[string][]string {
		},
		Body: pr,
		ContentLength: -1,
	}

	req.Header.Set("Proxy-Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(u.User.String())))

	var htConn net.Conn

	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext(ctx, network, addr)
			htConn = conn
			return conn, err
		},
	}
	client := &http.Client{
		Transport: tr,
	}

	res, err := client.Do(req)
	//res, err := httputil.DumpRequestOut(req, true)
	//fmt.Printf("res: %s\nerr: %s\n", res, err)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err %s\n", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "Code %s\n", 	res.StatusCode)
		os.Exit(1)
	} 

	// res.Close = false

	//conn := res.Body
	_ = res
	//fmt.Printf("conn: %s\n", conn)


	go func(c io.ReadCloser) {
		defer c.Close()
		for true {
			var b [40960]byte
			n, err := c.Read(b[:])
			fmt.Fprintf(db, "Read Result: %s %s\n", n, err)
			db.Sync()
			fmt.Fprintf(db, "Read: %s\n", b[:n])
			db.Sync()
			os.Stdout.Write(b[:n])
			_ = err
		}		
	}(res.Body)

	_ = pw

	for true {
		var b [2000]byte
		n, _ := os.Stdin.Read(b[:])
		if n == 0 {
			break
		}
		// fmt.Printf("Write: %s\n", b[:n])
		fmt.Fprintf(db, "Write: %s\n", b[:n])
		db.Sync()
		// n, err := htConn.Write(b[:n])
		n, err := pw.Write(b[:n])
		fmt.Fprintf(db, "Write Result: %s %s\n", n, err)
		db.Sync()
	}
}
