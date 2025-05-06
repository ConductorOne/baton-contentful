![Baton Logo](./baton-logo.png)

# `baton-contentful` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-contentful.svg)](https://pkg.go.dev/github.com/conductorone/baton-contentful) ![main ci](https://github.com/conductorone/baton-contentful/actions/workflows/main.yaml/badge.svg)

`baton-contentful` is a connector for built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-contentful
baton-contentful
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_DOMAIN_URL=domain_url -e BATON_API_KEY=apiKey -e BATON_USERNAME=username ghcr.io/conductorone/baton-contentful:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-contentful/cmd/baton-contentful@main

baton-contentful

baton resources
```

# Requirements
- organization id
guide to find your organization id [here](https://www.contentful.com/help/organizations/find-organization-id/)
- CMA Token
you can create the token in the settings menu.

# Data Model

`baton-contentful` will pull down information about the following resources:
- Organizations
- Spaces
- Teams
- Users

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-contentful` Command Line Usage

```
baton-contentful

Usage:
  baton-contentful [flags]
  baton-contentful [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  config             Get the connector config schema
  help               Help about any command

Flags:
      --client-id string                                 The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string                             The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
      --external-resource-c1z string                     The path to the c1z file to sync external baton resources with ($BATON_EXTERNAL_RESOURCE_C1Z)
      --external-resource-entitlement-id-filter string   The entitlement that external users, groups must have access to sync external baton resources ($BATON_EXTERNAL_RESOURCE_ENTITLEMENT_ID_FILTER)
  -f, --file string                                      The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                                             help for baton-contentful
      --log-format string                                The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string                                 The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --organization-id string                           required: The ID of the organization to use. ($BATON_ORGANIZATION_ID)
      --otel-collector-endpoint string                   The endpoint of the OpenTelemetry collector to send observability data to (used for both tracing and logging if specific endpoints are not provided) ($BATON_OTEL_COLLECTOR_ENDPOINT)
  -p, --provisioning                                     This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --skip-full-sync                                   This must be set to skip a full sync ($BATON_SKIP_FULL_SYNC)
      --ticketing                                        This must be set to enable ticketing support ($BATON_TICKETING)
      --token string                                     required: The API token used to authenticate with the service. ($BATON_TOKEN)
  -v, --version                                          version for baton-contentful

Use "baton-contentful [command] --help" for more information about a command.
```
