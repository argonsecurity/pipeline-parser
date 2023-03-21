package http

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestDefaultRetryPolicy(t *testing.T) {
	type args struct {
		resp *http.Response
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "nil response and nil error",
			args: args{
				resp: nil,
				err:  nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "nil response and not nil error",
			args: args{
				resp: nil,
				err:  errors.New("error"),
			},
			want:    true,
			wantErr: true,
		},
		{
			name: "response with status code 0 and nil error",
			args: args{
				resp: &http.Response{
					StatusCode: 0,
				},
				err: nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "response with 0 < status code < 500 and nil error",
			args: args{
				resp: &http.Response{
					StatusCode: 400,
				},
				err: nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "response with status code 500 and nil error",
			args: args{
				resp: &http.Response{
					StatusCode: 500,
				},
				err: nil,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DefaultRetryPolicy(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultRetryPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DefaultRetryPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBackoff(t *testing.T) {
	type args struct {
		min        time.Duration
		max        time.Duration
		attemptNum int
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "attemptNum = 0",
			args: args{
				min:        1 * time.Second,
				max:        10 * time.Second,
				attemptNum: 0,
			},
			want: 1 * time.Second,
		},
		{
			name: "not exceeding max",
			args: args{
				min:        1 * time.Second,
				max:        10 * time.Second,
				attemptNum: 2,
			},
			want: 4 * time.Second,
		},
		{
			name: "exceeding max",
			args: args{
				min:        1 * time.Second,
				max:        10 * time.Second,
				attemptNum: 5,
			},
			want: 10 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultBackoff(tt.args.min, tt.args.max, tt.args.attemptNum); got != tt.want {
				t.Errorf("DefaultBackoff() = %v, want %v", got, tt.want)
			}
		})
	}
}
