package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

func main() {
	portEnv, ok := os.LookupEnv("PORT")
	if !ok {
		portEnv = "8080"
	}

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		log.Fatal(err)
	}

	remoteEnv, ok := os.LookupEnv("REMOTE")
	if !ok {
		log.Fatal("REMOTE environment variable required")
	}

	remote, err := url.Parse(remoteEnv)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.Handle("/", &proxyHandler{proxy})

	log.Printf("proxying requests to %s on port %d ...", remote.String(), port)
	if err = http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Fatal(err)
	}
}

type proxyHandler struct {
	p *httputil.ReverseProxy
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	h.p.ServeHTTP(w, r)
}
