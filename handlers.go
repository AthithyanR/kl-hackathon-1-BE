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
