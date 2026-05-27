package execution

import (
	"context"
	"time"
)

type Method string

const (
	MethodGet    Method = "GET"
	MethodPost   Method = "POST"
	MethodPut    Method = "PUT"
	MethodPatch  Method = "PATCH"
	MethodDelete Method = "DELETE"
)

type Header struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type QueryParam struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type Request struct {
	Method      Method       `json:"method"`
	URL         string       `json:"url"`
	Headers     []Header     `json:"headers"`
	QueryParams []QueryParam `json:"queryParams"`
	Body        string       `json:"body"`
	Timeout     time.Duration `json:"timeout"`
}

type Response struct {
	Status          int           `json:"status"`
	StatusText      string        `json:"statusText"`
	Headers         []Header      `json:"headers"`
	Body            string        `json:"body"`
	Duration        float64       `json:"duration"` // Milliseconds
	Size            int64         `json:"size"`
	TTFB            float64       `json:"ttfb"` // Milliseconds
	Error           string        `json:"error,omitempty"`
	ExecutionTarget string        `json:"executionTarget"`
	RequestID       string        `json:"requestId"`
}

type Service interface {
	Execute(ctx context.Context, req Request) (Response, error)
}
