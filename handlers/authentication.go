package handlers

import (
	"encoding/json"

	"github.com/AthithyanR/kl-hackathon-1-BE/auth"
	"github.com/AthithyanR/kl-hackathon-1-BE/db"
	"github.com/AthithyanR/kl-hackathon-1-BE/models"
	"github.com/valyala/fasthttp"
)

func Authenticate(ctx *fasthttp.RequestCtx) {
	var requestBody models.User
	err := json.Unmarshal(ctx.PostBody(), &requestBody)
	if err != nil || requestBody.Email == "" || requestBody.Password == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	var existingUser models.User
	whereClause := &models.User{Email: requestBody.Email}
	db.DB.Where(whereClause).Find(&existingUser)

	if existingUser.Id == "" || existingUser.Password != requestBody.Password {
		sendFailureResponse(ctx, "Invalid username or password")
		return
	}

	token, err := auth.GenerateToken(&models.ClaimValues{Id: existingUser.Id, Email: existingUser.Email})

	if err != nil {
		sendFailureResponse(ctx, err.Error())
		return
	}

	sendSuccessResponse(ctx, token)
}
