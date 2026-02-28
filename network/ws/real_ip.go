package ws

import (
	"net"
	"net/http"
	"strings"
)

func resolveRealIP(r *http.Request, mode RealIPMode) string {
	if r == nil {
		return ""
	}

	if ip := parseXForwardedFor(r.Header.Get("X-Forwarded-For"), mode); ip != "" {
		return ip
	}

	if ip := parseIPToken(r.Header.Get("X-Real-IP")); ip != "" {
		return ip
	}

	return parseIPToken(r.RemoteAddr)
}

func parseXForwardedFor(value string, mode RealIPMode) string {
	tokens := strings.Split(value, ",")
	if len(tokens) == 0 {
		return ""
	}

	switch parseRealIPMode(string(mode)) {
	case RealIPModeLeft:
		for i := 0; i < len(tokens); i++ {
			if ip := parseIPToken(tokens[i]); ip != "" {
				return ip
			}
		}
	default:
		for i := len(tokens) - 1; i >= 0; i-- {
			if ip := parseIPToken(tokens[i]); ip != "" {
				return ip
			}
		}
	}

	return ""
}

func parseIPToken(token string) string {
	token = strings.TrimSpace(strings.Trim(token, `"`))
	if token == "" {
		return ""
	}

	if ip := net.ParseIP(token); ip != nil {
		return ip.String()
	}

	if host, _, err := net.SplitHostPort(token); err == nil {
		host = strings.TrimSpace(strings.Trim(host, "[]"))
		if ip := net.ParseIP(host); ip != nil {
			return ip.String()
		}
	}

	if ip := net.ParseIP(strings.Trim(token, "[]")); ip != nil {
		return ip.String()
	}

	return ""
}
