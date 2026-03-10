package salamoonder

import (
	"context"
	"encoding/json"
	"errors"
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

// ------------------------------------------------------------------
// New
// ------------------------------------------------------------------

func TestNew_NoApiKey(t *testing.T) {
	_, err := New("", nil)
	if err == nil {
		t.Fatal("New() expected error for empty key, got nil")
	}
	if !errors.Is(err, ErrNoApiKey) {
		t.Fatalf("errors.Is(err, ErrNoApiKey) = false, want true; err = %v", err)
	}
}

// ------------------------------------------------------------------
// Balance
// ------------------------------------------------------------------

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

// HTTP 200, but ErrorCode == 1
func TestBalance_APIError(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"error_code":1,"error_description":"insufficient funds","wallet":""}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	_, err := c.Balance(context.Background())
	if err == nil {
		t.Fatal("Balance() expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Balance() error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusOK {
		t.Errorf("APIError.StatusCode = %d, want %d", apiErr.StatusCode, http.StatusOK)
	}
	if apiErr.Msg == "" {
		t.Error("APIError.Msg is empty, want non-empty description")
	}
}

// ------------------------------------------------------------------
// CreateTask
// ------------------------------------------------------------------

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

	result, err := c.CreateTask(context.Background(), KasadaStandardOptions{
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

// HTTP 400
func TestCreateTask_HTTP400(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error_code":1,"error_description":"invalid task type"}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	_, err := c.CreateTask(context.Background(), KasadaStandardOptions{Pjs: "x"})
	if err == nil {
		t.Fatal("CreateTask() expected error on 400, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("CreateTask() error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("APIError.StatusCode = %d, want %d", apiErr.StatusCode, http.StatusBadRequest)
	}
	if apiErr.Msg == "" {
		t.Error("APIError.Msg is empty, want non-empty description")
	}
}

// HTTP 200, ErrorCode = 1
func TestCreateTask_APIError(t *testing.T) {
	const wantTaskId = "task-err-1"
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"error_code":1,"error_description":"solver unavailable","taskId":"` + wantTaskId + `"}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	_, err := c.CreateTask(context.Background(), KasadaStandardOptions{Pjs: "x"})
	if err == nil {
		t.Fatal("CreateTask() expected error on ErrorCode=1, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("CreateTask() error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusOK {
		t.Errorf("APIError.StatusCode = %d, want %d", apiErr.StatusCode, http.StatusOK)
	}
	if apiErr.TaskId != wantTaskId {
		t.Errorf("APIError.TaskId = %q, want %q", apiErr.TaskId, wantTaskId)
	}
	if apiErr.Msg == "" {
		t.Error("APIError.Msg is empty, want non-empty description")
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

func TestCreateTask_TwitchIntegrity_Success(t *testing.T) {
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

	result, err := c.CreateTask(context.Background(), TwitchIntegrityOptions{
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

	got, err := GetTaskResult[TwitchIntegritySolution](c, context.Background(), result.TaskId)
	if err != nil {
		t.Fatalf("GetTaskResult() error: %v", err)
	}
	if got.Solution.DeviceID != "dev" || got.Solution.ClientID != "CID" {
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

// Unsupported type:
//   - errors.Is(err, ErrUnsupportedTaskOptionsType)
//   - errors.As(err, &MethodError{})
func TestCreateTask_UnsupportedType(t *testing.T) {
	c, _ := New("test-api-key", nil)

	_, err := c.CreateTask(context.Background(), "invalid")
	if err == nil {
		t.Fatal("CreateTask() expected error, got nil")
	}

	if !errors.Is(err, ErrUnsupportedTaskOptionsType) {
		t.Fatalf("errors.Is(err, ErrUnsupportedTaskOptionsType) = false, want true; err = %v", err)
	}

	var methodErr *MethodError
	if !errors.As(err, &methodErr) {
		t.Fatalf("errors.As(*MethodError) = false, want true; err = %v", err)
	}
	if methodErr.OptionsValue != "invalid" {
		t.Errorf("MethodError.OptionsValue = %v, want \"invalid\"", methodErr.OptionsValue)
	}
}

// ------------------------------------------------------------------
// Task (raw)
// ------------------------------------------------------------------

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

	var solution KasadaStandardSolution
	if err := json.Unmarshal(got.Solution, &solution); err != nil {
		t.Fatalf("Failed to unmarshal solution: %v", err)
	}
	if solution.UserAgent != "UA" {
		t.Fatalf("Task() solution.UserAgent = %q, want UA", solution.UserAgent)
	}
}

// HTTP 200, ErrorCode = 1
func TestTask_Raw_APIError(t *testing.T) {
	const wantTaskId = "task-fail-1"
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"errorId":1,"status":"failed","solution":null}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	_, err := c.Task(context.Background(), wantTaskId)
	if err == nil {
		t.Fatal("Task() expected error on errorId=1, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Task() error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusOK {
		t.Errorf("APIError.StatusCode = %d, want %d", apiErr.StatusCode, http.StatusOK)
	}
	if apiErr.TaskId != wantTaskId {
		t.Errorf("APIError.TaskId = %q, want %q", apiErr.TaskId, wantTaskId)
	}
}

// HTTP 400
func TestTask_Raw_HTTP400(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error_code":1,"error_description":"task not found"}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	_, err := c.Task(context.Background(), "bad-task-id")
	if err == nil {
		t.Fatal("Task() expected error on 400, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Task() error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("APIError.StatusCode = %d, want %d", apiErr.StatusCode, http.StatusBadRequest)
	}
	if apiErr.Msg == "" {
		t.Error("APIError.Msg is empty, want non-empty description")
	}
}

// ------------------------------------------------------------------
// GetTaskResult (typed)
// ------------------------------------------------------------------

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

	got, err := GetTaskResult[KasadaStandardSolution](c, context.Background(), "task-1")
	if err != nil {
		t.Fatalf("GetTaskResult() error: %v", err)
	}
	if got.Solution.UserAgent != "UA" || got.Solution.XKpsdkCd != "cd" {
		t.Fatalf("GetTaskResult() unexpected solution: %+v", got.Solution)
	}
}

// HTTP 200, ErrorCode = 1
func TestGetTaskResult_APIError(t *testing.T) {
	const wantTaskId = "task-err-42"
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"errorId":1,"status":"failed","solution":null}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	_, err := GetTaskResult[KasadaStandardSolution](c, context.Background(), wantTaskId)
	if err == nil {
		t.Fatal("GetTaskResult() expected error on errorId=1, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("GetTaskResult() error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusOK {
		t.Errorf("APIError.StatusCode = %d, want %d", apiErr.StatusCode, http.StatusOK)
	}
	if apiErr.TaskId != wantTaskId {
		t.Errorf("APIError.TaskId = %q, want %q", apiErr.TaskId, wantTaskId)
	}
}

// HTTP 400
func TestGetTaskResult_HTTP400(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error_code":1,"error_description":"invalid taskId format"}`))
	}

	c, closeFn := newTestClient(t, h)
	defer closeFn()

	_, err := GetTaskResult[KasadaStandardSolution](c, context.Background(), "malformed")
	if err == nil {
		t.Fatal("GetTaskResult() expected error on 400, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("GetTaskResult() error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("APIError.StatusCode = %d, want %d", apiErr.StatusCode, http.StatusBadRequest)
	}
	if apiErr.Msg == "" {
		t.Error("APIError.Msg is empty, want non-empty description")
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
