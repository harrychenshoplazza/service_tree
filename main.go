package main

import (
	"github.com/harrychenshoplazza/service_tree/internal/models"
	. "github.com/harrychenshoplazza/service_tree/internal/services"
	"github.com/harrychenshoplazza/service_tree/routes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	// 初始化数据库
	var err error
	DB, err = gorm.Open(sqlite.Open("service-tree.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移表结构
	DB.AutoMigrate(&models.ServiceTreeNode{}, &models.Dependency{})
	serviceHandler := NewServiceHandler(DB)

	// 初始化路由
	r := routes.SetupRouter(serviceHandler)
	r.Run(":8080")
}
