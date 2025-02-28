package services

import (
	"github.com/gin-gonic/gin"
	"github.com/harrychenshoplazza/service_tree/internal/models"
	"github.com/harrychenshoplazza/service_tree/utils"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// 使用依赖注入方式
type ServiceHandler struct {
	db *gorm.DB
}

func NewServiceHandler(db *gorm.DB) *ServiceHandler {
	return &ServiceHandler{db: db}
}

// 创建服务节点
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var input struct {
		ServiceName  string  `json:"service_name" binding:"required"`
		ParentID     *string `json:"parent_id"`
		Status       string  `json:"status" binding:"required"`
		Dependencies []struct {
			TargetID string `json:"service_id"`
			Status   string `json:"status"`
		} `json:"dependencies"`
		Owner string `json:"owner"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serviceID := utils.GenerateServiceID()

	node := models.ServiceTreeNode{
		ServiceID:   serviceID,
		ServiceName: input.ServiceName,
		ParentID:    input.ParentID,
		Status:      input.Status,
		Owner:       input.Owner,
		LastUpdated: time.Now(),
	}

	// 处理路径
	tx := h.db.Begin()
	if input.ParentID != nil {
		var parent models.ServiceTreeNode
		if err := tx.Where("service_id = ?", *input.ParentID).First(&parent).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent node not found"})
			return
		}
		node.Path = parent.Path + "/" + node.ServiceID
	} else {
		node.Path = "/" + node.ServiceID
	}

	// 创建依赖关系
	for _, dep := range input.Dependencies {
		// 验证依赖服务是否存在
		var target models.ServiceTreeNode
		if err := tx.Where("service_id = ?", dep.TargetID).First(&target).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dependency service not found: " + dep.TargetID})
			return
		}

		dependency := models.Dependency{
			SourceID: node.ServiceID,
			TargetID: dep.TargetID,
			Status:   dep.Status,
		}
		if err := tx.Create(&dependency).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Create(&node).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, node)
}

// 获取服务节点信息
func (h *ServiceHandler) GetService(c *gin.Context) {
	var service models.ServiceTreeNode
	if err := h.db.First(&service, c.Param("service_id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	c.JSON(http.StatusOK, service)
}

// 获取子服务
func (h *ServiceHandler) GetChildren(c *gin.Context) {
	var services []models.ServiceTreeNode
	h.db.Where("parent_id = ?", c.Param("service_id")).Find(&services)
	c.JSON(http.StatusOK, services)
}

// 删除服务
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id := c.Param("service_id")
	var service models.ServiceTreeNode
	if err := h.db.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
	}
	h.db.Where("path LIKE ?", service.Path+"%").Delete(&models.ServiceTreeNode{})
	c.JSON(http.StatusNoContent, nil)
}
