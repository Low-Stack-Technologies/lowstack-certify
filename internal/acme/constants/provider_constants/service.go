package provider_constants

const (
	ProviderCloudflare = Provider("cloudflare")
	ProviderWebsupport = Provider("websupport")
	ProviderCPanel     = Provider("cpanel")
)

type Provider string
