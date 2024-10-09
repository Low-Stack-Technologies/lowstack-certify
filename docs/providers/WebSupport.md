# WebSupport Provider

If you are using [WebSupport](https://www.websupport.se/) as your DNS provider, you can use the WebSupport provider to obtain and renew certificates.

## Required Provider Options

```yaml
provider: websupport

provider_options:
  api_key: YOUR_API_KEY
  api_secret: YOUR_API_SECRET
```

- `provider: websupport` - This has to be set to `websupport` to use the WebSupport provider.
- `provider_options.api_key` - This is the API key that will be used to authenticate with the WebSupport API.
- `provider_options.api_secret` - This is the API secret that will be used to authenticate with the WebSupport API.

## How to Obtain the API Key and Secret

Head over to the [Security and login](https://admin.websupport.se/en/auth/security-settings) page and create a new API key.

1. Scroll down to "API Authentication & Dynamic DNS".
2. Click on "+ Generate new API access".
3. Choose "Standard" and click "+ Generate new API access".
4. (Optional) Give the API key a suitable name.
5. Copy the API key (Identifier) and API secret (Secret).