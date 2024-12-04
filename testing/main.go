package main

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
)

type TestResult struct {
	Package string `json:"package"`
	Output  string `json:"output"`
	Passed  bool   `json:"passed"`
}

var testResults []TestResult

func runTests() {
	// Run `go test ./...` and capture the output
	cmd := exec.Command("go", "test", "./...")
	out, _ := cmd.CombinedOutput()

	// Parse the results
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ok") || strings.HasPrefix(line, "FAIL") {
			parts := strings.Fields(line)
			passed := strings.HasPrefix(line, "ok")
			testResults = append(testResults, TestResult{
				Package: parts[1],
				Output:  strings.Join(parts[2:], " "),
				Passed:  passed,
			})
		}
	}
}

func serveResults(w http.ResponseWriter, r *http.Request) {
	// Serve the captured test results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testResults)
}

func main() {
	// Run tests before starting the server
	runTests()

	// Start an HTTP server to serve the results
	http.HandleFunc("/results", serveResults)

	println("Serving test results on http://localhost:8080/results")
	http.ListenAndServe(":8080", nil)
}
