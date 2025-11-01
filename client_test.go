package salamoonder

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, func()) {
	t.Helper()
	ts := httptest.NewServer(handler)
	c, err := New("test-api-key", nil)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	c.client.setBaseURL(ts.URL)
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

	got, err := c.Balance(context.Background())
	if err != nil {
		t.Fatalf("Balance() error: %v", err)
	}
	if got.Wallet != "123.45" {
		t.Fatalf("Balance() = %v, want 123.45", got.Wallet)
	}
}

func TestCreateTask_Kasada_Success(t *testing.T) {
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

	result, err := c.CreateTask(context.Background(), KasadaOptions{
		Pjs:    "pjs-code",
		CdOnly: true,
	})
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}
	if result.TaskId != "task-1" {
		t.Fatalf("CreateTask() = %q, want task-1", result.TaskId)
	}
}

func TestGetTaskResult_Kasada_Success(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getTaskResult" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"errorId":0,
			"status":"ready",
			"solution":{
				"user-agent":"UA",
				"x-is-human":"human",
				"x-kpsdk-cd":"cd",
				"x-kpsdk-cr":"cr",
				"x-kpsdk-ct":"ct",
				"x-kpsdk-r":"r",
				"x-kpsdk-st":"st"
			}
		}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	got, err := GetTaskResult[KasadaSolution](c, context.Background(), "task-1")
	if err != nil {
		t.Fatalf("GetTaskResult() error: %v", err)
	}
	if got.Solution.UserAgent != "UA" || got.Solution.XKpsdkCd != "cd" {
		t.Fatalf("GetTaskResult() unexpected solution: %+v", got.Solution)
	}
}

func TestTask_Raw_Success(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getTaskResult" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"errorId":0,
			"status":"ready",
			"solution":{
				"user-agent":"UA",
				"x-kpsdk-cd":"cd"
			}
		}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	got, err := c.Task(context.Background(), "task-1")
	if err != nil {
		t.Fatalf("Task() error: %v", err)
	}
	if got.Status != "ready" {
		t.Fatalf("Task() status = %q, want ready", got.Status)
	}

	var solution KasadaSolution
	if err := json.Unmarshal(got.Solution, &solution); err != nil {
		t.Fatalf("Failed to unmarshal solution: %v", err)
	}
	if solution.UserAgent != "UA" {
		t.Fatalf("Task() solution.UserAgent = %q, want UA", solution.UserAgent)
	}
}

func TestCreateTask_TwitchScraper_Success(t *testing.T) {
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

	result, err := c.CreateTask(context.Background(), TwitchScraperOptions{})
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}
	if result.TaskId != "ts-1" {
		t.Fatalf("CreateTask() = %q, want ts-1", result.TaskId)
	}
}

func TestGetTaskResult_TwitchScraper_Success(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getTaskResult" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"errorId":0,
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

	got, err := GetTaskResult[TwitchScraperSolution](c, context.Background(), "ts-1")
	if err != nil {
		t.Fatalf("GetTaskResult() error: %v", err)
	}
	if got.Solution.Username != "name" || got.Solution.ProfilePicture != "pic" {
		t.Fatalf("GetTaskResult() unexpected solution: %+v", got.Solution)
	}
}

func TestCreateTask_TwitchPublicIntegrity_Success(t *testing.T) {
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
				"errorId":0,
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

	result, err := c.CreateTask(context.Background(), TwitchPublicIntegrityOptions{
		Proxy:       "prx",
		AccessToken: "acc",
		DeviceID:    "dev",
		ClientID:    "cid",
	})
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}
	if result.TaskId != "pi-1" {
		t.Fatalf("CreateTask() = %q, want pi-1", result.TaskId)
	}

	got, err := GetTaskResult[TwitchPublicIntegritySolution](c, context.Background(), result.TaskId)
	if err != nil {
		t.Fatalf("GetTaskResult() error: %v", err)
	}
	if got.Solution.DeviceID != "dev" || got.Solution.ClientID != "CID" {
		t.Fatalf("GetTaskResult() unexpected solution: %+v", got.Solution)
	}
}

func TestCreateTask_TwitchLocalIntegrity_Success(t *testing.T) {
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
				"errorId":0,
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

	result, err := c.CreateTask(context.Background(), TwitchLocalIntegrityOptions{
		Proxy:    "prx",
		DeviceID: "dev",
		ClientID: "cid",
	})
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}
	if result.TaskId != "li-1" {
		t.Fatalf("CreateTask() = %q, want li-1", result.TaskId)
	}

	got, err := GetTaskResult[TwitchLocalIntegritySolution](c, context.Background(), result.TaskId)
	if err != nil {
		t.Fatalf("GetTaskResult() error: %v", err)
	}
	if got.Solution.DeviceID != "dev2" || got.Solution.ClientID != "CID2" {
		t.Fatalf("GetTaskResult() unexpected solution: %+v", got.Solution)
	}
}

func TestCreateTask_Reese84_Success(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/createTask" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"error_code":0,"error_description":"","taskId":"r84-1"}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	result, err := c.CreateTask(context.Background(), Reese84Options{
		Website:       "https://example.com",
		SubmitPayload: true,
	})
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}
	if result.TaskId != "r84-1" {
		t.Fatalf("CreateTask() = %q, want r84-1", result.TaskId)
	}
}

func TestCreateTask_Uutmvc_Success(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/createTask" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"error_code":0,"error_description":"","taskId":"utmvc-1"}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	result, err := c.CreateTask(context.Background(), UutmvcOptions{
		Website: "https://example.com",
	})
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}
	if result.TaskId != "utmvc-1" {
		t.Fatalf("CreateTask() = %q, want utmvc-1", result.TaskId)
	}
}

func TestCreateTask_UnsupportedType(t *testing.T) {
	c, _ := New("test-api-key", nil)

	_, err := c.CreateTask(context.Background(), "invalid")
	if err == nil {
		t.Fatal("CreateTask() expected error, got nil")
	}
	if err != ErrUnsupportedTaskOptionsType {
		t.Fatalf("CreateTask() error = %v, want ErrUnsupportedTaskOptionsType", err)
	}
}
