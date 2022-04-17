# example-golang

An example repository to demonstrate Pants's experimental Golang support.

See the [Golang blog post](https://blog.pantsbuild.org/golang-support-pants-28/) for some unique
benefits Pants brings to Golang repositories, and see
[pantsbuild.org/docs/go-overview](https://www.pantsbuild.org/v2.8/docs/go-overview) for more detailed
documentation.

This is only one possible way of laying out your project with Pants. See 
[pantsbuild.org/docs/source-roots#examples](https://www.pantsbuild.org/docs/source-roots#examples) 
for some other example layouts.

Note: for now, Pants only supports repositories using a single `go.mod`. Please comment on 
[#13114](https://github.com/pantsbuild/pants/issues/13114) if you need support for greater 
than one `go.mod` so that we can prioritize adding support.

# How to use this template?
- global search "gocode-template" and replace it with your service name
- refactor all Company structs in repository, service and handler to your business logic, by following MVC framework

# Prerequisites
- python 3
- direnv
  - for mac, `brew install direnv`
  - for debian linux, `sudo apt install direnv`
- go
  - for mac, `brew install go`
  - for debian linux, `wget -q -O - https://raw.githubusercontent.com/canha/golang-tools-install-script/master/goinstall.sh | bash`
- protobuf
  - for mac, `brew install protobuf`
  - for debian linux, `sudo apt install protobuf-compiler`
- `pip install pantsbuild.pants`

# Running Pants

You run Pants goals using the `./pants` wrapper script, which will bootstrap the
configured version of Pants if necessary.

# Goals

Pants commands are called _goals_. You can get a list of goals with

```
pants help goals
```

Most goals take arguments to run on. To run on a single directory, use the directory name with 
`:` at the end. To recursively run on a directory and all its subdirectories, add `::` to the 
end.

For example:

```
pants lint cmd: internal::
```

You can run on all changed files:

```
pants --changed-since=HEAD lint
```

You can run on all changed files, and any of their "dependees":

```
pants --changed-since=HEAD --changed-dependees=transitive test
```

# Common Goals

## Generate go code for all protobuf

```
pants proto-compile-go
```

## Create BUILD files

```
pants tailor
```

## Format BUILD files

```
pants update-build-files
```

## Generate mock files for all go code

```
pants gomock
```

## Add dependencies
You don't need to explicitly add new dependencies into go.mod. You can just use it in the go files and run

```
go mod tidy
```

and go will automatically fetch new dependencies and update go.mod. Then use the command below to compile.

## Check compilation

```
pants check ::  # Compile all packages.
pants check cmd/greeter_en:  # Compile only this package and its transitive dependencies.
```

## Run Gofmt

```
pants fmt ::  # Format all packages.
pants fmt cmd/greeter_en:  # Format only this package.
pants lint pkg::  # Check that all packages under `pkg` are formatted.
```

## Run tests

```
pants test ::  # Run all tests in the repository.
pants test pkg/uuid:  # Run all the tests in this package.
pants test pkg/uuid: -- -run TestGenerateUuid  # Run just this one test.
```

## Create a binary file

Writes the result to the `dist/` folder.

```
pants package cmd/greeter_en:
pants package cmd::  # Create all binaries.
```

## Run a binary

```
pants run cmd/greeter_en:
pants run cmd/greeter_es: -- --help
```

## Determine dependencies

```
pants dependencies cmd/greeter_en:
pants dependencies --transitive cmd/greeter_en:
```

## Determine dependees

That is, find what code depends on a particular package(s).

```
pants dependees pkg/uuid:
pants dependees --transitive pkg/uuid: 
```

## Count lines of code

```
pants count-loc '**/*'
```

## Update proto related binary

```
go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

and use `brew upgrade` (for mac) or `apt update` (for debian linux) to update protobuf on your OS

## Update gomock binary

```
go install github.com/golang/mock/mockgen@latest
```

# Local Run
Use `pants run cmd/company:` to start the server and run the following command in a new terminal

```
curl -X POST -d '{"instances": [{"name": "company1", "id": "1"}, {"name": "company2", "id": "2"}]}' -H 'Content-Type: application/json' http://localhost:8678/twirp/CompanyService/Add

```
This will add two companies into your local mongo DB.

```
curl -X POST -d '{"name": "company1"}' -H 'Content-Type: application/json' http://localhost:8678/twirp/CompanyService/Search
```

You are supposed to get company1 data from mongo DB.

# TODO
- complete tests
- add pre-commit hooks to check lint, format, tests etc.
- add documentation for this template
- add pants goal for code coverage
- add examples for integration tests and load tests
- add pants goal for shell scripts
- add pants goal for docker image build