package setting

const (
	KeyAppName             = "app_name"
	KeyAppPrivacyPolicy    = "app_privacy_policy"
	KeyAbout               = "about"
	KeyRules               = "rules"
	KeyAppPhone            = "app_phone"
	KeyAppWhatsappUrl      = "app_whatsapp_url" //nolint: gosec // don't worry
	KeyAppFacebookUrl      = "app_facebook_url"
	KeyAppTelegramUrl      = "app_telegram_url"
	KeyAppInstagramUrl     = "app_instagram_url"
	KeyAppWebsiteUrl       = "app_website_url" //nolint: gosec // don't worry
	KeyAppEmailUrl         = "app_email_url"
	KeyAppTwitterUrl       = "app_twitter_url"
	KeyAppLogo             = "app_logo"
	KeyAppColorPrimary     = "app_color_primary"
	KeyAppColorSecondary   = "app_color_secondary"
	KeyAppColorOnPrimary   = "app_color_on_primary"
	KeyAppColorOnSecondary = "app_color_on_secondary"
	KeyAppAppStoreUrl      = "app_app_store_url"
	KeyAppGooglePlayUrl    = "app_google_play_url"
	KeyOrderMinimum        = "order_minimum"
	KeyFreeDelivery        = "free_delivery"

	KeyTlyncEndpoint = "tlync_endpoint"
	KeyTlyncToken    = "tlync_token"
	KeyTlyncStoreID  = "tlync_store_id"
	KeyTlyncFrontUrl = "tlync_front_url"

	KeySas4Username = "sas4_username"
	KeySas4Password = "sas4_password"
	KeySas4Token    = "sas4_token"

	// app specific

	KeyMainHero         = "main_hero"
	KeyMainTitle        = "main_title"
	KeyMainSubtitle     = "main_subtitle"
	KeyMainSummary      = "main_summary"
	KeyMainCallToAction = "main_call_to_action"

	KeyPaymentGatewayEndpoint = "payment_gateway_endpoint"
	KeyPaymentGatewayAPIKey   = "payment_gateway_api_key"
)

var CoreKeys = []string{
	KeyAppName,
	KeyAppPrivacyPolicy,
	KeyAbout,
	KeyRules,
	KeyAppPhone,
	KeyAppWhatsappUrl,
	KeyAppFacebookUrl,
	KeyAppTelegramUrl,
	KeyAppInstagramUrl,
	KeyAppWebsiteUrl,
	KeyAppEmailUrl,
	KeyAppTwitterUrl,
	KeyAppLogo,
	KeyAppColorPrimary,
	KeyAppColorSecondary,
	KeyAppColorOnPrimary,
	KeyAppColorOnSecondary,
	KeyAppAppStoreUrl,
	KeyAppGooglePlayUrl,
	KeyOrderMinimum,
	KeyFreeDelivery,

	KeyTlyncEndpoint,
	KeyTlyncToken,
	KeyTlyncStoreID,
	KeyTlyncFrontUrl,

	KeySas4Username,
	KeySas4Password,
	KeySas4Token,

	KeyMainHero,
	KeyMainTitle,
	KeyMainSubtitle,
	KeyMainSummary,
	KeyMainCallToAction,

	KeyPaymentGatewayEndpoint,
	KeyPaymentGatewayAPIKey,
}
