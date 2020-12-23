package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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
	log.SetFlags(0)
	if engine := os.Getenv("ENGINE"); engine != "EDGE" {
		flag.StringVar(&flagWorkflowPath, "workflow-path", "", "Path to the Workflow you want to execute")
		flag.StringVar(&flagTaskRunnerURL, "tr-url", "", "URL of the Task Runner. Should look like `https://xxxxxxxx-yyyy-zz.a.run.app`")
		flag.StringVar(&flagTaskRunnerAuthJWT, "tr-auth-jwt", "", "JWT to use to connect to the Task Runner")

		flag.Parse()

		os.Setenv("TASK_RUNNER_URL", flagTaskRunnerURL)

		os.Setenv("WORKFLOW_TR_AUTH_JWT", flagTaskRunnerAuthJWT)

		if flagWorkflowPath == "" {
			log.Fatal("workflow-path is mandatory")
		}
		// Get the DAG file name from the path
		workflowPathSplit := strings.Split(flagWorkflowPath, "/")
		DAGfileName := workflowPathSplit[len(workflowPathSplit)-1]
		os.Setenv("DAG_FILE_NAME", strings.ToLower(regexp.MustCompile(`\.yaml`).ReplaceAllString(DAGfileName, "")))

		if flagTaskRunnerURL == "" {
			log.Fatal("`tr-url` parameter is mandatory")
		}
		log.Printf("Using the Task Runner url: %s", flagTaskRunnerURL)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	v := struct {
		TRURL        string `json:"tr-url"`
		WorkflowPath string `json:"workflow-path"`
	}{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		fmt.Fprintf(w, "Error when decoding the body %s!\n", err.Error())
		return
	}
	start := time.Now()
	flagWorkflowPath = v.WorkflowPath
	if err := exec(flagWorkflowPath); err != nil {
		fmt.Fprintf(w, "Flow %s failed in %s. %s.", flagWorkflowPath, time.Now().Sub(start), err)
		return
	}
	fmt.Fprintf(w, "Flow %s ran with success in %s.", flagWorkflowPath, time.Now().Sub(start))
}

func main() {
	if engine := os.Getenv("ENGINE"); engine == "EDGE" {
		log.Print("starting server...")
		http.HandleFunc("/", handler)

		// Determine port for HTTP service.
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
			log.Printf("defaulting to port %s", port)
		}

		// Start HTTP server.
		log.Printf("listening on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	} else {
		start := time.Now()
		if err := exec(flagWorkflowPath); err != nil {
			log.Fatalf("Flow %s failed in %s. %s.", flagWorkflowPath, time.Now().Sub(start), err)
		}
		log.Printf("Flow %s ran with success in %s.", flagWorkflowPath, time.Now().Sub(start))
	}
}
