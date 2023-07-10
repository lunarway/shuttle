<p align="center">
  <a href="https://github.com/lunarway/shuttle">
    <img src="docs/logo.png" alt="Shuttle logo">
  </a>

<p align="center">
    A CLI for handling shared build and deploy tools between many projects no matter what technologies the project is using.
    <br>
    <a href="https://github.com/lunarway/shuttle/issues/new?template=bug.md">Report bug</a>
    ·
    <a href="https://github.com/lunarway/shuttle/issues/new?template=feature.md&labels=feature">Request feature</a>
    ·
    <a href="https://github.com/lunarway/shuttle/releases">Releases</a>
    ·
    <a href="https://github.com/lunarway/shuttle/releases/latest">Latest release</a>
  </p>
</p>

## Table of contents

- [What is shuttle?](#what-is-shuttle)
- [Status](#status)
- [How?](#how)
- [Features](#features)
- [Documentation](#documentation)
- [Installing](#installing)
- [Functions](#functions)
- [Release History](#release-history)

## What is shuttle?

`shuttle` is a CLI for handling shared build and deploy tools between many
projects no matter what technologies the project is using.

## Status

_DISCLAIMER: shuttle is in beta, so stuff may change. However we are using
shuttle heavily at Lunar Way and we use it to deploy to production, so it is
pretty battle proven._

[![Go Report Card](https://goreportcard.com/badge/github.com/lunarway/shuttle)](https://goreportcard.com/report/github.com/lunarway/shuttle)

## How?

Projects that use `shuttle` are always referencing a `shuttle plan`. A plan
describes what can be done with shuttle. Fx:

```yaml
# plan.yaml file
scripts:
  build:
    description: Build the docker image
    args:
      - name: tag
        required: true
    actions:
      - shell: docker -f $plan/Dockerfile build -t $(shuttle get docker.image):$tag
  test:
    description: Run test for the project
    actions:
      - shell: go test
```

The `plan.yaml` is located at the root of the plan directory which is located
elsewhere of the actual project using it. The plan directory can be locally
stored or in a git repository. The directory structure could be something like:

```sh
workspace
│
└───moon-base          # project
│   │   shuttle.yaml   # project specific shuttle.yaml file
│   │   main.go
│
└───station-plan       # plan to be shared by projects
    │   plan.yaml
    │   Dockerfile
```

To use a plan a project must specify the `shuttle.yaml` file:

```yaml
plan: ../the-plan
vars:
  docker:
    image: earth-united/moon-base
```

With this in place a docker image can be built:

```sh
$ cd workspace/moon-base
$ shuttle run build tag=v1
```

You can also put your `plan.yaml` in a different git repository, and then
reference it from your `shuttle.yaml`:

```
# shuttle.yaml
plan: git://git@github.com:lunarway/station-plan.git
```

## Features

- Fetch shuttle plans from git repositories
- Create any script you like in the plan
- Overwrite scripts in local projects when they defer from the plan
- Write templates in plans and overwrite them in projects when they defer
- ...

### Telemetry

see [telemetry](./docs/features/telemetry.md)

## Documentation

Plan documentation can be inspected using the `shuttle documentation` command.

When writing shuttle plans you can hint as to where to find documentation for
the plan. Users of the plan will open the specified URL when requesting
documentation.

```yaml
# plan.yaml
documentation: https://docs.my-corp.com
```

If no specific `documentation` field is set in the plan it will be inferred from
the plan reference. In below example shuttle will open
`https://github.com/lunarway/shuttle-example-go-plan.git`

```yaml
# shuttle.yaml
plan: git://git@github.com:lunarway/shuttle-example-go-plan.git
```

### Git Plan

Specify the protocol to use for pulling the remote repository using either
`https://` for HTTPS, or `git://` for SSH:

- `https://github.com/lunarway/shuttle-example-go-plan.git`
- `git://git@github.com:lunarway/shuttle-example-go-plan.git`

Choose a specific branch to use:

- `https://github.com/lunarway/shuttle-example-go-plan.git#change-build`
- `git://git@github.com:lunarway/shuttle-example-go-plan.git#change-build`

The `#change-build` points the plan to a specific branch, which by default would
be `master`.

It can also be used to point to a tag or a git SHA, like this:

- `https://github.com/lunarway/shuttle-example-go-plan.git#v1.2.3`
- `git://git@github.com:lunarway/shuttle-example-go-plan.git#46ce3cc`

#### Caching

By default shuttle will pull the upstream plan on every pull. To prevent this
you can use

```bash
export SHUTTLE_CACHE_DURATION_MIN=60 # Cache a plan for 60 minutes
```

This feature caches pr. repo, as such the cache isn't shared between working
repositories.

### Overloading the plan

It is possible to overload the plan specified in `shuttle.yaml` file by using
the `--plan` argument or the `SHUTTLE_PLAN_OVERLOAD` environment variable.
Following arguments are supported

- A path to a local plan like `--plan ../local-checkout-of-plan`. Absolute paths
  is also supported
- Another git plan like `--plan git://github.com/some-org/some-plan`
- A git tag to append to the plan like `--plan #some-branch`, `--plan #some-tag`
  or a SHA `--plan #2b52c21`

## Installing

### Mac OS

```console
curl -LO https://github.com/lunarway/shuttle/releases/download/$(curl -Lso /dev/null -w %{url_effective} https://github.com/lunarway/shuttle/releases/latest | grep -o '[^/]*$')/shuttle-darwin-amd64
chmod +x shuttle-darwin-amd64
sudo mv shuttle-darwin-amd64 /usr/local/bin/shuttle
```

### Linux

```console
curl -LO https://github.com/lunarway/shuttle/releases/download/$(curl -Lso /dev/null -w %{url_effective} https://github.com/lunarway/shuttle/releases/latest | grep -o '[^/]*$')/shuttle-linux-amd64
chmod +x shuttle-linux-amd64
sudo mv shuttle-linux-amd64 /usr/local/bin/shuttle
```

### GitHub Actions

Shuttle can be installed on your GitHub Runner by adding this line to your
workflow:

```
- use: lunarway/shuttle
```

After this point you can use shuttle in the scripts in your workflow job.

## Functions

### `shuttle get <variable>`

Used to get a variable defined in shuttle.yaml

```console
$ shuttle get some.variable
> some variable content

$ shuttle get does.not.exist
> # nothing
```

### `shuttle plan`

Inspect the plan in use for a project. Use the `template` flag to customize the
output to your needs.

```console
$ shuttle plan
https://github.com/lunarway/shuttle-example-go-plan.git
```

### `shuttle has <variable>`

It is possible to easily check if a variable or script is defined

```console
shuttle has some.variable
shuttle has --script integration
```

Output is either statuscode=0 if variable is found or statuscode=1 if variables
isn't found. The output can also be a stdout boolean like

```console
$ shuttle has my.docker.image --stdout
> false
```

### Template functions

The `template` command along with commands taking a `--template` flag has
multiple templating functions available. The
[masterminds/sprig](http://masterminds.github.io/sprig) v3.2.2 functions are
available along with those described below.

Examples are based on the below `shuttle.yaml` file.

| Function                      | Description                                                                                                                                                             | Example                                                                 | Output                            |
| ----------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- | --------------------------------- |
| `array <path> <value>`        | Get array from path. If value is a map, the values of the map is returned in deterministic order.                                                                       | `array "args" .`                                                        | `helloworld`                      |
| `fileExists <file-path>`      | Returns whether a file exists.                                                                                                                                          | `fileExists ".gitignore"`                                               | `true`                            |
| `fromYaml <value>`            | Unmarshal YAML string to a `map[string]interface{}`. In case of YAML parsing errors the `Error` key in the result contains the error message. See notes below on usage. | `fromYaml "api: v1"`                                                    | `map[api:v1]`                     |
| `get <path> <value>`          | Get a value from a field path. `.` is read as nested nested objects                                                                                                     | `get "docker.image" .`                                                  | `earth-united/moon-base`          |
| `getFileContent <file-path>`  | Get raw contents of a file.                                                                                                                                             | `getFileContent ".gitignore"`                                           | `dist/`<br> `vendor/`<br> `...`   |
| `getFiles <directory-path>`   | Returns a slice of files in the provided directory as [`os.FileInfo`](https://golang.org/pkg/os/#FileInfo) structs.                                                     | `{{ range $i, $f := (getFiles "./") }}{{ .Name }} {{ end }}`            | `.git .gitignore ...`             |
| `int <path> <value>`          | Get int value without formatting. Note that this is a direct `int` cast ie. value `1.2` will generate an error.                                                         | `int "replicas" .`                                                      | `1`                               |
| `is <value-a> <value-b>`      | Equality indication by Go's `==` comparison.                                                                                                                            | `is "foo" "bar"`                                                        | `false`                           |
| `isnt <value-a> <value-b>`    | Inequality indication by Go's `!=` comparison.                                                                                                                          | `isnt "foo" "bar"`                                                      | `true`                            |
| `objectArray <path> <value>`  | Get object key-value pairs from path. Each key-value is returned in a `{ Key Value}` object.                                                                            | `{{ range objectArray "docker" . }}{{ .Key }} -> {{ .Value }}{{ end }}` | `image -> earth-united/moon-base` |
| `rightPad <string> <padding>` | Add space padding to the right of a string.                                                                                                                             | `{{ rightPad "padded" 10 }}string`                                      | `padded string`                   |
| `strConst <value>`            | Convert string to upper snake casing converting `.` to `_`.                                                                                                             | `strConst "a.value"`                                                    | `A_VALUE`                         |
| `string <path> <value>`       | Format any value as a string.                                                                                                                                           | `string "replicas" .`                                                   | `"1"`                             |
| `toYaml <value>`              | Marshal value to YAML. In case of non-parsable string an empty string is returned. See notes below on usage.                                                            | `toYaml (get "args" .)`                                                 | `- hello`<br> `- world`           |
| `trim <string>`               | Trim leading and trailing whitespaces. Uses [`strings.TrimSpace`](https://golang.org/pkg/strings/#TrimSpace).                                                           | `trim " a string "`                                                     | `a string`                        |
| `upperFirst <string>`         | Upper case first character in string.                                                                                                                                   | `upperFirst "a string"`                                                 | `A string`                        |

```yaml
plan: ../the-plan
vars:
  docker:
    image: earth-united/moon-base
  replicas: 1
  args:
    - hello
    - world
```

**Notes on YAML parsers**: The `toYaml` and `fromYaml` template functions are
intented to be used inside file templates (not `template` flags). Because of
this they ignore errors and some YAML documents canont be parsed.

## Release History

See the [releases](https://github.com/lunarway/shuttle/releases) for more
information on changes between releases.
