// Copyright 2024 itpey
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

// jet is a construct that lets us runs API tests with similar requirements
type jet struct {
	host      string
	clientID  string
	authToken string
	debug     bool
	client    *http.Client
	results   chan result
}

// jetOptions are the options to make a jet
type jetOptions struct {
	Host      string
	ClientID  string
	AuthToken string
	Timeout   uint
	Debug     bool
}

// results contains the count of passed, failed, and total tests
type results struct {
	Passed int
	Failed int
	Total  int
}

// result is the result from a test
type result struct {
	test   test
	errors []error
}

// test represents an API test
type test struct {
	Name    string  `yaml:"name"`
	Request request `yaml:"request"`
	Expect  expect  `yaml:"expect"`
}

// request contains everything required to make the test API request
type request struct {
	Method        string      `yaml:"method"`
	Path          string      `yaml:"path"`
	Params        url.Values  `yaml:"params"`
	Headers       http.Header `yaml:"headers"`
	WithClientID  bool        `yaml:"with_clientID"`
	WithAuthToken bool        `yaml:"with_token"`
	Body          string      `yaml:"body"`
}

// input represents the input configuration for running API tests.
type input struct {
	Tests []test `yaml:"tests"`
}

// Expect contains the expectations of the API response or behavior
type expect struct {
	StatusCode int           `yaml:"status_code"`
	MaxLatency time.Duration `yaml:"max_latnecy"`
	Body       []body        `yaml:"body"`
}

// body represents a segment of the response body that is expected
type body struct {
	Path          string `yaml:"path"`
	ExpectedValue string `yaml:"value"`
}

// main function for CLI

func Create() *cli.App {
	app := &cli.App{
		Name:        "JetTest",
		Usage:       "A command-line tool for testing APIs using YAML configuration",
		Version:     appVersion,
		Description: appDescription,
		Copyright:   appCopyright,
		Authors:     appAuthors,
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v", "ver", "about"},
				Usage:   "Print the version",
				Action: func(c *cli.Context) error {
					return showVersion()
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Usage:    "Specify the API host",
				EnvVars:  []string{"JETTEST_HOST"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Specify the path to the YAML test file",
				EnvVars:  []string{"JETTEST_FILE"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "clientID",
				Aliases: []string{"cid"},
				Usage:   "Set the client ID for requests",
				EnvVars: []string{"JETTEST_CLIENT_ID"},
			},
			&cli.StringFlag{
				Name:     "authToken",
				Aliases:  []string{"at"},
				Usage:    "Provide an authentication token for requests",
				EnvVars:  []string{"JETTEST_AUTH_TOKEN"},
				Required: false,
			},
			&cli.UintFlag{
				Name:     "timeout",
				Aliases:  []string{"t"},
				Usage:    "Set the request timeout duration in seconds (default: 30)",
				EnvVars:  []string{"JETTEST_TIMEOUT"},
				Required: false,
				Value:    30,
			},
			&cli.BoolFlag{
				Name:     "debug",
				Aliases:  []string{"d"},
				Usage:    "Enable debug mode (default: false)",
				EnvVars:  []string{"JETTEST_DEBUG"},
				Required: false,
				Value:    false,
			},
		},

		Action: func(c *cli.Context) error {
			fmt.Print(appNameArt)
			host := c.String("host")
			file := c.String("file")
			clientID := c.String("clientID")
			authToken := c.String("authToken")
			timeout := c.Uint("timeout")
			debug := c.Bool("debug")

			cfg, err := getConfiguration(file)
			if err != nil {
				return cli.Exit(err, 1)
			}

			jet := newJet(jetOptions{
				Host:      host,
				ClientID:  clientID,
				AuthToken: authToken,
				Debug:     debug,
				Timeout:   timeout,
			})
			defer jet.cleanup()

			results := jet.test(cfg.Tests...)

			fmt.Println()
			fmt.Println("test results:")
			fmt.Printf("total tests: %d\n", results.Total)
			fmt.Printf("tests passed: %d\n", results.Passed)
			fmt.Printf("tests failed: %d\n", results.Failed)
			fmt.Println()

			if results.Failed == 0 {
				fmt.Println("all tests passed successfully!")
			} else {
				fmt.Println("some tests failed.")
				if !debug {
					fmt.Println("hint: use '-d' for detailed request and response information.")
				}
			}

			return nil
		},
	}

	return app
}

// reads and parses the YAML configuration file
func getConfiguration(file string) (input, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return input{}, err
	}
	var cfg input
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return input{}, err
	}
	return cfg, nil
}

// creates a new Jet which allows us to run tests
func newJet(opt jetOptions) *jet {
	fmt.Printf("creating a new test jet for the host: %s\n", opt.Host)
	return &jet{
		host:      opt.Host,
		clientID:  opt.ClientID,
		authToken: opt.AuthToken,
		debug:     opt.Debug,
		client: &http.Client{
			Timeout: time.Second * time.Duration(opt.Timeout),
		},
		results: make(chan result),
	}
}

// test runs API tests in parallel and returns the results
func (b *jet) test(ts ...test) results {
	for _, t := range ts {
		go func(t test) {
			errs := b.run(t)
			b.results <- result{
				test:   t,
				errors: errs,
			}
		}(t)
	}
	timeout := 30 * time.Second
	results := results{Total: len(ts)}
	for i := 0; i < len(ts); i++ {
		t := time.NewTimer(timeout)
		select {
		case r := <-b.results:
			r.print(b.debug)
			if len(r.errors) == 0 {
				results.Passed = results.Passed + 1
			} else {
				results.Failed = results.Failed + 1
			}
		case <-t.C:
			panic(fmt.Sprintf("an unexpected delay has occurred. no result has been received for %v. only %v events have been processed.", timeout, i))
		}
	}
	return results
}

// cleans up the jet
func (b *jet) cleanup() {
	close(b.results)
}

// run makes the API call and checks the response
func (b *jet) run(t test) []error {
	if err := t.validateRequestMethod(); err != nil {
		return []error{err}
	}
	path := fmt.Sprintf("%s%s?%s", b.host, t.Request.Path, t.Request.Params.Encode())
	req, err := http.NewRequest(strings.ToUpper(t.Request.Method), path, strings.NewReader(t.Request.Body))
	if err != nil {
		return []error{err}
	}
	if t.Request.Headers != nil {
		req.Header = t.Request.Headers
	}
	if t.Request.WithClientID {
		req.Header.Add("client-id", b.clientID)
	}
	if t.Request.WithAuthToken {
		req.Header.Add("Authorization", b.authToken)
	}

	start := time.Now()
	resp, err := b.client.Do(req)
	if err != nil {
		return []error{err}
	}
	latency := time.Since(start)
	defer resp.Body.Close()
	return b.check(req, resp, latency, t.Expect)
}

// check runs through all the Expected response behaviors and records any that don't match
func (b *jet) check(req *http.Request, resp *http.Response, latency time.Duration, e expect) []error {
	errs := []error{}
	if resp.StatusCode != e.StatusCode {
		errs = append(errs, fmt.Errorf("status code failure: expected status code %v, but received %v", e.StatusCode, resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errs = append(errs, fmt.Errorf("body parsing failure: unable to read response body. skipping response body checks"))
		return errs
	}
	for _, b := range e.Body {
		err = b.checkResponseBody(string(body))
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 && b.debug {
		errs = append(errs, fmt.Errorf("body for debugging: %v", string(body)))
		errs = append(errs, fmt.Errorf("request for debugging: %+v", req))
		errs = append(errs, fmt.Errorf("response for debugging: %+v", resp))
	}
	if latency > e.MaxLatency {
		errs = append(errs, fmt.Errorf("latency failure: expected latency to be below %v, but it took %v", e.MaxLatency, latency))
	}
	return errs
}

// print prints the result to console but ensures that it will print one test's result at a time
func (r result) print(debug bool) {
	if len(r.errors) == 0 {
		if debug {
			fmt.Printf("Test Passed: %s (%s %s)\n", r.test.Name, strings.ToUpper(r.test.Request.Method), r.test.Request.Path)
		}
		return
	}
	fmt.Printf("Test Failed: %s (%s %s)\n", r.test.Name, strings.ToUpper(r.test.Request.Method), r.test.Request.Path)
	for _, e := range r.errors {
		fmt.Printf("\t- %s\n", e)
	}
}

// checks if the expected value matches the actual value in the response body
func (b body) checkResponseBody(body string) error {
	r := gjson.Get(body, b.Path)
	if b.ExpectedValue != r.String() {
		return fmt.Errorf("response body failure: expected %s to be %s, but got %s", b.Path, b.ExpectedValue, r.String())
	}
	return nil
}

// validates if the HTTP request method is supported
func (t test) validateRequestMethod() error {
	switch strings.ToUpper(t.Request.Method) {
	case "GET", "POST", "PUT":
		return nil
	}
	return fmt.Errorf("http method '%s' is not supported", t.Request.Method)
}
