package internal

// JsonResponse is the huma output envelope; huma serializes only the Body field.
type JsonResponse[T any] struct {
	Body T
}

func NewResponse[T any](body T) *JsonResponse[T] {
	return &JsonResponse[T]{Body: body}
}
