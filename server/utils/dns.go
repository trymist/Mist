package utils

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// returns the public IP, useful when no wildcard domain is configured
func GetServerIP() (string, error) {
	// TODO: remove this Getenv, not being used
	if serverIP := os.Getenv("SERVER_IP"); serverIP != "" {
		return serverIP, nil
	}

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// validate if dns is pointed to the server IP
func ValidateDNS(domain string) (bool, error) {
	serverIP, err := GetServerIP()
	if err != nil {
		return false, fmt.Errorf("failed to get server IP: %w", err)
	}

	ips, err := net.LookupIP(domain)
	if err != nil {
		return false, fmt.Errorf("DNS lookup failed: %w", err)
	}

	if len(ips) == 0 {
		return false, fmt.Errorf("no DNS records found")
	}

	for _, ip := range ips {
		if ip.String() == serverIP {
			return true, nil
		}
	}

	return false, fmt.Errorf("DNS does not point to server (expected: %s, found: %s)", serverIP, formatIPs(ips))
}

func ValidateDNSWithTimeout(domain string, timeout time.Duration) (bool, error) {
	type result struct {
		valid bool
		err   error
	}

	ch := make(chan result, 1)

	go func() {
		valid, err := ValidateDNS(domain)
		ch <- result{valid, err}
	}()

	select {
	case res := <-ch:
		return res.valid, res.err
	case <-time.After(timeout):
		return false, fmt.Errorf("DNS validation timeout after %v", timeout)
	}
}

func formatIPs(ips []net.IP) string {
	strs := make([]string, len(ips))
	for i, ip := range ips {
		strs[i] = ip.String()
	}
	return strings.Join(strs, ", ")
}
