package accountmanagement

import (
	"autumnomous-jobs-applicant-api/shared/services/security/encryption"
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

type ApplicantRepository struct {
	Database *sql.DB
}

type Applicant struct {
	FirstName   string  `json:"firstname"`
	LastName    string  `json:"lastname"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phonenumber"`
	Address     string  `json:"address"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Zipcode     string  `json:"zipcode"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	// MobileNumber     string `json:"mobilenumber"`
	// Role             string `json:"role"`
	// Facebook         string `json:"facebook"`
	// Twitter          string `json:"twitter"`
	// Instagram        string `json:"instagram"`
	// TotalPostsBought int    `json:"totalpostsbought"`
	RegistrationStep string `json:"registrationstep"`
	Password         string
	// CompanyPublicID  string `json:"companypublicid"`
	PublicID string `json:"publicid"`
}

// RegistrationStep represents which stage in the registration process the user is in
type RegistrationStep int64

const (
	// ChangePassword Registration Step 1
	ChangePassword RegistrationStep = iota

	// PersonalInformation Registration Step 2
	PersonalInformation

	// Job Preferences Registration Step 3
	JobPreferences

	// Complete Registration Step 4
	RegistrationComplete
)

func (rs RegistrationStep) String() string {
	return [...]string{"change-password", "personal-information", "job-preferences", "registration-complete"}[rs]
}

func NewApplicantRepository(db *sql.DB) *ApplicantRepository {
	return &ApplicantRepository{Database: db}
}

func (repository *ApplicantRepository) CreateApplicant(firstName, lastName, email, password string) (*Applicant, error) {

	if firstName == "" || lastName == "" || email == "" || password == "" {
		return nil, errors.New("data cannot be empty")
	}

	applicant := &Applicant{FirstName: firstName, LastName: lastName, Email: email}

	stmt, err := repository.Database.Prepare(`INSERT INTO applicants(email, firstname, lastname, password, registrationstep) VALUES ($1, $2, $3, $4, 'change-password') RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(email, firstName, lastName, password).Scan(&applicant.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return applicant, nil
}

func (repository *ApplicantRepository) GetApplicant(userID string) (*Applicant, error) {

	if userID == "" {
		return nil, errors.New("missing required value")
	}
	var applicant Applicant

	stmt, err := repository.Database.Prepare(`
		SELECT firstname, lastname, email, registrationstep, phonenumber, address, city, state, zipcode
		FROM applicants
		WHERE publicid=$1;`,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var app_phone_number, registrationstep, address, city, state, zipcode sql.NullString

	err = stmt.QueryRow(userID).Scan(&applicant.FirstName, &applicant.LastName, &applicant.Email, &registrationstep, &app_phone_number, &address, &city, &state, &zipcode)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if app_phone_number.Valid {
		applicant.PhoneNumber = app_phone_number.String
	}

	if registrationstep.Valid {
		applicant.RegistrationStep = registrationstep.String
	}

	if address.Valid {
		applicant.Address = address.String
	}

	if city.Valid {
		applicant.City = city.String
	}

	if state.Valid {
		applicant.State = state.String
	}

	if zipcode.Valid {
		applicant.Zipcode = zipcode.String
	}
	// if emp_facebook.Valid {
	// 	applicant.Facebook = emp_facebook.String
	// }

	// if emp_twitter.Valid {
	// 	applicant.Twitter = emp_twitter.String
	// }

	// if emp_instagram.Valid {
	// 	applicant.Instagram = emp_instagram.String
	// }
	applicant.PublicID = userID

	return &applicant, nil
}

func (repository *ApplicantRepository) AuthenticateApplicantPassword(email, password string) (bool, string, string, error) {

	if email == "" || password == "" {
		return false, "", "", nil
	}

	var databasePassword, registrationStep sql.NullString
	var publicID string
	stmt, err := repository.Database.Prepare(`SELECT password, registrationStep, publicid FROM applicants WHERE email=$1;`)

	if err != nil {
		log.Println(err)
		return false, "", "", err
	}

	err = stmt.QueryRow(email).Scan(&databasePassword, &registrationStep, &publicID)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			return false, "", "", nil
		} else {
			log.Println(err)
			return false, "", "", err
		}

	}

	if databasePassword.Valid {
		if encryption.CompareHashes([]byte(databasePassword.String), []byte(password)) {
			return true, registrationStep.String, publicID, nil
		}
	}

	return false, "", "", nil
}

func (repository *ApplicantRepository) UpdateApplicantPassword(publicID, password, newPassword string) (bool, error) {

	if publicID == "" || password == "" || newPassword == "" {
		return false, nil
	}
	var databasePassword, registrationStep sql.NullString

	stmt, err := repository.Database.Prepare(`SELECT password, registrationstep FROM applicants WHERE publicid=$1;`)

	if err != nil {
		log.Println(err)
		return false, err
	}

	err = stmt.QueryRow(publicID).Scan(&databasePassword, &registrationStep)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			return false, nil
		} else {
			log.Println(err)
			return false, err
		}

	}

	if databasePassword.Valid {
		if encryption.CompareHashes([]byte(databasePassword.String), []byte(password)) {

			if registrationStep.Valid {
				if registrationStep.String == ChangePassword.String() {
					stmt, err = repository.Database.Prepare(`UPDATE applicants SET registrationstep='personal-information' WHERE publicid=$1;`)

					if err != nil {
						log.Println(err)
						return false, err
					}

					_, err = stmt.Exec(publicID)

					if err != nil {
						log.Println(err)
						return false, err
					}

				}
			}

			stmt, err = repository.Database.Prepare(`UPDATE applicants SET password=$1 WHERE publicid=$2;`)

			if err != nil {
				log.Println(err)
				return false, err
			}

			hashedNewPassword, err := encryption.HashPassword([]byte(newPassword))

			if err != nil {
				log.Println(err)
				return false, err
			}

			_, err = stmt.Exec(hashedNewPassword, publicID)

			if err != nil {
				log.Println(err)
				return false, err
			}

			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}
}

func (repository *ApplicantRepository) UpdateApplicantAccount(publicID, firstName, lastName, email, phoneNumber, address, city, state, zipcode string, latitude, longitude float64) (*Applicant, error) {

	applicant := &Applicant{}

	stmt, err := repository.Database.Prepare(`SELECT firstname, lastname, email FROM applicants WHERE publicid=$1;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(publicID).Scan(&applicant.FirstName, &applicant.LastName, &applicant.Email)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if firstName != "" {
		applicant.FirstName = firstName
	}

	if lastName != "" {
		applicant.LastName = lastName
	}

	if email != "" {
		applicant.Email = email
	}

	if phoneNumber != "" {
		applicant.PhoneNumber = phoneNumber
	}

	if address != "" {
		applicant.Address = address
	}

	if city != "" {
		applicant.City = city
	}

	if state != "" {
		applicant.State = state
	}

	if zipcode != "" {
		applicant.Zipcode = zipcode
	}

	if latitude != 0 {
		applicant.Latitude = latitude
	}

	if longitude != 0 {
		applicant.Longitude = longitude
	}

	// if mobileNumber != "" {
	// 	applicant.MobileNumber = mobileNumber
	// }

	// applicant.Facebook = facebook
	// applicant.Twitter = twitter
	// applicant.Instagram = instagram
	applicant.PublicID = publicID
	stmt, err = repository.Database.Prepare(`UPDATE applicants SET firstname=$1, lastname=$2, email=$3, phonenumber=$4, address=$5, city=$6, state=$7, zipcode=$8, latitude=$9, longitude=$10 WHERE publicid=$11;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = stmt.Exec(applicant.FirstName, applicant.LastName, applicant.Email, applicant.PhoneNumber, applicant.Address, applicant.City, applicant.State, applicant.Zipcode, applicant.Latitude, applicant.Longitude, applicant.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	applicant, _ = repository.GetApplicant(publicID)

	if applicant.RegistrationStep == PersonalInformation.String() {
		stmt, _ = repository.Database.Prepare(`UPDATE applicants SET registrationstep='job-preferences' WHERE publicid=$1;`)

		stmt.Exec(publicID)

	}

	return applicant, nil
}

func (repository *ApplicantRepository) UpdateApplicantJobPreferences(publicID string, desiredCities []map[string]interface{}) error {

	for _, city := range desiredCities {

		stmt, err := repository.Database.Prepare(`
			INSERT INTO
			desiredcities(city, state, country, latitude, longitude, text, applicantid)
			VALUES ($1, $2, $3, $4, $5, $6, (SELECT id FROM applicants WHERE publicid=$7));
		`)

		if err != nil {
			log.Println(err)
			return err
		}

		err = stmt.QueryRow(city["city"], city["state"], city["country"], city["latitude"], city["longitude"], city["text"], publicID).Scan()

		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				log.Println(err)
				return err
			}
		}

	}

	applicant, _ := repository.GetApplicant(publicID)

	if applicant.RegistrationStep == JobPreferences.String() {
		stmt, _ := repository.Database.Prepare(`UPDATE applicants SET registrationstep='registration-complete' WHERE publicid=$1;`)

		stmt.Exec(publicID)

	}

	return nil
}
