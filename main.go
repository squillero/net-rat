//        ____()()     NetRat v0.1
//       /      @@     ~~~~~~~~~~~
// `~~~~~\_;m__m._>o   A Go hack
//
// Coded in July 2024, between Italy and Australia (34,138 km).
// Copyright © 2024 Giovanni Squillero <giovanni.squillero@polito.it>
// Released for educational and academic purposes under 0BSD.

package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
)

const NetRatVersion = "0.1"

func main() {
	log.SetPrefix("🐀 ") // 🐁
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lmsgprefix | log.LUTC)

	// Parse flags
	verbose := flag.Bool("v", false, "Verbose operations")
	clearCache := flag.Bool("c", false, "Clear cache")
	flag.Parse()
	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	slog.Debug("This is NetRat v" + NetRatVersion)
	if *clearCache {
		slog.Debug("Clearing cache", "file", CacheFile)
		os.Remove(CacheFile)
	}

	ni := getNetInfo()
	slog.Debug("NetInfo", "info", ni)
	loopback := ni.GetType(LoopbackIP)
	tunnel := ni.GetType(TunnelIP)
	local := ni.GetType(LocalIP)
	public := ni.GetType(PublicIP)

	slog.Debug("LoopbackIP", "val", loopback)
	slog.Debug("TunnelIP", "val", tunnel)
	slog.Debug("LocalIP", "val", local)
	slog.Debug("PublicIP", "val", public)

	// 	return strings.Join(r[:], "/")

	if hostname, err := os.Hostname(); err == nil {
		fmt.Printf("%s @ ", hostname)
	}

	// https://superuser.com/questions/730591/why-am-i-seeing-multiple-different-ip-addresses-reported-as-my-public-ip
	if len(local) == 0 && len(public) == 0 && len(loopback) == 0 {
		fmt.Print("No network connection")
	} else if len(local) == 0 && len(public) == 0 && len(loopback) == 1 {
		fmt.Printf("%s [Loopback Only]", loopback[0].Describe())
	} else if len(local) > 0 && len(public) == 0 {
		fmt.Printf("%s [Local Only]", local[0].Describe())
	} else if len(local) == 1 && len(public) == 1 && local[0].RawIp == public[0].RawIp {
		fmt.Print(public[0].Describe())
	} else if len(local) == 1 {
		fmt.Printf("%s ", local[0].Describe())

		if len(tunnel) > 0 {
			var r []string
			for _, v := range tunnel {
				r = append(r, v.Describe())
			}
			fmt.Printf("via %s ", strings.Join(r[:], "/"))
		}

		if len(public) == 1 {
			fmt.Printf("// %s", public[0].Describe())
		} else if len(public) > 1 {
			d := ""
			for _, v := range public {
				if len(v.Comment) > len(d) {
					d = v.Comment
				}
			}
			fmt.Printf("// NAT outbound [%s]", d)
		}
	}
	fmt.Println()
}
