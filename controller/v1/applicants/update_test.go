package applicants_test

import (
	"autumnomous-jobs-applicant-api/controller/v1/applicants"
	"autumnomous-jobs-applicant-api/shared/services/security/encryption"
	"autumnomous-jobs-applicant-api/shared/services/security/jwt"
	"autumnomous-jobs-applicant-api/shared/testhelper"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	testhelper.Init()
}

func Test_Applicant_UpdatePassword_Success(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(applicants.UpdatePassword))

	data := map[string]string{
		"password":    string(encryption.GeneratePassword(9)),
		"newpassword": string(encryption.GeneratePassword(9)),
	}

	applicant := &testhelper.TestApplicant{FirstName: "First", LastName: "Last", Email: fmt.Sprintf("email-%s@site.com", encryption.GeneratePassword(9))}

	hashedPassword, err := encryption.HashPassword([]byte(data["password"]))

	if err != nil {
		t.Fatal()
	}

	applicant.HashedPassword = hashedPassword

	applicant = testhelper.Helper_CreateApplicant(applicant, t)

	token, err := jwt.GenerateToken(applicant.PublicID)

	if err != nil {
		t.Fatal()
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal()
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))
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
	assert.Equal(int(http.StatusOK), response.StatusCode)
}

func Test_Applicant_UpdatePassword_IncorrectMethod(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(applicants.UpdatePassword))

	defer ts.Close()

	data := map[string]string{
		"password":    "",
		"newpassword": "",
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal(err)
	}

	methods := []string{"GET", "PUT", "DELETE"}

	for _, method := range methods {

		request, err := http.NewRequest(method, ts.URL, bytes.NewBuffer(requestBody))

		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}

		response, err := client.Do(request)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(http.StatusMethodNotAllowed, response.StatusCode)
	}
}
func Test_Applicant_UpdatePassword_IncorrectDataReceived_NoPassword(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(applicants.UpdatePassword))

	data := map[string]string{
		"password":    "",
		"newpassword": string(encryption.GeneratePassword(9)),
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal()
	}

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_Applicant_UpdatePassword_IncorrectDataReceived_NoNewPassword(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(applicants.UpdatePassword))

	data := map[string]string{
		"password":    string(encryption.GeneratePassword(9)),
		"newpassword": "",
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal()
	}

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_Applicant_UpdatePassword_IncorrectDataReceived_NoToken(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(applicants.UpdatePassword))

	data := map[string]string{
		"password":    string(encryption.GeneratePassword(9)),
		"newpassword": string(encryption.GeneratePassword(9)),
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal()
	}

	httpClient := &http.Client{}

	token, err := jwt.GenerateToken("")

	if err != nil {
		t.Fatal()
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_Applicant_UpdatePassword_IncorrectDataReceived_NoData(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(applicants.UpdatePassword))

	data := map[string]string{
		"password":    "",
		"newpassword": "",
	}

	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal()
	}

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(requestBody))

	if err != nil {
		t.Fatal()
	}

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		t.Fatal()
	}

	assert.Equal(int(http.StatusBadRequest), response.StatusCode)

}

func Test_Applicant_UpdateAccount_IncorrectDataReceived_NoToken(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(applicants.UpdateAccount))

	request, err := http.NewRequest("POST", ts.URL, nil)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	httpClient := http.Client{}

	token, err := jwt.GenerateToken("")

	if err != nil {
		t.Fatal()
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := httpClient.Do(request)

	assert.Nil(err)
	assert.Equal(int(http.StatusBadRequest), response.StatusCode)
}

func Test_Applicant_UpdateAccount_CorrectDataReceived(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer((http.HandlerFunc(applicants.UpdateAccount)))

	applicant := &testhelper.TestApplicant{FirstName: "First", LastName: "Last", Email: fmt.Sprintf("email-%s@site.com", encryption.GeneratePassword(9)), Password: string(encryption.GeneratePassword(9))}

	hashedPassword, err := encryption.HashPassword([]byte(applicant.Password))

	if err != nil {
		t.Fatal()
	}

	applicant.HashedPassword = hashedPassword

	applicant = testhelper.Helper_CreateApplicant(applicant, t)

	tests := map[string]map[string]string{
		"NewBio": {
			"firstname": "First",
			"lastname":  "Last",
			"email":     applicant.Email,
		},
		"New First Name": {
			"firstname": "NewFirst",
			"lastname":  "Last",
			"email":     applicant.Email,
		},
		"New Last Name": {
			"firstname": "NewFirst",
			"lastname":  "NewLast",
			"email":     applicant.Email,
		},
		"New Email": {
			"firstname": "NewFirst",
			"lastname":  "NewLast",
			"email":     fmt.Sprintf("new-email-%s@site.com", encryption.GeneratePassword(9)),
		},
	}

	for _, test := range tests {

		data, err := json.Marshal(test)

		if err != nil {
			log.Println(err)
			t.Fatal()
		}

		request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(data))
		if err != nil {
			t.Fatal()
		}

		token, err := jwt.GenerateToken(applicant.PublicID)

		if err != nil {
			t.Fatal()
		}

		token = base64.StdEncoding.EncodeToString([]byte(token))
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

		assert.Equal(int(http.StatusOK), response.StatusCode)
		assert.Equal(test["firstname"], result["firstname"])
		assert.Equal(test["lastname"], result["lastname"])
		assert.Equal(test["email"], result["email"])

	}

}

func Test_Applicant_UpdateJobPreferences_Correct(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(applicants.UpdateJobPreferences))

	data := map[string][]map[string]string{
		"desiredcities": {
			{"city": "Cleveland",
				"state":     "Alabama",
				"country":   "US",
				"latitude":  "34.0165",
				"longitude": "-86.5538",
				"id":        "1",
				"text":      "Cleveland, Alabama, US",
			},
		},
	}

	requestBody, err := json.Marshal(data)

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

	assert.Equal(int(http.StatusOK), response.StatusCode)

}

func Test_Applicant_UpdateJobPreferences_IncorrectMethod(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(applicants.UpdateJobPreferences))

	defer ts.Close()

	data := map[string][]map[string]string{
		"desiredcities": {
			{"city": "Cleveland",
				"state":     "Alabama",
				"country":   "US",
				"latitude":  "34.0165",
				"longitude": "-86.5538",
				"id":        "1",
				"text":      "Cleveland, Alabama, US",
			},
		},
	}
	requestBody, err := json.Marshal(data)

	if err != nil {
		t.Fatal(err)
	}

	methods := []string{"GET", "PUT", "DELETE"}

	for _, method := range methods {

		request, err := http.NewRequest(method, ts.URL, bytes.NewBuffer(requestBody))

		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}

		response, err := client.Do(request)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(http.StatusMethodNotAllowed, response.StatusCode)
	}
}
