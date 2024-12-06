package controllers

import (
	"errors"
	"log"
	"strconv"
	"tender_management/constants"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Response struct {
	Message interface{} `json:"message"`
}

func HandleResponse(c *gin.Context, statusCode int, message interface{}) {
	c.JSON(statusCode, Response{Message: message})

}

func handleError(c *gin.Context, statuscode int, message string, err error) {
	log.Printf("%s:%v", message, err)
	HandleResponse(c, statuscode, message)
}

func parseValidateDeliveryTime(deliveryTimeStr string) (*time.Time, error) {
	deliveryTime, err := time.Parse(constants.Layout, deliveryTimeStr)
	if err != nil {
		return nil, errors.New(constants.ErrFormatInput)
	}
	if !deliveryTime.After(time.Now()) {
		return nil, errors.New(constants.ErrDeadlinePassed)
	}
	return &deliveryTime, nil
}

func getPaginationParams(c *gin.Context) (int, int) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	return page, pageSize
}

func getByID(db *gorm.DB, id string, model interface{}) error {
	err := db.Where("id = ? AND deleted_at IS NULL", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("record not found")
		}
		return err
	}
	return nil
}
