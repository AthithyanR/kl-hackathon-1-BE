package handlers

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type BaseResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func sendSuccessResponse(ctx *fasthttp.RequestCtx, data any) {
	json.NewEncoder(ctx).Encode(&BaseResponse{Success: true, Data: data})
}

func sendFailureResponse(ctx *fasthttp.RequestCtx, data any) {
	json.NewEncoder(ctx).Encode(&BaseResponse{Success: false, Data: data})
}

func HealthCheck(ctx *fasthttp.RequestCtx) {
	sendSuccessResponse(ctx, nil)
}
