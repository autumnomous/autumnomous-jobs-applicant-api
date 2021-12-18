package applicants

import (
	"encoding/json"
	"jobs-applicant-api/shared/repository/applicants"
	"jobs-applicant-api/shared/repository/jobs"
	"jobs-applicant-api/shared/response"
	"jobs-applicant-api/shared/services/security/jwt"
	"jobs-applicant-api/shared/services/zipcode"
	"log"
	"net/http"
	"os"
)

type AutocompleteLocationData struct {
	Characters string `json:"chars"`
}

func GetApplicant(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	publicID := jwt.GetUserClaim(r)

	if publicID == "" {

		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	employer, err := repository.GetApplicant(publicID)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, employer)
}

func GetAutocompleteLocationData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var details AutocompleteLocationData

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&details)

	if details.Characters == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.MissingRequiredValue)
		return
	}

	gateway := zipcode.NewZipCodeGateway(os.Getenv("ZIPCODESERVICES_API_KEY"))

	data, err := gateway.GetAutoComplete(details.Characters)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, data)

}

func GetJobs(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	repository := jobs.NewJobRegistry().GetJobRepository()

	jobs, err := repository.GetJobs()

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, jobs)

}

func GetJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var details map[string]string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&details)

	if details["publicid"] == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	repository := jobs.NewJobRegistry().GetJobRepository()

	job, err := repository.GetJob(details["publicid"])

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, job)

}
