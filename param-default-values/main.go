package main

import "fmt"

type ServerConfig struct {
	Port    int
	Timeout int
}

type ServerOptionFunc func(cfg *ServerConfig)

func WithPort(port int) ServerOptionFunc {
	return func(cfg *ServerConfig) {
		cfg.Port = port
	}
}

func WithTimeout(timeout int) ServerOptionFunc {
	return func(cfg *ServerConfig) {
		cfg.Timeout = timeout
	}
}

func NewServer(opts ...ServerOptionFunc) *ServerConfig {
	cfg := &ServerConfig{
		Port:    8081,
		Timeout: 120,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func main() {
	s1 := NewServer()
	fmt.Printf("S1 Server Port: %d, Timeout: %d\n", s1.Port, s1.Timeout)

	s2 := NewServer(WithPort(8082), WithTimeout(600))
	fmt.Printf("S2 Server Port: %d, Timeout: %d\n", s2.Port, s2.Timeout)
}
