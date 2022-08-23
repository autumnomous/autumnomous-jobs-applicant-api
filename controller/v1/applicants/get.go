package applicants

import (
	"autumnomous-jobs-applicant-api/shared/repository/applicants"
	"autumnomous-jobs-applicant-api/shared/repository/jobs"
	"autumnomous-jobs-applicant-api/shared/response"
	"autumnomous-jobs-applicant-api/shared/services/security/jwt"
	"autumnomous-jobs-applicant-api/shared/services/zipcode"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type AutocompleteLocationData struct {
	Characters string `json:"chars"`
}

type GetJobsByRadiusData struct {
	Zipcode string  `json:"zipcode"`
	Radius  float64 `json:"radius"`
}

type getBookmarkJobData struct {
	JobID string `json:"jobid"`
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

func GetJobsByRadius(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var details GetJobsByRadiusData

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&details)

	// if details["zipcode"] == ""{

	// }

	// if details["radius"] == 0 {

	// }

	gateway := zipcode.NewZipCodeGateway(os.Getenv("ZIPCODESERVICES_API_KEY"))
	data, err := gateway.GetZipCodesInRadius(details.Zipcode, details.Radius)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	var result []*jobs.Job
	repository := jobs.NewJobRegistry().GetJobRepository()

	for i := 0; i < len(data); i++ {

		jobsInZipcode, err := repository.GetJobsByZipcode(data[i].ZipCode)

		if err != nil {
			log.Println(err)
			response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
			return
		}

		result = append(result, jobsInZipcode...)

	}
	// log.Println(data)
	// repository := jobs.NewJobRegistry().GetJobRepository()

	// jobs, err := repository.GetJobs()

	// if err != nil {
	// 	log.Println(err)
	// 	response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
	// 	return
	// }

	response.SendJSON(w, result)

}

func GetApplicantJobBookmark(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var data getBookmarkJobData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	if data.JobID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	publicID := jwt.GetUserClaim(r)

	repository := applicants.NewApplicantRegistry().GetBookmarkRepository()

	bookmark, err := repository.GetApplicantJobBookmark(publicID, data.JobID)

	if err != nil {
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, bookmark)
	return

}

func GetApplicantBookmarkedJobs(w http.ResponseWriter, r *http.Request) {

	publicID := jwt.GetUserClaim(r)

	repository := jobs.NewJobRegistry().GetJobRepository()

	jobs, err := repository.GetApplicantBookmarkedJobs(publicID)

	if err != nil {
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	response.SendJSON(w, jobs)
	return
}
