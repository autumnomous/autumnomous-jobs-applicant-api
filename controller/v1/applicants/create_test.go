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

func Test_Applicant_Create_BookmarkJob_Create_Success(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.BookmarkJob))

	defer ts.Close()

	employer := testhelper.Helper_RandomEmployer(t)
	job := testhelper.Helper_RandomJob(employer, t)

	data := map[string]string{
		"jobid": job.PublicID,
	}

	applicant := testhelper.Helper_RandomApplicant(t)

	token, err := jwt.GenerateToken(applicant.PublicID)

	if err != nil {
		t.Fatal()
	}
	token = base64.StdEncoding.EncodeToString([]byte(token))
	requestBody, err := json.Marshal(data)

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

	result, err := httpClient.Do(request)

	assert.Nil(err)
	assert.Equal(int(http.StatusOK), result.StatusCode)
}

func Test_Applicant_Create_BookmarkJob_IncorrectMethod(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.BookmarkJob))

	defer ts.Close()

	data := map[string]string{
		"jobid": "",
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal(err)
	}

	methods := []string{"GET", "PUT"}

	for _, method := range methods {
		request, err := http.NewRequest(method, ts.URL, bytes.NewBuffer(requestBody))

		if err != nil {
			t.Fatal()
		}

		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(int(http.StatusMethodNotAllowed), response.StatusCode)
	}
}

func Test_Applicant_Create_BookmarkJob_NoJobID(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.BookmarkJob))

	defer ts.Close()

	applicant := testhelper.Helper_RandomApplicant(t)

	token, err := jwt.GenerateToken(applicant.PublicID)

	if err != nil {
		t.Fatal(err)
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))
	data := map[string]string{
		"jobid": "",
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(request)

	assert.Nil(err)
	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_Applicant_Create_BookmarkJob_Delete_Success(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(applicants.BookmarkJob))

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

	request := httptest.NewRequest("DELETE", ts.URL, bytes.NewBuffer(requestBody))

	token, err := jwt.GenerateToken(applicant.PublicID)

	if err != nil {
		t.Fatal(err)
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}

	response, err := client.Do(request)

	assert.Nil(err)
	assert.Equal(int(http.StatusOK), response.StatusCode)

}
