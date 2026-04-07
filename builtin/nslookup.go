package builtin

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func init() {
	Register("nslookup", NSLookup, "Resolve DNS records (A/AAAA, CNAME, NS, TXT)")
}

func NSLookup(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "nslookup: usage: nslookup <hostname>")
		return
	}
	name := strings.TrimSpace(args[0])
	if name == "" {
		fmt.Fprintln(os.Stderr, "nslookup: empty hostname")
		return
	}

	start := time.Now()

	// CNAME (alias) if present
	if cname, err := net.LookupCNAME(name); err == nil && !strings.EqualFold(cname, name) {
		fmt.Printf("Name:\t%s\nAlias:\t%s\n\n", name, cname)
	}

	// A/AAAA
	if ips, err := net.LookupIP(name); err == nil && len(ips) > 0 {
		fmt.Printf("Non-authoritative answer:\n")
		for _, ip := range ips {
			fam := "A"
			if ip.To4() == nil {
				fam = "AAAA"
			}
			fmt.Printf("%s\t%s\n", fam, ip.String())
		}
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "nslookup: lookup A/AAAA: %v\n", err)
	}

	// NS
	if nss, err := net.LookupNS(name); err == nil && len(nss) > 0 {
		fmt.Printf("\nNS records:\n")
		for _, ns := range nss {
			fmt.Println(ns.Host)
		}
	}

	// TXT
	if txts, err := net.LookupTXT(name); err == nil && len(txts) > 0 {
		fmt.Printf("\nTXT records:\n")
		for _, t := range txts {
			fmt.Println(t)
		}
	}

	fmt.Printf("\nQuery time: %v\n", time.Since(start).Truncate(time.Millisecond))
}
