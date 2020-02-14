package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestDefaultGet - checks if the root "/" GET endpoint returns "Default Get"
func TestDefaultGet(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handleGet(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Default Get", rec.Body.String())
	}
}

// TestJsonPost - checks if POST "/json" will return the json body it is sent
func TestJsonPost(t *testing.T) {
	mcPostBody := map[string]interface{}{
		"question_text": "Is this a test post for MultiQuestion?",
		"first_name":    "Olin",
	}
	body, _ := json.Marshal(mcPostBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/json")
	c := e.NewContext(req, rec)

	h := handleJSON(c)

	recResponse, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		log.Println("Error reading Request Body into ByteArray", err)
	}

	message := map[string]interface{}{}
	err = json.Unmarshal(recResponse, &message)
	if err != nil {
		log.Println("Error decoding Response Body", err)
	}

	if assert.NoError(t, h) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, mcPostBody, message)
	}
}

// TestJsonPost - checks if POST "/donation_jsonapi" will return the json:api
// request data it was sent using the Donation type
func TestJsonApiPost(t *testing.T) {
	testDonation := Donation{
		Name:  "Jeff",
		Value: 88.95,
	}

	out := bytes.NewBuffer(nil)
	jsonapi.MarshalPayload(out, &testDonation)

	req := httptest.NewRequest(http.MethodPost, "/donation_jsonapi", out)
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/json")

	e := echo.New()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handleDonationJSONAPI(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		donation := new(Donation)
		err := jsonapi.UnmarshalPayload(rec.Body, donation)
		if err != nil {
			log.Println("Error decoding Response Body for testing", err)
		}

		assert.Equal(t, testDonation, *donation)
	}
}
