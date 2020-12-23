# Atom Workflows

## Task
### Description
Workflow task is taking as inputs a Worklow representation (DAG) and a Cloud Run
hostname (scheme + hotname) and will execute it.

### Parameters
| Parameter | Description | Default Value | Example | Optional |
| ------------- |-------------|-----|-----|:-----:|
| --workflow-path | Path to the Workflow you want to execute |  |gs://your-bucket/workflows/your-workflow.yaml | N |
| --tr-url | URL of the Task Runner |  |https://taskrunnermock-rps3r5yvgq-ew.a.run.app | N |
| --tr-auth-jwt | JWT to use to connect to the Task Runner |  |xxx.yyy.zzz | Y |

### Environment variables
| Variable | Description | Default Value |
| ------------- |-------------|-----|-----|:-----:|
| WORKFLOW_CHECK_STATE_DURATION_SEC | Number of seconds to wait between 2 calls to the Task Runner Check | 10 |
| GOOGLE_APPLICATION_CREDENTIAL | Path to your Google Credential Key | - |

## Development environment
### How to use a local Mock for the Task Runner

#### Clone the repository
```
git clone https://github.com/lvmh-group/atom-container-kernel-task-runner.git
cd atom-container-kernel-task-runner
git checkout ae
```

#### Build the Docker image
```
docker build -t taskrunner:mock .
```

#### Run the sever
```
docker run \
    --rm -it \
    -e PORT=8080 \
    -v $HOME/.config:/root/.config \
    -p 8080:8080 \
    taskrunner:mock
```
`-v $HOME/.config:/root/.config` will let you share your own Gcloud authentification credentials.

#### Run your workflow task
```
go run main.go \
    --workflow-path=./dags/fork_join_dag.yaml \
    --tr-url=http://localhost:8080 \
    --tr-auth-jwt=none
```

### How to use a Mock for the Task Runner on Cloud Run
#### Build and deploy your image
```
export PROJECT_ID=grp-sta-lab-prj-cm
docker build -t task_runner:mock .
docker tag task_runner:mock eu.gcr.io/${PROJECT_ID}/run/task_runner:mock
docker push eu.gcr.io/${PROJECT_ID}/run/task_runner:mock
gcloud run deploy \
    --image eu.gcr.io/${PROJECT_ID}/run/task_runner:mock \
    --platform managed \
    --update-env-vars PROJECT_ID=${PROJECT_ID} \
    --project=${PROJECT_ID}
# it will return your endpoint (i.e. https://taskrunnermock-rps3r5yvgq-ew.a.run.app)
```
#### Run your workflow task
```
go run main.go \
    --workflow-path=./dags/fork_join_dag.yaml \
    --tr-url=https://taskrunnermock-rps3r5yvgq-ew.a.run.app \
    --tr-auth-jwt=$(gcloud auth print-identity-token)
```

## Internal documentation

### Data structure
![Data Structure](doc/data-structures.jpg "Data Structure")

### Limitations
#### Credentials
The Workflow task required a Service Account (and cannot use gcloud authentification). This limitation
comes from the `google.golang.org/api/idtoken` modules. You will get the following error message otherwise:
```
error message: idtoken: credential must be service_account, found "authorized_user".
```

The SA used for the Task should have the following permission
```
storage.objects.get
run.routes.invoke
```
Includes in the following roles :
```
roles/storage.objectViewer
roles/run.invoker
```

The ML service Agent already has the storage permission but need to have the Cloud Run one. The default SA used by the AI Platform is the following one: `service-<PROJECT_NUMBER>@cloud-ml.google.com.iam.gserviceaccount.com`

## [DEV] Deployment
```
PROJECT_ID=grp-sta-lab-prj-cm
docker build -t workflow .
docker tag workflow eu.gcr.io/$PROJECT_ID/images/workflow          
docker push eu.gcr.io/$PROJECT_ID/images/workflow
```

```
PROJECT_ID=grp-sta-lab-prj-cm
docker build -t atom_tasks_workflow -f ./Dockerfile_edge .
docker tag atom_tasks_workflow eu.gcr.io/$PROJECT_ID/edge/workflow
docker push eu.gcr.io/$PROJECT_ID/edge/workflow
gcloud run deploy task-edge-workflow \
    --image eu.gcr.io/$PROJECT_ID/edge/workflow \
    --platform managed \
    --project=$PROJECT_ID \
    --set-env-vars=TASK_RUNNER_URL=https://taskrunner-j7s6l5wnca-ew.a.run.app \
    --region=europe-west1 \
    --no-allow-unauthenticated
```