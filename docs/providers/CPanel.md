# CPanel / WHM Provider

If your hosting service is using [CPanel / WHM](https://www.cpanel.net/) as their interface, you are very likely to be able to use this provider to obtain and renew certificates.

## Required Provider Options

```yaml
provider: cpanel

provider_options:
  username: YOUR_USERNAME # This is the CPanel username
  token: YOUR_API_TOKEN
  base_url: https://cpanel.example.com # This is the base URL of your CPanel installation
  # mode: whm # Optional, defaults to cpanel
```

- `provider: cpanel` - This has to be set to `cpanel` to use the CPanel / WHM provider.
- `provider_options.username` - This is the CPanel username that will be used to authenticate with the CPanel API.
- `provider_options.token` - This is the API token that will be used to authenticate with the CPanel API.
- `provider_options.base_url` - This is the base URL of your CPanel installation.
- `provider_options.mode` - This is the mode that will be used to authenticate with the CPanel API.

## How to Obtain the API Token

Follow the instructions in their documentation called [How to Use cPanel API Tokens](https://docs.cpanel.net/knowledge-base/security/how-to-use-cpanel-api-tokens/) to obtain the API token.