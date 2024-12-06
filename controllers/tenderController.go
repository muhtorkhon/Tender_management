package controllers

import (
	"net/http"
	"tender_management/constants"
	"tender_management/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenderController struct {
	Storage *gorm.DB
}

func NewTenderController(storage *gorm.DB) *TenderController {
	return &TenderController{
		Storage: storage,
	}
}

// CreateTender 	godoc
// @Summary 		Create a new tender
// @Description 	Creates a new tender with the provided details
// @Security 		BearerAuth
// @Tags 			tender
// @Accept 			json
// @Produce 		json
// @Param 			body body models.TenderRequest true "Tender Request Body"
// @Success 		201 {object} models.Tenders
// @Failure 		400 {object} Response "Bad Request"
// @Failure 		500 {object} Response "Internal Server Error"
// @Router 			/tenders [post]
func (t *TenderController) CreateTender(c *gin.Context) {
	var body models.TenderRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Failed to parse request body", err)
		return
	}

	deadline, err := parseValidateDeliveryTime(body.Deadline)
	if err != nil {
		handleError(c, http.StatusBadRequest, "invalid deadline time format", err)
		return
	}

	tender := models.Tenders{
		Title:       body.Title,
		Description: body.Description,
		Deadline:    deadline,
		Budget:      body.Budget,
		FileURL:     body.FileURL,
		ClientID:    body.ClientID,
	}

	if err := t.Storage.Create(&tender).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create tender", err)
		return
	}

	HandleResponse(c, http.StatusCreated, tender)
}

// GetTenders 	godoc
// @Summary 	Get tenders by client ID
// @Description Retrieve all tenders for a specific client
// @Security 	BearerAuth
// @Tags 		tender
// @Produce 	json
// @Param 		client_id path string true "Client ID"
// @Success 	200 {object} models.Tenders
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/tenders/{client_id} [get]
func (t *TenderController) GetTenders(c *gin.Context) {
	clientID := c.Param("client_id")

	var tenders models.Tenders
	if err := t.Storage.Where("client_id = ? and deleted_at IS NULL", clientID).First(&tenders).Error; err != nil {
		handleError(c, http.StatusNotFound, "Failed to fetch tenders", err)
		return
	}

	HandleResponse(c, http.StatusOK, tenders)
}

// GetAllTenders 	godoc
// @Summary 		Get all tenders with pagination
// @Description 	Retrieve all tenders with pagination support
// @Tags 			tender
// @Produce 		json
// @Param 			page query int false "Page number"
// @Param 			pageSize query int false "Page size"
// @Success 		200 {array} models.Tenders
// @Failure 		500 {object} Response "Internal Server Error"
// @Router 			/tenders [get]
func (t *TenderController) GetAllTenders(c *gin.Context) {

	page, pageSize := getPaginationParams(c)

	offset := (page - 1) * pageSize

	var tenders []models.Tenders
	if err := t.Storage.Where("deleted_at IS NULL").Limit(pageSize).Offset(offset).Find(&tenders).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to fetch tenders", err)
		return
	}

	HandleResponse(c, http.StatusOK, gin.H{
		"count": len(tenders),
		"offers": tenders,
	})
}

// UpdateTender 	godoc
// @Summary 		Update an existing tender
// @Description 	Updates the details of an existing tender
// @Tags 			tender
// @Security 		BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Tender ID"
// @Param 			body body models.TenderRequest true "Updated Tender Body"
// @Success 		200 {object} Response
// @Failure 		400 {object} Response "Bad Request"
// @Failure 		404 {object} Response "Tender not found"
// @Failure 		500 {object} Response "Internal Server Error"
// @Router 			/tenders/{id} [put]
func (t *TenderController) UpdateTender(c *gin.Context) {
	id := c.Param("id")

	var newtender models.TenderRequest

	if err := c.ShouldBindJSON(&newtender); err != nil {
		handleError(c, http.StatusBadRequest, "Failed to parse request newtender", err)
		return
	}

	deadline, err := parseValidateDeliveryTime(newtender.Deadline)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid deadline format", err)
		return
	}

	updatefields := map[string]interface{}{
		"title":       newtender.Title,
		"description": newtender.Description,
		"deadline":    deadline,
		"budget":      newtender.Budget,
		"file_url":    newtender.FileURL,
		"client_id":   newtender.ClientID,
	}

	if err := t.Storage.Model(&models.Tenders{}).Where("id = ?", id).Updates(updatefields).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to update offer", err)
		return
	}

	HandleResponse(c, http.StatusOK, updatefields)
}

// DeleteTender 	godoc
// @Summary 		Soft delete a tender by ID
// @Description 	Marks a tender as deleted by setting the `DeletedAt` timestamp.
// @Tags 			tender
// @Security 		BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Tender ID"
// @Success 		200 {string} string "Tender deleted successfully"
// @Failure 		400 {object} Response "Invalid request"
// @Failure 		404 {object} Response "Tender not found or already deleted"
// @Failure 		500 {object} Response "Internal server error"
// @Router 			/tenders/{id} [delete]
func (t *TenderController) DeleteTender(c *gin.Context) {
	id := c.Param("id")

	var tender models.Tenders

	if err := getByID(t.Storage, id, &tender); err != nil {
		handleError(c, http.StatusNotFound, constants.ErrRecordNotFound, err)
		return
	}

	if tender.DeletedAt != nil {
		handleError(c, http.StatusNotFound, "Tender has been deleted", nil)
		return
	}

	new := time.Now()
	tender.DeletedAt = &new

	if err := t.Storage.Save(&tender).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to soft delete tender", err)
		return
	}

	HandleResponse(c, http.StatusOK, "Tender deleted successfully")
}

// RestoreTender 	godoc
// @Summary 		Restore a soft deleted tender by ID
// @Description 	Restores a tender by removing the `DeletedAt` timestamp, making it active again.
// @Tags 			tender
// @Security 		BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Tender ID"
// @Success 		200 {string} Response "Tender restored successfully"
// @Failure 		404 {object} Response "Tender not found or already active"
// @Failure 		500 {object} Response "Internal server error"
// @Router 			/tenders/restore/{id} [patch]
func (t *TenderController) RestoreTender(c *gin.Context) {
	id := c.Param("id")

	var tender models.Tenders

	if err := t.Storage.Unscoped().Where("id = ?", id).First(&tender).Error; err != nil {
		handleError(c, http.StatusNotFound, "Failed to find tender or it may not be soft deleted", err)
		return
	}

	if tender.DeletedAt == nil {
		handleError(c, http.StatusBadRequest, "Tender is not soft deleted", nil)
		return
	}

	tender.DeletedAt = nil

	if err := t.Storage.Save(&tender).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to restore tender", err)
		return
	}

	HandleResponse(c, http.StatusOK, "Tender restored successfully")
}
