package errors_test

import (
	"errors"
	"testing"

	skylineerrors "github.com/github/gh-skyline/internal/errors"
)

func TestSkylineError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *skylineerrors.SkylineError
		want string
	}{
		{
			name: "error with underlying error",
			err: &skylineerrors.SkylineError{
				Type:    skylineerrors.ValidationError,
				Message: "invalid input",
				Err:     errors.New("value out of range"),
			},
			want: "[VALIDATION] invalid input: value out of range",
		},
		{
			name: "error without underlying error",
			err: &skylineerrors.SkylineError{
				Type:    skylineerrors.STLError,
				Message: "failed to process STL",
			},
			want: "[STL] failed to process STL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("SkylineError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		message string
		want    string
		wantNil bool
	}{
		{
			name:    "nil error returns nil",
			err:     nil,
			message: "test message",
			wantNil: true,
		},
		{
			name:    "wrap standard error",
			err:     errors.New("original error"),
			message: "wrapped message",
			want:    "[STL] wrapped message: original error",
		},
		{
			name: "wrap SkylineError preserves type",
			err: &skylineerrors.SkylineError{
				Type:    skylineerrors.ValidationError,
				Message: "original message",
				Err:     errors.New("base error"),
			},
			message: "wrapped message",
			want:    "[VALIDATION] wrapped message: original message: base error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := skylineerrors.Wrap(tt.err, tt.message)
			if tt.wantNil {
				if got != nil {
					t.Errorf("Wrap() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("Wrap() returned nil, want error")
			}
			if got.Error() != tt.want {
				t.Errorf("Wrap() = %v, want %v", got.Error(), tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		errType skylineerrors.ErrorType
		message string
		err     error
		want    string
	}{
		{
			name:    "new error without underlying error",
			errType: skylineerrors.ValidationError,
			message: "validation failed",
			err:     nil,
			want:    "[VALIDATION] validation failed",
		},
		{
			name:    "new error with underlying error",
			errType: skylineerrors.NetworkError,
			message: "network timeout",
			err:     errors.New("connection refused"),
			want:    "[NETWORK] network timeout: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := skylineerrors.New(tt.errType, tt.message, tt.err)
			if got.Error() != tt.want {
				t.Errorf("New() error = %v, want %v", got.Error(), tt.want)
			}
			if got.Type != tt.errType {
				t.Errorf("New() type = %v, want %v", got.Type, tt.errType)
			}
			if got.Message != tt.message {
				t.Errorf("New() message = %v, want %v", got.Message, tt.message)
			}
			if got.Err != tt.err {
				t.Errorf("New() underlying error = %v, want %v", got.Err, tt.err)
			}
		})
	}
}

func TestSkylineError_Is(t *testing.T) {
	tests := []struct {
		name   string
		err    *skylineerrors.SkylineError
		target error
		want   bool
	}{
		{
			name: "matching error types",
			err: &skylineerrors.SkylineError{
				Type: skylineerrors.ValidationError,
			},
			target: &skylineerrors.SkylineError{
				Type: skylineerrors.ValidationError,
			},
			want: true,
		},
		{
			name: "different error types",
			err: &skylineerrors.SkylineError{
				Type: skylineerrors.ValidationError,
			},
			target: &skylineerrors.SkylineError{
				Type: skylineerrors.NetworkError,
			},
			want: false,
		},
		{
			name: "non-SkylineError target",
			err: &skylineerrors.SkylineError{
				Type: skylineerrors.ValidationError,
			},
			target: errors.New("standard error"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Is(tt.target); got != tt.want {
				t.Errorf("SkylineError.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkylineError_Unwrap(t *testing.T) {
	baseErr := errors.New("base error")
	tests := []struct {
		name    string
		err     *skylineerrors.SkylineError
		wantErr error
	}{
		{
			name: "with underlying error",
			err: &skylineerrors.SkylineError{
				Type:    skylineerrors.ValidationError,
				Message: "test message",
				Err:     baseErr,
			},
			wantErr: baseErr,
		},
		{
			name: "without underlying error",
			err: &skylineerrors.SkylineError{
				Type:    skylineerrors.ValidationError,
				Message: "test message",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.err.Unwrap(); err != tt.wantErr {
				t.Errorf("SkylineError.Unwrap() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
