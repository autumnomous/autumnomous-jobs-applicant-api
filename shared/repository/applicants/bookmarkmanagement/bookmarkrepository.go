package bookmarkmanagement

import (
	"database/sql"
	"errors"
	"log"
)

type BookmarkRepository struct {
	Database *sql.DB
}

type ApplicantJobBookmark struct {
	PublicID          string `json:"publicid"`
	ApplicantPublicID string `json:"applicantpublicid"`
	JobPublicID       string `json:"jobpublicid"`
}

func NewBookmarkRepository(db *sql.DB) *BookmarkRepository {
	return &BookmarkRepository{Database: db}
}

func (repository *BookmarkRepository) CreateApplicantJobBookmark(userPublicID, jobID string) (*ApplicantJobBookmark, error) {

	if userPublicID == "" || jobID == "" {
		return nil, errors.New("missing required value")
	}

	var bookmark ApplicantJobBookmark

	bookmark.ApplicantPublicID = userPublicID
	bookmark.JobPublicID = jobID

	stmt, err := repository.Database.Prepare(`INSERT INTO applicantjobbookmarks(applicantpublicid, jobpublicid) VALUES ($1, $2) RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(bookmark.ApplicantPublicID, bookmark.JobPublicID).Scan(&bookmark.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &bookmark, nil

}

func (repository *BookmarkRepository) GetApplicantJobBookmark(userPublicID, jobID string) (*ApplicantJobBookmark, error) {

	if userPublicID == "" || jobID == "" {
		return nil, errors.New("missing required value")
	}

	var bookmark ApplicantJobBookmark

	bookmark.ApplicantPublicID = userPublicID
	bookmark.JobPublicID = jobID

	stmt, err := repository.Database.Prepare(`SELECT publicid from applicantjobbookmarks WHERE applicantpublicid=$1 AND jobpublicid=$2;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(bookmark.ApplicantPublicID, bookmark.JobPublicID).Scan(&bookmark.PublicID)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {

			return nil, nil

		} else {
			return nil, err
		}
	}

	return &bookmark, nil

}

func (repository *BookmarkRepository) DeleteApplicantJobBookmark(userPublicID, jobID string) (*ApplicantJobBookmark, error) {

	if userPublicID == "" || jobID == "" {
		return nil, errors.New("missing required value")
	}

	var bookmark ApplicantJobBookmark

	bookmark.ApplicantPublicID = userPublicID
	bookmark.JobPublicID = jobID

	stmt, err := repository.Database.Prepare(`DELETE from applicantjobbookmarks WHERE applicantpublicid=$1 AND jobpublicid=$2 RETURNING publicid;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = stmt.QueryRow(bookmark.ApplicantPublicID, bookmark.JobPublicID).Scan(&bookmark.PublicID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &bookmark, nil

}
