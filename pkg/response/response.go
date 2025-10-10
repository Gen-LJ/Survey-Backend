package response

type Response[T any] struct {
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
	Data    *T      `json:"data,omitempty"` // use pointer to allow nil
}

// Success method
func (r *Response[T]) OK(data T) *Response[T] {
	return &Response[T]{
		Success: true,
		Data:    &data, // pointer now
	}
}

// Success with no data method
func (r *Response[T]) StatusOK(msg string) *Response[T] {
	return &Response[T]{
		Success: true,
		Message: msg,
		Data:    nil,
	}
}

// Error method
func (r *Response[T]) Fail(msg string) *Response[T] {
	return &Response[T]{
		Success: false,
		Message: msg,
		Data:    nil, // data is nil on failure
	}
}
