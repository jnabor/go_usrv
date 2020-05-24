package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"
)

const message = "Hello Gopher 2020!"

var (
	// GoCerFile environment variable
	GoCerFile = os.Getenv("GO_CERT_FILE")
	// GoKeyFile environment variable
	GoKeyFile = os.Getenv("GO_KEY_FILE")
	// GoSrvAddr environment variable
	//GoSrvAddr = ":8080"
	GoSrvAddr = os.Getenv("GO_SRV_ADDR")
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(message))
	})

	srv := NewServer(mux, GoSrvAddr)
	err := srv.ListenAndServeTLS(GoCerFile, GoKeyFile)
	//err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// NewServer returns *http.Server
func NewServer(mux *http.ServeMux, serverAddress string) *http.Server {
	tlsConfig := &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
		},

		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	srv := &http.Server{
		Addr:         serverAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      mux,
	}
	return srv
}
