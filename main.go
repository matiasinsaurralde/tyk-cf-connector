package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

const (
	// defaultPort is the TCP port the connector will bind to, hopefully it will be provided by the CF environment.
	defaultPort = 7000
	// cfURLHeader is the HTTP header that contains the target URL.
	cfURLHeader = "X-Cf-Forwarded-Url"
)

func main() {
	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			forwardedURL := req.Header.Get(cfURLHeader)

			if forwardedURL == "" {
				log.Fatalln("No header")
			}

			url, err := url.Parse(forwardedURL)
			if err != nil {
				log.Fatalln(err.Error())
			}

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatalln(err.Error())
			}
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			req.URL = url
			req.Host = url.Host
		},
		Transport: http.DefaultTransport,
	}
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = strconv.Itoa(defaultPort)
	}
	http.ListenAndServe(":"+port, reverseProxy)
}
