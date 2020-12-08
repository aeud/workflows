# Atom Workflows

## Data Structure
![Data Structure](doc/data-structures.jpg "Data Structure")


## [DEV] Deployment
```
PROJECT_ID=grp-sta-atom-prj-aelab
docker build -t workflow .
docker tag workflow eu.gcr.io/$PROJECT_ID/images/workflow          
docker push eu.gcr.io/$PROJECT_ID/images/workflow
```
