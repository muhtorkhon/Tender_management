package controllers

import (
	"net/http"
	"tender_management/constants"
	"tender_management/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OfferController struct {
	Storage *gorm.DB
}

func NewOfferController(storage *gorm.DB) *OfferController {
	return &OfferController{
		Storage: storage,
	}
}

// @Summary      Create a new offer
// @Description  This endpoint creates a new offer with the provided details.
// @Tags         offers
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      models.OffersRequest  true  "Offer Request Body"
// @Success      201   {object}  models.Offers
// @Failure      400   {object}  Response  "Failed to parse request body"
// @Failure      500   {object}  Response "Failed to create offer"
// @Router       /offers [post]
func (o *OfferController) CreateOffer(c *gin.Context) {
	var body models.OffersRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		handleError(c, http.StatusBadRequest, "Failed to parse request body", err)
		return
	}

	deliveryTime, err := parseValidateDeliveryTime(body.DeliveryTime)
	if err != nil {
		handleError(c, http.StatusBadRequest, "invalid delivery time format", err)
		return
	}
	offer := models.Offers{
		TenderID:     body.TenderID,
		ContractorID: body.ContractorID,
		Price:        body.Price,
		DeliveryTime: deliveryTime,
		Comments:     body.Comments,
		Status:       body.Status,
	}

	if err := o.Storage.Create(&offer).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create offer", err)
		return
	}

	HandleResponse(c, http.StatusCreated, offer)
}

// @Summary      Get all offers
// @Description  Retrieve a paginated list of all offers.
// @Tags         offers
// @Security 	 BearerAuth
// @Produce      json
// @Param        page      query     int  false  "Page number"
// @Param        pageSize  query     int  false  "Page size"
// @Success      200       {array}   models.Offers
// @Failure      500       {object}  Response  "Failed to fetch offers"
// @Router       /offers [get]
func (o *OfferController) GetAllOffers(c *gin.Context) {

	page, pageSize := getPaginationParams(c)

	offset := (page - 1) * pageSize

	var offer []models.Offers

	if err := o.Storage.Where("deleted_at IS NULL").Limit(pageSize).Offset(offset).Find(&offer).Error; err != nil {
		handleError(c, http.StatusInternalServerError, constants.ErrRecordNotFound, err)
		return
	}

	HandleResponse(c, http.StatusOK, gin.H{
		"count": len(offer),
		"offers": offer,
	})
}

// @Summary      Get a specific offer
// @Description  Retrieve an offer by its ID.
// @Tags         offers
// @Security 	 BearerAuth
// @Produce      json
// @Param        contractor_id  path      string  true  "Contractor ID"
// @Success      200            {object}  models.Offers
// @Failure      404            {object}  Response "Offer not found"
// @Router       /offers/{contractor_id} [get]
func (o *OfferController) GetOffer(c *gin.Context) {
	contractorID := c.Param("contractor_id")
	var offer models.Offers

	if err := o.Storage.Where("contractor_id = ? AND deleted_at IS NULL", contractorID).First(&offer).Error; err != nil {
		handleError(c, http.StatusNotFound, constants.ErrRecordNotFound, err)
		return
	}
	HandleResponse(c, http.StatusOK, offer)
}

// GetFilterSort    godoc
// @Summary 		Get filtered and sorted offers with pagination
// @Description 	This endpoint retrieves a list of offers with pagination, sorted by price and delivery time.
// It also provides the total number of offers matching the filters (excluding deleted offers).
// @Tags            offers
// @Security 		BearerAuth
// @Accept  		json
// @Produce 		json
// @Param 			page query int false "Page number"
// @Param 			pageSize query int false "Number of offers per page"
// @Success 		200 {object} Response "Successful response with offers and total count"
// @Failure 		500 {object} Response "Internal server error"
// @Router          /offers/sorted [get]
func (o *OfferController) GetFilterSort(c *gin.Context) {
	var offers []models.Offers
	var totalRecords int64

	page, pageSize := getPaginationParams(c)

	offset := (page - 1) * pageSize

	if err := o.Storage.Model(&offers).
		Where("deleted_at IS NULL").Count(&totalRecords).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to count offers", err)
		return
	}

	if err := o.Storage.Model(&offers).
		Where("deleted_at IS NULL").Limit(pageSize).
		Offset(offset).Order("price ASC").Order("delivery_time ASC").
		Find(&offers).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to fetch offers", err)
		return
	}
		
	HandleResponse(c, http.StatusOK, gin.H{
		"totalRecords": totalRecords,
		"currentPage":  page,
		"pageSize":     pageSize,
		"offers":       offers,
	})
}

// GetMaxMinfilter godoc
// @Summary     Get min, max prices and delivery times with filtered count
// @Description Retrieves offers with minimum and maximum prices and delivery times, along with their counts.
// @Tags        offers
// @Security 	BearerAuth
// @Accept      json
// @Produce     json
// @Success     200 {object} Response "Details of min/max offers and counts"
// @Failure     500 {object} Response "Error message"
// @Router      /offers/filter [get]
func (o *OfferController) GetMaxMinFilter(c *gin.Context) {
	var offers []models.Offers
	var stats models.Stats

	query := `
		SELECT 
			MIN(price) AS min_price, 
			MAX(price) AS max_price, 
			MIN(delivery_time) AS min_delivery, 
			MAX(delivery_time) AS max_delivery,
			COUNT(*) AS total_records
		FROM offers
		WHERE deleted_at IS NULL
	`
	if err := o.Storage.Raw(query).Scan(&stats).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to fetch statistics", err)
		return
	}

	if err := o.Storage.Model(&offers).
		Where("deleted_at IS NULL").
		Where("price = ? OR price = ? OR delivery_time = ? OR delivery_time = ?",
			stats.MinPrice, stats.MaxPrice, stats.MinDelivery, stats.MaxDelivery).
		Find(&offers).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to fetch filtered offers", err)
		return
	}

	HandleResponse(c, http.StatusOK, gin.H{
		"stats":  stats,
		"count":  len(offers),
		"offers": offers,
	})
}

// @Summary      Update an existing offer
// @Description  Update the details of an offer by its ID.
// @Tags         offers
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string              true  "Offer ID"
// @Param        body body      models.OffersRequest true  "Updated offer details"
// @Success      200  {object}  models.OffersRequest
// @Failure      400  {object}  Response  "Failed to parse request body"
// @Failure      404  {object}  Response  "Offer not found"
// @Router       /offers/{id} [put]
func (o *OfferController) UpdateOffer(c *gin.Context) {
	id := c.Param("id")

	var newOffer models.OffersRequest

	if err := c.ShouldBindJSON(&newOffer); err != nil {
		handleError(c, http.StatusBadRequest, "Failed to parse request newOffer", err)
		return
	}

	deliveryTime, err := parseValidateDeliveryTime(newOffer.DeliveryTime)
	if err != nil {
		handleError(c, http.StatusBadRequest, "invalid delivery time format", err)
		return
	}

	updatefields := map[string]interface{}{
		"tender_id":     newOffer.TenderID,
		"contractor_id": newOffer.ContractorID,
		"price":         newOffer.Price,
		"delivery_time": deliveryTime,
		"comments":      newOffer.Comments,
		"status":        newOffer.Status,
	}

	if err := o.Storage.Model(&models.Offers{}).Where("id = ?", id).Updates(updatefields).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to update offer", err)
		return
	}

	HandleResponse(c, http.StatusOK, updatefields)
}

// @Summary      Soft delete an offer
// @Description  Mark an offer as deleted by setting the DeletedAt field.
// @Tags         offers
// @Security 	 BearerAuth
// @Param        id   path      string  true  "Offer ID"
// @Success      200  {string}  string  "Offer deleted successfully"
// @Failure      404  {object}  Response  "Offer not found"
// @Router       /offers/{id} [delete]
func (o *OfferController) DeleteOffer(c *gin.Context) {
	id := c.Param("id")

	var offer models.Offers

	if err := getByID(o.Storage, id, &offer); err != nil {
		handleError(c, http.StatusNotFound, constants.ErrRecordNotFound, err)
		return
	}

	if offer.DeletedAt != nil {
		handleError(c, http.StatusNotFound, "Offer already deleted", nil)
		return
	}

	new := time.Now()
	offer.DeletedAt = &new

	if err := o.Storage.Save(&offer).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to soft delete offer", err)
		return
	}

	HandleResponse(c, http.StatusOK, "Offer deleted successfully")
}

// @Summary      Restore a soft-deleted offer
// @Description  Restore an offer that was previously soft deleted.
// @Tags         offers
// @Security 	 BearerAuth
// @Param        id   path      string  true  "Offer ID"
// @Success      200  {string}  string  "Offer restored successfully"
// @Failure      404  {object}  Response  "Offer not found"
// @Router       /offers/{id}/restore [patch]
func (t *OfferController) RestoreOffer(c *gin.Context) {
	id := c.Param("id")

	var offer models.Offers

	if err := t.Storage.Unscoped().Where("id = ?", id).First(&offer).Error; err != nil {
		handleError(c, http.StatusNotFound, "Failed to find offer or it may not be soft deleted", err)
		return
	}

	if offer.DeletedAt == nil {
		handleError(c, http.StatusBadRequest, "Offer is not soft deleted", nil)
		return
	}

	offer.DeletedAt = nil

	if err := t.Storage.Save(&offer).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to restore offer", err)
		return
	}

	HandleResponse(c, http.StatusOK, "Offer restored successfully")
}
