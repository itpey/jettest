[//]: # "Title: jettest"
[//]: # "Author: itpey"
[//]: # "Attendees: itpey"
[//]: # "Tags: #itpey #go #test #golang #go-lang #cli #api #http #jettest"

<p align="center">
<img alt= "itpey JetTest jettest" src="https://raw.githubusercontent.com/itpey/jettest/main/static/images/jet_test_icon.png" width="200" height="200" border="10"/>
</p>

<h1 align="center">
 JetTest
</h1>

<p align="center">
JetTest is a versatile command-line tool designed for effortless API testing using YAML-based configuration files.
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/itpey/jettest">
    <img src="https://pkg.go.dev/badge/github.com/itpey/jettest.svg" alt="itpey JetTest Go Reference">
  </a>
  <a href="https://github.com/itpey/jettest/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/itpey/jettest" alt="itpey JetTest license">
  </a>
</p>

# Features

- **Simple YAML Configuration**: Define API tests easily using YAML files.
- **Flexible HTTP Methods**: Supports GET, POST, and PUT requests.
- **Customizable Parameters**: Specify API host, client ID, authentication token, and request timeout.
- **Debug Mode**: Enable detailed logging for request and response information.
- **Response Validation**: Validate expected status codes, response latency, and specific response body segments.

# Installation

Make sure you have Go installed and configured on your system. Use go install to install JetTest:

```bash
go install github.com/itpey/jettest@latest
```

Ensure that your `GOBIN` directory is in your `PATH` for the installed binary to be accessible globally.

# Usage

```bash
jettest --host <API_HOST> --file <YAML_TEST_FILE> [OPTIONS]
```

### Options

- `--host <API_HOST>`: Specify the API host (required).
- `--file, -f <YAML_TEST_FILE>`: Path to the YAML test file (required).
- `--clientID, --cid <CLIENT_ID>`: Set the client ID for requests.
- `--authToken, --at <AUTH_TOKEN>`: Provide an authentication token for requests.
- `--timeout, -t <SECONDS>`: Set the request timeout duration in seconds (default: 30).
- `--debug, -d`: Enable debug mode to print detailed request and response information.

## Example

Execute JetTest to run API tests against a target host using a specific YAML configuration file:

```bash
jettest --host https://api.example.com --file file.yaml --clientID myclient --authToken mytoken --timeout 60 --debug
```

# YAML Test Configuration

Define API tests in a YAML configuration file `file.yaml` with the following structure:

```yaml
tests:
  - name: Get user profile
    description: Retrieve user profile information
    request:
      method: GET
      path: /api/users/profile
      params:
        user_id: [12345]
      with_token: true
    expect:
      status_code: 200
      max_latency: 5000ms
      body:
        - path: data.username
          value: "john_doe"
        - path: data.email
          value: "john.doe@example.com"
```

### Explanation

- `name`: Name of the API test.
- `request`: Details of the API request.
- `method`: HTTP method (GET, POST, PUT).
- `path`: Endpoint path.
- `params`: Query parameters.
- `headers`: Request headers.
- `with_client_id`: Include client ID in headers (true/false).
- `with_token`: Include authentication token in headers (true/false).
- `body`: Request body (if applicable).
- `expect`: Expected response conditions.
- `status_code`: Expected HTTP status code.
- `max_latency`: Maximum acceptable response latency.
- `body`: List of expected response body conditions.

For more examples and detailed usage, refer to the [examples](https://github.com/itpey/jettest/tree/main/examples) examples provided in the repository.

# Feedback and Contributions

If you encounter any issues or have suggestions for improvement, please [open an issue](https://github.com/itpey/jettest/issues) on GitHub.

We welcome contributions! Fork the repository, make your changes, and submit a pull request.

# License

JetTest is open-source software released under the Apache License, Version 2.0. You can find a copy of the license in the [LICENSE](https://github.com/itpey/jettest/blob/main/LICENSE) file.

# Author

JetTest was created by [itpey](https://github.com/itpey)
