package salamoonder

import (
	"testing"
)

func TestFindPJS(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantErr   bool
		expectMsg string
	}{
		{
			name:    "nike.com has p.js",
			url:     "https://www.nike.com/",
			wantErr: false,
		},
		{
			name:      "wakatime no p.js",
			url:       "https://wakatime.com/",
			wantErr:   true,
			expectMsg: "p.js script src not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindPJS(tt.url)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (url=%s)", tt.url)
				}
				if tt.expectMsg != "" && err.Error() != tt.expectMsg {
					t.Errorf("unexpected error message, got %q, want %q", err.Error(), tt.expectMsg)
				}
				if got != "" {
					t.Errorf("expected empty result, got %q", got)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v (url=%s)", err, tt.url)
				}
				if got == "" {
					t.Errorf("expected non-empty result, got empty (url=%s)", tt.url)
				} else {
					t.Logf("found p.js: %s", got)
				}
			}
		})
	}
}
