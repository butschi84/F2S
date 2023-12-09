# F2S
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![F2S FunctionsAsAService](https://img.shields.io/badge/f2s-Functions_as_a_Service-blue)](https://opensight.ch/licenses/MIT)
[![Beta](https://img.shields.io/badge/status-beta-green)](https://opensight.ch/licenses/MIT)

An Open Source Function as a Service (FaaS) Platform. <br />
[Get started in 5 Minutes](#quick-start-guide) on your Kubernetes platform.

![](docs/highlights.png)
## Quick Start Guide
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
## Features

* **Gitops**<br/>
  Define your Functions as K8S CRDs
* **Scale to Zero**<br/>
  Can scale Deployments to zero when there is no activity
* **Autoscaling**<br/>
  F2S continuously measures the performance of your containers and uses the data for autoscaling
* **Authentication**<br/>
  Right now, F2S Supports Token (JWT) Authentication and Basic Auth
  * None
  * Token (JWT)
  * Basic Auth
  * <font color=orange>TO DO</font> Security (OAuth)
* **Authorization**<br/>
  Authorization (RBAC)
* **Kafka**<br/>
  <font color=orange>TO DO</font> Kafka Message Bus Integration
## Core Concept

* Keep it as simple as can be
* Run out of the box with as few dependencies as possible. <br/>
  No service meshes or other dependencies
* Simple start. Up and running in default config in 1 minute
* Lightweight. Use the features of vanilla kubernetes where ever possible
* Intuitive. No steep learning curve<br/>
  Beginners can use a UI to manage the soultion (i.e. create the CRD’s using the UI)
* No "enterprise only" features

# Content
- [F2S](#f2s)
  - [Quick Start Guide](#quick-start-guide)
  - [Features](#features)
  - [Core Concept](#core-concept)
- [Content](#content)
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
    - [Authorization](#authorization)
    - [Kafka](#kafka)
    - [Debugging](#debugging)
- [Building Custom Functions](#building-custom-functions)
  - [Example - NodeJS](#example---nodejs)
  - [Example - Python](#example---python)

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
All Metrics go to the prometheus instance. Prometheus is used to collect the metrics of all f2s operators and source for scaling decisions.

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

**none**<br />
Just allow all requests. Authorization Controls are also disabled in Authentication Mode 'none'.

```
f2s:
  auth:
    global_config:
      type: none
```

**basic**<br />
HTTP basic auth.

```
f2s:
  auth:
    global_config:
      type: basic
    basic:
      - username: roman
        password: helloworld
        group: admins
```

**token**<br />
Authentication with a jwt bearer token.

```
f2s:
  auth:
    global_config:
      type: token
    token:
      tokens:
        - token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQyMzc4OTgsImdyb3VwIjoiZ3JvdXAxIiwic3ViIjoicm9tYW4ifQ.xQOtzG2cNa4eg97qidR-YN7v3qyJ18qjShWYLFUs_bU
      jwt_secret: test
```

### Authorization
Each user Account can be assigned to a group. Global Privileges are then assigned to the group via F2S Configmap

```
f2s:
  ...
  auth:
    ...
    authorization:
      - group: admins
        privileges:
          - functions:list
          - functions:invoke
          - functions:create
          - functions:delete
          - functions:update
          - settings:view
          - settings:update

```

### Kafka
F2S can receive Function Invocations from kafka. "Listeners" and "Actions" have to be declared in the kafka config block in the f2s configmap. Here is an example how to configure. The Action will react to kafka messages with key 'test-key' and invoke the f2sfunction with the specified uid.

```
kafka:
    enabled: true
    consumergroup: f2s-consumer
    brokers:
      - kafka-1.kafka
    listeners:
      - topic: test
        actions:
          - name: test-action
            triggers:
              - type: key
                filter: equal
                value: test-key
            f2sfunctions:
              - xxxxxxxxxx
            response:
              key: test-response
```

Fitler values will have to be one of `equal`, `not equal`, `contains` or `not contains`

### Debugging
Environment Variables take precedence over the configmap and are useful for local testing / debugging.

```
export Prometheus_URL=localhost:9090
export KUBECONFIG=~/.kube/config
```

# Building Custom Functions

Each Function that you like to run on F2S needs to be containerized and instrumented with a 'f2s fizzlet', a process that listens on 9092 and accepts incoming requests.<br/>

F2SFizzlet can run in 'reverseproxy mode' ([example - nodejs](#example---nodejs)) or in 'command line mode'.

We have some example containers that should get you started:
* [NodeJS](testing/container_nodejs_test/)
* [Python](testing/container_python_test/)

## Example - NodeJS

An example container for a NodeJS application is under [testing/container_nodejs_test](testing/container_nodejs_test/). If your NodeJS application listens on 8080, fizzlet will notice and switch to "reverseproxy mode" automatically. 

![](docs/fizzlet_reverseproxy_mode.png)

## Example - Python
An example for a Python application is under [testing/container_python_test](testing/container_python_test). When fizzlet receives a POST,PUT request, which is in 'application/json' Content-Type, it will write the request body to a 'input.json' and invoke the start.sh script ('start.sh input.json')

![](docs/fizzlet_cmd_mode.png)