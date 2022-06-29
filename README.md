Health Checker
==============

Tiny (< 9MB) CLI utility to perform health checks everywhere (container, local, etc)

## TL;DR

```shell
# from source
make build

# check if github.com is up
$ build/bin/healthcheck probe '{{ .Assert.HTTPStatusCode .HTTP.Handler "GET" "https://github.com" nil 200 }}'
INFO   [16:47:40] starting probe                                probe=p0 status=STARTED
INFO   [16:47:40] finished probe                                probe=p0 status=PASS
INFO   [16:47:40] 1 of 1 probes passed 👍
EXIT 0
```

## Usage

```
NAME:
   healthcheck - HealthCheck Utility

USAGE:
   healthcheck [global options] command [command options] [arguments...]

COMMANDS:
   probe    Performs a HTTP connection check
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h             show help (default: false)
   --log-level value      Set the log output verbosity (default: "info") [$HC_DEBUG]
   --no-fail              If set failed checks won't cause the program to exit with an error code (default: false) [$HC_NO_FAIL]
   --no-fail-affirmative  If set the program won't exit with an error code as long as at least one probe passed (default: false) [$HC_NO_FAIL_AFFIRMATIVE]
   --output               If set rendered output will be printed to stdout (default: false) [$HC_OUTPUT]
   --timeout value        Maximum time in seconds after an unfinished probe will fail (default: 30) [$HC_TIMEOUT]
```

### Examples

```shell
# play with env vars
$ MY_NAME=Earl ./healthcheck probe '{{ .Assert.Equal "My name is Earl!" ("My name is $MY_NAME!"|expandenv) }}'
INFO   [16:19:52] starting probe                                probe=p0 status=STARTED
INFO   [16:19:52] finished probe                                probe=p0 status=PASS
INFO   [16:19:52] 1 of 1 probes passed 👍
EXIT 0

# run multiple probes
$ MY_NAME=Earl ./healthcheck probe '{{ .Assert.Equal "Me llamo Earl!" ("My name is $MY_NAME!"|expandenv) }}' \
                                   '{{ .Assert.Equal "My name is Earl!" ("My name is $MY_NAME!"|expandenv) }}'
INFO   [16:41:53] starting probe                                probe=p0 status=STARTED
WARNING[16:41:53]

	Error:      	Not equal:
	            	expected: "Me llamo Earl!"
	            	actual  : "My name is Earl!"

	            	Diff:
	            	--- Expected
	            	+++ Actual
	            	@@ -1 +1 @@
	            	-Me llamo Earl!
	            	+My name is Earl!  probe=p0
WARNING[16:41:53] finished probe                                probe=p0 status=FAIL
INFO   [16:41:53] starting probe                                probe=p1 status=STARTED
INFO   [16:41:53] finished probe                                probe=p1 status=PASS
ERROR  [16:41:53] 1 of 2 probes failed 😭
FAIL
EXIT 1
```

### Probes

Probes are basically go-templates (see "https://pkg.go.dev/text/template" to learn how utilize them) where you can
perform assertions on static or dynamic values

The renderer comes with the Sprig extension to add some more convenience. See
"http://masterminds.github.io/sprig/" for documentation and the sources on Github (https://github.com/Masterminds/sprig)

Here are some examples on how sprig can help you:

#### Environment variables

```bash
export HTTP_PORT=8080 healthcheck probe '{{ .Assert.Equal "http://localhost:8080" ("http://localhost:$HTTP_PORT"|expandenv) }}'
export HTTP_URL=http://localhost:8080 healthcheck probe '{{ .Assert.Equal "http://localhost:8080" ("HTTP_URL"|env) }}'
```

### Assertions

Most assertions available in the probes are provided by Testify (see https://github.com/stretchr/testify).
Have a look at the API Doc to see detailed information about it:
https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Assertions


Here are some examples on how you can use assertions in probes:
#### Basics

```gotemplate
{{ .Assert.True false }}
{{ .Assert.Equal 200 (.HTTP.Request "GET" "http://localhost:8080").StatusCode }}
{{ .Assert.Regexp "^start"|regexp "its not starting" }}
{{ .Assert.NotRegexp (call .Regexp "^start") "its not starting" }}
```

#### HTTP

```gotemplate
{{ .Assert.HTTPStatusCode .HTTP.Handler "GET" "http://google.com" nil }}
{{ .Assert.HTTPError .HTTP.Handler "GET" "http://localhost:8080" nil }}
{{ .Assert.HTTPSuccess .HTTP.Handler "GET" "http://localhost:8080" nil }}
{{ .Assert.HTTPBodyContains .HTTP.Handler "GET" "http://localhost:8080" nil }}
{{ .Assert.HTTPBodyNotContains .HTTP.Handler "GET" "http://localhost:8080" nil }}
```

## Credits

This tool is base on multiple awesome open source libraries shout-outs to:
* https://github.com/stretchr/testify
* https://github.com/Masterminds/sprig
* https://github.com/urfave/cli

## License

```
MIT License

Copyright (c) 2022 Maximilian Reichel

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
