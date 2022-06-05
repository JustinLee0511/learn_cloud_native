package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
)

/*
1.接收客户端 request，并将 request 中带的 header 写入 response header
2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4.当访问 {url}/healthz 时，应返回200
*/

func main() {
	mux := http.NewServeMux() // 06. debug
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/", index)
	mux.HandleFunc("/healthz", healthz)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("start http server failed, error: %s\n", err.Error())
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	os.Setenv("VERSION", "0.1")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION,", version)
	fmt.Printf("os version:%s\n", version)
	w.Write([]byte("<h1>Welcome to cloudnative</h1>"))
	for k, v := range r.Header {
		for _, vv := range v {
			fmt.Printf("Header key: %s,Header value: %s", k, v)
			w.Header().Set(k, vv)
		}
	}
	clientIP := ClientIP(r)
	log.Printf("Response code: %d", 200)
	log.Printf("clientIP: %s", clientIP)

}

func getCurrentIP(r *http.Request) string {
	ip := r.Header.Get("X-REAL-IP")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
		//ip:port
	}
	return ip
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "working")
}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
