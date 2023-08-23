# F2S
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![F2S FunctionsAsAService](https://img.shields.io/badge/f2s-Functions_as_a_Service-blue)](https://opensight.ch/licenses/MIT)
[![Beta](https://img.shields.io/badge/status-beta-green)](https://opensight.ch/licenses/MIT)

An Open Source Function as a Service (FaaS) Platform. <br />
[Get started in 5 Minutes](#quick-start-guide) on your Kubernetes platform.

![](docs/highlights.png)

## Features

* **Gitops**<br/>
  Define your Functions as K8S CRDs
* **Scale to Zero**<br/>
  Can scale Deployments to zero when there is no activity
* **Authentication**<br/>
  Right now, F2S Supports Token (JWT) Authentication and Basic Auth

# Content
- [F2S](#f2s)
  - [Features](#features)
- [Content](#content)
- [Quick Start Guide](#quick-start-guide)
- [Core Concept](#core-concept)
  - [All Features included](#all-features-included)
- [Architecture](#architecture)
  - [Namespaces](#namespaces)
  - [Gitops (CRDs Config)](#gitops-crds-config)
  - [High Availability](#high-availability)
  - [Autoscaling](#autoscaling)
- [Configuration](#configuration)
  - [CRDs functions.f2s.opensight.ch](#crds-functionsf2sopensightch)
  - [Configmap - config.yaml](#configmap---configyaml)
    - [Timeouts](#timeouts)
    - [Authentication](#authentication)
    - [Debugging](#debugging)

# Quick Start Guide
This will install f2s on your kubernetes cluster. 

* "Functions" CRD
* Namespaces "f2s", "f2s-containers"
* ClusterRoles and Bindings
* Deployments for F2S, Grafana, Prometheus
* Services

```
helm repo add f2s https://butschi84.github.io/F2S/helm-release
helm repo update

# install crds
kubectl apply -f https://butschi84.github.io/F2S/helm-release/crds/crds.yaml

# install f2s
helm install f2s f2s/f2s
```

# Core Concept

* Keep it as simple as can be
* Run out of the box with as few dependencies as possible. <br/>
  No service meshes or other dependencies
* Simple start. Up and running in default config in 1 minute
* Lightweight. Use the features of vanilla kubernetes where ever possible
* Intuitive. No steep learning curve<br/>
  Beginners can use a UI to manage the soultion (i.e. create the CRD’s using the UI)

## All Features included
No "enterprise only" features
* <font color=orange>TO DO</font> Kafka Message Bus Integration
* Scale to Zero
* Authentication
  * None
  * Token (JWT)
  * Basic Auth
  * <font color=orange>TO DO</font> Security (OAuth)
* <font color=orange>TO DO</font> Authorization (RBAC)
# Architecture

![](docs/architecture.png)

* **API Server**<br/>
  'API Server' is the REST API Interface. Functions can be defined and invoked via API. Authentication and Authorization are part of 'API Server'. All function invocations are then sent to 'dispatcher'.
* **Metrics**<br/>
  F2S Exposes Metrics in order to be able to make scaling decisions and have insight in all activity
* **Config**<br/>
  The Config package observes CRD's (F2SFunction Declarations) in "f2s" namespace on kubernetes
* **Kafka**<br/>
  Package for Interaction with Kafka Event Bus (planned). All invocations are sent to dispatcher.
* **Operator**<br/>
  Operator reacts to Config Changes and creates or deletes deployments and services in f2s-containers namespace. It is also responsible for scaling functions up and down.
* **Dispatcher**<br/>
  Dispatcher picks up all incoming requests (REST, Kafka) and decides which function pod should serve the request.

## Namespaces
F2S uses 2 fixed namespaces in kubernetes
* **F2S**<br/>
contains the F2S operational components
* **F2S-Containers**<br/>
contains the running pods managed by F2S
## Gitops (CRDs Config)
F2SFunctions are managed by CRDs (bring your own Gitops)

## High Availability
We use a redundant setup of 2 F2S Pods. 

## Autoscaling
All Metrics go to the prometheus instance. The plan is to use it with autoscaler on k8s

# Configuration
## CRDs functions.f2s.opensight.ch
Initial Datamodel is for testing and will certainly change
![](docs/datamodel.png)

## Configmap - config.yaml

F2S Operator is configured by using [a configmap](helm/templates/f2s/f2s-configmap.yaml). Here is an example config for f2s with all parameters.

```
# enable debug output
debug: true

prometheus:
  url: prometheus-service.f2s:9090

f2s:
  timeouts:
    request_timeout: 120000
    http_timeout: 60000
    scaling_timeout: 45000
   auth:
    global_config:
      type: token
    basic:
      - username: roman
        password: helloworld
        group: group1
    token:
      tokens:
        - token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQyMzc4OTgsImdyb3VwIjoiZ3JvdXAxIiwic3ViIjoicm9tYW4ifQ.xQOtzG2cNa4eg97qidR-YN7v3qyJ18qjShWYLFUs_bU
      jwt_secret: test
```

### Timeouts
Timeouts can be configured in the configmap. F2S will abort those requests that exceed a timeout period.

* scaling_timeout<br />
  When zero replicas of a function are available, f2s will scale up the deployment from zero to one.
* http_timeout<br />
  How long should f2s wait at maximum for completion of a backend function
* request timeout<br />
  timeout for completion of the whole request


![](docs/timeouts.png)

### Authentication
Right now, F2S Supports the Authentication Modes:

* none<br />
  just allow all requests
* basic<br />
  http basic auth
* token<br />
  authentication with a jwt bearer token
### Debugging
Environment Variables take precedence over the configmap and are useful for local testing / debugging.

```
export Prometheus_URL=localhost:9090
export KUBECONFIG=~/.kube/config
```