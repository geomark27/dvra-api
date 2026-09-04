package models

type CompanySettings struct {
	BaseModel
	CompanyID             uint     `gorm:"not null;uniqueIndex" json:"company_id"`
	DisplayName           string   `gorm:"type:varchar(255)" json:"display_name,omitempty"`
	LogoURL               string   `gorm:"type:text" json:"logo_url,omitempty"`
	PrimaryColor          string   `gorm:"type:varchar(7);not null;default:'#2563eb'" json:"primary_color"`
	SecondaryColor        string   `gorm:"type:varchar(7);not null;default:'#0f172a'" json:"secondary_color"`
	Locale                string   `gorm:"type:varchar(20);not null;default:'es'" json:"locale"`
	CareerPageTitle       string   `gorm:"type:varchar(255)" json:"career_page_title,omitempty"`
	CareerPageDescription string   `gorm:"type:text" json:"career_page_description,omitempty"`
	CareerPageEnabled     bool     `gorm:"not null;default:true" json:"career_page_enabled"`
	CustomDomain          string   `gorm:"type:varchar(255)" json:"custom_domain,omitempty"`
	ReplyToEmail          string   `gorm:"type:varchar(255)" json:"reply_to_email,omitempty"`
	Company               *Company `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"company,omitempty"`
}

func (CompanySettings) TableName() string { return "company_settings" }
