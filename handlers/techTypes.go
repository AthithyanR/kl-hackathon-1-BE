package handlers

import (
	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/valyala/fasthttp"
)

func GetAllTechTypes(ctx *fasthttp.RequestCtx) {
	var techTypes []models.TechType
	db.DB.Find(&techTypes)
	sendSuccessResponse(ctx, techTypes)
}
