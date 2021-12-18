package applicants

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	applicants "jobs-applicant-api/shared/repository/applicants"

	"jobs-applicant-api/shared/repository/applicants/accountmanagement"
	"jobs-applicant-api/shared/response"
	"jobs-applicant-api/shared/services/security/encryption"
	"jobs-applicant-api/shared/services/security/jwt"

	mailgun "github.com/mailgun/mailgun-go/v4"
)

var SendWelcomeMessageFunction = SendWelcomeMessage
var AuthenticationFunction = AuthenticatePassword

type SignUpCredentials struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var credentials SignUpCredentials
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	if err != nil {
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	if credentials.FirstName == "" || credentials.LastName == "" || credentials.Email == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, response.MissingRequiredValue)
		return
	}

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	password := encryption.GeneratePassword(9)
	hashedPassword, err := encryption.HashPassword([]byte(password))

	if err != nil {
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	employer, err := repository.CreateApplicant(credentials.FirstName, credentials.LastName, credentials.Email, string(hashedPassword))

	if err != nil {
		response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
		return
	}

	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")
	_, err = SendWelcomeMessageFunction(domain, apiKey, string(password), employer)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, "JSON error:"+err.Error())
		return
	}

	response.SendJSON(w, "")
}

//SendWelcomeMessage Sends a welcome message
func SendWelcomeMessage(domain, apiKey, password string, employer *accountmanagement.Applicant) (string, error) {

	message := fmt.Sprintf("Thank you for joining BiT Jobs, %s!\nYour temporary password is %s", employer.FirstName, password)
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
		"BiT Jobs Support <admin@autumnomous.git.beanstalkapp.com/jobs-applicant-api>",
		"Welcome to BiT Jobs!",
		message,
		employer.Email,
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, id, err := mg.Send(ctx, m)

	return id, err
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.SendJSONMessage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	decoder := json.NewDecoder(r.Body)
	credentials := LoginCredentials{}
	decoder.Decode(&credentials)

	if credentials.Email == "" || credentials.Password == "" {
		response.SendJSONMessage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	match, registrationStep, publicID, err := AuthenticationFunction(credentials.Email, credentials.Password)

	if err != nil {
		log.Println(err)
		response.SendJSONMessage(w, http.StatusInternalServerError, "Internal Server Error")
	}

	if match {

		// use   to create JWT
		tokenStr, err := jwt.GenerateToken(publicID)

		if err != nil {
			log.Println(err)
			response.SendJSONMessage(w, http.StatusInternalServerError, response.FriendlyError)
			return
		}

		encodedTokenStr := base64.StdEncoding.EncodeToString([]byte(tokenStr))

		token := map[string]interface{}{
			"token":            encodedTokenStr,
			"registrationstep": registrationStep,
		}
		response.SendJSON(w, token)
		return
	} else {
		response.SendJSONMessage(w, http.StatusUnauthorized, "Login failed")
		return
	}

}

func AuthenticatePassword(email, password string) (bool, string, string, error) {

	repository := applicants.NewApplicantRegistry().GetApplicantRepository()

	match, registrationStep, publicID, err := repository.AuthenticateApplicantPassword(email, password)

	if err != nil {
		log.Println(err)
		return false, "", "", err
	}

	return match, registrationStep, publicID, nil

}
