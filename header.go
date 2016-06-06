package hpkp

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Header holds a domain's hpkp information
type Header struct {
	Created           int64
	MaxAge            int64
	IncludeSubDomains bool
	Permanent         bool
	Sha256Pins        []string
}

// Matches checks whether the provided pin is in the header list
func (h *Header) Matches(pin string) bool {
	for i := range h.Sha256Pins {
		if h.Sha256Pins[i] == pin {
			return true
		}
	}
	return false
}

// ParseHeader parses the hpkp information from an http.Response. It should only
// be used on HTTPS connections.
func ParseHeader(resp *http.Response) *Header {
	header := &Header{
		Sha256Pins: []string{},
	}

	v, ok := resp.Header["Public-Key-Pins"]
	if !ok {
		return header
	}

	for _, field := range strings.Split(v[0], ";") {
		field = strings.TrimSpace(field)

		i := strings.Index(field, "pin-sha256")
		if i >= 0 {
			header.Sha256Pins = append(header.Sha256Pins, field[i+12:len(field)-1])
			continue
		}

		i = strings.Index(field, "max-age=")
		if i >= 0 {
			ma, err := strconv.Atoi(field[i+8:])
			if err == nil {
				header.MaxAge = int64(ma)
			}
			continue
		}

		if strings.Contains(field, "includeSubDomains") {
			header.IncludeSubDomains = true
			continue
		}
	}

	header.Created = time.Now().Unix()
	return header
}
