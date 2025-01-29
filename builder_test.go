package url_builder

import (
	"fmt"
	"strings"
	"testing"
)

func TestBuilder_Build(t *testing.T) {
	tests := []struct {
		name    string
		builder *Builder
		want    string
		wantErr bool
	}{
		{
			name:    "Domain",
			builder: New().WithDomain("example.com"),
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name:    "Empty domain",
			builder: New().WithDomain(""),
			want:    "",
			wantErr: true,
		},
		{
			name:    "Missing domain",
			builder: New().WithScheme("https"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "Scheme",
			builder: New().WithScheme("https").WithDomain("example.com"),
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name:    "Extra scheme",
			builder: New().WithScheme("https").WithDomain("http://example.com"),
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name:    "Port",
			builder: New().WithDomain("example.com").WithPort(8080),
			want:    "http://example.com:8080",
			wantErr: false,
		},
		{
			name:    "Extra port",
			builder: New().WithDomain("example.com:80").WithPort(8080),
			want:    "http://example.com:8080",
			wantErr: false,
		},
		{
			name:    "Invalid port",
			builder: New().WithDomain("example.com").WithPort(70000),
			want:    "",
			wantErr: true,
		},
		{
			name:    "Extra scheme and port",
			builder: New().WithDomain("ftp://example.com:80").WithPort(8080),
			want:    "http://example.com:8080",
			wantErr: false,
		},
		{
			name:    "Credentials",
			builder: New().WithDomain("example.com").WithCredentials("user", "pass"),
			want:    "http://user:pass@example.com",
			wantErr: false,
		},
		{
			name:    "Empty credentials",
			builder: New().WithDomain("example.com").WithCredentials("", ""),
			want:    "",
			wantErr: true,
		},
		{
			name:    "Path",
			builder: New().WithDomain("example.com").WithPath("users", "123"),
			want:    "http://example.com/users/123",
			wantErr: false,
		},
		{
			name:    "Empty path",
			builder: New().WithDomain("example.com").WithPath("", ""),
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name:    "Multiple path calls",
			builder: New().WithDomain("example.com").WithPath("users").WithPath("123"),
			want:    "http://example.com/users/123",
			wantErr: false,
		},
		{
			name:    "Escape characters in path",
			builder: New().WithDomain("example.com").WithPath("hello world"),
			want:    "http://example.com/hello%20world",
			wantErr: false,
		},
		{
			name:    "Long path",
			builder: New().WithDomain("example.com").WithPath(strings.Repeat("a", 1000)),
			want:    fmt.Sprintf("http://example.com/%s", strings.Repeat("a", 1000)),
			wantErr: false,
		},
		{
			name:    "Query params",
			builder: New().WithDomain("example.com").WithQuery("key", "value"),
			want:    "http://example.com?key=value",
			wantErr: false,
		},
		{
			name:    "Single query key and multiple values in a single call",
			builder: New().WithDomain("example.com").WithQuery("key", "val1", "val2"),
			want:    "http://example.com?key=val1&key=val2",
			wantErr: false,
		},
		{
			name:    "Single query key and multiple query values in multiple calls",
			builder: New().WithDomain("example.com").WithQuery("key", "val1").WithQuery("key", "val2"),
			want:    "http://example.com?key=val1&key=val2",
			wantErr: false,
		},
		{
			name:    "Multiple query key and multiple query values",
			builder: New().WithDomain("example.com").WithQuery("key1", "val1").WithQuery("key2", "val2"),
			want:    "http://example.com?key1=val1&key2=val2",
			wantErr: false,
		},
		{
			name:    "Escape characters in query",
			builder: New().WithDomain("example.com").WithQuery("key", "a b"),
			want:    "http://example.com?key=a+b",
			wantErr: false,
		},
		{
			name:    "Empty query key",
			builder: New().WithDomain("example.com").WithQuery("", "val"),
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name:    "Empty query value",
			builder: New().WithDomain("example.com").WithQuery("key", ""),
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name:    "Anchor",
			builder: New().WithDomain("example.com").WithAnchor("Anchor"),
			want:    "http://example.com#Anchor",
			wantErr: false,
		},
		{
			name:    "Empty anchor",
			builder: New().WithDomain("example.com").WithAnchor(""),
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name: "Combined",
			builder: New().
				WithScheme("https://").
				WithCredentials("user", "pass").
				WithDomain("test.example.com").
				WithPort(1234).
				WithPath("path1", "path2").
				WithQuery("key1", "val1").
				WithQuery("key2", "val2").
				WithAnchor("#Anchor"),
			want:    "https://user:pass@test.example.com:1234/path1/path2?key1=val1&key2=val2#Anchor",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}
