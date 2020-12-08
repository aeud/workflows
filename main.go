package main

import (
	"flag"
	"log"
	"os"
	"time"
	"workflows/flows"
)

var (
	flagDAGPath            string
	flagTaskRunnerHostname string
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
	flag.StringVar(&flagDAGPath, "dag-path", "", "Path to the DAG you want to execute")
	flag.StringVar(&flagTaskRunnerHostname, "tr-hostname", "", "Hostname of the Task Runner")

	flag.Parse()

	log.Printf("Using the Task Runner hostname: %s", flagTaskRunnerHostname)

	os.Setenv("TASK_RUNNER_HOSTNAME", flagTaskRunnerHostname)

	if flagDAGPath == "" {
		log.Fatal("dag-path is mandatory")
	}
	if flagTaskRunnerHostname == "" {
		log.Fatal("tr-hostname is mandatory")
	}
}

func main() {
	start := time.Now()
	if err := exec(flagDAGPath); err != nil {
		log.Fatalf("Flow %s failed in %s. %s.", flagDAGPath, time.Now().Sub(start), err)
	}
	log.Printf("Flow %s ran with success in %s.", flagDAGPath, time.Now().Sub(start))
}
