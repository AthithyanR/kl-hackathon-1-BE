package handlers

import (
	"encoding/json"

	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/AthithyanR/kl-hackathon-1-BE/utils"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm/clause"
)

func GetAllTechTypes(ctx *fasthttp.RequestCtx) {
	var techTypes []models.TechType
	if err := db.DB.Find(&techTypes).Error; err != nil {
		sendFailureResponse(ctx, nil)
	}
	sendSuccessResponse(ctx, techTypes)
}

func AddTechTypes(ctx *fasthttp.RequestCtx) {
	var techTypes []models.TechType
	var response []string
	err := json.Unmarshal(ctx.PostBody(), &techTypes)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if len(techTypes) == 0 {
		sendFailureResponse(ctx, "No resource provided")
		return
	}
	for i := 0; i < len(techTypes); i++ {
		id := utils.CanonicId()
		techTypes[i].Id = id
		response = append(response, id)
	}
	result := db.DB.Create(&techTypes)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, response)
}

func UpdateTechTypes(ctx *fasthttp.RequestCtx) {
	var techTypes []models.TechType
	err := json.Unmarshal(ctx.PostBody(), &techTypes)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if len(techTypes) == 0 {
		sendFailureResponse(ctx, "No resource provided")
		return
	}
	result := db.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&techTypes)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, nil)
}

func DeleteTechTypes(ctx *fasthttp.RequestCtx) {
	var techTypeIds []string
	err := json.Unmarshal(ctx.PostBody(), &techTypeIds)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if len(techTypeIds) == 0 {
		sendFailureResponse(ctx, "No ids provided")
		return
	}
	result := db.DB.Delete(&models.TechType{}, techTypeIds)
	if result.Error != nil {
		sendFailureResponse(ctx, result.Error.Error())
		return
	}
	sendSuccessResponse(ctx, nil)
}
