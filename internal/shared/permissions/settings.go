package permissions

const (
	SettingsView   = "settings.view"
	SettingsUpdate = "settings.update"
	BillingView    = "billing.view"
)

func init() { define(SettingsView, SettingsUpdate, BillingView) }
