# Shuttle Extensions

Shuttle extensions aims to provide a holistic developer experience for your
entire organisation. It was created to support our users, developers and people
in between. It allows you to quickly and easily express your manual workflows,
or ineffecienies using proper engineering techniques and easily make those
available to your peers.

Shuttle extensions are a certain set of tools you'd like to provide your users,
or for yourself, which enables workflows that span multiple projects, or
entirely external to your projects.

A few examples could be:

```bash
# Will release the current project into dev
shuttle release --env dev 

# Will call the current service in dev according to its own OpenAPI spec.
shuttle call myServiceEndpoint --env dev --someArg someArg

# Search for code
shuttle search "FindMyDependency"

# Open the current docs page for the given service
shuttle docs
shuttle docs --service lunarway/some-service

# Whatever workflow you can imagine
shuttle lunar --transfer 500 --from someone --to someoneElse
```

Given the generic nature of these extensions shuttle itself acts as a registry
and package manager for your extensions, but doesn't provide much for them to
use. Other than the variables used in the project

Shuttle will not reserve keywords for its own commands but, will always execute
its own first, this effectively means that we will not respect your list of
extensions.

Notice: Extensions will be turned off if shuttle is run within plan actions
(shuttle run build etc.). This is because extensions are not made to be
distributed to those parts of your CI and in this case it is your plans job to
make sure you have good developer experience. An extension can still use a plan,
but not the other way around

## Usage

To enable shuttle extensions one or more registries are required. A registry is
an opinionated repository with a shuttle-extensions.json file containing all the
extensions, their binary paths, checksums and so on.

To add a registry simply set variable

```bash
export SHUTTLE_EXTENSIONS_REGISTRY="git=github.com/lunarway/example-registry"
```

The host machine needs to have native access to said registry otherwise an error
is returned. As long you can `git clone` said repository you should be good.

Shuttle extensions needs to be installed and updated manually or via. a
background job. Shuttle itself at the moment make no such attempt at runtime to
not incur a runtime cost. But will periodically fetch said git repository and
print a notice that an update should be done via. `shuttle ext update`.

The notice is opt in, so as to not pollute machine usage

```bash
export SHUTTLE_EXTENSIONS_NOTICE="warn" # none/warn/error
```

To run install the extensions:

```bash
# for installing current extensions from the upstream repository, is a noop if the extensions are already installed
shuttle ext install 

# for updating and installing missing extensions, is nearly never a noop as it will fetch the upstream registries first
shuttle ext update
```

## Building an extension

Currently only golang shuttle extensions are supported, but other types of
binaries are technically supported but will need their own sdks to work.

To bootstrap an extension simply run

```bash
shuttle ext init  
# Prompt: name (my-example-extensions)
# Prompt: path (.) 
# Prompt: provider: (github)
```

This will create a repository matching the template found in: TBD

Such as a:

- shuttle.extension.yaml: describing certain properties of shuttle, mainly type
  of extensions and name of final binary
- main.go: Basic extension containing the shuttle extension sdk
- .github/workflows/ci.yaml: contains ci required for building binaries for
  various operating systems, as well as publishing the extension to the registry
  described in shuttle.extension.yaml

## Architecture

Shuttle itself doesn't contain much more than the `ext` subcommand, along with
command calling facilities.

It will however provide an opinionated layout for where it wants to store files.

Shuttle will place whatever files it needs in `~/.shuttle`, or whatever
`SHUTTLE_CONFIG` points to.

```
- ~/.shuttle/registry: will contain a folder for each of the registries listed in SHUTTLE_EXTENSIONS_REGISTRY  
- ~/.shuttle/extensions/cache: will contain all the extensions as binaries fetched for the registries and providers
```

That is it for now.
