package controllers

import (
	"net/http"
	"tender_management/models"

	"github.com/gin-gonic/gin"
	"errors"
	"gorm.io/gorm"
)

type NotifController struct {
	Storage *gorm.DB
}

func NewNotifController(storage *gorm.DB) *NotifController {
	return &NotifController{
		Storage: storage,
	}
}

// @Summary      Create Notification
// @Description  Yangi xabar yaratish (Client yoki Contractor uchun)
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        body  body  models.NotifRequest  true  "Notification Body"
// @Success      200   {object}  models.Notif
// @Failure      400   {object} Response  "Bad Request"
// @Failure      500   {object} Response  "Internal Server Error"
// @Router       /notifs [post]
func (n *NotifController) CreateNotif(c *gin.Context) {

	var body models.NotifRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Failed to parse request body", err)
		return
	}

	notif := models.Notif{
		UserID: body.UserID,
		Message: body.Message,
		RelationID: body.RelationID,
		Type: body.Type,
	}

	if err := n.Storage.Create(&notif).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create notification", err)
		return
	}

	HandleResponse(c, http.StatusOK, notif)
} 

// @Summary      Get User Notification
// @Description  User uchun o‘ziga tegishli xabarni olish
// @Tags         Notifications
// @Produce      json
// @Security     BearerAuth
// @Param        user_id      path  string  true  "User ID"
// @Param        relation_id  path  string  true  "Relation ID"
// @Success      200  {object}  models.Notif
// @Failure      404  {object}  Response  "Notif not found"
// @Failure      400  {object}  Response  "Bad Request"
// @Router       /notifs/{user_id}/{relation_id} [get]
func (n *NotifController) GetNotifsUser(c *gin.Context) {
	userID := c.Param("user_id")
	relationID := c.Param("relation_id")

	var notif models.Notif

	if err := n.Storage.Where("user_id = ? AND relation_id = ?", userID, relationID).First(&notif).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			handleError(c, http.StatusNotFound, "Notif not found", err)
			return
		}
		handleError(c, http.StatusBadRequest, "Something went wrong while fetching the notification", err)
		return
	}
	
	HandleResponse(c, http.StatusOK, notif)
}



