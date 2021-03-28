package goscrobble

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// decodeJson - Returns a map[string]interface{}
func decodeJson(body io.ReadCloser) (map[string]interface{}, error) {
	var jsonInput map[string]interface{}
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&jsonInput)

	return jsonInput, err
}

// isEmailValid - checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

// contains - Check if string is in list
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// getUserIp - Returns IP that isn't set in REVERSE_PROXY
func getUserIp(r *http.Request) net.IP {
	var ip net.IP
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	if contains(ReverseProxies, host) {
		forwardedFor := r.Header.Get("X-Forwarded-For")
		if forwardedFor != "" {
			host = forwardedFor
		}
		// 	realIp := r.Header.Get("X-Real-IP")
	}

	ip = net.ParseIP(host)
	log.Printf("%+v", ip)
	return ip
}

// Inet_Aton converts an IPv4 net.IP object to a 64 bit integer.
func Inet_Aton(ip net.IP) int64 {
	ipv4Int := big.NewInt(0)
	ipv4Int.SetBytes(ip.To4())
	return ipv4Int.Int64()
}

// Inet6_Aton converts an IP Address (IPv4 or IPv6) net.IP object to a hexadecimal
// representaiton. This function is the equivalent of
// inet6_aton({{ ip address }}) in MySQL.
func Inet6_Aton(ip net.IP) string {
	ipv4 := false
	if ip.To4() != nil {
		ipv4 = true
	}

	ipInt := big.NewInt(0)
	if ipv4 {
		ipInt.SetBytes(ip.To4())
		ipHex := hex.EncodeToString(ipInt.Bytes())
		return ipHex
	}

	ipInt.SetBytes(ip.To16())
	ipHex := hex.EncodeToString(ipInt.Bytes())
	return ipHex
}
