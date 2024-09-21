package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

func Get[T any](host, uri, cookie string, params map[string]string, t []T) ([]T, error) {
	param := url.Values{}
	for k, v := range params {
		param.Set(k, v)
	}
	u, _ := url.ParseRequestURI(host)
	u.Path = uri
	u.RawQuery = param.Encode()
	addr := u.Scheme + "://" + u.Host + "/" + u.Path + "?" + u.RawQuery

	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		slog.Error("http new request failed", slog.String("error", err.Error()), slog.String("addr", addr))
		return nil, err
	}
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("http new request failed", slog.String("error", err.Error()), slog.String("addr", addr))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, nil
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&t); err != nil {
		slog.Error("http response json decode failed", slog.String("error", err.Error()), slog.String("addr", addr))
		return nil, err
	}
	return t, nil
}

func GetOne[T any](host, uri, cookie string, params map[string]string, t T) (T, error) {
	param := url.Values{}
	for k, v := range params {
		param.Set(k, v)
	}
	u, _ := url.ParseRequestURI(host)
	u.Path = uri
	u.RawQuery = param.Encode()
	addr := u.Scheme + "://" + u.Host + "/" + u.Path + "?" + u.RawQuery

	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		slog.Error("http new request failed", slog.String("error", err.Error()), slog.String("addr", addr))
		return t, err
	}
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("http new request failed", slog.String("error", err.Error()), slog.String("addr", addr))
		return t, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return t, errors.New(fmt.Sprintf("http response status is %d", resp.StatusCode))
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&t); err != nil {
		slog.Error("http response json decode failed", slog.String("error", err.Error()))
		return t, err
	}
	return t, nil
}
