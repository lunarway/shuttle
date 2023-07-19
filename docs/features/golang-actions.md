# Golang Actions

Shuttle supports multiple types of actions. Actions are the things command you
do on the command line via. shuttle which does stuff.

That is `shuttle run build`, `shuttle run test` etc.

These actions can be defined in multiple ways, and in multiple formats. As the
`examples` show they can be `docker`, `shell` or `golang`.

As defined in either the `shuttle.yaml` or `plan.yaml` when building a plan.

A script section can be defined.

```yaml
scripts:
  build:
	  actions:
		  - shell: echo "build...
```

That of course is a shell action, it can also call a file via:
`- shell: $scripts/build.sh`

However, we can also use golang. To use golang you need a few prerequisites:

- golang 1.19+
- or docker

To create an action in your plan, you don't need either the `plan.yaml` or the
`shuttle.yaml` script sections. You just create a folder:
`mkdir -p actions && cd actions`

The actions folder is the place our golang code actions will live. Each regular
golang file in this folder will be treated as an action, this doesn't include
tests or folders. As such if you need files that aren't actions, create a folder
and put the code in there.

```bash
go mod init actions
go get github.com/lunarway/shuttle
```

Create a golang submodule, and add shuttle to it.

```bash
echo <--GOEOF
package main

import (
	"context",
	_ "github.com/lunarway/shuttle" // default the base shuttle so that it doesn't disappear from go.mod
)

func Build(ctx context.Context) error {
	println("build")

	return nil
}
GOEOF > build.go
```

This will create a file and a function with the same name. Note, the file name
and the function name have to match. golang file names support snake_case, and
golang func names support PascalCase. As such they have to match. Otherwise the
build will fail.

1. file: "build_production.go" -> func BuildProduction
2. file: "build.go" -> func Build

Lower case functions are ignored, only 1 public function is allowed pr actions
file.

1. func handleBuild -> ignored
2. func Build -> allowed
3. func build, func Build -> allowed
4. func Build, func TestBuild -> error (only 1 public function pr file is
   allowed)

Now you can run the command via. shuttle.

```
shuttle run build
stdout: build
```

## Why

Why would you want such a feature?

- First of all software engineering. Using golang over bash, enables a few
  things.
- It is easier to build solid tools in golang than bash. Golang code can be
  shared and distributed.
- You can include packages, snippets and whatnot in golang.
- Golang can more easily be tested and benchmarked.
- Binaries can be saved and reused to improve performance and startup time.
- A more thorough user experience can be designed.
- No longer bound to what is installed on a client machine, you don't need curl,
  wget, uname, grep, yq, jq etc. installed
