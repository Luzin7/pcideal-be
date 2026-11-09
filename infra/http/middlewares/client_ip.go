package middlewares

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

var ErrInvalidClientIPHeader = errors.New("invalid CLIENT_IP_HEADER")

type ClientIPResolver struct {
	header string
}

func NewClientIPResolver(header string) (ClientIPResolver, error) {
	name := strings.TrimSpace(header)
	if name == "" {
		return ClientIPResolver{}, nil
	}

	if !isValidHeaderName(name) {
		return ClientIPResolver{}, fmt.Errorf("%w: %q", ErrInvalidClientIPHeader, header)
	}

	return ClientIPResolver{header: name}, nil
}

func isValidHeaderName(name string) bool {
	if name == "" {
		return false
	}

	for _, char := range name {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' {
			continue
		}

		if strings.ContainsRune("!#$%&'*+-.^_`|~", char) {
			continue
		}

		return false
	}

	return true
}

func (r ClientIPResolver) Resolve(c *gin.Context) string {
	if r.header == "" {
		return c.ClientIP()
	}

	address := net.ParseIP(strings.TrimSpace(c.GetHeader(r.header)))
	if address == nil {
		return c.ClientIP()
	}

	return address.String()
}
