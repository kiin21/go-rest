package httputil

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type RequestURLResolver struct{}

func NewRequestURLResolver() *RequestURLResolver {
	return &RequestURLResolver{}
}

func (r RequestURLResolver) Scheme(c *gin.Context) string {
	if c != nil && c.Request != nil {
		if s := getForwardedProto(c); s != "" {
			return s
		}
		if c.Request.TLS != nil {
			return "https"
		}
	}
	return "http"
}

func (r RequestURLResolver) Host(c *gin.Context) string {
	if c != nil && c.Request != nil {
		if c.Request.Host != "" {
			return c.Request.Host
		}
		if h := c.GetHeader("X-Forwarded-Host"); h != "" {
			return strings.TrimSpace(strings.Split(h, ",")[0])
		}
	}
	return ""
}

func (r RequestURLResolver) AbsoluteURL(c *gin.Context, path string, query url.Values) string {
	u := url.URL{
		Scheme:   r.Scheme(c),
		Host:     r.Host(c),
		Path:     path,
		RawQuery: query.Encode(),
	}
	return u.String()
}

func getForwardedProto(c *gin.Context) string {
	if s := c.GetHeader("X-Forwarded-Proto"); s != "" {
		return strings.ToLower(strings.TrimSpace(strings.Split(s, ",")[0]))
	}
	if s := c.GetHeader("X-Forwarded-Scheme"); s != "" {
		return strings.ToLower(strings.TrimSpace(strings.Split(s, ",")[0]))
	}
	return ""
}
