package controllers

import (
	"net/http"
	"urlshortner/utils"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	utils.FormatResponse("true", "Success is not final; failure is not fatal: It is the courage to continue that counts.", w)
}
