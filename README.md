<div align="center">
 
<img src="https://github.com/ragrag/envoy/assets/35541698/84ca78d8-ed2b-404f-b50d-63bcd4da6493" width="500" />
 
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

### Supported Languages
This is the current list of supported languages, contributions to add or request more language support is always welcome 

`c`,
`c++`,
`csharp`,
`elixir`,
`erlang`,
`go`,
`haskell`,
`java`,
`javascript`,
`kotlin`,
`php`,
`python2`,
`python3`,
`rust`,
`scala`,
`swift`,
`typescript`,
`zig`

## Getting Started

### Deploying
The recommanded way to deploy Envoy is via Docker. Envoy comes with a  [prebuilt docker image on dockerhub](https://hub.docker.com/r/ragrag/envoy) that you can directly use. a customized Docker image can always be built from the provided Dockerfile in the repository.

:warning: For Envoy to function correctly, the Docker container must be run in [privileged mode](https://docs.docker.com/engine/reference/commandline/run/#privileged).
 
Running Docker in privileged mode gives all capabilities to the container, effectively disabling the security boundaries between the Docker container and the host system. This is necessary for [isolate](https://github.com/ioi/isolate) to run properly, which is the underlying sandboxing technology Envoy uses (more on that below)

How to enable privileged mode depends on the environment where you are running your Docker image:
- For running a single Docker container using the Docker CLI, you can use the `--privileged` flag. 
- When deploying in a Kubernetes environment, you can set `privileged: true` in the `securityContext` section of your pod specification.
- For cloud hosted providers, the method to enable privileged mode may vary, and you should consult the specific provider's documentation.


### Running Locally

To start local development with Envoy, it is recommended to use the provided VSCode devcontainer. This provides a complete, pre-configured development environment with all the necessary dependencies for Envoy.

Before you start, ensure the following prerequisites are met:

1. **Install Docker Desktop:** VSCode devcontainers require Docker to operate. [Docker Desktop](https://www.docker.com/products/docker-desktop/) is the easiest way to get Docker on your machine and it must be installed and running before proceeding.

2. **Prepare the .env file:** Copy the `.env.example` file to a new file named `.env` and change the configuration values if needed.

With the prerequisites met, you can now run Envoy:

1. Open the project in VSCode. 
2. VSCode should automatically suggest opening the project in a devcontainer. If not, you can manually launch the devcontainer by clicking on the green '><' button in the bottom-left corner and selecting 'Remote-Containers: Open Folder in Container...' or doing the same via the command-pallete.
3. Once the devcontainer is running, open a terminal in VSCode and run Envoy with the following command:

    ```bash
    go run cmd/main.go
    ```

## API Reference

### Authorization

By default, all API requests are publicly accessible.

Bearer Authentication can be added on all requests by providing an auth token in the server configuration by setting the `SERVER_AUTH_TOKEN` environment variable, as mentioned in the [Configuration](#configuration) section. 

Once an auth token is set, it must be included as a Bearer token in the Authorization header of all requests.

Example of a request with the Authorization header:

```shell
curl -i -H 'Accept: application/json' -H 'Authorization: Bearer <your_token_here>' http://localhost:8080/runtimes
```

### Get List of Available Runtimes

#### Request

`GET /runtimes`

#### Response
An array of available programming languages that can be used
- `id`: the id of the runtime, this is a unique identifier for a language and is used for other requests to reference that language
- `language`: the programming language
- `version`: the current version of the language

```http
HTTP/1.1 200 OK
Status: 200 OK

[
    {
        "id": "go",
        "language": "Go",
        "version": "1.20.3"
    },
    {
        "id": "javascript",
        "language": "JavaScript",
        "version": "18.16.0 (Node.js)"
    },
    {
        "id": "rust",
        "language": "Rust",
        "version": "1.69.0"
    }
]
```
