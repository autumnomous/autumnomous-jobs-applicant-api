package jobs

import (
	"database/sql"
	"log"
)

type JobRepository struct {
	Database *sql.DB
}

// CREATE TABLE jobs (
//     id SERIAL PRIMARY KEY,
//     publicid text NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
//     title text NOT NULL,
//     jobtype text,
//     category text,
//     description text,
//     minsalary integer,
//     maxsalary integer,
//     payperiod text,
//     poststartdatetime timestamp without time zone,
//     postenddatetime timestamp without time zone,
//     employerid integer NOT NULL REFERENCES employers(id) ON DELETE CASCADE ON UPDATE CASCADE,
//     iscustomized boolean DEFAULT false,
//     createdate timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
//     updatedate timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
//     slug text,
//     remote boolean DEFAULT false
// );

type Job struct {
	ID              int
	PublicID        string `json:"publicid"`
	Title           string `json:"title"`
	JobType         string `json:"jobtype"`
	Category        string `json:"category"`
	Description     string `json:"description"`
	VisibleDate     string `json:"visibledate"`
	Remote          bool   `json:"remote"`
	IsCustomized    bool   `json:"iscustomized"`
	CreateDate      string `json:"createdate"`
	UpdateDate      string `json:"updatedate"`
	EmployerID      string `json:"employerid"`
	CompanyName     string `json:"companyname"`
	CompanyURL      string `json:"companyurl"`
	CompanyLogo     string `json:"companylogo"`
	CompanyPublicID string `json:"companypublicid"`
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{Database: db}
}

func (repository *JobRepository) GetJobs() ([]*Job, error) {

	var jobs []*Job

	stmt, err := repository.Database.Prepare(`
		SELECT 
			jobs.title, jobs.jobtype, jobs.category, jobs.description, jobs.visibledate, jobs.remote, jobs.publicid,
			employers.companyid, companies.url, companies.name, companies.logo
		FROM 
			jobs 
		JOIN employers ON employers.id=jobs.employerid
		JOIN companies ON companies.id=employers.companyid
		WHERE 
			now() >= visibledate AND now() <= (visibledate + '30 days'::interval);`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		job := &Job{}

		var visibleDate sql.NullString
		err := rows.Scan(&job.Title, &job.JobType, &job.Category, &job.Description, &visibleDate, &job.Remote, &job.PublicID, &job.EmployerID, &job.CompanyURL, &job.CompanyName, &job.CompanyLogo)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		if visibleDate.Valid {
			job.VisibleDate = visibleDate.String
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (repository *JobRepository) GetJob(publicid string) (*Job, error) {

	var job Job

	stmt, err := repository.Database.Prepare(`
		SELECT 
			jobs.title, jobs.jobtype, jobs.category, jobs.description, jobs.visibledate, jobs.remote, jobs.publicid,
			employers.companyid, companies.url, companies.name, companies.logo
		FROM 
			jobs 
		JOIN employers ON employers.id=jobs.employerid
		JOIN companies ON companies.id=employers.companyid
		WHERE 
			jobs.publicid=$1;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var visibleDate sql.NullString

	err = stmt.QueryRow(publicid).Scan(&job.Title, &job.JobType, &job.Category, &job.Description, &visibleDate, &job.Remote, &job.PublicID, &job.EmployerID, &job.CompanyURL, &job.CompanyName, &job.CompanyLogo)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if visibleDate.Valid {
		job.VisibleDate = visibleDate.String
	}

	return &job, nil
}
