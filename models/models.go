package models

// Role Database model info
// @Description App type information
type Role struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
}

// RolePost model info
// @Description RolePost type information
type RolePost struct {
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
}

// RoleGet model info
// @Description RoleGet type information
type RoleGet struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
}

// RolePut model info
// @Description RolePut type information
type RolePut struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
}

// RolePatch model info
// @Description RolePatch type information
type RolePatch struct {
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
}
