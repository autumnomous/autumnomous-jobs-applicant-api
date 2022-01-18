package jobs

import (
	"autumnomous-jobs-applicant-api/shared/database"
)

type JobRegistry struct {
}

func NewJobRegistry() *JobRegistry {
	return &JobRegistry{}
}

func (*JobRegistry) GetJobRepository() *JobRepository {
	return NewJobRepository(database.DB)
}
