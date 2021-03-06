[![Build Status](https://snap-ci.com/ashwanthkumar/marathonctl/branch/master/build_image)](https://snap-ci.com/ashwanthkumar/marathonctl/branch/master)

# Marathon CLI
CLI tool to access and deploy apps and services to [Marathon](https://mesosphere.github.io/marathon/).

## NOTE
The master branch is a WIP of package repository feature. Marathonctl as a simple commandline tool to deploy apps to marathon is implemented and you can find the binaries [here](https://github.com/ashwanthkumar/marathonctl/releases). Please use the latest `v0.0.x`. The package manager changes will come in `v0.1.x` series as we're making some non-backward compatible changes.

## Usage
You can download a binary distribution from the [releases](https://github.com/ashwanthkumar/marathonctl/releases).

```
$ marathonctl
Command line client to Marathon

Usage:
  marathonctl [command]

Available Commands:
  deploy      Deploy an app using Marathon's app definition
  package     Manage packages which needs to be installed on Marathon
  repo        Manage remote repositories where packages can be installed
  version     Version of the Marathon CLI

Flags:
  -h, --help                   help for marathonctl
      --marathon.host string   Marathon host in http://host:port form. (default "http://localhost:8080")
      --mesos.master string    Mesos host in host:port form. (default "localhost:5050")
      --zk.host string         ZK host in host:port form. (default "localhost:2181")

Use "marathonctl [command] --help" for more information about a command.
```

## Configuration
You can optionally create a configuration file `$HOME/.marathonctl/config.json` with the following contents, which can be overriden using the above flags.

```
{
  "marathon": {
    "host": "http://marathon.host:8080"
  },
  "mesos": {
    "master": "mesos.master:5050"
  },
  "zk": {
    "host": "zk01:2181"
  }
}
```


## Deploy Apps
`marathonctl deploy` helps you deploy applications to your Marathon setup from command line. It takes an app definition and tries to deploy it.

```
$ marathonctl deploy -h
Deploy an app using Marathon's app definition

Usage:
  marathonctl deploy <app.json> [flags]

Flags:
  -d, --dry-run              Print the final application configuration but don't deploy
  -e, --environment string   Environment to deploy (default "test")
  -f, --force                Force deploy the app
  -t, --timeout int          timeout in seconds for deployment to complete, else we'll fail (default 900)
```

### Application Definition
The application definition (`app.json`) that's passed it treated as a [Go Template](https://golang.org/pkg/text/template/) and rendered. The available variables for the template is `{{ .DEPLOY_ENV }}`.

#### Using Environment Variables
You can also access environment variables in your app.json using the convention `{{ .Env.GO_PIPELINE_LABEL }}`, where `GO_PIPELINE_LABEL` is an environment variable.

Example app.json file could be something like
```
{
  "id": "{{ .DEPLOY_ENV }}.http",
  "cpus": 0.1,
  "mem": 10,
  "instances": 1,
  "ports": [
    0
  ],
  "cmd": "python -m SimpleHTTPServer $PORT0",
  "uris": [
    "https://github.com/ashwanthkumar/wasp-cli/releases/download/v{{ .Env.WASP_CLI_VERSION }}/wasp-linux-amd64"
  ],
  "upgradeStrategy": {
    "minimumHealthCapacity": 0.9,
    "maximumOverCapacity": 0.1
  },
  "env": {
    "DEPLOY_ENV": "{{ .DEPLOY_ENV }}"
  },
  "healthChecks": [
    {
      "protocol": "COMMAND",
      "command": { "value": "curl -f http://$HOST:$PORT0/" },
      "gracePeriodSeconds": 60,
      "intervalSeconds": 30,
      "maxConsecutiveFailures": 3,
      "timeoutSeconds": 10
    }
  ]
}
```
