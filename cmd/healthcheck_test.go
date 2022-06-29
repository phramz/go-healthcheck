package main

var basicAssertions = []struct {
	Name             string
	Probe            string
	ExpectedExitCode int
	ExpectedResult   bool
}{
	{Name: "True", Probe: `{{ .Assert.True false }}`},
	{Name: "Equal", Probe: `{{ .Assert.Equal 200 (.HTTP.Request "GET" "http://localhost:8080").StatusCode }}`},
	{Name: "Regexp", Probe: `{{ .Assert.Regexp "^start"|regexp "its not starting" }}`},
	{Name: "NotRegexp", Probe: `{{ .Assert.NotRegexp (call .Regexp "^start") "its not starting" }}`},
	{Name: "HTTPSuccess", Probe: `{{ .Assert.HTTPSuccess .HTTP.Handler "GET" "http://localhost:8080" nil }}`},
	{Name: "HTTPBodyContains", Probe: `{{ .Assert.HTTPBodyContains .HTTP.Handler "GET" "http://localhost:8080" nil }}`},
	{Name: "HTTPBodyNotContains", Probe: `{{ .Assert.HTTPBodyNotContains .HTTP.Handler "GET" "http://localhost:8080" nil }}`},
}

//func TestMain(t *testing.T) {
//	for _, tc := range basicAssertions {
//		tc := tc
//		t.Run(tc.Name, func(t *testing.T) {
//			t.Parallel()
//			assert.True(t, assert.ObjectsAreEqualValues(tc.Expected, res))
//		})
//	}
//
//	for _, tc := range basicAssertions {
//		exitCode := m.Run()
//
//		if tc.ExpectedExitCode != m.Run() {
//			m.
//		}
//	}
//
//	assert.Equal(t, 4, out.Score(identities["identityI"], identities["identityJ"]))
//}
