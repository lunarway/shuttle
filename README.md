# shuttle [![Build Status](https://travis-ci.com/lunarway/shuttle.svg?branch=master)](https://travis-ci.com/lunarway/shuttle)

*DISCLAIMER: shuttle is in its alpha stage and is not yet production ready. Expect the APIs to change.*

## What is shuttle?
`shuttle` is a CLI for handling shared build and deploy tools between many projects no matter what technologies the project is using.

## How?

Projects that use `shuttle` are always referencing a `shuttle plan`. A plan describes what can be done with shuttle. Fx:

```yaml
# plan.yaml file
scripts:
  build:
    description: Build the docker image
    args:
    - name: tag
      required: true
    actions:
    - shell: docker -f $plan/Dockerfile build -t (shuttle get docker.image):$tag
  test:
    description: Run test for the project
    actions:
    - shell: go test
```

The `plan.yaml` is located at the root of the plan directory which is located elsewhere of the actual project using it. The plan directory can be locally stored or in a git repository. The directory structure could be something like:

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

## Features
* Fetch shuttle plans from git repositories
* Create any script you like in the plan
* Overwrite scripts in local projects when they defer from the plan
* Write templates in plans and overwrite them in projects when they defer
* ...

## Documentation
*Documention is coming*


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


## Release History

See the [releases](https://github.com/lunarway/shuttle/releases) for more
information on changes between releases.
