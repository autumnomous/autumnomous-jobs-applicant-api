package applicants

import (
	"autumnomous-jobs-applicant-api/shared/repository/applicants"
	"autumnomous-jobs-applicant-api/shared/response"
	"autumnomous-jobs-applicant-api/shared/services/security/jwt"
	"encoding/json"
	"log"
	"net/http"
)

type bookmarkJobData struct {
	JobID string `json:"jobid"`
}

func BookmarkJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, response.FriendlyError)
		return
	}

	var data bookmarkJobData
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

	if publicID == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.FriendlyError)
		return
	}

	repository := applicants.NewApplicantRegistry().GetBookmarkRepository()

	if r.Method == http.MethodPost {

		bookmark, err := repository.GetApplicantJobBookmark(publicID, data.JobID)

		if err != nil {
			log.Println(err)
			response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
			return
		}

		if bookmark == nil {

			_, err = repository.CreateApplicantJobBookmark(publicID, data.JobID)

			if err != nil {
				log.Println(err)
				response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
				return
			}
		}

		response.SendJSONMessage(w, http.StatusOK, "Bookmarked")
		return
	}

	if r.Method == http.MethodDelete {

		_, err = repository.DeleteApplicantJobBookmark(publicID, data.JobID)

		if err != nil {
			log.Println(err)
			response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
			return
		}

		response.SendJSONMessage(w, http.StatusOK, "Deleted")
		return

	}

}
