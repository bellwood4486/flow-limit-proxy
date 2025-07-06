package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/sync/semaphore"
)

// Config holds the proxy configuration
type Config struct {
	FromPort   uint  // Source port to listen on (1-65535)
	ToPort     uint  // Target port to forward requests to (1-65535)
	MaxConns   int64 // Maximum number of concurrent connections
}

// NewConfig creates a new Config with validation
func NewConfig(fromPort, toPort int, limit int64) (*Config, error) {
	if err := validatePort(fromPort); err != nil {
		return nil, fmt.Errorf("invalid fromPort: %w", err)
	}
	
	if err := validatePort(toPort); err != nil {
		return nil, fmt.Errorf("invalid toPort: %w", err)
	}
	
	return &Config{
		FromPort: uint(fromPort),
		ToPort:   uint(toPort),
		MaxConns: limit,
	}, nil
}

// validatePort validates that a port number is within the valid range
func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", port)
	}
	return nil
}

func ListenProxy(config *Config) error {
	proxy, err := newReverseProxy(config.ToPort, config.MaxConns)
	if err != nil {
		return fmt.Errorf("failed to new proxy: %w", err)
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.FromPort),
		Handler: proxy,
	}

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	go func() {
		<-quit
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("failed to gracefully shutdown: %v\n", err)
		}
	}()

	log.Printf("start proxy...(limit:%d)", config.MaxConns)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to ListenAndServ: %w", err)
	}
	log.Printf("shutdown")

	return nil
}

func newReverseProxy(targetPort uint, limit int64) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(fmt.Sprintf("http://localhost:%d", targetPort))
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = newCustomTransport(limit)
	proxy.ErrorHandler = func(_ http.ResponseWriter, r *http.Request, err error) {
		log.Printf("fail request: %s %s: %v", r.Method, r.URL, err)
	}

	return proxy, nil
}

// customTransport は、プロキシするHTTP通信を制御するための構造体です。
// 以下の機能を持ちます。
// - 同時通信数の制御
// - 通信エラー時のリトライ
type customTransport struct {
	base http.RoundTripper
	sem  *semaphore.Weighted
}

func newCustomTransport(concurrentLimit int64) http.RoundTripper {
	return &customTransport{
		base: http.DefaultTransport,
		sem:  semaphore.NewWeighted(concurrentLimit),
	}
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 同時通信数の制御
	if err := t.sem.Acquire(req.Context(), 1); err != nil {
		return nil, fmt.Errorf("failed to acquire semaphore: %v", err)
	}
	defer t.sem.Release(1)

	// 指数バックオフしながらリクエストを送る
	var res *http.Response
	tryCount := 0
	err := backoff.Retry(func() error {
		tryCount++
		var err error
		res, err = t.base.RoundTrip(req)
		// エラーのときだけリトライ。errがnilでステータスコード500は成功とみなす。
		if err != nil {
			log.Printf("retry%d: %s %s", tryCount, req.Method, req.URL)
			return err
		}
		return nil
	}, newBackOffConfig())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func newBackOffConfig() *backoff.ExponentialBackOff {
	config := backoff.NewExponentialBackOff()
	config.MaxInterval = 3 * time.Second
	config.MaxElapsedTime = 10 * time.Second
	return config
}
