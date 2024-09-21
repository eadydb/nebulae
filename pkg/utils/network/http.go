package network

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type HttpClinet struct {
	Ctx  context.Context
	Host string
	Uri  string
}

func NewHttpClient(host, uri string) *HttpClinet {
	return &HttpClinet{
		Host: host,
		Uri:  uri,
	}
}

func (h *HttpClinet) Delete(ctx context.Context, host, uri, cookie string, params map[string]string) error {
	return h.Do(cookie, params, http.MethodDelete)
}

func (h *HttpClinet) Post(ctx context.Context, host, uri, cookie string, params map[string]string) error {
	return h.Do(cookie, params, http.MethodPost)
}

func (h *HttpClinet) PostJson(ctx context.Context, host, uri, cookie string, params map[string]string) error {
	return h.DoJson(cookie, params, http.MethodPost)
}

func (h *HttpClinet) Put(ctx context.Context, host, uri, cookie string, params map[string]string) error {
	return h.Do(cookie, params, http.MethodPut)
}

func (h *HttpClinet) Do(cookie string, params map[string]string, model string) error {
	param := url.Values{}
	for k, v := range params {
		param.Set(k, v)
	}
	u, _ := url.ParseRequestURI(h.Host)
	u.Path = h.Uri
	u.RawQuery = param.Encode()
	addr := u.Scheme + "://" + u.Host + "/" + u.Path + "?" + u.RawQuery

	req, err := http.NewRequest(model, addr, nil)
	if err != nil {
		slog.Error("http new request failed", slog.String("error", err.Error()), slog.String("addr", addr))
		return err
	}
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("http new request failed", slog.String("error", err.Error()), slog.String("addr", addr))
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (h *HttpClinet) DoJson(cookie string, params map[string]string, model string) error {
	u, _ := url.ParseRequestURI(h.Host)
	u.Path = h.Uri
	addr := u.Scheme + "://" + u.Host + "/" + u.Path
	body, _ := json.Marshal(params)
	req, err := http.NewRequest(model, addr, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		slog.Error("http new request failed", slog.String("error", err.Error()), slog.String("addr", addr))
		return err
	}
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("http new request failed", slog.String("error", err.Error()), slog.String("addr", addr))
		return err
	}
	defer resp.Body.Close()
	return nil
}
