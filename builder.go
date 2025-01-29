package url_builder

import (
	"fmt"
	"net/url"
	"strings"
)

// Builder constructs URLs step by step. Domain name is the only required value.
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

// WithDomain sets the domain name of the URL. This part is required.
func (b *Builder) WithDomain(domain string) *Builder {
	b.domain = domain
	return b
}

// WithPort sets the port number. Must be in the range [1, 65535].
func (b *Builder) WithPort(port int) *Builder {
	b.port = port
	return b
}

// WithCredentials sets the username and password for authentication.
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
// If the same key is added multiple times, values are appended rather than replaced.
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

	rawUrl := b.domain
	if !strings.Contains(b.domain, "//") {
		rawUrl = fmt.Sprintf("%s://%s", b.scheme, b.domain)
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	u.Host = u.Hostname()
	u.Scheme = b.scheme

	if b.port > 0 {
		if b.port > 65535 {
			return "", fmt.Errorf("port must be in range [1, 65535]")
		}

		u.Host = fmt.Sprintf("%s:%d", u.Host, b.port)
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

	qVal := url.Values{}
	for k, v := range b.query {
		if k == "" || len(v) == 0 {
			continue
		}
		for i := range v {
			if v[i] == "" {
				continue
			}
			qVal.Add(k, v[i])
		}
	}

	if len(qVal) > 0 {
		u.RawQuery = qVal.Encode()
	}

	if b.anchor != "" {
		u.Fragment = b.anchor
	}

	return u.String(), nil
}
