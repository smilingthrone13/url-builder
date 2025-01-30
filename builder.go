package url_builder

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

type Builder struct {
	scheme      string
	domain      string
	port        int
	credentials *credentials
	path        []string
	query       map[string][]string
	anchor      string
}

type credentials struct {
	user     string
	password string
}

// New creates a new empty Builder instance. Scheme set to "http" by default.
func New() *Builder {
	return &Builder{
		scheme: "http",
		query:  make(map[string][]string),
	}
}

// WithScheme sets the URL scheme (e.g., "http", "https", "ftp").
// Passed value is not validated.
func (b *Builder) WithScheme(scheme string) *Builder {
	b.scheme = strings.Trim(scheme, ":/")
	return b
}

// WithSchemeHTTP sets the URL scheme to HTTP.
func (b *Builder) WithSchemeHTTP() *Builder {
	b.scheme = "http"
	return b
}

// WithSchemeHTTPS sets the URL scheme to HTTPS.
func (b *Builder) WithSchemeHTTPS() *Builder {
	b.scheme = "https"
	return b
}

// WithDomain sets the domain name of the URL.
// Build will return an error if input string contains slashes or colons ("/", ":").
func (b *Builder) WithDomain(domain string) *Builder {
	b.domain = strings.TrimSuffix(domain, "/")
	return b
}

// WithIPv4 sets given IPv4 address as the domain of the URL.
// Build will return an error if input string contains slashes or colons ("/", ":").
func (b *Builder) WithIPv4(address string) *Builder {
	b.domain = strings.TrimSuffix(address, "/")
	return b
}

// WithIPv6 sets given IPv6 address as the domain of the URL.
// Build will return an error if input string contains slashes ("/").
func (b *Builder) WithIPv6(address string) *Builder {
	address = strings.TrimSuffix(strings.Trim(address, "[]"), "/")
	b.domain = fmt.Sprintf("[%s]", address)
	return b
}

// WithPort sets the port number.
// Build will return an error if port not in [1, 65535] range.
func (b *Builder) WithPort(port int) *Builder {
	b.port = port
	return b
}

// WithCredentials sets the username and password for authentication.
// Build will return an error if one of the parameters is empty.
func (b *Builder) WithCredentials(user, password string) *Builder {
	b.credentials = &credentials{
		user:     user,
		password: password,
	}
	return b
}

// WithPath appends path segments to the URL path.
// Multiple calls to this method add more segments instead of replacing them.
func (b *Builder) WithPath(elements ...string) *Builder {
	b.path = append(b.path, elements...)
	return b
}

// WithQuery adds query parameters to the URL.
// If the same key is added multiple times, values are appended.
// Build will return an error on empty keys or values.
func (b *Builder) WithQuery(key string, values ...string) *Builder {
	b.query[key] = append(b.query[key], values...)
	return b
}

// WithAnchor sets the fragment (anchor) part of the URL.
func (b *Builder) WithAnchor(anchor string) *Builder {
	b.anchor = strings.Trim(anchor, "#/")
	return b
}

// Build constructs the final URL string based on the provided data.
// Returns an error if data is invalid.
func (b *Builder) Build() (string, error) {
	if b.domain == "" {
		return "", fmt.Errorf("domain is required")
	}

	// check given domain
	// todo: can't detect if given ipv6 contains port, so result string might be broken.
	if strings.Contains(b.domain, "/") || // assume domain contains scheme
		strings.Count(b.domain, ":") == 1 { // assume domain contains port (valid ipv6 have at least 2 colons)
		return "", fmt.Errorf("domain contains forbidden symbols")
	}

	rawBaseUrl := fmt.Sprintf("%s://%s", b.scheme, b.domain)

	if b.port > 0 {
		if b.port > 65535 {
			return "", fmt.Errorf("port must be in range [1, 65535]")
		}
		rawBaseUrl = fmt.Sprintf("%s:%d", rawBaseUrl, b.port)
	}

	u, err := url.Parse(rawBaseUrl)
	if err != nil {
		return "", err
	}

	if b.credentials != nil {
		if b.credentials.user == "" {
			return "", fmt.Errorf("user not set")
		}

		if b.credentials.password == "" {
			return "", fmt.Errorf("password not set")
		}

		u.User = url.UserPassword(b.credentials.user, b.credentials.password)
	}

	if len(b.path) > 0 {
		u = u.JoinPath(b.path...)
	}

	for k, v := range b.query {
		if k == "" {
			return "", fmt.Errorf("query key is empty")
		}
		if i := slices.Index(v, ""); i != -1 {
			return "", fmt.Errorf("query query value for key %s", k)
		}
	}

	u.RawQuery = url.Values(b.query).Encode()

	if b.anchor != "" {
		u.Fragment = b.anchor
	}

	return u.String(), nil
}
