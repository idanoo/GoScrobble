package goscrobble

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var usernameRegex = regexp.MustCompile("^[a-zA-Z0-9_\\.]+$")

// decodeJson - Returns a map[string]interface{}
func decodeJson(body io.ReadCloser) (map[string]interface{}, error) {
	var jsonInput map[string]interface{}
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&jsonInput)

	return jsonInput, err
}

// isEmailValid - checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 5 && len(e) > 254 {
		return false
	}

	if !emailRegex.MatchString(e) {
		return false
	}

	// Do MX lookup
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}

	return true
}

// isUsernameValid - Checks if username is alphanumeric+underscores+dots
func isUsernameValid(e string) bool {
	if len(e) > 64 {
		return false
	}

	return usernameRegex.MatchString(e)
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

// getUserIp - Returns IP that isn't set in REVERSE_PROXIES
func getUserIp(r *http.Request) net.IP {
	var ip net.IP
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	if contains(ReverseProxies, host) {
		forwardedFor := r.Header.Get("X-Forwarded-For")
		if !contains(ReverseProxies, forwardedFor) {
			host = forwardedFor
		} else {
			realIp := r.Header.Get("X-Real-IP")
			if !contains(ReverseProxies, realIp) {
				host = realIp
			}
		}
	}

	if host == "" {
		host = "0.0.0.0"
	}

	ip = net.ParseIP(host)
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

// calcPageOffsetString - Used to SQL paging
func calcPageOffsetString(page int, offset int) string {
	return fmt.Sprintf("%d", page*offset)
}

// timestampToSeconds - Converts HH:MM:SS to (int)seconds
func timestampToSeconds(timestamp string) int {
	var h, m, s int
	n, err := fmt.Sscanf(timestamp, "%d:%d:%d", &h, &m, &s)
	if err != nil || n != 3 {
		return 0
	}
	return h*3600 + m*60 + s
}

func filterSlice(s []string) []string {
	m := make(map[string]bool)
	for _, item := range s {
		if item != "" {
			if _, ok := m[strings.TrimSpace(item)]; !ok {
				m[strings.TrimSpace(item)] = true
			}
		}
	}

	var result []string
	for item, _ := range m {
		result = append(result, item)
	}

	fmt.Printf("RESTULS: %+v", result)
	return result
}

func isValidTimezone(tz string) bool {
	_, err := time.LoadLocation(tz)
	return err == nil
}
