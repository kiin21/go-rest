package httpctx

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequestURLResolver centralizes logic for inferring request metadata such as scheme and host.
// It inspects the current HTTP request (including proxy headers) and falls back to
// optional defaults derived from configuration so the same resolver can be reused
// in background jobs or tests that do not have full HTTP context.
type RequestURLResolver struct {
	defaultScheme string
	defaultHost   string
}

// NewRequestURLResolver parses the provided baseURL (for example, "https://example.com")
// and stores its scheme/host as the fallback values. If baseURL is empty or invalid,
// the resolver simply keeps empty defaults and will later fall back to "http".
func NewRequestURLResolver(baseURL string) RequestURLResolver {
	if baseURL == "" {
		return RequestURLResolver{}
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return RequestURLResolver{}
	}

	return RequestURLResolver{
		defaultScheme: parsed.Scheme,
		defaultHost:   parsed.Host,
	}
}

// Scheme returns the best-effort scheme for the current request.
// Priority:
//  1. Forwarded proto headers (trusted proxies must be configured at the router level)
//  2. Presence of TLS on the active request
//  3. Scheme stored in the request URL (useful in tests)
//  4. Default scheme derived from configuration
//  5. Plain "http" as the ultimate fallback
func (r RequestURLResolver) Scheme(ctx *gin.Context) string {
	if ctx != nil && ctx.Request != nil {
		if scheme := detectSchemeFromHeaders(ctx); scheme != "" {
			return scheme
		}

		if ctx.Request.TLS != nil {
			return "https"
		}

		if ctx.Request.URL != nil && ctx.Request.URL.Scheme != "" {
			return ctx.Request.URL.Scheme
		}
	}

	if r.defaultScheme != "" {
		return r.defaultScheme
	}
	return "http"
}

// Host returns the host for the current request, preferring the value provided
// by the reverse proxy headers if Gin considers the proxy trustworthy. If the
// request does not expose a host, the resolver falls back to the configured default.
func (r RequestURLResolver) Host(ctx *gin.Context) string {
	if ctx != nil && ctx.Request != nil {
		if host := ctx.Request.Host; host != "" {
			return host
		}

		if host := firstHeaderValue(ctx, "X-Forwarded-Host"); host != "" {
			return host
		}

		if host := parseForwardedDirective(ctx.GetHeader("Forwarded"), "host"); host != "" {
			return host
		}
	}

	return r.defaultHost
}

// AbsoluteURL builds an absolute URL using the inferred scheme and host and the supplied
// path/query. The query values are encoded without mutating the original map.
func (r RequestURLResolver) AbsoluteURL(ctx *gin.Context, path string, query url.Values) string {
	clone := url.Values{}
	for key, values := range query {
		clone[key] = append([]string(nil), values...)
	}

	u := url.URL{
		Scheme:   r.Scheme(ctx),
		Host:     r.Host(ctx),
		Path:     path,
		RawQuery: clone.Encode(),
	}

	return u.String()
}

func detectSchemeFromHeaders(ctx *gin.Context) string {
	if scheme := firstHeaderValue(ctx, "X-Forwarded-Proto"); scheme != "" {
		return scheme
	}

	if scheme := firstHeaderValue(ctx, "X-Forwarded-Scheme"); scheme != "" {
		return scheme
	}

	if scheme := parseForwardedDirective(ctx.GetHeader("Forwarded"), "proto"); scheme != "" {
		return scheme
	}

	return ""
}

func firstHeaderValue(ctx *gin.Context, key string) string {
	if ctx == nil {
		return ""
	}

	header := ctx.GetHeader(key)
	if header == "" {
		return ""
	}

	// According to the spec, multiple values are comma-separated; we only need the first one.
	parts := strings.Split(header, ",")
	if len(parts) == 0 {
		return ""
	}

	value := strings.TrimSpace(parts[0])
	return strings.ToLower(value)
}

func parseForwardedDirective(header string, directive string) string {
	if header == "" {
		return ""
	}

	entries := strings.Split(header, ",")
	for _, entry := range entries {
		pairs := strings.Split(entry, ";")
		for _, pair := range pairs {
			if pair == "" {
				continue
			}
			kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
			if len(kv) != 2 {
				continue
			}
			if strings.EqualFold(strings.TrimSpace(kv[0]), directive) {
				return strings.ToLower(strings.Trim(strings.TrimSpace(kv[1]), `"`))
			}
		}
	}

	return ""
}
