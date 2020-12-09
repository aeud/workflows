package main

import (
	"flag"
	"log"
	"os"
	"time"
	"workflows/flows"
)

var (
	flagWorkflowPath      string
	flagTaskRunnerURL     string
	flagTaskRunnerAuthJWT string
)

func exec(path string) error {
	flow, err := flows.NewFlowFromYAMLFile(path)
	if err != nil {
		return err
	}
	log.Printf("Running the following Flow:\n%s", flow.YAML())
	if err := flow.Walk(); err != nil {
		return err
	}
	return nil
}

func init() {
	flag.StringVar(&flagWorkflowPath, "workflow-path", "", "Path to the Workflow you want to execute")
	flag.StringVar(&flagTaskRunnerURL, "tr-url", "", "URL of the Task Runner. Should look like `https://xxxxxxxx-yyyy-zz.a.run.app`")
	flag.StringVar(&flagTaskRunnerAuthJWT, "tr-auth-jwt", "", "JWT to use to connect to the Task Runner")

	flag.Parse()

	log.Printf("Using the Task Runner url: %s", flagTaskRunnerURL)

	os.Setenv("TASK_RUNNER_URL", flagTaskRunnerURL)

	os.Setenv("WORKFLOW_TR_AUTH_JWT", flagTaskRunnerAuthJWT)

	if flagWorkflowPath == "" {
		log.Fatal("workflow-path is mandatory")
	}
	if flagTaskRunnerURL == "" {
		log.Fatal("`tr-url` parameter is mandatory")
	}
}

func main() {
	start := time.Now()
	if err := exec(flagWorkflowPath); err != nil {
		log.Fatalf("Flow %s failed in %s. %s.", flagWorkflowPath, time.Now().Sub(start), err)
	}
	log.Printf("Flow %s ran with success in %s.", flagWorkflowPath, time.Now().Sub(start))
}
