package testhelper

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"autumnomous-jobs-applicant-api/shared/database"
	"autumnomous-jobs-applicant-api/shared/services/security/encryption"

	"github.com/joho/godotenv"
)

type TestApplicant struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	// CompanyPublicID  string
	HashedPassword   []byte
	PublicID         string
	RegistrationStep string
	PhoneNumber      string
	// MobileNumber     string
	// Role             string
}

type TestJob struct {
	PublicID         string `jaon:"publicid"`
	Title            string `json:"title"`
	JobType          string `json:"jobtype"`
	Category         string `json:"category"`
	Description      string `json:"description"` // make required?
	EmployerPublicID string `json:"employerpublicid"`
	Remote           bool   `json:"remote"`
	VisibleDate      string `json:"visibledate"`
}

type TestApplication struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Email       string `json:"email"`
	ID          int
	ApplicantID string `json:"applicantid"`
	PublicID    string `json:"publicid"`
}

type TestJobPackage struct {
	ID           int     `json:"id"`
	TypeID       string  `json:"typeid"`
	IsActive     bool    `json:"isactive"`
	Title        string  `json:"title"`
	NumberOfJobs int     `json:"numberofjobs"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
}

type TestCompany struct {
	Name         string `json:"name"`
	Domain       string `json:"domain"`
	Location     string `json:"location"`
	URL          string `json:"url"`
	Facebook     string `json:"facebook"`
	Twitter      string `json:"twitter"`
	Instagram    string `json:"instagram"`
	Description  string `json:"description"`
	Logo         string `json:"logo"`
	ExtraDetails string `json:"extradetails"`
	PublicID     string `json:"publicid"`
	ID           string `json:"id"`
	Zipcode      string `json:"zipcode"`
}

type TestEmployer struct {
	FirstName        string
	LastName         string
	Email            string
	Password         string
	TotalPostsBought int
	CompanyPublicID  string
	HashedPassword   []byte
	PublicID         string
	RegistrationStep string
	PhoneNumber      string
	MobileNumber     string
	Role             string
}

type TestJobBookmark struct {
	PublicID          string `json:"publicid"`
	ApplicantPublicID string `json:"applicantpublicid"`
	JobPublicID       string `json:"jobpublicid"`
}

func Init() {
	os.Clearenv()

	err := godotenv.Load("test.env")

	if err != nil {
		log.Println(err)
		log.Fatal("Error loading test.env file:", err)
	}

	database.Connect("DATABASE_URL")

}

func Helper_CreateApplicant(applicant *TestApplicant, t *testing.T) *TestApplicant {

	stmt, err := database.DB.Prepare(`INSERT INTO applicants (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	err = stmt.QueryRow(applicant.FirstName, applicant.LastName, applicant.Email, applicant.HashedPassword).Scan(&applicant.PublicID)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	return applicant
}

func Helper_RandomApplicant(t *testing.T) *TestApplicant {
	applicant := &TestApplicant{FirstName: string(encryption.GeneratePassword(5)),
		LastName: string(encryption.GeneratePassword(5)),
		Email:    fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		Password: string(encryption.GeneratePassword(9))}

	hashedPassword, err := encryption.HashPassword([]byte(applicant.Password))

	if err != nil {
		t.Fatal()
	}

	applicant.HashedPassword = hashedPassword

	return Helper_CreateApplicant(applicant, t)
}

func Helper_GetApplicant(publicID string, t *testing.T) *TestApplicant {

	stmt, err := database.DB.Prepare(`SELECT 
			email, firstname, lastname
		FROM applicants 
		WHERE publicid=$1;`)

	if err != nil {
		t.Fatal()
	}

	var applicant TestApplicant
	err = stmt.QueryRow(publicID).Scan(&applicant.Email, &applicant.FirstName, &applicant.LastName)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	return &applicant
}

func Helper_CreateJob(job *TestJob, t *testing.T) *TestJob {
	stmt, err := database.DB.Prepare(`INSERT INTO
											jobs(title, jobtype, category, description, visibledate,remote, employerid)
											VALUES ($1, $2, $3, $4, $5, $6, (SELECT id FROM employers WHERE publicid=$7))
											RETURNING publicid;`)
	if err != nil {
		log.Println(err)
		return nil
	}

	err = stmt.QueryRow(job.Title, job.JobType, job.Category, job.Description, job.VisibleDate, job.Remote, job.EmployerPublicID).Scan(&job.PublicID)

	if err != nil {
		log.Println(err)
		return nil
	}

	return job
}

func Helper_RandomJob(employer *TestEmployer, t *testing.T) *TestJob {

	job := &TestJob{
		Title:            string(encryption.GeneratePassword(5)),
		JobType:          string(encryption.GeneratePassword(5)),
		Category:         string(encryption.GeneratePassword(5)),
		Description:      string(encryption.GeneratePassword(5)),
		Remote:           true,
		VisibleDate:      time.Now().Format(time.RFC3339),
		EmployerPublicID: employer.PublicID,
	}

	return Helper_CreateJob(job, t)
}

func Helper_CreateEmployer(employer *TestEmployer, t *testing.T) *TestEmployer {

	stmt, err := database.DB.Prepare(`INSERT INTO employers (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	err = stmt.QueryRow(employer.FirstName, employer.LastName, employer.Email, employer.HashedPassword).Scan(&employer.PublicID)

	if err != nil {
		log.Println(err)
		t.Fatal()
	}

	return employer
}

func Helper_RandomEmployer(t *testing.T) *TestEmployer {
	employer := &TestEmployer{FirstName: string(encryption.GeneratePassword(5)),
		LastName: string(encryption.GeneratePassword(5)),
		Email:    fmt.Sprintf("email-%s@site.com", string(encryption.GeneratePassword(9))),
		Password: string(encryption.GeneratePassword(9))}

	hashedPassword, err := encryption.HashPassword([]byte(employer.Password))

	if err != nil {
		t.Fatal()
	}

	employer.HashedPassword = hashedPassword

	return Helper_CreateEmployer(employer, t)
}

func Helper_CreateCompany(company *TestCompany, t *testing.T) *TestCompany {

	stmt, err := database.DB.Prepare(`INSERT INTO 
			companies(domain, name, location, url, facebook, twitter, instagram, description, logo, extradetails, zipcode) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
			RETURNING publicid;`)
	if err != nil {
		log.Println(err)
		return nil
	}

	err = stmt.QueryRow(company.Domain, company.Name, company.Location, company.URL, company.Facebook, company.Twitter, company.Instagram, company.Description, company.Logo, company.ExtraDetails, company.Zipcode).Scan(&company.PublicID)

	if err != nil {
		log.Println(err)
		return nil
	}

	return company
}

func Helper_RandomCompany(t *testing.T) *TestCompany {

	company := &TestCompany{
		Domain:    string(encryption.GeneratePassword(5)),
		Name:      string(encryption.GeneratePassword(5)),
		Location:  string(encryption.GeneratePassword(5)),
		URL:       string(encryption.GeneratePassword(5)),
		Facebook:  string(encryption.GeneratePassword(5)),
		Twitter:   string(encryption.GeneratePassword(5)),
		Instagram: string(encryption.GeneratePassword(5)),
		Zipcode:   string(encryption.GeneratePassword(5)),
	}

	return Helper_CreateCompany(company, t)
}

func Helper_SetEmployerCompany(employerPublicID, companyPublicID string) error {

	if employerPublicID == "" || companyPublicID == "" {
		return errors.New("missing required value")
	}

	stmt, err := database.DB.Prepare(`UPDATE employers SET companyid=(SELECT id FROM companies WHERE publicid=$1) WHERE publicid=$2;`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(companyPublicID, employerPublicID)

	if err != nil {
		return err
	}

	return nil

}

func Helper_CreateApplicantJobBookmark(applicant *TestApplicant, job *TestJob, t *testing.T) *TestJobBookmark {

	if applicant.PublicID == "" || job.PublicID == "" {
		return nil
	}

	var bookmark TestJobBookmark
	stmt, err := database.DB.Prepare(`INSERT INTO applicantjobbookmarks(applicantpublicid, jobpublicid) VALUES ($1, $2) RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		return nil
	}

	err = stmt.QueryRow(applicant.PublicID, job.PublicID).Scan(&bookmark.PublicID)

	if err != nil {
		log.Println(err)
		return nil
	}

	bookmark.ApplicantPublicID = applicant.PublicID
	bookmark.JobPublicID = job.PublicID

	return &bookmark
}

// func Helper_CreateJobPackage(pack *TestJobPackage, t *testing.T) *TestJobPackage {

// 	stmt, err := database.DB.Prepare(`INSERT INTO
// 											jobpackages(typeid, isactive, title, numberofjobs, description, price)
// 											VALUES ($1, $2, $3, $4, $5, $6)
// 											RETURNING id;`)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	err = stmt.QueryRow(pack.TypeID, pack.IsActive, pack.Title, pack.NumberOfJobs, pack.Description, pack.Price).Scan(&pack.ID)

// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	return pack
// }

// func Helper_RandomJobPackage(t *testing.T) *TestJobPackage {

// 	pack := &TestJobPackage{
// 		TypeID:       string(encryption.GeneratePassword(5)),
// 		IsActive:     true,
// 		Title:        string(encryption.GeneratePassword(5)),
// 		NumberOfJobs: 3,
// 		Description:  string(encryption.GeneratePassword(5)),
// 		Price:        100.00,
// 	}

// 	return Helper_CreateJobPackage(pack, t)
// }

// func Helper_CreateCompany(company *TestCompany, t *testing.T) *TestCompany {

// 	stmt, err := database.DB.Prepare(`INSERT INTO
// 			companies(domain, name, location, url, facebook, twitter, instagram, description, logo, extradetails)
// 			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
// 			RETURNING publicid;`)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	err = stmt.QueryRow(company.Domain, company.Name, company.Location, company.URL, company.Facebook, company.Twitter, company.Instagram, company.Description, company.Logo, company.ExtraDetails).Scan(&company.PublicID)

// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	return company
// }

// func Helper_RandomCompany(t *testing.T) *TestCompany {

// 	company := &TestCompany{
// 		Domain:    string(encryption.GeneratePassword(5)),
// 		Name:      string(encryption.GeneratePassword(5)),
// 		Location:  string(encryption.GeneratePassword(5)),
// 		URL:       string(encryption.GeneratePassword(5)),
// 		Facebook:  string(encryption.GeneratePassword(5)),
// 		Twitter:   string(encryption.GeneratePassword(5)),
// 		Instagram: string(encryption.GeneratePassword(5)),
// 	}

// 	return Helper_CreateCompany(company, t)
// }

// func Helper_SetEmployerCompany(employerPublicID, companyPublicID string) error {

// 	if employerPublicID == "" || companyPublicID == "" {
// 		return errors.New("missing required value")
// 	}

// 	stmt, err := database.DB.Prepare(`UPDATE employers SET companyid=(SELECT id FROM companies WHERE publicid=$1) WHERE publicid=$2;`)

// 	if err != nil {
// 		return err
// 	}

// 	_, err = stmt.Exec(companyPublicID, employerPublicID)

// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

// func Helper_ChangeRegistrationStep(step string, Applicant *TestUser, t *testing.T) error {

// 	stmt, err := database.DB.Prepare(`UPDATE applications SET registrationstep=$1 WHERE Applicantid=(SELECT id FROM Applicants WHERE publicid=$2);`)

// 	if err != nil {
// 		t.Fatal()
// 	}

// 	err = stmt.QueryRow(step, Applicant.PublicID).Err()

// 	return err

// }
