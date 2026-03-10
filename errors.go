package salamoonder

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNoApiKey = errors.New("no api key")

	/*
		ErrUnsupportedTaskOptionsType is returned from CreateTask if the provided
		options type is not supported. Use errors.Is for checking,
		and errors.As(*MethodError) to get details (the actual type).
	*/
	ErrUnsupportedTaskOptionsType = errors.New("unsupported task options type")

	_ error = (*APIError)(nil)
	_ error = (*MethodError)(nil)
)

var allowedTaskTypes = make([]reflect.Type, 0)

type (
	APIError struct {
		TaskId     string
		StatusCode int
		Msg        string
	}

	MethodError struct {
		OptionsValue any
	}
)

func (a *APIError) Error() string {
	switch a.StatusCode {
	case 200:
		if a.TaskId != "" {
			return fmt.Sprintf("api error for task [%s]: %s", a.TaskId, a.Msg)
		}
		return fmt.Sprintf("api error: %s", a.Msg)
	case 400:
		if a.TaskId != "" {
			return fmt.Sprintf("invalid request for task [%s]: %s", a.TaskId, a.Msg)
		}
		return fmt.Sprintf("invalid request: %s", a.Msg)
	default:
		return fmt.Sprintf("unexpected api error (status %d): %s", a.StatusCode, a.Msg)
	}
}

func (m *MethodError) Error() string {
	actualType := reflect.TypeOf(m.OptionsValue)
	allowed := make([]string, len(allowedTaskTypes))
	for i, t := range allowedTaskTypes {
		allowed[i] = t.Name()
	}

	return fmt.Sprintf(
		"invalid Task type %v; allowed: %v",
		actualType.Name(),
		allowed,
	)
}

/*
Is allows using errors.Is(err, ErrUnsupportedTaskOptionsType).
To get details (actual type), use errors.As(*MethodError).
*/
func (m *MethodError) Is(target error) bool {
	return target == ErrUnsupportedTaskOptionsType
}
