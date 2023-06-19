# F2S
An Open Source Function as a Service (FaaS) Platform

Status <font color=red>Early Alpha</font>

# Content
- [F2S](#f2s)
- [Content](#content)
- [Core Concept](#core-concept)
  - [KISS - Simple and Stupid](#kiss---simple-and-stupid)
  - [All Features included](#all-features-included)
- [Architecture](#architecture)
  - [Namespaces](#namespaces)
  - [Gitops (CRDs Config)](#gitops-crds-config)
- [Configuration](#configuration)
  - [CRDs functions.f2s.opensight.ch](#crds-functionsf2sopensightch)
# Core Concept
Personally, I work on the project mainly to learn more Golang and because I have a usecase for a F2S Platform currently.

These will be the core concepts of this F2S Platform
## KISS - Simple and Stupid

* Keep it as simple as can be
* Run out of the box with as few dependencies as possible. <br/>
  No service meshes or other dependencies
* Simple start. Up and running in default config in 1 minute
* Lightweight. Use the features of vanilla kubernetes where ever possible
* Intuitive. No steep learning curve<br/>
  Beginners can use a UI to manage the soultion (i.e. create the CRD’s using the UI)

## All Features included

* No "enterprise only" features
* <font color=orange>TO DO</font> Kafka Message Bus Integration
* <font color=orange>TO DO</font> Scale to Zero
* <font color=orange>TO DO</font> Security (OAuth)
# Architecture
This is a first draft of the architecture and can still change.

![](docs/architecture.png)

## Namespaces
F2S uses 2 fixed namespaces in kubernetes
* **F2S**<br/>
contains the F2S operational components
* **F2S-Containers**<br/>
contains the running pods managed by F2S
## Gitops (CRDs Config)
F2SFunctions are managed by CRDs (bring your own Gitops)
We use a redundant setup of 2 F2S Pods. Metrics go to a prometheus instance

# Configuration
## CRDs functions.f2s.opensight.ch
Initial Datamodel is for testing and will certainly change
![](docs/datamodel.png)
