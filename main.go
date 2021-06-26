package main

import (
	"context"
	"os"

	"github.com/ghchinoy/cx-go-examples/testcases"
)

func main() {

	ctx := context.Background()

	location := "global"
	projectID := os.Getenv("PROJECT_ID")
	agentID := os.Getenv("AGENT_ID")
	testCaseID := os.Getenv("TEST_CASE_ID")

	testcases.RunTestCase(ctx, projectID, location, agentID, testCaseID)
}
