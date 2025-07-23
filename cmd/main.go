package main

import (
	"golang.org/x/net/publicsuffix"
	"log"
	"net/http"
)

func init() {
	// 规整日志格式，包含日期、时间、微秒
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

func infoLog(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

func errorLog(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	rootDomain, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		errorLog("Failed to parse root domain from host %q: %v", host, err)
		rootDomain = host
	}

	target := "https://www." + rootDomain + r.URL.RequestURI()
	infoLog("Redirecting %s%s to %s", host, r.URL.RequestURI(), target)
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}

func main() {
	http.HandleFunc("/", redirectHandler)

	port := "80"
	infoLog("Starting redirect gateway on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		errorLog("ListenAndServe failed: %v", err)
	}
}
