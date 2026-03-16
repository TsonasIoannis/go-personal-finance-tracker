package apperrors

import "errors"

type Kind string

const (
	KindValidation   Kind = "validation"
	KindUnauthorized Kind = "unauthorized"
	KindForbidden    Kind = "forbidden"
	KindNotFound     Kind = "not_found"
	KindConflict     Kind = "conflict"
	KindUnavailable  Kind = "unavailable"
	KindInternal     Kind = "internal"
)

type Error struct {
	Kind    Kind
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.Message != "" {
		return e.Message
	}

	if e.Err != nil {
		return e.Err.Error()
	}

	return "application error"
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func New(kind Kind, code, message string) *Error {
	return &Error{
		Kind:    kind,
		Code:    code,
		Message: message,
	}
}

func Wrap(kind Kind, code, message string, err error) *Error {
	return &Error{
		Kind:    kind,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func Validation(code, message string) *Error {
	return New(KindValidation, code, message)
}

func Unauthorized(code, message string) *Error {
	return New(KindUnauthorized, code, message)
}

func Forbidden(code, message string) *Error {
	return New(KindForbidden, code, message)
}

func NotFound(code, message string) *Error {
	return New(KindNotFound, code, message)
}

func Conflict(code, message string) *Error {
	return New(KindConflict, code, message)
}

func Unavailable(code, message string) *Error {
	return New(KindUnavailable, code, message)
}

func Internal(code, message string, err error) *Error {
	return Wrap(KindInternal, code, message, err)
}

func As(err error) (*Error, bool) {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr, true
	}

	return nil, false
}
