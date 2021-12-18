package applicants_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"jobs-applicant-api/controller/v1/applicants"
	"jobs-applicant-api/shared/services/security/jwt"
	"jobs-applicant-api/shared/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	testhelper.Init()
}

func Test_Applicant_GetApplicant_Correct(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.GetApplicant))

	defer ts.Close()

	applicant := testhelper.Helper_RandomApplicant(t)

	token, err := jwt.GenerateToken(applicant.PublicID)

	if err != nil {
		t.Fatal()
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))

	request, err := http.NewRequest("GET", ts.URL, nil)

	if err != nil {
		t.Fatal()
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	var result map[string]interface{}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	if err != nil {
		t.Fatal()
	}

	assert.Nil(err)
	assert.Equal(int(http.StatusOK), response.StatusCode)
	// assert.Contains(result, "publicid")

}

func Test_Applicant_GetApplicant_Incorrect(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.GetApplicant))

	defer ts.Close()

	request, err := http.NewRequest("GET", ts.URL, nil)

	if err != nil {
		t.Fatal()
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer ")

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	var result map[string]interface{}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	if err != nil {
		t.Fatal()
	}

	assert.Nil(err)
	assert.Equal(int(http.StatusBadRequest), response.StatusCode)
	// assert.Contains(result, "publicid")

}

func Test_Applicant_GetAutocompleteLocationData_Correct(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.GetAutocompleteLocationData))

	defer ts.Close()

	test := map[string]string{
		"chars": "Cleve",
	}

	requestBody, err := json.Marshal(test)

	if err != nil {
		t.Fatal()
	}

	applicant := testhelper.Helper_RandomApplicant(t)

	token, err := jwt.GenerateToken(applicant.PublicID)

	if err != nil {
		t.Fatal()
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal()
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	var result map[string]interface{}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	if err != nil {
		t.Fatal()
	}

	assert.Nil(err)
	assert.Equal(int(http.StatusOK), response.StatusCode)
	// assert.Contains(result, "publicid")

}

func Test_Applicant_GetAutocompleteLocationData_IncorrectMethod(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.GetAutocompleteLocationData))

	defer ts.Close()

	methods := []string{"GET", "DELETE", "PUT"}

	for _, method := range methods {
		request, err := http.NewRequest(method, ts.URL, nil)

		if err != nil {
			t.Fatal()
		}

		httpClient := &http.Client{}

		response, err := httpClient.Do(request)

		assert.Nil(err)
		assert.Equal(int(http.StatusMethodNotAllowed), response.StatusCode)
	}

}

func Test_Applicant_GetJobs_Correct(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(applicants.GetJobs))

	defer ts.Close()

	applicant := testhelper.Helper_RandomApplicant(t)

	token, err := jwt.GenerateToken(applicant.PublicID)

	if err != nil {
		t.Fatal()
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))

	request, err := http.NewRequest("POST", ts.URL, nil)

	if err != nil {
		t.Fatal()
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	var result map[string]interface{}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	if err != nil {
		t.Fatal()
	}

	assert.Nil(err)
	assert.Equal(int(http.StatusOK), response.StatusCode)

}

func Test_Applicant_GetJobs_IncorrectMethod(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.GetJobs))

	defer ts.Close()

	methods := []string{"GET", "DELETE", "PUT"}

	for _, method := range methods {
		request, err := http.NewRequest(method, ts.URL, nil)

		if err != nil {
			t.Fatal()
		}

		httpClient := &http.Client{}

		response, err := httpClient.Do(request)

		assert.Nil(err)
		assert.Equal(int(http.StatusMethodNotAllowed), response.StatusCode)
	}

}
