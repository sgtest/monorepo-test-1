/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tcp

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sourcegraph/monorepo-test-1/kubernetes-4/pkg/probe"
)

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

func TestTcpHealthChecker(t *testing.T) {
	// Setup a test server that responds to probing correctly
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	tHost, tPortStr, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	tPort, err := strconv.Atoi(tPortStr)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tests := []struct {
		host string
		port int

		expectedStatus probe.Result
		expectedError  error
		// Some errors are different depending on your system.
		// The test passes as long as the output matches one of them.
		expectedOutputs []string
	}{
		// A connection is made and probing would succeed
		{tHost, tPort, probe.Success, nil, []string{""}},
		// No connection can be made and probing would fail
		{tHost, -1, probe.Failure, nil, []string{
			"unknown port",
			"Servname not supported for ai_socktype",
			"nodename nor servname provided, or not known",
			"dial tcp: invalid port",
		}},
	}

	prober := New()
	for i, tt := range tests {
		status, output, err := prober.Probe(tt.host, tt.port, 1*time.Second)
		if status != tt.expectedStatus {
			t.Errorf("#%d: expected status=%v, get=%v", i, tt.expectedStatus, status)
		}
		if err != tt.expectedError {
			t.Errorf("#%d: expected error=%v, get=%v", i, tt.expectedError, err)
		}
		if !containsAny(output, tt.expectedOutputs) {
			t.Errorf("#%d: expected output=one of %#v, get=%s", i, tt.expectedOutputs, output)
		}
	}
}
