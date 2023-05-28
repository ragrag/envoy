<div align="center">
 
<img src="https://github.com/ragrag/envoy/assets/35541698/84ca78d8-ed2b-404f-b50d-63bcd4da6493" width="500" />
 
![envoi-steins-gate-machine](https://github.com/ragrag/envoy/assets/35541698/65307737-18d5-4e25-8830-a8acefde02ca)

</div>

## Table of Contents

- [About](#about)
- [Features](#features)
- [Getting Started](#getting-started)
- [API Reference](#api-reference)
- [Configuration](#api-reference)

## About

**Envoy** _(aan.voy)_ is a simple, performant and secure code execution engine designed from the ground up for running code as well as judging multiple test case submissions.
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

The recommanded way to deploy Envoy is via Docker. Envoy comes with a [prebuilt docker image on dockerhub](https://hub.docker.com/r/ragrag/envoy) that you can directly use. a customized Docker image can always be built from the provided Dockerfile in the repository.

:warning: For Envoy to function correctly, the Docker container must be run in [privileged mode](https://docs.docker.com/engine/reference/commandline/run/#privileged).

Running Docker in privileged mode gives all capabilities to the container, effectively disabling the security boundaries between the Docker container and the host system. This is necessary for [isolate](https://github.com/ioi/isolate) to run properly, which is the underlying sandboxing technology Envoy uses.

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

### Run Code

#### Request

`POST /run`

body:

- `language`: The id of the programming language.
- `code`: The code to be executed.
- `options` (optional): Options that can be provided to the execution environment
  - `timeLimit` (optional): Maximum time (in seconds) for the program execution, exceeding this limit will result in `TIME_LIMIT_EXCEEDED` status
  - `memoryLimit` (optional): Maximum memory (in KB) that the program can use, exceeding this limit will result in `MEMORY_LIMIT_EXCEEDED` status

```json
{
  "language": "javascript",
  "code": "console.log('Hello, world!');",
  "options": {
    "timeLimit": 2,
    "memoryLimit": 64000
  }
}
```

#### Response

- `output`: The output generated by the executed code.
- `status`: The status of the code execution, and can be one of:
  - `SUCCESS`: the code executed successfully without any errors
  - `COMPILATION_ERROR`: a compilation error occured
  - `TIME_LIMIT_EXCEEDED`: the program running time exceeded that set by the preconfigured limit or the provided limit in the request
  - `MEMORY_LIMIT_EXCEEDED`: the program used memory exceeded that set by the preconfigured limit or the provided limit in the request
  - `RUNTIME_ERROR`: a runtime error occured
- `time` (optional): The time (in seconds) the program took to execute, in case of `COMPILATION_ERROR` this is not provided
- `memory` (optional): The memory (in KB) the program used, in case of `COMPILATION_ERROR` this is not provided

```http
HTTP/1.1 200 OK
Status: 200 OK

{
    "output": "Hello, world!\n",
    "status": "SUCCESS",
    "time": 0.029,
    "memory": 11512
}
```

```http
HTTP/1.1 200 OK
Status: 200 OK

{
    "output": "error: cannot find macro `invalid_println` in this scope\n --> main.rs:1:13\n  |\n1 | fn main() { invalid_println!(\"Hello World!\"); }\n  |             ^^^^^^^^^^^^^^^\n\nerror: aborting due to previous error\n\n",
    "status": "COMPILATION_ERROR"
}
```

### Judge Code

#### Request

`POST /judge`

body:

- `language`: The id of the programming language.
- `code`: The code to be executed.
- `testCases`: An array containing test cases
  - `input`: input of the test case (stdin)
  - `expectedOutput`: expected output of the test case
- `options` (optional): Options that can be provided to the execution environment
  - `timeLimit` (optional): Maximum time (in seconds) for the program execution, exceeding this limit will result in `TIME_LIMIT_EXCEEDED` status
  - `memoryLimit` (optional): Maximum memory (in KB) that the program can use, exceeding this limit will result in `MEMORY_LIMIT_EXCEEDED` status

```json
{
  "language": "cpp",
  "code": "#include <iostream>; \n using namespace std; int main(){string i; cin>>i;cout<< i<<endl; return 0;}",
  "testCases": [
    {
      "input": "1",
      "expectedOutput": "1\n"
    },
    {
      "input": "2",
      "expectedOutput": "2\n"
    },
    {
      "input": "3",
      "expectedOutput": "3\n"
    },
    {
      "input": "4",
      "expectedOutput": "4\n"
    }
  ]
}
```

#### Response

- `verdict`: object containing the result of judging the submission
  - `status`: can be one of:
    - `SUCCESS`: all test cases passed
    - `COMPILATION_ERROR`: a compilation error occured
    - `TIME_LIMIT_EXCEEDED`: the program running time exceeded that set by the preconfigured limit or the provided limit in the request
    - `MEMORY_LIMIT_EXCEEDED`: the program used memory exceeded that set by the preconfigured limit or the provided limit in the request
    - `RUNTIME_ERROR`: a runtime error occured
    - `WRONG_ANSWER`: a runtime error occured
  - `totalTime (optional)`: The aggregated total time (in seconds) took to by the program while running test cases, in case of `COMPILATION_ERROR` this is not provided
  - `totalMemory (optional)`: The aggregated total memory (in KB) used by the program while running all test cases, in case of `COMPILATION_ERROR` this is not provided
  - `input (optional)`: The input (stdin) for the failed test case, only provided on `WRONG_ANSWER`
  - `output (optional)`: The actual output (stdout) by the program for the failed test case, only provided on `WRONG_ANSWER`
  - `expectedOutput (optional)`: The expected output (stdout) for the failed test case, only provided on `WRONG_ANSWER`
- `results`: an array containing the result of each test case, empty on `COMPILATION_ERROR`
  - `status`: can be one of:
    - `SUCCESS`: all test cases passed
    - `TIME_LIMIT_EXCEEDED`: the program running time for the test case exceeded that set by the preconfigured limit or the provided limit in the request
    - `MEMORY_LIMIT_EXCEEDED`: the program used memory for the test case exceeded that set by the preconfigured limit or the provided limit in the request
    - `RUNTIME_ERROR`: a runtime error occured
    - `WRONG_ANSWER`: a runtime error occured
  - `time`: The time (in seconds) the program took to execute the test case
  - `memory`: The memory (in KB) the program used to execute the test case
  - `input`: The input (stdin) for the test case
  - `expectedOutput`: The expected output (stdout) for the test case
  - `output (optional)`: The actual output (stdout) by the program for the test case

```http
HTTP/1.1 200 OK
Status: 200 OK

{
    "verdict": {
        "status": "SUCCESS",
        "totalTime": 0.006,
        "totalMemory": 5432
    },
    "results": [
        {
            "output": "1\n",
            "expectedOutput": "1\n",
            "input": "1",
            "status": "SUCCESS",
            "time": 0.002,
            "memory": 1120
        },
        {
            "output": "2\n",
            "expectedOutput": "2\n",
            "input": "2",
            "status": "SUCCESS",
            "time": 0.002,
            "memory": 1292
        },
        {
            "output": "3\n",
            "expectedOutput": "3\n",
            "input": "3",
            "status": "SUCCESS",
            "time": 0.001,
            "memory": 1412
        }
    ]
}
```

```http
HTTP/1.1 200 OK
Status: 200 OK

{
    "verdict": {
        "output": "2\n",
        "expectedOutput": "1\n",
        "input": "2",
        "status": "WRONG_ANSWER",
        "totalTime": 0.006,
        "totalMemory": 4956
    },
    "results": [
        {
            "output": "2\n",
            "expectedOutput": "1\n",
            "input": "2",
            "status": "WRONG_ANSWER",
            "time": 0.002,
            "memory": 932
        },
        {
            "output": "2\n",
            "expectedOutput": "2\n",
            "input": "2",
            "status": "SUCCESS",
            "time": 0.002,
            "memory": 1156
        },
        {
            "output": "3\n",
            "expectedOutput": "3\n",
            "input": "3",
            "status": "SUCCESS",
            "time": 0.001,
            "memory": 1336
        }
    ]
}
```

## Configuration

All configurations are set with environment variables. Below is a table of all available configuration parameters:

| Environment Variable           | Description                                                                                                                    | Default Value |
| :----------------------------- | :----------------------------------------------------------------------------------------------------------------------------- | :------------ |
| `LOG_LEVEL`                    | log level. Available options are trace, debug, info, warn, error, fatal, panic                                                 | `info`        |
| `SERVER_PORT`                  | port number that the server will use                                                                                           | `4000`        |
| `SERVER_AUTH_TOKEN` (optional) | If set, all requests must include this token in the Authorization header (Bearer token).                                       | `None`        |
| `ENGINE_WORKER_COUNT`          | This defines the maximum number of concurrent workers running code or juding submissions. The value must be between 1 and 999. | `100`         |
| `ENGINE_TIME_LIMIT_SECONDS`    | The default time limit (in seconds) for a single program run or a single test case run                                         | `2`           |
| `ENGINE_MEMORY_LIMIT_KB`       | The default memory limit (in KB) for a single program run or a single test case run                                            | `128000`      |
