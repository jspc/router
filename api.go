package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type API struct {
	Docker Docker
	Domain string
}

func NewAPI(docker Docker, domain string) API {
	if !strings.HasPrefix(domain, ".") {
		domain = fmt.Sprintf(".%s", domain)
	}

	return API{
		Docker: docker,
		Domain: domain,
	}
}

func (a API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	containerName := a.InferContainerName(req.Host)

	scheme, addr, err := a.Docker.GetContainerAddress(containerName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)

		return
	}

	req.URL.Scheme = *scheme
	req.URL.Host = addr

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)

		return
	}

	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (a API) InferContainerName(host string) string {
	return strings.TrimSuffix(host, a.Domain)
}
