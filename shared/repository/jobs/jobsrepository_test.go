package jobs_test

import (
	"jobs-applicant-api/shared/repository/jobs"
	"jobs-applicant-api/shared/testhelper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	testhelper.Init()
}
func Test_JobsRepository_GetJobs(t *testing.T) {

	assert := assert.New(t)

	company := testhelper.Helper_RandomCompany(t)
	employer := testhelper.Helper_RandomEmployer(t)

	testhelper.Helper_SetEmployerCompany(employer.PublicID, company.PublicID)
	testhelper.Helper_RandomJob(employer, t)
	testhelper.Helper_RandomJob(employer, t)
	testhelper.Helper_RandomJob(employer, t)

	repository := jobs.NewJobRegistry().GetJobRepository()
	jobs, err := repository.GetJobs()

	assert.Nil(err)
	assert.GreaterOrEqual(len(jobs), 3)
}

func Test_JobsRepository_GetJob(t *testing.T) {

	assert := assert.New(t)

	company := testhelper.Helper_RandomCompany(t)
	employer := testhelper.Helper_RandomEmployer(t)

	testhelper.Helper_SetEmployerCompany(employer.PublicID, company.PublicID)
	job := testhelper.Helper_RandomJob(employer, t)

	repository := jobs.NewJobRegistry().GetJobRepository()
	result, err := repository.GetJob(job.PublicID)

	assert.Nil(err)
	assert.Equal(result.PublicID, job.PublicID)

}
