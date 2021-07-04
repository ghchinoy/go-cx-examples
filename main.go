package main

import (
	"context"
	"os"

	"github.com/ghchinoy/cx-examples/testcases"
)

func main() {

	ctx := context.Background()

	projectID := os.Getenv("PROJECT_ID")
	location := os.Getenv("LOCATION")
	agentID := os.Getenv("AGENT_ID")
	testCaseID := os.Getenv("TEST_CASE_ID")

	err := testcases.RunTestCase(ctx, projectID, location, agentID, testCaseID)
	if err != nil {
		os.Exit(1)
	}
}
