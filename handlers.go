package main

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type BaseResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func healthCheck(ctx *fasthttp.RequestCtx) {
	response := BaseResponse{true, nil}
	json.NewEncoder(ctx).Encode(response)
}

// questions
func getQuestionsByTechType(ctx *fasthttp.RequestCtx) {
	var questions []Question
	DB.Where(&Question{TechType: ctx.UserValue("techType").(string)}).Find(&questions)
	json.NewEncoder(ctx).Encode(&BaseResponse{Success: true, Data: questions})
}
