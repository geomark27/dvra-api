package service

type UpdateCompanySettingsDTO struct {
	DisplayName           *string `json:"display_name,omitempty" binding:"omitempty,max=255"`
	LogoURL               *string `json:"logo_url,omitempty"`
	PrimaryColor          *string `json:"primary_color,omitempty" binding:"omitempty,hexcolor"`
	SecondaryColor        *string `json:"secondary_color,omitempty" binding:"omitempty,hexcolor"`
	Locale                *string `json:"locale,omitempty" binding:"omitempty,max=20"`
	CareerPageTitle       *string `json:"career_page_title,omitempty" binding:"omitempty,max=255"`
	CareerPageDescription *string `json:"career_page_description,omitempty"`
	CareerPageEnabled     *bool   `json:"career_page_enabled,omitempty"`
	CustomDomain          *string `json:"custom_domain,omitempty" binding:"omitempty,max=255"`
	ReplyToEmail          *string `json:"reply_to_email,omitempty" binding:"omitempty,email"`
}

type SubscriptionUsageDTO struct {
	Users        int64 `json:"users"`
	Jobs         int64 `json:"jobs"`
	Candidates   int64 `json:"candidates"`
	Applications int64 `json:"applications"`
	StorageBytes int64 `json:"storage_bytes"`
}
