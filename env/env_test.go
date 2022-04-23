package env_test

import (
	"os"
	"testing"

	"github.com/trivelaapp/go-kit/env"
)

func TestGetEnvironmentByLabel(t *testing.T) {
	tt := []struct {
		desc  string
		input string
		env   env.Environment
	}{
		{
			desc:  "should be Dev #1",
			input: "development",
			env:   env.DevelopmentEnvironment,
		},
		{
			desc:  "should be Dev #2",
			input: "dev",
			env:   env.DevelopmentEnvironment,
		},
		{
			desc:  "should be Dev #3",
			input: "DEVELOPMENT",
			env:   env.DevelopmentEnvironment,
		},
		{
			desc:  "should be Test #1",
			input: "test",
			env:   env.TestEnvironment,
		},
		{
			desc:  "should be Test #2",
			input: "homolog",
			env:   env.TestEnvironment,
		},
		{
			desc:  "should be Test #3",
			input: "staging",
			env:   env.TestEnvironment,
		},
		{
			desc:  "should be Prod #1",
			input: "production",
			env:   env.ProductionEnvironment,
		},
		{
			desc:  "should be Prod #2",
			input: "prod",
			env:   env.ProductionEnvironment,
		},
		{
			desc:  "should be Prod #3",
			input: "PROD",
			env:   env.ProductionEnvironment,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			e := env.GetEnvironmentByLabel(tc.input)
			if e != tc.env {
				t.Errorf("Mismatch response, got %d", e)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "fake-env-1",
		"FAKE_ENV_2": "fake-env-2",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc          string
		env           string
		defaultValues []string
		expectedValue string
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			defaultValues: []string{"default-value-1", "default-value-2"},
			expectedValue: "fake-env-1",
		},
		{
			desc:          "should access an unknown environment variable with default value successfully",
			env:           "FAKE_ENV_3",
			defaultValues: []string{"default-value-1", "default-value-2"},
			expectedValue: "default-value-1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			res := env.GetString(tc.env, tc.defaultValues...)
			if tc.expectedValue != res {
				t.Errorf("Mismatch response. Expected '%s', got '%s'.", tc.expectedValue, res)
			}
		})
	}
}

func TestMustGetString(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "fake-env-1",
		"FAKE_ENV_2": "fake-env-2",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc             string
		env              string
		expectedPanicMsg string
		expectedValue    string
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			expectedValue: "fake-env-1",
		},
		{
			desc:             "should panic when accessing an unknown environment variable",
			env:              "FAKE_ENV_3",
			expectedPanicMsg: "FAKE_ENV_3 can't be empty",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				e := recover()
				if e == nil && tc.expectedPanicMsg == "" {
					return
				}

				if e != tc.expectedPanicMsg {
					t.Errorf("Mismatch panic msg response. Expected '%s', got '%s'.", tc.expectedPanicMsg, e)
				}
			}()

			res := env.MustGetString(tc.env)
			if tc.expectedValue != res {
				t.Errorf("Mismatch response. Expected '%s', got '%s'.", tc.expectedValue, res)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "1",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc          string
		env           string
		defaultValues []int
		expectedValue int
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			defaultValues: []int{1, 2},
			expectedValue: 1,
		},
		{
			desc:          "should access an unknown environment variable with default value successfully",
			env:           "FAKE_ENV_3",
			defaultValues: []int{1, 2},
			expectedValue: 1,
		},
		{
			desc:          "should use default value when environment variable is not an Int",
			env:           "FAKE_ENV_2",
			defaultValues: []int{1, 2},
			expectedValue: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			res := env.GetInt(tc.env, tc.defaultValues...)
			if tc.expectedValue != res {
				t.Errorf("Mismatch response. Expected '%d', got '%d'.", tc.expectedValue, res)
			}
		})
	}
}

func TestMustGetInt(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "1",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc             string
		env              string
		expectedPanicMsg string
		expectedValue    int
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			expectedValue: 1,
		},
		{
			desc:             "should panic when environment variable value is not and int",
			env:              "FAKE_ENV_2",
			expectedPanicMsg: "FAKE_ENV_2 must contain an int value",
		},
		{
			desc:             "should panic when accessing an unknown environment variable",
			env:              "FAKE_ENV_3",
			expectedPanicMsg: "FAKE_ENV_3 must contain an int value",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				e := recover()
				if e == nil && tc.expectedPanicMsg == "" {
					return
				}

				if e != tc.expectedPanicMsg {
					t.Errorf("Mismatch panic msg response. Expected '%s', got '%s'.", tc.expectedPanicMsg, e)
				}
			}()

			res := env.MustGetInt(tc.env)
			if tc.expectedValue != res {
				t.Errorf("Mismatch response. Expected '%d', got '%d'.", tc.expectedValue, res)
			}
		})
	}
}

func TestGetFloat(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "1.7",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc          string
		env           string
		defaultValues []float64
		expectedValue float64
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			defaultValues: []float64{1.5, 2.3},
			expectedValue: 1.7,
		},
		{
			desc:          "should access an unknown environment variable with default value successfully",
			env:           "FAKE_ENV_3",
			defaultValues: []float64{1.5, 2.3},
			expectedValue: 1.5,
		},
		{
			desc:          "should use default value when environment variable is not a Float",
			env:           "FAKE_ENV_2",
			defaultValues: []float64{1.5, 2.3},
			expectedValue: 1.5,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			res := env.GetFloat(tc.env, tc.defaultValues...)
			if tc.expectedValue != res {
				t.Errorf("Mismatch response. Expected '%f', got '%f'.", tc.expectedValue, res)
			}
		})
	}
}

func TestMustGetFloat(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "1.7",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc             string
		env              string
		expectedPanicMsg string
		expectedValue    float64
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			expectedValue: 1.7,
		},
		{
			desc:             "should panic when accessing an environment variable that isn't a Float",
			env:              "FAKE_ENV_2",
			expectedPanicMsg: "FAKE_ENV_2 must contain a float value",
		},
		{
			desc:             "should panic when accessing an unknown environment variable",
			env:              "FAKE_ENV_3",
			expectedPanicMsg: "FAKE_ENV_3 must contain a float value",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				e := recover()
				if e == nil && tc.expectedPanicMsg == "" {
					return
				}

				if e != tc.expectedPanicMsg {
					t.Errorf("Mismatch panic msg response. Expected '%s', got '%s'.", tc.expectedPanicMsg, e)
				}
			}()

			res := env.MustGetFloat(tc.env)
			if tc.expectedValue != res {
				t.Errorf("Mismatch response. Expected '%f', got '%f'.", tc.expectedValue, res)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "true",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc          string
		env           string
		defaultValues []bool
		expectedValue bool
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			defaultValues: []bool{true, false},
			expectedValue: true,
		},
		{
			desc:          "should use default value when environment variable is not a Float",
			env:           "FAKE_ENV_2",
			defaultValues: []bool{true, false},
			expectedValue: true,
		},
		{
			desc:          "should access an unknown environment variable with default value successfully",
			env:           "FAKE_ENV_3",
			defaultValues: []bool{true, false},
			expectedValue: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			res := env.GetBool(tc.env, tc.defaultValues...)
			if tc.expectedValue != res {
				t.Errorf("Mismatch response. Expected '%t', got '%t'.", tc.expectedValue, res)
			}
		})
	}
}

func TestMustGetBool(t *testing.T) {
	envs := map[string]string{
		"FAKE_ENV_1": "true",
		"FAKE_ENV_2": "fake-invalid-env",
	}

	for key, value := range envs {
		_ = os.Setenv(key, value)
	}

	tt := []struct {
		desc             string
		env              string
		expectedPanicMsg string
		expectedValue    bool
	}{
		{
			desc:          "should access a valid environment variable successfully",
			env:           "FAKE_ENV_1",
			expectedValue: true,
		},
		{
			desc:             "should panic when accessing an environment variable that isn't a Boolean",
			env:              "FAKE_ENV_2",
			expectedPanicMsg: "FAKE_ENV_2 must contain a boolean value",
		},
		{
			desc:             "should panic when accessing an unknown environment variable",
			env:              "FAKE_ENV_3",
			expectedPanicMsg: "FAKE_ENV_3 must contain a boolean value",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				e := recover()
				if e == nil && tc.expectedPanicMsg == "" {
					return
				}

				if e != tc.expectedPanicMsg {
					t.Errorf("Mismatch panic msg response. Expected '%s', got '%s'.", tc.expectedPanicMsg, e)
				}
			}()

			res := env.MustGetBool(tc.env)
			if tc.expectedValue != res {
				t.Errorf("Mismatch response. Expected '%t', got '%t'.", tc.expectedValue, res)
			}
		})
	}
}
