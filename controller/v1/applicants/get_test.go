package applicants_test

import (
	"autumnomous-jobs-applicant-api/controller/v1/applicants"
	"autumnomous-jobs-applicant-api/shared/services/security/jwt"
	"autumnomous-jobs-applicant-api/shared/testhelper"
	"bytes"
	"encoding/base64"
	"encoding/json"
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

	var result []map[string]interface{}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)
	assert.Nil(err)
	assert.Equal(int(http.StatusOK), response.StatusCode)

}

func Test_Applicant_GetJobs_IncorrectMethod(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.GetJobs))

	defer ts.Close()

	methods := []string{"POST", "DELETE", "PUT"}

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

func Test_Applicant_GetJobsByRadius_Correct(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(applicants.GetJobsByRadius))

	defer ts.Close()

	applicant := testhelper.Helper_RandomApplicant(t)

	token, err := jwt.GenerateToken(applicant.PublicID)

	if err != nil {
		t.Fatal()
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))

	test := map[string]string{
		"zipcode": "44145",
		"radius":  "10",
	}

	requestBody, err := json.Marshal(test)

	if err != nil {
		t.Fatal()
	}

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

}

func Test_Applicant_GetJobsByRadius_IncorrectMethod(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.GetJobsByRadius))

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

func Test_Applicant_GetApplicantJobBookmark_IncorrectMethod(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.GetApplicantJobBookmark))

	defer ts.Close()

	for _, method := range []string{"GET", "DELETE", "PUT"} {

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

func Test_Applicant_GetApplicantJobBookmark_Success(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.GetApplicantJobBookmark))

	defer ts.Close()

	applicant := testhelper.Helper_RandomApplicant(t)
	employer := testhelper.Helper_RandomEmployer(t)
	job := testhelper.Helper_RandomJob(employer, t)

	data := map[string]string{
		"jobid": job.PublicID,
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	token, err := jwt.GenerateToken(applicant.PublicID)
	token = base64.StdEncoding.EncodeToString([]byte(token))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	assert.Nil(err)
	assert.NotNil(response)
}
