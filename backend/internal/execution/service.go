package execution

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

type httpService struct {
	client *http.Client
}

func NewService() Service {
	return &httpService{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (s *httpService) Execute(ctx context.Context, req Request) (Response, error) {
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return Response{}, err
	}

	query := parsedURL.Query()
	for _, p := range req.QueryParams {
		if p.Enabled {
			query.Add(p.Key, p.Value)
		}
	}
	parsedURL.RawQuery = query.Encode()

	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, string(req.Method), parsedURL.String(), bodyReader)
	if err != nil {
		return Response{}, err
	}

	for _, h := range req.Headers {
		if h.Enabled && h.Key != "" {
			httpReq.Header.Add(h.Key, h.Value)
		}
	}

	var start time.Time
	var ttfb time.Duration

	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			ttfb = time.Since(start)
		},
	}
	httpReq = httpReq.WithContext(httptrace.WithClientTrace(httpReq.Context(), trace))

	start = time.Now()
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{
			Status:     resp.StatusCode,
			StatusText: resp.Status,
			Error:      err.Error(),
		}, nil
	}
	duration := time.Since(start)

	headers := make([]Header, 0, len(resp.Header))
	for k, v := range resp.Header {
		headers = append(headers, Header{
			Key:     k,
			Value:   strings.Join(v, ", "),
			Enabled: true,
		})
	}

	return Response{
		Status:          resp.StatusCode,
		StatusText:      resp.Status,
		Headers:         headers,
		Body:            string(bodyBytes),
		Duration:        float64(duration.Nanoseconds()) / 1e6,
		Size:            int64(len(bodyBytes)),
		TTFB:            float64(ttfb.Nanoseconds()) / 1e6,
		ExecutionTarget: parsedURL.Host,
		RequestID:       "req_" + strings.ToLower(strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")),
	}, nil
}
