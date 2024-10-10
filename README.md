# Low-Stack Certify

[![GitHub last commit (branch)](https://img.shields.io/github/last-commit/Low-Stack-Technologies/lowstack-certify/main)](https://github.com/Low-Stack-Technologies/lowstack-certify)

## Introduction

Low-Stack Certify is a tool that automates the process of obtaining and renewing SSL/TLS certificates using the ACME protocol. It is designed to work with different DNS providers, and be easily expanded upon. The tool can be deployed using Docker Compose, allowing users to configure and run it with custom settings for certificate storage, zone management, and periodic execution. The configuration is highly customizable, supporting various key types and file permissions, and is intended to simplify the management of SSL/TLS certificates for multiple domains.

## Table of Contents

- [Introduction](#introduction)
- [Getting Started](#getting-started)
  - [Using Docker](#using-docker)
  - [Using Docker Compose](#using-docker-compose)
  - [Running the binary (TBD)](#running-the-binary)
- [Configuration](#configuration)
- [Providers](#providers)

## Getting Started

### Using Docker

To get started using Docker, you can use the following command:

```bash
$ docker run -d --name certify \
  -v /path/to/config:/config \
  -v /path/to/zones:/zones \
  -v /path/to/certificates:/certificates \
  ghcr.io/low-stack-technologies/lowstack-certify:latest
```

This will start the application and mount the configuration, zones, and certificates directories to the container.

Then you can continue to the [Configuration](#configuration) section to configure the application.

### Using Docker Compose

To get started using Docker Compose, you can use the following configuration:

```yaml
services:
  certify:
    image: ghcr.io/low-stack-technologies/lowstack-certify:latest
    container_name: certify
    user: "1000:1000" # This has to be set to 0:0 if you want to set the file permissions
    volumes:
      - ./config:/config
      - ./zones:/zones
      - ./certificates:/certificates
    environment:
      - CUSTOM_CONFIGURATION_PATH=/config/config.yaml
    restart: unless-stopped
```

This will start the application and mount the configuration, zones, and certificates directories to the container.

Then you can continue to the [Configuration](#configuration) section to configure the application.

### Running the binary

To be written

## Configuration

Configuration is split into two parts: [application config](#application-config) and [zone configuration](#zone-configuration).

### Application Config

The application config is located at the path specified by the `CUSTOM_CONFIGURATION_PATH` environment variable, or wherever you mounted the configuration directory.

#### Runtime

The runtime configuration is used to control the behavior of the application. It is used to enable or disable the periodic execution of the application, and to set the interval between executions.

```yaml
runtime:
  run_periodically: true
  period_in_minutes: 15
```

`run_periodically` is a boolean value that determines whether the application should run periodically or not. If set to `false`, the application will only run once and exit. This is useful for testing purposes, running the application manually, or if you wish to set up a cron job to run the application periodically.

`period_in_minutes` is an integer value that determines the interval between executions of the application in minutes. This is the interval that will be used if `run_periodically` is set to `true`. If `run_periodably` is set to `false`, this value will be ignored.

#### Paths

The paths are used to specify the paths where the zones and certificates directories are located. These paths can be absolute or relative.

```yaml
zones_path: "/zones"
certificates_path: "/certificates"
```

`zones_path` is the path where the zones directory is located. This is the directory where the zone configuration files are stored. See the [Zone Configuration](#zone-configuration) section for more information.

`certificates_path` is the path where the certificates directory is located. Each zone will have its own directory within this directory, where the certificates will be stored. The zone `unique_identifier` will be used as the directory name, see the [Zone Configuration](#zone-configuration) section for more information.

#### CA URL

The CA URL is used to specify the URL of the CA directory that will be used to sign certificates. This is used by the ACME client to fetch the CA certificate. This is by default set to Let's Encrypt's production CA URL, but can be changed to use a different CA.

```yaml
ca_url: "https://acme-v02.api.letsencrypt.org/directory"
```

`ca_url` is the URL of the CA directory that will be used to sign certificates. This is used by the ACME client to fetch the CA certificate. Self-hosted CAs can be used by setting this to the URL of the CA directory.

#### Zones

Zones are the configuration files that define the domains that will be managed by the application.

### Zone Configuration

A zone is a group of subdomains for a domain. This can include specific subdomains, wildcard subdomains, or a combination of both. Zones are defined in YAML files, and are located in the `zones` directory specified by the `zones_path` in the [Application Config](#application-config).

A minimal zone configuration file for Cloudflare looks like this:

```yaml
unique_identifier: example.com

hostnames:
  - example.com
  - "*.example.com"

identity_email: example@example.com

renewal_days: 15

provider: cloudflare
provider_options:
  api_token: YOUR_API_TOKEN

key_type: 2048

file_permissions:
  enabled: false
  uid: 0
  gid: 0
  private_key_mode: 0600
  full_chain_mode: 0644
```

#### Unique Identifier

The unique identifier is used to identify the zone. This is used to create the directory where the certificates will be stored. This can be anything, but it is recommended to use a domain name to make it easier to identify the zone.

#### Hostnames

The hostnames are the hostnames that will be managed by the application. This can be a root domain, a subdomain, or a wildcard subdomain. If a wildcard subdomain is used, it has to be enclosed in quotes.

#### Identity Email

This is the email address that will be used to register the user with the ACME server.

#### Renewal Days

This is the number of days before the certificate expires on which to renew the certificate. Normally, a new certificate expires every 90 days.

#### Provider and Provider Options

The provider is the DNS provider that will be used to obtain and renew the certificate. The provider options are the options that will be passed to the provider when requesting a certificate. These options are provider-specific, and can be found in the [Providers](#providers) section.

#### Key Type

The key type is the type of key that will be generated for the certificate. This can be a RSA key (2048, 3072, or 4096), or an EC key (P256, P384, or P521).

#### File Permissions

To automatically set the permissions of the certificate files you can enable this feature. This will update the file permissions to match the specified permissions. Otherwise, the certificate files will be left as-is. **Important!** This does require root permissions to update the file permissions!

```yaml
file_permissions:
  enabled: false
  uid: 0
  gid: 0
  private_key_mode: 0600
  full_chain_mode: 0644
```

`enabled` is a boolean value that determines whether the file permissions should be updated or not.

`uid` is the user ID that will be used to update the file permissions.

`gid` is the group ID that will be used to update the file permissions.

`private_key_mode` is the file mode that will be used to update the private key file permissions.

`full_chain_mode` is the file mode that will be used to update the full chain file permissions.

## Providers

The following providers are currently supported:

- [Cloudflare](docs/providers/Cloudflare.md)
- [WebSupport](docs/providers/WebSupport.md)
- [CPanel / WHM](docs/providers/CPanel.md)

## Contributing

If you would like to contribute to Low-Stack Certify, please read the [CONTRIBUTING.md](CONTRIBUTING.md) file for more information.

## Acknowledgements

This project is based on the following projects:

- [github.com/go-acme/lego](https://github.com/go-acme/lego)

## License

Low-Stack Certify is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.