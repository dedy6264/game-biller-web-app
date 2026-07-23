package utils

import (
	"bytes"
	"fmt"
	"gamebiller/models"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func WorkerRequestPOSTFromData(
	urlApi string,
	requestBody interface{},
	requestHeader models.ReqHeader,
	timeout time.Duration,
) (result []byte, statusCode int, err error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	val := reflect.ValueOf(requestBody)

	// Handle pointer
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {

	case reflect.Map:
		for _, key := range val.MapKeys() {

			strKey := fmt.Sprintf("%v", key.Interface())
			strVal := fmt.Sprintf("%v", val.MapIndex(key).Interface())

			if err := writer.WriteField(strKey, strVal); err != nil {
				return nil, 0, err
			}
		}

	case reflect.Struct:

		typ := val.Type()

		for i := 0; i < val.NumField(); i++ {

			field := typ.Field(i)

			fieldName := field.Name

			tag := field.Tag.Get("form")

			if tag == "" {
				tag = field.Tag.Get("json")
			}

			if tag != "" && tag != "-" {

				// handle json:"data_no,omitempty"
				tag = strings.Split(tag, ",")[0]

				fieldName = tag
			}
			strVal := fmt.Sprintf("%v", val.Field(i).Interface())

			if err := writer.WriteField(fieldName, strVal); err != nil {
				return nil, 0, err
			}
		}

	default:
		return nil, 0, fmt.Errorf("requestBody must be map or struct")
	}

	writer.Close()

	reqHTTP, err := http.NewRequest(
		http.MethodPost,
		urlApi,
		&b,
	)
	if err != nil {
		return nil, 0, err
	}

	reqHTTP.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	reqHTTP = GenRequestHeader(reqHTTP, requestHeader)

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(reqHTTP)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	result, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, statusCode, err
	}

	log.Println("Request URL :", urlApi)
	log.Println("Request Body :", b.String())
	log.Println("Response :", string(result))

	return result, statusCode, nil
}
