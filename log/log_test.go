package log

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/trivelaapp/go-kit/log/format"
)

func mockedTimmer() func() time.Time {
	return func() time.Time {
		return time.Date(2021, 0, 1, 12, 0, 0, 0, time.UTC)
	}
}

func TestNewLogger(t *testing.T) {
	t.Run("should use LogLevel INFO when not specified", func(t *testing.T) {
		ctx := context.Background()

		logger := NewLogger(LoggerParams{})
		logger.now = mockedTimmer()

		out := captureOutput(func() {
			logger.Debug(ctx, "random message")
		})

		if diff := cmp.Diff("", out); diff != "" {
			t.Errorf("mismatch (-want, +got):\n%s", diff)
		}

		out = captureOutput(func() {
			logger.Info(ctx, "random message")
		})

		if diff := cmp.Diff(`{"level":"INFO","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`, out); diff != "" {
			t.Errorf("mismatch (-want, +got):\n%s", diff)
		}
	})
}

func TestDebug(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		attrs       format.LogAttributeSet
		msg         string
		msgArgs     []interface{}
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message",
			expectedLog: `{"level":"DEBUG","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message with dynamic data %d",
			msgArgs:     []interface{}{1},
			expectedLog: `{"level":"DEBUG","message":"random message with dynamic data 1","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       "DEBUG",
			msg:         "random message",
			attrs:       format.LogAttributeSet{"attr1": true},
			expectedLog: `{"attributes":{"attr1":"value1"},"level":"DEBUG","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := NewLogger(LoggerParams{Level: tc.level, Attributes: tc.attrs})
			logger.now = mockedTimmer()

			out := captureOutput(func() {
				logger.Debug(tc.ctx, tc.msg, tc.msgArgs...)
			})

			if diff := cmp.Diff(tc.expectedLog, out); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		attrs       format.LogAttributeSet
		msg         string
		msgArgs     []interface{}
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message",
			expectedLog: `{"level":"INFO","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			msg:         "random message",
			expectedLog: `{"level":"INFO","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message with dynamic data %d",
			msgArgs:     []interface{}{1},
			expectedLog: `{"level":"INFO","message":"random message with dynamic data 1","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       "DEBUG",
			msg:         "random message",
			attrs:       format.LogAttributeSet{"attr1": true},
			expectedLog: `{"attributes":{"attr1":"value1"},"level":"INFO","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := NewLogger(LoggerParams{Level: tc.level, Attributes: tc.attrs})
			logger.now = mockedTimmer()

			out := captureOutput(func() {
				logger.Info(tc.ctx, tc.msg, tc.msgArgs...)
			})

			if diff := cmp.Diff(tc.expectedLog, out); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestWarning(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		attrs       format.LogAttributeSet
		msg         string
		msgArgs     []interface{}
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message",
			expectedLog: `{"level":"WARNING","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			msg:         "random message",
			expectedLog: `{"level":"WARNING","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			msg:         "random message",
			expectedLog: `{"level":"WARNING","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			msg:         "random message",
			expectedLog: "",
		},
		{
			desc:        "should log with dynamic message",
			ctx:         ctx,
			level:       "DEBUG",
			msg:         "random message with dynamic data %d",
			msgArgs:     []interface{}{1},
			expectedLog: `{"level":"WARNING","message":"random message with dynamic data 1","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       "DEBUG",
			msg:         "random message",
			attrs:       format.LogAttributeSet{"attr1": true},
			expectedLog: `{"attributes":{"attr1":"value1"},"level":"WARNING","message":"random message","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := NewLogger(LoggerParams{Level: tc.level, Attributes: tc.attrs})
			logger.now = mockedTimmer()

			out := captureOutput(func() {
				logger.Warning(tc.ctx, tc.msg, tc.msgArgs...)
			})

			if diff := cmp.Diff(tc.expectedLog, out); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestError(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		attrs       format.LogAttributeSet
		err         error
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"random error"},"level":"ERROR","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"random error"},"level":"ERROR","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"random error"},"level":"ERROR","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"random error"},"level":"ERROR","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			err:         errors.New("random error"),
			expectedLog: "",
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       "DEBUG",
			err:         errors.New("random error"),
			attrs:       format.LogAttributeSet{"attr1": true},
			expectedLog: `{"attributes":{"attr1":"value1","err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"random error"},"level":"ERROR","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := NewLogger(LoggerParams{Level: tc.level, Attributes: tc.attrs})
			logger.now = mockedTimmer()

			out := captureOutput(func() {
				logger.Error(tc.ctx, tc.err)
			})

			if diff := cmp.Diff(tc.expectedLog, out); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestCritical(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		level       string
		attrs       format.LogAttributeSet
		err         error
		expectedLog string
	}{
		{
			desc:        "should log when LogLevel is DEBUG",
			ctx:         ctx,
			level:       "DEBUG",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is INFO",
			ctx:         ctx,
			level:       "INFO",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is WARNING",
			ctx:         ctx,
			level:       "WARNING",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log when LogLevel is ERROR",
			ctx:         ctx,
			level:       "ERROR",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log when LogLevel is CRITICAL",
			ctx:         ctx,
			level:       "CRITICAL",
			err:         errors.New("random error"),
			expectedLog: `{"attributes":{"err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log with attributes",
			ctx:         context.WithValue(ctx, "attr1", "value1"),
			level:       "DEBUG",
			err:         errors.New("random error"),
			attrs:       format.LogAttributeSet{"attr1": true},
			expectedLog: `{"attributes":{"attr1":"value1","err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"random error"},"level":"CRITICAL","message":"random error","timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := NewLogger(LoggerParams{Level: tc.level, Attributes: tc.attrs})
			logger.now = mockedTimmer()

			out := captureOutput(func() {
				logger.Critical(tc.ctx, tc.err)
			})

			if diff := cmp.Diff(tc.expectedLog, out); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	ctx := context.Background()

	tt := []struct {
		desc        string
		ctx         context.Context
		sysLogLevel string
		level       []Level
		attrs       format.LogAttributeSet
		data        any
		expectedLog string
	}{
		{
			desc:        "should log when system LogLevel is DEBUG and uses default JSON log Level",
			ctx:         ctx,
			sysLogLevel: "DEBUG",
			data:        map[string]string{"foo": "bar"},
			expectedLog: `{"level":"DEBUG","message":"JSON data logged","payload":{"foo":"bar"},"timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should log an error when data can not be JSON marshalled",
			ctx:         ctx,
			sysLogLevel: "DEBUG",
			data:        make(chan int),
			expectedLog: `{"attributes":{"err_code":"UNKNOWN","err_kind":"UNEXPECTED","root_error":"json: unsupported type: chan int"},"level":"ERROR","message":"could not marshal payload to JSON format","timestamp":"2020-12-01T12:00:00Z"}`,
		},
		{
			desc:        "should not log anything when system LogLevel is INFO and uses default JSON log Level",
			ctx:         ctx,
			sysLogLevel: "INFO",
			data:        map[string]string{"foo": "bar"},
			expectedLog: "",
		},
		{
			desc:        "should log when system LogLevel is INFO and JSON log Level is WARNING",
			ctx:         ctx,
			sysLogLevel: "INFO",
			level:       []Level{LevelWarning},
			data:        map[string]string{"foo": "bar"},
			expectedLog: `{"level":"WARNING","message":"JSON data logged","payload":{"foo":"bar"},"timestamp":"2020-12-01T12:00:00Z"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			logger := NewLogger(LoggerParams{Level: tc.sysLogLevel, Attributes: tc.attrs})
			logger.now = mockedTimmer()

			out := captureOutput(func() {
				logger.JSON(tc.ctx, tc.data, tc.level...)
			})

			if diff := cmp.Diff(tc.expectedLog, out); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func captureOutput(output func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	output()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	return strings.TrimRight(string(out), "\n")
}
