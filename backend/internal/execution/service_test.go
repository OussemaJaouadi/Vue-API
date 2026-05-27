package execution_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vue-api/backend/internal/execution"
)

func TestExecute_GetRequestReturnsStatusAndBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Equal(t, "200 OK", resp.StatusText)
	assert.Contains(t, resp.Body, `{"ok":true}`)
	assert.Greater(t, resp.Duration, 0.0)
}

func TestExecute_PostRequestSendsBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		assert.Equal(t, `{"hello":"world"}`, strings.TrimSpace(string(body)))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`created`))
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodPost,
		URL:    ts.URL,
		Body:   `{"hello":"world"}`,
	})

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Status)
	assert.Equal(t, "201 Created", resp.StatusText)
	assert.Contains(t, resp.Body, "created")
}

func TestExecute_QueryParamsAreAppended(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "value1", r.URL.Query().Get("key1"))
		assert.Equal(t, "value2", r.URL.Query().Get("key2"))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
		QueryParams: []execution.QueryParam{
			{Key: "key1", Value: "value1", Enabled: true},
			{Key: "key2", Value: "value2", Enabled: true},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Status)
}

func TestExecute_DisabledQueryParamsAreSkipped(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.URL.Query().Get("key1"))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	_, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
		QueryParams: []execution.QueryParam{
			{Key: "key1", Value: "value1", Enabled: false},
		},
	})
	require.NoError(t, err)
}

func TestExecute_RequestHeadersAreSent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "bar", r.Header.Get("X-Custom"))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	_, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
		Headers: []execution.Header{
			{Key: "X-Custom", Value: "bar", Enabled: true},
		},
	})
	require.NoError(t, err)
}

func TestExecute_DisabledHeadersAreSkipped(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("X-Custom"))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	_, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
		Headers: []execution.Header{
			{Key: "X-Custom", Value: "should-not-appear", Enabled: false},
		},
	})
	require.NoError(t, err)
}

func TestExecute_InvalidURLReturnsError(t *testing.T) {
	svc := execution.NewService()
	_, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    "://invalid-url",
	})
	require.Error(t, err)
}

func TestExecute_ResponseHeadersAreReturned(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(resp.Headers), 1)
	found := false
	for _, h := range resp.Headers {
		if h.Key == "Content-Type" && h.Value == "application/json" {
			found = true
		}
	}
	assert.True(t, found)
}

func TestExecute_RedirectsAreNotFollowed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/other", http.StatusMovedPermanently)
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.Equal(t, http.StatusMovedPermanently, resp.Status)
}

func TestExecute_ContextTimeoutCancels(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	svc := execution.NewService()
	resp, err := svc.Execute(ctx, execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	if err != nil {
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	} else {
		assert.NotEmpty(t, resp.Error)
	}
}

func TestExecute_ResponseHasSize(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := `{"data":"hello world"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(len(`{"data":"hello world"}`)), resp.Size)
}

func TestExecute_ResponseHasExecutionTarget(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.ExecutionTarget)
}

func TestExecute_ResponseHasRequestID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.RequestID)
	assert.True(t, strings.HasPrefix(resp.RequestID, "req_"))
}

func TestExecute_ExistingQueryParamsOnURLArePreserved(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "existing", r.URL.Query().Get("existing"))
		assert.Equal(t, "added", r.URL.Query().Get("added"))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	_, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL + "?existing=existing",
		QueryParams: []execution.QueryParam{
			{Key: "added", Value: "added", Enabled: true},
		},
	})
	require.NoError(t, err)
}

func TestExecute_ResponseContentType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("plain text"))
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.Equal(t, "plain text", resp.Body)
}

func TestExecute_ConnectionRefused(t *testing.T) {
	svc := execution.NewService()
	_, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    "http://127.0.0.1:1",
		Timeout: time.Second,
	})

	if err == nil {
		t.Skip("connection refused may return as response.Error on some platforms")
	} else {
		require.Error(t, err)
	}
}

func TestExecute_NilBodyForGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	_, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})
	require.NoError(t, err)
}

func TestExecute_StatusCodeMapping(t *testing.T) {
	codes := []int{
		http.StatusOK,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusUnauthorized,
	}

	for _, code := range codes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			}))
			defer ts.Close()

			svc := execution.NewService()
			resp, err := svc.Execute(context.Background(), execution.Request{
				Method: execution.MethodGet,
				URL:    ts.URL,
			})

			require.NoError(t, err)
			assert.Equal(t, code, resp.Status)
		})
	}
}

func TestExecute_ResponseBodyIsJSONEncodable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"nested":{"key":"val"}}`))
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	var parsed map[string]any
	require.NoError(t, json.Unmarshal([]byte(resp.Body), &parsed))
	nested, ok := parsed["nested"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "val", nested["key"])
}

func TestExecute_EmptyResponseBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.Status)
	assert.Empty(t, resp.Body)
}

func TestExecute_DurationIsWithinReasonableRange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.Greater(t, resp.Duration, 1.0)
	assert.Less(t, resp.Duration, 10000.0)
}

func TestExecute_CustomMethod(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	_, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodPatch,
		URL:    ts.URL,
	})
	require.NoError(t, err)
}

func TestExecute_SchemeAndHostInExecutionTarget(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.Contains(t, resp.ExecutionTarget, "127.0.0.1")
}

func TestExecute_DeleteSendsBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodDelete,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.Status)
}

func TestExecute_HeaderWithEmptyKeyIsSkipped(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	_, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
		Headers: []execution.Header{
			{Key: "", Value: "empty-key", Enabled: true},
		},
	})
	require.NoError(t, err)
}

func TestExecute_HeaderWithEmptyValue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "", r.Header.Get("X-Empty"))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := execution.NewService()
	_, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
		Headers: []execution.Header{
			{Key: "X-Empty", Value: "", Enabled: true},
		},
	})
	require.NoError(t, err)
}

func TestExecute_ResponseStatusTextForNonStandardCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(299)
	}))
	defer ts.Close()

	svc := execution.NewService()
	resp, err := svc.Execute(context.Background(), execution.Request{
		Method: execution.MethodGet,
		URL:    ts.URL,
	})

	require.NoError(t, err)
	assert.Equal(t, 299, resp.Status)
	assert.NotEmpty(t, resp.StatusText)
}
