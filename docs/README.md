# Low-Stack Certify

## Table of Contents

- [Getting Started](#getting-started)
 - [Docker Compose](#docker-compose)

## Getting Started

### Docker Compose

To get started using Docker Compose you can use the following configuration:

```yaml
services:
  certify:
    image: lowstack/certify:latest
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

### Configuration

The configuration file can be found at the path specified by the `CUSTOM_CONFIGURATION_PATH` environment variable,
or wherever you mounted the configuration directory.

Here is the default configuration file:

```yaml
# This is the default configuration file for the Low-Stack Certify application.
# You can change all values here to customize the application to your needs.

# This is the runtime configuration for the application.
runtime:
  # This is whether the application should run periodically or not.
  # If set to false, the application will only run once and exit.
  run_periodically: true

  # This is the number of minutes between each run of the application.
  period_in_minutes: 15

# The path to the directory containing the zone configuration files.
# This can be a relative or absolute path.
zones_path: "/zones"

# The path to the directory where certificates will be stored.
# This can be a relative or absolute path.
certificates_path: "/certificates"

# This is the URL of the CA directory that will be used to sign certificates.
#ca_url: "https://acme-staging-v02.api.letsencrypt.org/directory" # Use this for testing
ca_url: "https://acme-v02.api.letsencrypt.org/directory"
```

#### Zones

Zones are the configuration files that define the domains that will be managed by the application.

Here is an example zone configuration file:

```yaml
# This is an example zone configuration file,
# you can copy this file to zones/<zone>.yaml and edit it.

# This is a unique identifier for the zone.
#
# IMPORTANT! This has to be safe to use as a directory
# name, because it will be used as the directory name
# for storing the certificates!
unique_identifier: example.com

# The hostnames that this zone will be responsible for.
# You can specify multiple hostnames, including wildcard
# hostnames using *.example.com.
hostnames:
  - example.com
  - www.example.com

# This is the email used when registering with the ACME
# server. This can be shared between multiple zones.
#
# You might get notifications from the ACME server
# when certificates are about to expire, so it's
# important to keep this email up to date.
identity_email: example@example.com

# This is the number of days before the certificate
# expires on which to renew the certificate.
#
# Normally, a new certificate expires every 90 days.
renewal_days: 15

# This is the provider that will be used to challenge
# the certificate.
#
# Available providers: cloudflare, websupport
provider: cloudflare

# This is the options that will be passed to the
# provider when requesting a certificate.
# These options are provider-specific.
provider_options:
  # Cloudflare
  api_token: YOUR_API_TOKEN # Requires DNS:Edit & ZONE:Read permissions
  
  # WebSupport
  #api_key: YOUR_API_KEY
  #api_secret: YOUR_API_SECRET

# This is the type of key that will be generated for the
# certificate.
#
# Available key types: EC256 (P256), EC384 (P384),
# RSA2048 (2048), RSA3072 (3072), RSA4096 (4096),
# RSA8192 (8192)
key_type: 2048

# This is the file permissions that will be used for the
# certificate files.
#
# If enabled, the certificate files will be updated to
# match the specified permissions. Otherwise, the
# certificate files will be left as-is.
#
# IMPORTANT! This does require root permissions to
# update the file permissions!
file_permissions:
  enabled: false
  uid: 0
  gid: 0
  private_key_mode: 0600
  full_chain_mode: 0644
```

##### Providers

The following providers are currently supported:

- [Cloudflare](providers/Cloudflare.md)
- [WebSupport](providers/WebSupport.md)
