package builtin

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func init() {
	Register("ping", Ping, "Send ICMP echo requests (requires CAP_NET_RAW)")
}

func Ping(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "ping: usage: ping <host> [-c N] [-i SEC] [-W SEC]")
		return
	}

	host := args[0]
	count := 4
	interval := time.Second
	timeout := 2 * time.Second

	// Parse simple flags
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-c":
			if i+1 < len(args) {
				if v, err := strconv.Atoi(args[i+1]); err == nil && v > 0 {
					count = v
					i++
				}
			}
		case "-i":
			if i+1 < len(args) {
				if v, err := strconv.ParseFloat(args[i+1], 64); err == nil && v > 0 {
					interval = time.Duration(v * float64(time.Second))
					i++
				}
			}
		case "-W":
			if i+1 < len(args) {
				if v, err := strconv.ParseFloat(args[i+1], 64); err == nil && v > 0 {
					timeout = time.Duration(v * float64(time.Second))
					i++
				}
			}
		}
	}

	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		fmt.Fprintf(os.Stderr, "ping: resolve %s: %v\n", host, err)
		return
	}

	var dst net.IP
	for _, ip := range ips {
		if ip = ip.To4(); ip != nil {
			dst = ip
			break
		}
	}

	if dst == nil {
		fmt.Fprintf(os.Stderr, "ping: %s has no IPv4 address (IPv6 not supported in this minimal ping)\n", host)
		return
	}

	fmt.Printf("PING %s (%s): %d data bytes\n", host, dst.String(), 56)

	// Open raw ICMP socket (IPv4)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		if errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EACCES) {
			fmt.Fprintln(os.Stderr, "ping: operation not permitted (need CAP_NET_RAW or root)")
		} else {
			fmt.Fprintf(os.Stderr, "ping: socket error: %v\n", err)
		}
		return
	}
	defer syscall.Close(fd)

	// Set receive timeout
	tv := syscall.Timeval{
		Sec:  int64(timeout / time.Second),
		Usec: int64((timeout % time.Second) / time.Microsecond),
	}
	_ = syscall.SetsockoptTimeval(fd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv)

	// Destination sockaddr
	sa := &syscall.SockaddrInet4{Port: 0}
	copy(sa.Addr[:], dst.To4())

	pid := os.Getpid() & 0xffff
	seq := 0

	var received int
	var minRTT, maxRTT, sumRTT time.Duration

	for i := 0; i < count; i++ {
		seq++
		// Build ICMP Echo (type 8, code 0)
		payload := []byte("ctenterd-ping-payload")
		echo := buildICMPEcho(uint16(pid), uint16(seq), payload)

		sendTime := time.Now()
		if err := syscall.Sendto(fd, echo, 0, sa); err != nil {
			fmt.Fprintf(os.Stderr, "ping: send: %v\n", err)
			time.Sleep(interval)
			continue
		}

		// Receive
		buf := make([]byte, 1500)
		n, from, rerr := syscall.Recvfrom(fd, buf, 0)
		rtt := time.Since(sendTime)

		if rerr != nil {
			if isTimeout(rerr) {
				fmt.Printf("Request timeout for icmp_seq %d\n", seq)
			} else {
				fmt.Fprintf(os.Stderr, "ping: recv: %v\n", rerr)
			}
			time.Sleep(interval)
			continue
		}

		// Parse reply (skip IP header; assume minimal 20 bytes)
		if n < 28 {
			time.Sleep(interval)
			continue
		}
		icmp := buf[20:n]
		if icmp[0] != 0 || icmp[1] != 0 { // type 0 code 0 = Echo Reply
			time.Sleep(interval)
			continue
		}
		// id/seq
		id := binary.BigEndian.Uint16(icmp[4:6])
		sq := binary.BigEndian.Uint16(icmp[6:8])
		if int(id) != pid || int(sq) != seq {
			time.Sleep(interval)
			continue
		}

		received++
		if received == 1 || rtt < minRTT {
			minRTT = rtt
		}
		if rtt > maxRTT {
			maxRTT = rtt
		}
		sumRTT += rtt

		src := ""
		switch a := from.(type) {
		case *syscall.SockaddrInet4:
			src = net.IP(a.Addr[:]).String()
		default:
			src = dst.String()
		}

		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", n-20, src, seq, rtt.Truncate(time.Microsecond))
		time.Sleep(interval)
	}

	lost := count - received
	var avg time.Duration
	if received > 0 {
		avg = sumRTT / time.Duration(received)
	}
	fmt.Printf("\n--- %s ping statistics ---\n", host)
	fmt.Printf("%d packets transmitted, %d received, %.1f%% packet loss\n",
		count, received, 100*float64(lost)/float64(count))
	if received > 0 {
		fmt.Printf("round-trip min/avg/max = %v/%v/%v\n",
			minRTT.Truncate(time.Microsecond),
			avg.Truncate(time.Microsecond),
			maxRTT.Truncate(time.Microsecond))
	}
}

func isTimeout(err error) bool {
	// On Go+Linux, timeouts come as EAGAIN/EWOULDBLOCK when SO_RCVTIMEO is set.
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "deadline") ||
		strings.Contains(msg, "resource temporarily unavailable")
}

func checksum(b []byte) uint16 {
	var sum uint32
	for i := 0; i+1 < len(b); i += 2 {
		sum += uint32(binary.BigEndian.Uint16(b[i : i+2]))
	}
	if len(b)%2 == 1 {
		sum += uint32(b[len(b)-1]) << 8
	}
	for (sum >> 16) > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}

func buildICMPEcho(id, seq uint16, payload []byte) []byte {
	hdr := make([]byte, 8+len(payload))
	hdr[0] = 8 // Echo request
	hdr[1] = 0 // code
	// checksum later
	binary.BigEndian.PutUint16(hdr[4:6], id)
	binary.BigEndian.PutUint16(hdr[6:8], seq)
	copy(hdr[8:], payload)
	sum := checksum(hdr)
	binary.BigEndian.PutUint16(hdr[2:4], sum)
	return hdr
}
