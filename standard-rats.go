//        ____()()     NetRat v0.1
//       /      @@     ~~~~~~~~~~~
// `~~~~~\_;m__m._>o   A Go hack
//
// Coded in July 2024, between Italy and Australia (34,138 km).
// Copyright © 2024 Giovanni Squillero <giovanni.squillero@polito.it>
// Released for educational and academic purposes under 0BSD.

package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

func fetchJson(out chan IpInfo, url, tag string) {
	result, err := http.Get(url)
	if err != nil {
		return
	}
	raw, err := io.ReadAll(result.Body)
	if err != nil || raw == nil {
		return
	}
	var cooked map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &cooked); err != nil {
		slog.Debug("fetchJson: ", "err", err)
	}
	info := IpInfo{
		RawIp:     cooked[tag].(string),
		Source:    url,
		Flags:     PublicIP,
		Timestamp: time.Now(),
	}
	out <- info
}

func fetchRaw(out chan IpInfo, url string) {
	result, err := http.Get(url)
	if err != nil {
		return
	}
	cooked, err := io.ReadAll(result.Body)
	if err != nil {
		return
	}
	info := IpInfo{
		RawIp:     strings.TrimSpace(string(cooked)),
		Source:    url,
		Flags:     PublicIP,
		Timestamp: time.Now(),
	}
	out <- info
}

func getLocalIpUDP(out chan IpInfo) {
	if conn, err := net.Dial("udp", "8.8.8.8:80"); err == nil {
		defer conn.Close()
		localAddress := conn.LocalAddr().(*net.UDPAddr)
		info := IpInfo{
			RawIp:     localAddress.IP.String(),
			Source:    "udp",
			Flags:     LocalIP,
			Timestamp: time.Now(),
		}
		out <- info
	}
}

func getLocalIpIFACE(out chan IpInfo) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		//if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				var t IpFlags
				if ipnet.IP.IsLoopback() {
					t = LoopbackIP
				} else if ipnet.IP.IsPrivate() {
					t = LocalIP
				} else if ipnet.IP.IsGlobalUnicast() {
					t = PublicIP
				} else {
					slog.Debug("***IP", "ip", ipnet.IP)
				}
				if t > 0 {
					out <- IpInfo{
						RawIp:     ipnet.IP.String(),
						Source:    "IFace",
						Flags:     t,
						Timestamp: time.Now(),
					}
				}
			}
		}
	}
}
