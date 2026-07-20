package models

// Department represents one node in a tenant's organizational structure.
// Subdepartments use the same table through ParentDepartmentID.
type Department struct {
	BaseModel

	CompanyID          uint   `gorm:"not null;index" json:"company_id"`
	ParentDepartmentID *uint  `gorm:"index" json:"parent_department_id,omitempty"`
	Name               string `gorm:"type:varchar(150);not null" json:"name"`
	Code               string `gorm:"type:varchar(50)" json:"code,omitempty"`
	Description        string `gorm:"type:text" json:"description,omitempty"`
	ManagerUserID      *uint  `gorm:"index" json:"manager_user_id,omitempty"`
	IsActive           bool   `gorm:"not null;default:true;index" json:"is_active"`

	Company          *Company             `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ParentDepartment *Department          `gorm:"foreignKey:ParentDepartmentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Children         []Department         `gorm:"foreignKey:ParentDepartmentID" json:"-"`
	Manager          *User                `gorm:"foreignKey:ManagerUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Path             []DepartmentPathItem `gorm:"-" json:"-"`
}

func (Department) TableName() string { return "departments" }

// DepartmentPathItem is one ancestor in a root-to-leaf department path.
type DepartmentPathItem struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code,omitempty"`
}
