package handlers

import (
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/models"
	"github.com/BPJS-Hackathon/Blockchain-API-Gateway-BPJS-Claim-Backend/utils"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminServ models.AdminService
}

func NewAdminHandler(engine *gin.Engine, adminService models.AdminService, jwtManager *utils.JWTManager) {
	handler := &AdminHandler{adminServ: adminService}
	// public routes
	admin := engine.Group("/admin")

}
