package route

import (
	"net/http"

	"jobs-applicant-api/controller/v1/applicants"
	"jobs-applicant-api/controller/v1/utilities"
	"jobs-applicant-api/route/middleware/acl"
	"jobs-applicant-api/route/middleware/cors"
	hr "jobs-applicant-api/route/middleware/httprouterwrapper"
	"jobs-applicant-api/route/middleware/logrequest"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// LoadRoutes returns the routes and middleware
func LoadRoutes() http.Handler {
	//return routes()
	return middleware(routes())
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes() *httprouter.Router {
	r := httprouter.New()

	r.POST("/upload/image", hr.Handler(alice.New(acl.AllowAPIKey).ThenFunc(utilities.UploadImage)))

	r.POST("/applicant/signup", hr.Handler(alice.New(acl.AllowAPIKey).ThenFunc(applicants.SignUp)))
	r.POST("/applicant/login", hr.Handler(alice.New(acl.AllowAPIKey).ThenFunc(applicants.Login)))

	r.POST("/applicant/update-password", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(applicants.UpdatePassword)))
	r.POST("/applicant/update-account", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(applicants.UpdateAccount)))
	r.POST("/applicant/update-job-preferences", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(applicants.UpdateJobPreferences)))
	r.GET("/applicant/get", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(applicants.GetApplicant)))
	r.POST("/applicant/get/location/autocomplete", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(applicants.GetAutocompleteLocationData)))
	r.GET("/applicant/get/jobs", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(applicants.GetJobs)))
	r.POST("/applicant/get/job", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(applicants.GetJob)))

	// r.POST("/employer/update-company", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.UpdateCompany)))
	// r.POST("/employer/update-payment-method", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.UpdatePaymentMethod)))
	// r.POST("/employer/update-payment-details", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.UpdatePaymentDetails)))
	// r.GET("/employer/get/company", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.GetEmployerCompany)))
	// r.POST("/employer/create/job", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.CreateJob)))
	// r.POST("/employer/edit/job", hr.Handler(alice.New(acl.ValidateJWT, acl.ValidateJWT).ThenFunc(employers.EditJob)))
	// r.GET("/employer/get/jobs", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.GetJobs)))
	// r.POST("/employer/get/job", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.GetJob)))
	// r.DELETE("/employer/delete/job", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.DeleteJob)))
	// r.GET("/employer/get/jobpackages/active", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.GetActiveJobPackages)))

	// r.POST("/employer/buy/job-package", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.PurchaseJobPackage)))

	// r.POST("/get-user", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(users.GetUser)))

	// r.GET("/get/client/registration", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(clients.CheckRegistration)))
	// r.POST("/set/client/registration", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(clients.SetRegistration)))

	// r.GET("/customers/:id", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(clients.GetClientCustomers)))
	// Enable Pprof
	// r.GET("/debug/pprof/*pprof", hr.Handler(alice.
	// 	New(acl.ValidateJWT).
	// 	ThenFunc(pprofhandler.Handler)))

	return r
}

// *****************************************************************************
// Middleware
// *****************************************************************************

func middleware(h http.Handler) http.Handler {
	// Log every request
	h = logrequest.Handler(h)

	// Cors for swagger-ui
	h = cors.Handler(h)

	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)

	return h
}
