package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os/user"
	"rsc.io/letsencrypt"
)

func run(register, listenhttp, listenhttps, state, serve string) {
	var letsEncryptManager letsencrypt.Manager
	if err := letsEncryptManager.CacheFile(state); err != nil {
		log.Fatal(err)
	}
	if register != "" && !letsEncryptManager.Registered() {
		letsEncryptManager.Register(register, func(terms string) bool {
			log.Printf("Agreeing to %s ...", terms)
			return true
		})
	}

	srv := &http.Server{
		Addr:      listenhttps,
		TLSConfig: &tls.Config{GetCertificate: letsEncryptManager.GetCertificate},
		Handler:   http.FileServer(http.Dir(serve)),
	}

	go func() { log.Fatal(http.ListenAndServe(listenhttp, http.HandlerFunc(letsencrypt.RedirectHTTP))) }()
	log.Fatal(srv.ListenAndServeTLS("", ""))
}

var register = flag.String("register", "", "(optional) email address for letsencrypt registration")
var listenhttp = flag.String("listenhttp", ":http", "host:port to listen for HTTP on")
var listenhttps = flag.String("listenhttps", ":https", "host:port to listen for HTTPS on")
var serve = flag.String("serve", ".", "directory to serve files from")

var default_statefile = "~/.letsencrypt.cache" // Go does NOT understand '~' paths, correct in main()
var state = flag.String("state", default_statefile, "path at which the letsencrypt server state will be recorded")

func main() {
	flag.Parse()
	state := *state
	if state == default_statefile {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		state = usr.HomeDir + default_statefile[1:]
	}
	run(*register, *listenhttp, *listenhttps, state, *serve)
}