package ws

import (
	"net/http"
	"testing"
)

func TestResolveRealIP(t *testing.T) {
	tests := []struct {
		name       string
		mode       RealIPMode
		xff        string
		xRealIP    string
		remoteAddr string
		want       string
	}{
		{
			name:       "xff right mode",
			mode:       RealIPModeRight,
			xff:        "1.1.1.1, 2.2.2.2",
			remoteAddr: "10.0.0.1:12345",
			want:       "2.2.2.2",
		},
		{
			name:       "xff left mode",
			mode:       RealIPModeLeft,
			xff:        "1.1.1.1, 2.2.2.2",
			remoteAddr: "10.0.0.1:12345",
			want:       "1.1.1.1",
		},
		{
			name:       "fallback x-real-ip",
			mode:       RealIPModeRight,
			xff:        "unknown, -,",
			xRealIP:    "3.3.3.3",
			remoteAddr: "10.0.0.1:12345",
			want:       "3.3.3.3",
		},
		{
			name:       "fallback remote addr",
			mode:       RealIPModeRight,
			remoteAddr: "4.4.4.4:56789",
			want:       "4.4.4.4",
		},
		{
			name:       "token trim host port ipv6",
			mode:       RealIPModeRight,
			xff:        "  bad-token  , \"5.5.5.5:1234\", [2001:db8::1]:443 ",
			remoteAddr: "10.0.0.1:12345",
			want:       "2001:db8::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header:     make(http.Header),
				RemoteAddr: tt.remoteAddr,
			}
			if tt.xff != "" {
				req.Header.Set("X-Forwarded-For", tt.xff)
			}
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			if got := resolveRealIP(req, tt.mode); got != tt.want {
				t.Fatalf("resolveRealIP() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveRealIPNilRequest(t *testing.T) {
	if got := resolveRealIP(nil, RealIPModeRight); got != "" {
		t.Fatalf("resolveRealIP(nil) = %q, want empty", got)
	}
}
