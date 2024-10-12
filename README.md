# CCJP - Another JSON Parser

This project is a custom implementation of a JSON parser. It was created as part of a coding challenge [here](https://codingchallenges.fyi/challenges/challenge-json-parser). The utility is designed to validate and parse JSON files, as well as run predefined tests on JSON datasets to ensure their correctness.

## Features

- Validate JSON files for syntax and structural correctness.
- Identify and report detailed errors in JSON files.
- Run predefined tests on JSON files to verify their validity.
- Efficiently parse and validate JSON from files up to 100MB in size.
- Command-line interface (CLI) powered by Cobra for easy usage.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- You need to have Go installed on your machine (Go 1.15 or later is recommended).
- You can download and install Go from [https://golang.org/dl/](https://golang.org/dl/).

### Installing

Clone the repository to your local machine:

```bash
git clone https://github.com/nullsploit01/cc-json-parser.git
cd cc-json-parser
```

### Building

Compile the project using:

```bash
go build -o ccjp
```

### Testing

```bash
go test ./...
```

### Usage

To run the utility, you can either validate a JSON file or run predefined tests. The command-line interface (CLI) provides flags for different operations:

#### Validate a JSON file:

```bash
./ccjp --json-parser path/to/file.json
```

#### Run predefined tests:

```bash
./ccjp --run-tests
```

This command runs predefined tests on JSON files located in the test directories (./test_data/pass and ./test_data/fail). The tests check whether the files pass or fail as expected.

### Examples

```bash
./ccjp --json-parser test_data/pass/valid_3.json
```
