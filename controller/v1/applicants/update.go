package applicants

import (
	"encoding/json"
	applicants "jobs-applicant-api/shared/repository/applicants"
	"jobs-applicant-api/shared/response"
	"jobs-applicant-api/shared/services/security/jwt"
	"jobs-applicant-api/shared/services/zipcode"
	"log"
	"net/http"
	"os"
)

type updatePasswordCredentials struct {
	Password    string `json:"password"`
	NewPassword string `json:"newpassword"`
}

type updateAccountData struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	Zipcode     string `json:"zipcode"`
	// Facebook     string `json:"facebook"`
	// Twitter      string `json:"twitter"`
	// Instagram    string `json:"instagram"`
	// Bio          string `json:"bio"`
}

type updateJobPreferencesData struct {
	DesiredCities []map[string]interface{} `json:"desiredcities"`
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, response.FriendlyError)
		return
	}

	var credentials updatePasswordCredentials
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	if credentials.Password == "" || credentials.NewPassword == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	publicID := jwt.GetUserClaim(r)

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	updated, err := repository.UpdateApplicantPassword(publicID, credentials.Password, credentials.NewPassword)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	if updated {
		response.SendJSONMessage(w, http.StatusOK, response.Success)
		return
	} else {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, response.FriendlyError)
		return
	}

	publicID := jwt.GetUserClaim(r)

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	var data updateAccountData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	gateway := zipcode.NewZipCodeGateway(os.Getenv("ZIPCODESERVICES_API_KEY"))

	zip_code, err := gateway.GetZipCode(data.Zipcode)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	applicant, err := repository.UpdateApplicantAccount(publicID, data.FirstName, data.LastName, data.Email, data.PhoneNumber, data.Address, data.City, data.State, data.Zipcode, zip_code.Latitude, zip_code.Longitude)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, applicant)
}

func UpdateJobPreferences(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, response.FriendlyError)
		return
	}

	publicID := jwt.GetUserClaim(r)

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	var data updateJobPreferencesData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	err = repository.UpdateApplicantJobPreferences(publicID, data.DesiredCities)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

}
