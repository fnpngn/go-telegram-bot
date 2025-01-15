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

func runTests() error {
	// Run `go test ./...` and capture the output
	cmd := exec.Command("go", "test", "./...")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err // Return the error if `go test` fails
	}

	// Parse the results
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ok") || strings.HasPrefix(line, "FAIL") {
			parts := strings.Fields(line)
			passed := strings.HasPrefix(line, "ok")

			var output, packageStr string
			if len(parts) > 2 {
				packageStr = parts[1]
				output = strings.Join(parts[2:], " ")
			} else {
				packageStr = parts[0]
				output = line
			}
			testResults = append(testResults, TestResult{
				Package: packageStr,
				Output:  output,
				Passed:  passed,
			})
		}
	}
	return nil
}

func serveResults(w http.ResponseWriter, r *http.Request) {
	// Serve the captured test results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testResults)
}

func main() {
	// Run tests before starting the server
	if err := runTests(); err != nil {
		panic(err) // Stop execution if tests fail to run
	}

	// Start an HTTP server to serve the results
	http.HandleFunc("/results", serveResults)

	println("Serving test results on http://localhost:8080/results")
	http.ListenAndServe(":8080", nil)
}
