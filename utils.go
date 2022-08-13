package main

import (
	"fmt"
	"os"

	"github.com/valyala/fasthttp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getDb() *gorm.DB {
	template := "%s:%s@tcp(%s:3306)/%s?charset=utf8mb4"

	dsn := fmt.Sprintf(template, "root", getenv("password", "atr"), getenv("host", "athi.fun"), getenv("dbname", "entretien"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(`Unable to establish database connection :-(`)
	}

	return db
}

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func middleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
		fmt.Println(ctx)
		h(ctx)
	}
}
