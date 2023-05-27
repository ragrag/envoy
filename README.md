<div align="center">

<img src="https://github.com/ragrag/envoy/assets/35541698/84ca78d8-ed2b-404f-b50d-63bcd4da6493" width="500" />

---

![envoi-steins-gate-machine](https://github.com/ragrag/envoy/assets/35541698/65307737-18d5-4e25-8830-a8acefde02ca)

</div>


## Table of Contents
* [About](#about)
* [Features](#features)
* [Getting Started](#getting-started)
* [API Reference](#api-reference)
* [Configuration](#api-reference)
* [Contributing](#contributing)


## About
**Envoy** *(aan.voy)* is a simple, performant and secure code execution engine designed from the ground up for running code as well as judging multiple test case submissions.
Envoy provides a straight-to-the-point API surface, configuration with sensible defaults and out-of-the-box support for a multitude of programming languages.

## Features
- Secure sandboxed code execution
- Stupid simple configuration for time and memory limits
- Minimal straight to the point API Surface 
- Batched test case support

## Getting Started

### Deploying
The recommanded way to deploy Envoy is via Docker. Envoy comes with a  [prebuilt docker image on dockerhub](https://hub.docker.com/r/ragrag/envoy) that you can directly use. a customized Docker image can always be built from the provided Dockerfile in the repository.

:warning: For Envoy to function correctly, the Docker container must be run in [privileged mode](https://docs.docker.com/engine/reference/commandline/run/#privileged).
 
Running Docker in privileged mode gives all capabilities to the container, effectively disabling the security boundaries between the Docker container and the host system. This is necessary for [isolate](https://github.com/ioi/isolate) to run properly, which is the underlying sandboxing technology Envoy uses (more on that below)

How to enable privileged mode depends on the environment where you are running your Docker image:
- For running a single Docker container using the Docker CLI, you can use the `--privileged` flag. 
- When deploying in a Kubernetes environment, you can set `privileged: true` in the securityContext section of your pod specification.
- For cloud hosted providers, the method to enable privileged mode may vary, and you should consult the specific provider's documentation.



### Running Locally
