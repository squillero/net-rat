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
	"log"
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
	if err != nil {
		return
	}
	var cooked map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &cooked); err != nil {
		log.Println("fetchJson: ", err)
	}
	info := IpInfo{
		RawIp:     cooked[tag].(string),
		CookedIp:  cooked[tag].(string),
		Source:    url,
		Flags:     PublicIP,
		Timestamp: time.Time{},
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
		CookedIp:  strings.TrimSpace(string(cooked)),
		Source:    url,
		Flags:     PublicIP,
		Timestamp: time.Time{},
	}
	out <- info
}

func getLocalIpUDP(out chan IpInfo) {
	if conn, err := net.Dial("udp", "8.8.8.8:80"); err == nil {
		defer conn.Close()
		localAddress := conn.LocalAddr().(*net.UDPAddr)
		info := IpInfo{
			RawIp:     localAddress.IP.String(),
			CookedIp:  localAddress.IP.String(),
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
	var ips []string
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				slog.Debug("ipnet", "ipnet", ipnet)
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	info := IpInfo{
		RawIp:     strings.Join(ips[:], "/"),
		CookedIp:  strings.Join(ips[:], "/"),
		Source:    "IFace",
		Flags:     LocalIP | CoolIP,
		Timestamp: time.Now(),
	}
	out <- info
}
