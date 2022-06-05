package enhancers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
	"golang.org/x/net/html"
)

// https://oembed.com/#section4
type OEmbedStrategy struct {
	responseCache sync.Map
	once          sync.Once
}

func (s *OEmbedStrategy) Name() string {
	return "oembed"
}

func (s *OEmbedStrategy) Accepts(u *url.URL) bool {
	resp, err := s.getBody(u)
	if err != nil {
		return false
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	respStr := string(respBody)
	return strings.Contains(respStr, "application/json+oembed")
}

func (s *OEmbedStrategy) Provide(u *url.URL) (*EnhancedLink, error) {
	resp, err := s.getBody(u)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	tk := html.NewTokenizer(resp.Body)
	for {
		ntk := tk.Next()
		switch {
		case ntk == html.ErrorToken:
			goto done
		case ntk == html.StartTagToken:
			t := tk.Token()
			if t.Data != "link" {
				continue
			}

			data := getOEmbedData(t.Attr)
			if data != nil {
				return &EnhancedLink{
					Link:  u.String(),
					Title: data.Title,
				}, nil
			}
		}
	}

done:
	return nil, fmt.Errorf("no data")
}

func (s *OEmbedStrategy) cacheCleanup() {
	var keysToDelete []string

	// Collect keys to delete
	s.responseCache.Range(func(key, value any) bool {
		if !strings.HasPrefix(key.(string), "expiry.") {
			return true
		}

		expiredAt := value.(time.Time)
		logging.V(1).Info("expiry check", "key", key, "expiry", expiredAt)
		if expiredAt.Before(time.Now()) {
			keysToDelete = append(keysToDelete, key.(string))
			keysToDelete = append(keysToDelete, strings.TrimPrefix(key.(string), "expiry."))
		}

		return true
	})

	// Use LoadAndDelete
	for _, k := range keysToDelete {
		s.responseCache.LoadAndDelete(k)
		logging.V(1).Info("deleted key", "key", k)
	}
}

func (s *OEmbedStrategy) getBody(u *url.URL) (*http.Response, error) {
	s.once.Do(func() {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			for range ticker.C {
				s.cacheCleanup()
			}
		}()
	})

	key := u.String()
	urlLogger := logging.WithValues("key", key)
	v, ok := s.responseCache.Load(key)
	if ok {
		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}

		cachedReader := bufio.NewReader(bytes.NewReader(v.([]byte)))
		cachedResponse, err := http.ReadResponse(cachedReader, req)
		if err == nil {
			return cachedResponse, err
		} else {
			urlLogger.Error(err, "failed to load cached response", "cachedResp", string(v.([]byte)))
		}
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(30 * time.Second)
	s.responseCache.Store(key, respBytes)
	s.responseCache.Store(fmt.Sprintf("expiry.%s", key), expiresAt)

	newResponse, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(respBytes)), req)
	if err != nil {
		return nil, err
	}

	resp.Body = io.NopCloser(newResponse.Body)
	logging.Info("storing into cache", "key", key, "expiry", expiresAt, "response_bytes", len(respBytes))
	return resp, err
}

type oembedData struct {
	Type  string
	Title string
}

func getOEmbedData(attrs []html.Attribute) *oembedData {
	var data oembedData
	for _, attr := range attrs {
		if attr.Key == "type" && attr.Val == "application/json+oembed" {
			data.Type = attr.Val
		}

		if attr.Key == "title" {
			data.Title = attr.Val
		}
	}

	if data.Type == "" || data.Title == "" {
		return nil
	}

	return &data
}
