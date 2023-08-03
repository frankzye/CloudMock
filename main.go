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
	"path/filepath"
)

//goland:noinspection GoUnhandledErrorResult
func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func NewTLSServer(handler http.Handler) *httptest.Server {
	caCert, caKey, _ := LoadX509KeyPair(filepath.Join(GetRoot(), "public.pem"), filepath.Join(GetRoot(), "private.pem"))

	ts := httptest.NewUnstartedServer(handler)
	ts.TLS = &tls.Config{
		ClientAuth:         0,
		InsecureSkipVerify: true,
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			pemCert, pemKey := CreateCert([]string{info.ServerName}, caCert, caKey, 1)
			cert, err := tls.X509KeyPair(pemCert, pemKey)
			return &cert, err
		},
	}
	ts.StartTLS()
	return ts
}

func main() {
	requests := ReadMapping()

	ts := NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, ExpressRule(requests, &Request{
			Host:   r.Host,
			Path:   r.RequestURI,
			Method: r.Method,
		}))
	}))
	defer ts.Close()

	server := &http.Server{
		Addr: ":9999",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			localUrl, _ := url.Parse(ts.URL)
			if r.Method == http.MethodConnect {
				destConn, err := net.Dial("tcp", localUrl.Host)
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
