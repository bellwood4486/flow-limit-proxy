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

func ListenProxy(portFrom, portTo uint, concurrentLimit int64) error {
	proxy, err := newReverseProxy(portTo, concurrentLimit)
	if err != nil {
		return fmt.Errorf("failed to new proxy: %w", err)
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", portFrom),
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

	log.Printf("start proxy...(limit:%d)", concurrentLimit)
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
	ctx := context.Background()
	if err := t.sem.Acquire(ctx, 1); err != nil {
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
