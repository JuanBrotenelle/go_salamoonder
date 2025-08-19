package salamoonder

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

func boolToString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

// FindPJS fetches the page at pageURL and returns the first script src that ends with "/p.js".
// It returns an error if the request fails or no such script tag is found.
func FindPJS(pageURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	pattern := `(?i)<script[^>]*\bsrc=["']([^"'>]+/p\.js(?:\?[^"'>]*)?)["'][^>]*>`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) >= 2 {
		return matches[1], nil
	}
	return "", fmt.Errorf("p.js script src not found")
}
