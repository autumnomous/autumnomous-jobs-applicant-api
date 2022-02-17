package applicants

import (
	"autumnomous-jobs-applicant-api/shared/database"
	"autumnomous-jobs-applicant-api/shared/repository/applicants/accountmanagement"
	"autumnomous-jobs-applicant-api/shared/repository/applicants/bookmarkmanagement"
)

type ApplicantRegistry struct {
}

func NewApplicantRegistry() *ApplicantRegistry {
	return &ApplicantRegistry{}
}

func (*ApplicantRegistry) GetApplicantRepository() *accountmanagement.ApplicantRepository {
	return accountmanagement.NewApplicantRepository(database.DB)
}

func (*ApplicantRegistry) GetBookmarkRepository() *bookmarkmanagement.BookmarkRepository {
	return bookmarkmanagement.NewBookmarkRepository(database.DB)
}
