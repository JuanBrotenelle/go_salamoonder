package salamoonder

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (Salamoonder, func()) {
	t.Helper()
	ts := httptest.NewServer(handler)
	c := New("test-api-key").(*client)
	c.setBaseURL(ts.URL)
	return c, ts.Close
}

func TestBalance_Success(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/createTask" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"error_code":0,"error_description":"","wallet":"123.45"}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	got, err := c.Balance()
	if err != nil {
		t.Fatalf("Balance() error: %v", err)
	}
	if got != 123.45 {
		t.Fatalf("Balance() = %v, want 123.45", got)
	}
}

func TestKasadaCreate_Success(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/createTask" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"error_code":0,"error_description":"","taskId":"task-1"}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	id, err := c.KasadaCreate("pjs-code", true)
	if err != nil {
		t.Fatalf("KasadaCreate() error: %v", err)
	}
	if id != "task-1" {
		t.Fatalf("KasadaCreate() = %q, want task-1", id)
	}
}

func TestKasada_SuccessReady(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getTaskResult" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"error_id":0,
			"status":"ready",
			"solution":{
				"user-agent":"UA",
				"x-kpsdk-cd":"cd",
				"x-kpsdk-cr":"cr",
				"x-kpsdk-ct":"ct",
				"x-kpsdk-r":"r",
				"x-kpsdk-v":"v",
				"x-kpsdk-st":"st"
			}
		}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	got, err := c.Kasada("task-1")
	if err != nil {
		t.Fatalf("Kasada() error: %v", err)
	}
	if got.UserAgent != "UA" || got.XCd != "cd" {
		t.Fatalf("Kasada() unexpected solution: %+v", got)
	}
}

func TestTwitchScraperCreate_Success(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/createTask" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"error_code":0,"error_description":"","taskId":"ts-1"}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	id, err := c.TwitchScraperCreate()
	if err != nil {
		t.Fatalf("TwitchScraperCreate() error: %v", err)
	}
	if id != "ts-1" {
		t.Fatalf("TwitchScraperCreate() = %q, want ts-1", id)
	}
}

func TestTwitchScraper_Success(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getTaskResult" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"error_id":0,
			"status":"ready",
			"solution":{
				"biography":"bio",
				"profile_picture":"pic",
				"username":"name"
			}
		}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	got, err := c.TwitchScraper("ts-1")
	if err != nil {
		t.Fatalf("TwitchScraper() error: %v", err)
	}
	if got.Username != "name" || got.ProfilePicture != "pic" {
		t.Fatalf("TwitchScraper() unexpected solution: %+v", got)
	}
}

func TestTwitchPI_CreateAndGet_Success(t *testing.T) {
	createHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/createTask" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"error_code":0,"error_description":"","taskId":"pi-1"}`))
			return
		}
		if r.URL.Path == "/getTaskResult" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"error_id":0,
				"status":"ready",
				"solution":{
					"device_id":"dev",
					"integrity_token":"token",
					"proxy":"prx",
					"user-agent":"UA",
					"client-id":"CID"
				}
			}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}

	c, closeFn := newTestClient(t, createHandler)
	defer closeFn()

	id, err := c.TwitchPICreate("acc", "prx")
	if err != nil {
		t.Fatalf("TwitchPICreate() error: %v", err)
	}
	if id != "pi-1" {
		t.Fatalf("TwitchPICreate() = %q, want pi-1", id)
	}

	got, err := c.TwitchPI(id)
	if err != nil {
		t.Fatalf("TwitchPI() error: %v", err)
	}
	if got.DeviceID != "dev" || got.ClientID != "CID" {
		t.Fatalf("TwitchPI() unexpected solution: %+v", got)
	}
}

func TestTwitchLI_CreateAndGet_Success(t *testing.T) {
	createHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/createTask" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"error_code":0,"error_description":"","taskId":"li-1"}`))
			return
		}
		if r.URL.Path == "/getTaskResult" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"error_id":0,
				"status":"ready",
				"solution":{
					"device_id":"dev2",
					"integrity_token":"token2",
					"proxy":"prx2",
					"user-agent":"UA2",
					"client-id":"CID2"
				}
			}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}

	c, closeFn := newTestClient(t, createHandler)
	defer closeFn()

	id, err := c.TwitchLICreate("prx")
	if err != nil {
		t.Fatalf("TwitchLICreate() error: %v", err)
	}
	if id != "li-1" {
		t.Fatalf("TwitchLICreate() = %q, want li-1", id)
	}

	got, err := c.TwitchLI(id)
	if err != nil {
		t.Fatalf("TwitchLI() error: %v", err)
	}
	if got.DeviceID != "dev2" || got.ClientID != "CID2" {
		t.Fatalf("TwitchLI() unexpected solution: %+v", got)
	}
}
