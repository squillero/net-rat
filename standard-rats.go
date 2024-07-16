//        ____()()     NetRat v0.2
//       /      @@     ~~~~~~~~~~~
// `~~~~~\_;m__m._>o   A Go experiment
//
// Coded in July 2024, between Italy and Australia (34,138 km).
// Copyright © 2024 Giovanni Squillero <giovanni.squillero@polito.it>
// Released under 0BSD (see LICENSE).

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
			Source:    "net.Dial()",
			Flags:     LocalIP,
			Timestamp: time.Now(),
		}
		out <- info
	}
}

func getLocalIpIFACE(out chan IpInfo) {
	if ifaces, err := net.Interfaces(); err == nil {
		for _, v := range ifaces {
			if addrs, err := v.Addrs(); err == nil {
				for _, addr := range addrs {
					//slog.Debug("Address", "net", addr.Network(), "str", addr.String())
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					if ip != nil && ip.To4() != nil {
						var f IpFlags

						if v.Flags&net.FlagLoopback != 0 {
							f = LoopbackIP
						} else if v.Flags&net.FlagBroadcast != 0 {
							f = LocalIP
						} else if v.Flags&net.FlagPointToPoint != 0 {
							f = TunnelIP
						}
						out <- IpInfo{
							RawIp:     ip.String(),
							Comment:   "",
							Source:    "Interfaces()",
							Flags:     f,
							Timestamp: time.Now(),
						}
					}
				}
			}
		}
	}
}
