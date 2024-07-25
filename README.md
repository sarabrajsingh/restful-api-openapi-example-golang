# API developed by Sarabraj Singh for the take-home assignment.

## Table of Contents

- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
  - [Implementation](#implementation)
  - [Endpoints](#endpoints)
    - [Home Page](#home-page)
    - [Post Data Object](#post-temp)
    - [Get Errors](#get-errors)
    - [Delete Errors](#delete-errors)
- [OpenAPI Specification](#openapi-specification)
- [Testing](#testing)
- [Deployment](#deployment)
- [Application Logs](#application-logs)

## Project Structure

```bash
.
├── api # contains the openapi contract
├── config # contains the configuration helper package
├── internal
│   ├── global_errors # contains the in-memory global error handler for the API
│   ├── handlers # contains the API handlers
│   ├── logging # contains the custom logging stack
│   ├── models # contains the data models used in the API
│   ├── server # contains the API server implementation
│   └── utils # contains helpful utils that I developed when creating this API
└── swaggerui # contains the OpenAPI Swagger Frontend UI
│   └── dist
└── main.go
```

## Getting Started
Ensure that you have at least Golang `v1.21.6` installed locally, as it the minimim golang version needed to compile and execute this project locally. From here on out, I will be providing command line instructions for RHEL/Fedora based operating systems with the `dnf` package manager.


```bash
sudo dnf install golang
sudo dnf install make
```

If you plan on running the binary locally, [Docker](https://www.docker.com/) is an excellent choice to deploy and run the binary. Take a look at their documentation for installation instructions. The [Deployment](#deployment) steps in this documentation assume that the runtime user has docker properly installed.

## Running The Application

Download the necessary project dependencies

```bash
go mod download
```

Compile the application
```bash
go build -o ./app-api-server .
```

Run the API server. The default port that the server exposes itself on is `8080`.

```bash
$ ./app-api-server 
INFO: 2024/07/27 01:59:12 logger.go:20: Server started
INFO: 2024/07/27 01:59:12 logger.go:20: Validating Contract
INFO: 2024/07/27 01:59:12 logger.go:20: Successfully validated the contract
INFO: 2024/07/27 01:59:12 logger.go:20: API server is running on port :8080

```
## API Documentation

### Implementation

The core of the API was generated from the [OpenAPI](api/openapi.yaml) contract via the Swagger UI, and its functionalities extended through custom code. A mock-first approach was used to develop the rest of the API and some middleware layers were added. Please refer to the diagram below.

![High-Level Architecture](docs/architecture.drawio.svg)

- The request-to-response stack has two middleware layers:
  - A logging middleware layer that logs every interaction with an endpoint handler
  - An OpenAPI validation middleware layer that:
    - Verifies if an incoming request has a valid endpoint in the API server as defined by the contract
    - Ensures the request parameters satisfy the specification in the OpenAPI contract

If there is an incoming request that the router does not recognize, the router will drop the request on the ground, and return a `404`.

```bash
$ curl -X GET https://localhost:8080/api/v1/foobar

404 page not found
```

### Endpoints

All endpoints are prefixed with `/api/v1` as per canonical norms for API versioning.

### GET /

**Summary**: Home Page

**Description**: Serves the Swagger documentation.

**Responses**:
- **200 OK**: Returns the home page content.

**Example Response**:
```html
<!DOCTYPE html>
<html lang="en">
  // some HTML Content
</html>
```

### POST /temp

**Summary**: An endpoint that accepts a user request for interpretation

**Description**: This endpoint validates JSON blobs being sent from a client. There are two possible good responses, which are determined based on the JSON blurb coming into the endpoint. If the request is invalid, or fails the internal validation machinery, an error is returned.

**Responses**:
- **200 OK**: Returns either an `overtemp: true` response or a `overtemp: false` response.
- **400 BAD REQUEST**: The client provided an invalid request to this endpoint.

**Expected Headers**:
- `Content-Type: application/json`

**Example Request Body**:
```json
{
  "data": "365951380:1722089835:'Temperature':98.48256793121914"
}
```

**Example Response (Overtemp=true)**:
```json
{
  "overtemp": true,
  "device_id": 365951380,
  "formatted_time": "2024/07/27 14:17:15"
}
```
**Example Response (Overtemp=false)**:
```json
{
  "overtemp": false
}
```

**Example Curl Request**:

The application is hosted in Google's App Engine platform and requires the following headers being provided with the payload.

- `Content-Type: application/json`

Your client should calculate `Content-Length` and the `Host` headers as well, which the middleware and Google's App Engine expect.

##### Request:
```bash
$ curl -X POST --location 'https://localhost:8080/api/v1/temp' \
--header 'Content-Type: application/json' \
--data '{"data": "365951380:1722089835:'\''Temperature'\'':98.48256793121914"}'
```
##### Response:
```json
{
  "overtemp": true,
  "device_id": 365951380,
  "formatted_time": "2024/07/27 14:17:15"
}
```

Here is another example of a possible return response from the API server when the temperature is below 90 degrees celcius.

##### Request:
```bash
$ curl -X POST --location 'https://localhost:8080/api/v1/temp' \
--header 'Content-Type: application/json' \
--data '{"data": "365951380:1722089835:'\''Temperature'\'':89.48256793121914"}'
```
##### Response:
```json
{
  "overtemp": false
}
```
If there are errors in the payload or the validation of the payload fails past the middleware machinery, a `400` will be returned and the entry will get stored in the in-memory errors array within the API server. The array of errors can be queried by issuing a `GET` on the `/api/v1/errors` endpoint.

##### Request:
```bash
$ curl -X POST --location 'https://localhost:8080/api/v1/temp' \
--header 'Content-Type: application/json' \
--data '{"data": "365951380:1722089835:'\''Foobar'\'':89.48256793121914"}'
```
##### Response:
```json
{
  "error": "bad request"
}
```
Let's know peek at the errors.
##### Request:
```bash
$ curl -X GET --location 'https://localhost:8080/api/v1/errors'
```
##### Response:
```json
{
  "errors": [
    "365951380:1722089835:'Foobar':89.48256793121914"
  ]
}

```

### GET /errors

**Summary**: An endpoint that returns a list of errors that are known by the API server.

**Description**: The errors are collected and stored into an in-memory buffer within the API server. Mutual exclusion is used to add objects to the errors array with a capacity of 512. There is overflow protection, such that if the errors array is full, the array gets reset with the last element as the first element in the array.

**Responses**:
- **200 OK**: Returns the in-memory errors array from the API server.

**Example Curl Request**:
##### Request:
```bash
$ curl -X GET --location 'https://localhost:8080/api/v1/errors'
```
##### Response:
```json
{
  "errors": []
}
```
If there are errors in the in-memory array, the response will look like this:
##### Response that shows API errors:
```json
{
  "errors": [
    "365951380:1722089835:'Temperaure':89.48256793121914",
    "365951380:1722089835:'Foobar':89.48256793121914",
    "365951380:1722089835:'Foobar':89.48256793121914",
    "365951380:1722089835:'Foobar':89.48256793121914",
    "not_a_device_id:1722089835:'Temperature':89.48256793121914"
  ]
}
```

### DELETE /errors

**Summary**: An endpoint that clears the in-memory error array in the API server.

**Description**: This endpoint flushes the in-memory error array in the API server. This is a destructive command and is guarenteed to always work if the `DELETE` gets received to the `/errors` endpoint.

**Responses**:
- **200 OK**

**Example Curl Request**:

#### Errors Array Before `DELETE`:
```bash
$ curl -X GET --location 'https://localhost:8080/api/v1/errors'

{
    "errors": [
        "365951380:1722089835:'Temperaure':89.48256793121914",
        "365951380:1722089835:'Foobar':89.48256793121914",
        "365951380:1722089835:'Foobar':89.48256793121914",
        "365951380:1722089835:'Foobar':89.48256793121914",
        "not_a_device_id:1722089835:'Temperature':89.48256793121914"
    ]
}
```

##### Request:
```bash
$ curl --location --request DELETE 'https://localhost:8080/api/v1/errors'

# 200 OK received
```

Now lets check the errors endpoint to see if the buffer got flushed.
#### Errors Array After `DELETE`:
```bash
$ curl -X GET --location 'https://localhost:8080/api/v1/errors'

{
    "errors": []
}
```
## OpenAPI Specification

The [homepage](https://localhost:8080/) of the application displays the `openapi` contract through the Swagger UI for easy visualization.

The middleware layer is pretty powerful in ensuring that payloads for the `POST` endpoint at `/api/v1/temp` are enforced. This allows the API server to really use the [OpenAPI](./api/openapi.yaml) contract as the source of truth for the API. Lets explore an example of sending a JSON blurb that doesn't conform to the expected request body as defined in the contract for the `POST` endpoint.

```bash
curl -X POST --location 'https://localhost:8080/api/v1/temp' \
--header 'Content-Type: application/json' \
--data '{"foobar": "36595138029567120956:1722089835:'\''Temperature'\'':89.48256793121914"}'
```

In the example above, we swapped the `data` key for `foobar`. This request is to a validate endpoint, which will pass the router validation steps. However, when the middleware examines the request, it will be able to determine that the `foobar` key is invalid and that we only accept `data` as the payload key. Here's the response from the API server, which returns a middleware message:

```json
{
  "error": "OpenAPI Middleware: Request validation failed: request body has an error: doesn't match schema #/components/schemas/TempPostBody: Error at \"/data\": property \"data\" is missing\nSchema:\n  {\n    \"properties\": {\n      \"data\": {\n        \"example\": \"365951380:1640995229697:'Temperature':58.48256793121914\",\n        \"type\": \"string\"\n      }\n    },\n    \"required\": [\n      \"data\"\n    ],\n    \"type\": \"object\"\n  }\n\nValue:\n  {\n    \"foobar\": \"36595138029567120956:1722089835:'Temperature':89.48256793121914\"\n  }\n\n"
}
```

The middleware message is a bit hard to read with human eyes but it explicitly states the violation and remedation steps.

## Testing

Mocks are heavily used in the unit and integration tests for this project. Ensure that the mocks are generated before running the test suites.

The `make` build tool is used to generate the mocks.

```bash
make mockgen
```
After the mocks have been generated, you can run the `go` test suite. Run this command from the root of the project directory.
```bash
go test ./...
```
## Deployment

Please take a look at the [Dockerfile](./Dockerfile) if you want to see how the container image is built. Docker was chosen as the container engine because it is portable and supported by many deployment methods.

##### Build the Docker Container
```bash
docker build -t test:latest .
```

##### Depoy the Docker Container
```bash
docker run -d -p 8080:8080 test:latest
```

You will now be able to access the API at `localhost:8080`.

##### Example Request at `localhost:8080`:
```bash
$ curl --location 'http://localhost:8080/api/v1/temp' --header 'Content-Type: application/json' --data '{"data": "12345678:1722089835:'\''Temperature'\'':99.48256793121914"}' | jq

# response
{
  "overtemp": true,
  "device_id": 12345678,
  "formatted_time": "2024/07/27 14:17:15"
}
```

## Application Logs

#### Example Payload with Strange Values:
##### Scenario 1: Device ID is supposed to be an `int32` but lets give it a number larger than an `int32`:
```bash
$ curl --location 'http://localhost:8080/api/v1/temp' \
--header 'Content-Type: application/json' \
--data '{"data": "36595138029567120956:1722089835:'\''Temperature'\'':89.48256793121914"}'
```
###### Application Logs:
```bash
docker logs <container_runtime_id> #for example 843897f07326

# truncated output
INFO: 2024/07/27 15:26:35 logger.go:20: POST	/api/v1/temp	TempPost	172.17.0.1:53002
INFO: 2024/07/27 15:26:35 logger.go:24: POST /api/v1/temp - Malformed data string received. Error:  could not parse device_id=36595138029567120956 to an int32
INFO: 2024/07/27 15:26:35 logger.go:20: appending [36595138029567120956:1722089835:'Temperature':89.48256793121914] to errorBuffer
```