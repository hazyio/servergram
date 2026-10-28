package internal

import (
	"errors"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type BodyRequest[T any] struct {
	Body T
}

// JsonResponse is the huma output envelope; huma serializes only the Body field.
type JsonResponse[T any] struct {
	Body T
}

// Sentinel errors your business logic can return (optionally wrapped with %w).
var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrNotFound         = errors.New("not found")
	ErrConflict         = errors.New("conflict")
	InternalServerError = errors.New("internal server error")
	BadRequest          = errors.New("bad request")
)

func NewResponse[T any](body T) *JsonResponse[T] {
	return &JsonResponse[T]{Body: body}
}
func NewHumaError(status int, message string) error {
	return &huma.ErrorModel{
		Status: status,
		Title:  http.StatusText(status),
		Detail: message,
	}
}

func ToHumaError(err error) error {
	// Already a Huma error: pass through untouched.
	var se huma.StatusError
	if errors.As(err, &se) {
		return se
	}

	switch {
	case errors.Is(err, ErrUnauthorized):
		return huma.Error401Unauthorized("unauthorized")
	case errors.Is(err, ErrNotFound):
		return huma.Error404NotFound("not found")
	case errors.Is(err, ErrConflict):
		return huma.Error409Conflict("conflict")
	case errors.Is(err, InternalServerError):
		return huma.Error500InternalServerError("internal server error")
	case errors.Is(err, BadRequest):
		return huma.Error400BadRequest("bad request")
	}

	// Unknown: log the real error, return a generic one.
	log.Printf("unhandled error: %v", err)
	return huma.Error500InternalServerError("internal error")
}
