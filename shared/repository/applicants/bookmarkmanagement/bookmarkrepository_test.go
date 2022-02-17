package bookmarkmanagement_test

import (
	"autumnomous-jobs-applicant-api/shared/database"
	"autumnomous-jobs-applicant-api/shared/repository/applicants/bookmarkmanagement"
	"autumnomous-jobs-applicant-api/shared/testhelper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BookmarkRepository_CreateApplicantJobBookmark_Fail_EmptyData(t *testing.T) {

	assert := assert.New(t)

	repository := bookmarkmanagement.NewBookmarkRepository(database.DB)

	result, err := repository.CreateApplicantJobBookmark("", "")

	assert.NotNil(err)
	assert.Nil(result)
}

func Test_BookmarkRepository_CreateApplicantJobBookmark_Success(t *testing.T) {

	assert := assert.New(t)

	applicant := testhelper.Helper_RandomApplicant(t)
	employer := testhelper.Helper_RandomEmployer(t)
	job := testhelper.Helper_RandomJob(employer, t)

	repository := bookmarkmanagement.NewBookmarkRepository((database.DB))

	result, err := repository.CreateApplicantJobBookmark(applicant.PublicID, job.PublicID)

	assert.Nil(err)
	assert.NotNil(result)

}

func Test_BookmarkRepository_GetApplicantJobBookmark_Fail_EmptyData(t *testing.T) {

	assert := assert.New(t)

	repository := bookmarkmanagement.NewBookmarkRepository(database.DB)

	result, err := repository.GetApplicantJobBookmark("", "")

	assert.NotNil(err)
	assert.Nil(result)
}

func Test_BookmarkRepository_GetApplicantJobBookmark_Success(t *testing.T) {

	assert := assert.New(t)

	applicant := testhelper.Helper_RandomApplicant(t)
	employer := testhelper.Helper_RandomEmployer(t)
	job := testhelper.Helper_RandomJob(employer, t)

	repository := bookmarkmanagement.NewBookmarkRepository((database.DB))

	result, err := repository.GetApplicantJobBookmark(applicant.PublicID, job.PublicID)

	assert.Nil(err)
	assert.NotNil(result)

}

func Test_BookmarkRepository_DeleteApplicantJobBookmark_Success(t *testing.T) {

	assert := assert.New(t)

	applicant := testhelper.Helper_RandomApplicant(t)
	employer := testhelper.Helper_RandomEmployer(t)
	job := testhelper.Helper_RandomJob(employer, t)

	repository := bookmarkmanagement.NewBookmarkRepository((database.DB))

	result, err := repository.DeleteApplicantJobBookmark(applicant.PublicID, job.PublicID)

	assert.Nil(err)
	assert.NotNil(result)
}

func Test_BookmarkRepository_DeleteApplicantJobBookmark_Fail_EmptyData(t *testing.T) {

	assert := assert.New(t)

	repository := bookmarkmanagement.NewBookmarkRepository(database.DB)

	result, err := repository.DeleteApplicantJobBookmark("", "")

	assert.NotNil(err)
	assert.Nil(result)
}
