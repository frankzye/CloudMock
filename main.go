package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"
)

//goland:noinspection GoUnhandledErrorResult
func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func main() {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, r.Host)
	}))
	defer ts.Close()

	server := &http.Server{
		Addr: ":9999",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			localUrl, _ := url.Parse(ts.URL)
			if r.Method == http.MethodConnect {

				destConn, err := net.DialTimeout("tcp", localUrl.Host, 10*time.Second)
				w.WriteHeader(http.StatusOK)
				hijacker, ok := w.(http.Hijacker)
				if !ok {
					http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
					return
				}
				clientConn, _, err := hijacker.Hijack()
				if err != nil {
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
				}
				go transfer(destConn, clientConn)
				go transfer(clientConn, destConn)

			} else {
				fmt.Println("unexpected request")
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Fatal(server.ListenAndServe())
}
