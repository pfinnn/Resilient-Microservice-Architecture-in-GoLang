# Resilient MicroService

### A microservice architecture with resilience patterns and ToxiProxy

![image](https://drive.google.com/uc?export=view&id=12RFDgG_Jg21lF3adyCLSATOnGQ9QnDex)

## Features
- Two microservices communicating via HTTP
- Deployed in Docker containers, Orchestrated with kubernetes and running in the Google Cloud

- Resilience Patterns implemented: Retry, Bulkhead, Circuit Breaker
- Unit, Integration and Resilience Tests
- Continuous Integration Pipeline
- Toxi Proxy to simulate failures, resilience tests to prove the usefulness of the Resilience Patterns

## Wiki

[Introduction](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Main/001-Introduction)

[Specification](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Main/002-Specification)

[Solution Strategy](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Main/003-Solution-Strategy)

[Runtime View](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Main/004-Runtime-View)

[Quality Assurance](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Main/005-Quality-Assurance)


### Project Guidelines

[Goals Definition](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Project-Guidelines/Goals-Defintion)

[Planning and Progress Tracking](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Project-Guidelines/Planning-and-Progress-Tracking)

[Journal](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Journal)


### Research, Notes and Literature

<p>This section contains relevant links, images and descriptions that I found useful and worth saving</p>

<details>

<summary><strong>Research</strong></summary>

[Notes and Literature](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Research/00-Notes-and-Literature)

[Resilience](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Research/01-Resilience)

[Toxi Proxy](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Research/02-ToxiProxy)

[Circuit Breaker](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Research/03-Circuit-Breaker)

[Bulkhead](https://git.haw-hamburg.de/acm746/resilient-microservice/-/wikis/Research/04-Bulkhead)
</details>

## Requirements

### Go

You need Go installed

https://golang.org/doc/install

### Docker

You need Docker installed and running, when developing and testing this application locally

https://docs.docker.com/get-docker/

### Google Cloud SDK

https://cloud.google.com/sdk/install

### ToxiProxy

https://github.com/Shopify/toxiproxy/releases/tag/v2.1.4

If you are not under Windows, you have to install the ToxiProxy like so:
[text](https://github.com/Shopify/toxiproxy#1-installing-toxiproxy)
Make sure the ToxiProxy HTTP Server is running on your machine, when you are running in the test environment.

## Setup 
You need Go installed to compile this project

To run and test the application locally you need Docker installed and running

To deploy the application in the cloud, there need to be a cluster running

In Terminal of Goland in project root you can run this to start clusters with the right settings:

`make clusterUp`

and delete them with

`make clusterDown`

Use Port Forwarding to test services via Postman, when everything is deployed in the google cloud:

Run in Google Cloud SDK shell
```
kubectl port-forward service/forgeservice 8080:8080
kubectl port-forward service/smelterservice 8081:8081
```

## Resilience Tests
**Before running the resilience tests** make sure you have an instance of the toxiproxy running listening on the default port 8474

In the root of the project run following command in terminal to recursively run all integration and resilience tests:

`make test-global-resilience`´

alternatively you can run the tests in your IDE


## Google Cloud Platform (GCP)

### Preparation 

- set `GCP_PROJECTID` in `Makefile` (found on your GCP dashboard)
- set `GCP_CLUSTERID` in `Makefile` (choose something, this will be the name of cluster created by `$ make clusterUp`)
- set `GCP_PROJECTID`, `GCP_CLUSTER` and `GCP_GITLAB_TOKENFILE` in GitLab CI variables

### Deployment with GitLabCI

- create cluster with `$ make clusterUp`
- push zo gitlab / run pipeline
- check if pods are running with `$ kubectl get pods`
- portforwarding to receptionservice: `$ kubectl port-forward <name from previous cmd> 8080:8080`
- when done run `$ make undeploy` followed by `$ make clusterDown`