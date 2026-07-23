package utils

import (
	"gamebiller/models"
	"net/http"
	"strings"
)

func GenRequestHeader(req *http.Request, reqHeader models.ReqHeader) *http.Request {

	for _, v := range reqHeader.Header {
		if v.IsUpCase {
			req.Header.Set(strings.ToUpper(v.Key), v.Val)
			continue
		}
		req.Header.Set(v.Key, v.Val)
	}

	return req
}
