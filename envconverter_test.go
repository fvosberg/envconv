package envconv

import (
	"flag"
	"testing"
)

func TestTable(t *testing.T) {
	tests := []struct {
		title    string
		envv     map[string]string
		argv     []string
		expected string
	}{

		{
			title:    "Test flag default",
			envv:     map[string]string{},
			argv:     []string{},
			expected: "default-value",
		},
		{
			title:    "Set flag has precedence",
			envv:     map[string]string{},
			argv:     []string{"--foobar", "flag-value"},
			expected: "flag-value",
		},
		{
			title: "Env overwrites default",
			envv: map[string]string{
				"EXAMPLE_FOOBAR": "ENV-VALUE",
			},
			argv:     []string{},
			expected: "ENV-VALUE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			t.Run("plain", func(t *testing.T) {
				testPlainApproach(t, tt.envv, tt.argv, tt.expected)
			})
			t.Run("prepend", func(t *testing.T) {
				testPrepend(t, tt.envv, tt.argv, tt.expected)
			})
			t.Run("double", func(t *testing.T) {
				testDoubleParse(t, tt.envv, tt.argv, tt.expected)
			})
		})
	}
}

func testPlainApproach(t *testing.T, envv map[string]string, argv []string, expected string) {
	fs := flag.NewFlagSet("command", flag.ContinueOnError)
	act := fs.String("foobar", "default-value", "the foobar flag")

	var flagsFromEnv []string
	if envv["EXAMPLE_FOOBAR"] != "" {
		flagsFromEnv = append(flagsFromEnv, "--foobar", envv["EXAMPLE_FOOBAR"])
	}
	err := fs.Parse(append(flagsFromEnv, argv...))
	if err != nil {
		t.Fatalf("Parsing command line flags failes: %s", err)
	}
	if *act != expected {
		t.Errorf("Expected value: %q; got %q", expected, *act)
	}
}

func testPrepend(t *testing.T, envv map[string]string, argv []string, expected string) {
	ec := envConverter{
		getEnv: func(name string) string {
			return envv[name]
		},
	}

	fs := flag.NewFlagSet("command", flag.ContinueOnError)
	act := fs.String("foobar", "default-value", "the foobar flag")
	err := fs.Parse(ec.PrependToFlags(argv, map[string]string{"--foobar": "EXAMPLE_FOOBAR"}))

	if err != nil {
		t.Fatalf("Parsing command line flags failes: %s", err)
	}
	if *act != expected {
		t.Errorf("Expected value: %q; got %q", expected, *act)
	}
}

func testDoubleParse(t *testing.T, envv map[string]string, argv []string, expected string) {
	ec := envConverter{
		getEnv: func(name string) string {
			return envv[name]
		},
	}

	fs := flag.NewFlagSet("command", flag.ContinueOnError)
	act := fs.String("foobar", "default-value", "the foobar flag")
	err := fs.Parse(ec.FlagsFromEnv(map[string]string{"--foobar": "EXAMPLE_FOOBAR"}))
	if err != nil {
		t.Fatalf("Parsing env flags failed: %s", err)
	}
	err = fs.Parse(argv)
	if err != nil {
		t.Fatalf("Parsing command line flags failed: %s", err)
	}
	if *act != expected {
		t.Errorf("Expected value: %q; got %q", expected, *act)
	}
}
