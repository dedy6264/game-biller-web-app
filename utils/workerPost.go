package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"gamebiller/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func WorkerRequestPOST(tipeRequest, urlApi string, requestBody interface{}, requestHeader models.ReqHeader, timeout time.Duration) (result []byte, statusCode int, err error) {

	bodyRequest, _ := json.Marshal(requestBody)

	// CREATING REQUEST HTTP
	reqHTTP, err := http.NewRequest("POST", urlApi, bytes.NewBuffer(bodyRequest))
	if err != nil {
		return result, statusCode, err
	}
	// END CREATING REQUEST HTTP

	reqHTTP = GenRequestHeader(reqHTTP, requestHeader)

	// Set Content-type header
	if tipeRequest == "urlencoded" {
		reqHTTP.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if tipeRequest == "json" {
		reqHTTP.Header.Add("Content-Type", "application/json")
	}
	reqHTTP.Header.Add("Content-Length", strconv.FormatInt(reqHTTP.ContentLength, 10))
	reqHTTP.Header.Set("Connection", "close")

	// log.Printf("Request Header: %v\n", reqHTTP.Header)

	if bodyRequest != nil {
		defer reqHTTP.Body.Close()
	}
	reqHTTP.Close = true
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		return result, statusCode, err
	}
	resp.Header.Set("Connection", "close")
	defer resp.Body.Close()
	resp.Close = true

	statusCode = resp.StatusCode

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, statusCode, err
	}
	log.Println("Worker Request Url : ", urlApi)
	log.Println("Worker Request Data : ", string(bodyRequest))
	log.Println("Worker Response Data : ", string(result))
	return result, statusCode, nil
}
