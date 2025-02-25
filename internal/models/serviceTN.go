package models

import "time"

type ServiceTreeNode struct {
	ServiceID    string       `gorm:"primaryKey;size:50" json:"service_id"`
	ServiceName  string       `json:"service_name" binding:"required"`
	ParentID     *string      `gorm:"size:50" json:"parent_id"`
	Status       string       `gorm:"size:20" json:"status" binding:"required,oneof=active warning inactive"`
	Dependencies []Dependency `gorm:"foreignKey:SourceID" json:"dependencies"`
	Owner        string       `gorm:"size:100" json:"owner"`
	LastUpdated  time.Time    `json:"last_updated"`
	Path         string       `json:"path"` // 路径追踪
}

type Dependency struct {
	ID          uint   `gorm:"primaryKey"`
	SourceID    string `gorm:"size:50"` // 当前服务ID
	TargetID    string `gorm:"size:50"` // 依赖服务ID
	Status      string `gorm:"size:20"`
	Description string `gorm:"size:200"`
}
