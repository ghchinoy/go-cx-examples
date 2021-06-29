package testcases

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	cx "cloud.google.com/go/dialogflow/cx/apiv3beta1"
	"google.golang.org/api/option"
	cxpb "google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3beta1"
)

func RunTestCase(ctx context.Context, projectID, location, agentID, testCaseID string) error {
	agent := fmt.Sprintf("projects/%s/locations/%s/agents/%s", projectID, location, agentID)

	var apiEndpoint string
	if location == "global" {
		apiEndpoint = "dialogflow.googleapis.com:443"
	} else {
		apiEndpoint = fmt.Sprintf("%s-dialogflow.googleapis.com:443", location)
	}

	var testCaseResults []*cxpb.TestCaseResult

	// TestCases client
	tcc, err := cx.NewTestCasesClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		log.Fatalf("Unable to obtain Dialogflow CX client: %v", err)
		return err
	}
	defer tcc.Close()

	testCaseName := fmt.Sprintf("%s/testCases/%s", agent, testCaseID)

	tc, err := GetTestCaseDetails(ctx, location, testCaseName)
	if err != nil {
		log.Printf("can't get test case details for %s: %s", testCaseName, err)
		return err
	}
	log.Printf("running test case '%s' (%s)", tc.DisplayName, testCaseName)

	runTestCaseRequest := &cxpb.RunTestCaseRequest{
		Name: testCaseName,
	}

	op, err := tcc.RunTestCase(ctx, runTestCaseRequest)
	if err != nil {
		log.Printf("error creating test case run request: %v", err)
		return err
	}

	result, err := op.Wait(ctx)
	if err != nil {
		log.Printf("error waiting test case run request: %v", err)
		return err
	}
	testCaseResults = append(testCaseResults, result.Result)

	for _, result := range testCaseResults {
		//log.Printf("%#v", resp)
		name := strings.Split(result.Name, "/")
		testCaseResultID := name[len(name)-1]
		log.Printf("Name: %s", testCaseResultID)
		env := "draft"
		if result.Environment != "" {
			env = result.Environment
		}
		log.Printf("Env: %s", env)
		log.Printf("Turns: %d", len(result.ConversationTurns))
		log.Printf("Result: %s", result.TestResult.String())
		log.Printf("Time: %s", result.TestTime.AsTime().Format(time.UnixDate))

	}

	/* single result
	data := [][]string{}
	for _, turn := range result.Result.ConversationTurns {
		data = append(data, []string{
			fmt.Sprintf("%s", turn.UserInput.Input.GetText()),
			fmt.Sprintf("%s", turn.VirtualAgentOutput.TextResponses),
		})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"User", "Agent"})
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(data)
	table.Render()
	*/

}

// GetTestCaseDetails returns a TestCase given an ID
func GetTestCaseDetails(ctx context.Context, location, testCaseName string) (*cxpb.TestCase, error) {
	var testcase *cxpb.TestCase
	var apiEndpoint string
	if location == "global" {
		apiEndpoint = "dialogflow.googleapis.com:443"
	} else {
		apiEndpoint = fmt.Sprintf("%s-dialogflow.googleapis.com:443", location)
	}

	// TestCases client
	tcc, err := cx.NewTestCasesClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		log.Fatalf("Unable to obtain Dialogflow CX client: %v", err)
		//return err
		os.Exit(1)
	}
	defer tcc.Close()

	req := &cxpb.GetTestCaseRequest{
		Name: testCaseName,
	}
	resp, err := tcc.GetTestCase(ctx, req)
	if err != nil {
		return testcase, err
	}
	// TODO: Use resp.
	return resp, nil
}
