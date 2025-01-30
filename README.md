# url-builder

![Go Version](https://img.shields.io/badge/Go-1.23-blue)
![License](https://img.shields.io/badge/license-MIT-green)

A simple and flexible Go library for constructing URLs in a structured way.

## Features

- Supports custom schemes, domains, ports, paths, query parameters, credentials, and anchors
- Ensures correct URL encoding
- Chainable API for easy usage

## Installation

```sh
go get github.com/smilingthrone13/url-builder
```

## Usage

### Creating a Basic URL

```go
url, err := url_builder.New().WithDomain("example.com").Build()
// url: "http://example.com"
```

### Creating a Basic URL from IPv4 Address

```go
url, err := url_builder.New().WithIPv4("192.168.1.1").Build()
// url: "http://192.168.1.1"
```

### Using HTTPS

```go
url, err := url_builder.New().WithSchemeHTTPS().WithDomain("example.com").Build()
// url: "https://example.com"
```

### Adding a Port

```go
url, err := url_builder.New().WithDomain("example.com").WithPort(8080).Build()
// url: "http://example.com:8080"
```

### Adding Credentials

```go
url, err := url_builder.New().WithDomain("example.com").WithCredentials("user", "pass").Build()
// url: "http://user:pass@example.com"
```

### Adding a Path

```go
url, err := url_builder.New().WithDomain("example.com").WithPath("users", "123").Build()
// url: "http://example.com/users/123"
```

### Adding Query Parameters

```go
url, err := url_builder.New().WithDomain("example.com").WithQuery("key", "value").Build()
// url: "http://example.com?key=value"
```

### Adding an Anchor

```go
url, err := url_builder.New().WithDomain("example.com").WithAnchor("section").Build()
// url: "http://example.com#section"
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
