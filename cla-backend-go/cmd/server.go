// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"

	ini "github.com/communitybridge/easycla/cla-backend-go/init"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"

	"github.com/communitybridge/easycla/cla-backend-go/auth"
	"github.com/communitybridge/easycla/cla-backend-go/company"
	"github.com/communitybridge/easycla/cla-backend-go/config"
	"github.com/communitybridge/easycla/cla-backend-go/docraptor"
	"github.com/communitybridge/easycla/cla-backend-go/docs"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi"
	"github.com/communitybridge/easycla/cla-backend-go/gen/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/github"
	"github.com/communitybridge/easycla/cla-backend-go/health"
	"github.com/communitybridge/easycla/cla-backend-go/template"
	"github.com/communitybridge/easycla/cla-backend-go/user"
	"github.com/communitybridge/easycla/cla-backend-go/whitelist"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-openapi/loads"
	"github.com/lytics/logrus"
	"github.com/rs/cors"
	"github.com/savaki/dynastore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Version is the application version - either a git SHA or tag value
	Version string

	// Commit is the application commit hash
	Commit string

	// Branch the build branch
	Branch string

	// BuildDate is the date of the build
	BuildDate string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the backend server",
	Long:  `Run the backend server which listens for http requests over a given port.`,
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// server function called by environment specific server functions
func server(localMode bool) http.Handler {

	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("unable to get hostname. Error: %v", err)
	}

	log.Infof("Service %s starting...", ini.ServiceName)

	// Show the version and build info
	log.Infof("Name                  : %s", ini.ServiceName)
	log.Infof("Version               : %s", Version)
	log.Infof("Git commit hash       : %s", Commit)
	log.Infof("Branch                : %s", Branch)
	log.Infof("Build date            : %s", BuildDate)
	log.Infof("Golang OS             : %s", runtime.GOOS)
	log.Infof("Golang Arch           : %s", runtime.GOARCH)
	log.Infof("Service Host          : %s", host)
	log.Infof("Service Port          : %d", *portFlag)

	awsSession, err := ini.GetAWSSession()
	if err != nil {
		log.Panicf("Unable to load AWS session - Error: %v", err)
	}

	configFile, err := config.LoadConfig(configFile, awsSession, viper.GetString("STAGE"))
	if err != nil {
		log.Panicf("Unable to load config - Error: %v", err)
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logrus.Panicf("Invalid swagger file for initializing cla - Error: %v", err)
	}

	api := operations.NewClaAPI(swaggerSpec)
	docraptorClient, err := docraptor.NewDocraptorClient(configFile.Docraptor.APIKey, configFile.Docraptor.TestMode)
	if err != nil {
		logrus.Panicf("Unable to setup docraptor client - Error: %v", err)
	}

	authValidator, err := auth.NewAuthValidator(
		configFile.Auth0.Domain,
		configFile.Auth0.ClientID,
		configFile.Auth0.UsernameClaim,
		configFile.Auth0.Algorithm)
	if err != nil {
		logrus.Panic(err)
	}

	userRepo := user.NewDynamoRepository(awsSession, viper.GetString("STAGE"), configFile.SenderEmailAddress)
	templateRepo := template.NewRepository(awsSession, viper.GetString("STAGE"))
	whitelistRepo := whitelist.NewRepository(awsSession, viper.GetString("STAGE"))
	companyRepo := company.NewRepository(awsSession, viper.GetString("STAGE"))

	healthService := health.New(Version, Commit, Branch, BuildDate)
	templateService := template.NewService(viper.GetString("STAGE"), templateRepo, docraptorClient, awsSession)
	whitelistService := whitelist.NewService(whitelistRepo, http.DefaultClient)
	companyService := company.NewService(companyRepo, awsSession, configFile.SenderEmailAddress, configFile.CorporateConsoleURL, userRepo)
	authorizer := auth.NewAuthorizer(authValidator, userRepo)

	sessionStore, err := dynastore.New(dynastore.Path("/"), dynastore.HTTPOnly(), dynastore.TableName(configFile.SessionStoreTableName), dynastore.DynamoDB(dynamodb.New(awsSession)))
	if err != nil {
		log.Fatalf("Unable to create new Dynastore session - Error: %v", err)
	}

	api.OauthSecurityAuth = authorizer.SecurityAuth
	health.Configure(api, healthService)
	template.Configure(api, templateService)
	github.Configure(api, configFile.Github.ClientID, configFile.Github.ClientSecret, sessionStore)
	whitelist.Configure(api, whitelistService, sessionStore)
	docs.Configure(api)

	company.Configure(api, companyService)

	// For local mode - we allow anything, otherwise we use the value specified in the config (e.g. AWS SSM)
	var apiHandler http.Handler
	if localMode {
		apiHandler = setupGlobalMiddlewareLocal(api.Serve(setupMiddlewares))
	} else {
		apiHandler = setupGlobalMiddleware(api.Serve(setupMiddlewares), configFile.AllowedOrigins)
	}

	return apiHandler
}

// setupMiddlewares The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return responseLoggingMiddleware(handler)
}

// setupGlobalMiddleware sets up the CORS logic and creates the middleware HTTP handler
func setupGlobalMiddleware(handler http.Handler, allowedOrigins []string) http.Handler {

	log.Debugf("Allowed origins: %v", allowedOrigins)
	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			u, err := url.Parse(origin)
			if err != nil {
				log.Warnf("cors parse origin issue: %v", err)
				return false
			}

			if u.Scheme != "https" {
				log.Warnf("non-https scheme - blocking origin: %s", origin)
				return false
			}

			// Ensure the origin is in our allowed list
			allowedOrigin := stringInSlice(origin, allowedOrigins)
			if allowedOrigin {
				log.Debugf("origin %s is allowed", origin)
			} else {
				log.Warnf("origin %s is NOT allowed - not in allowed list: %v", origin, allowedOrigins)
			}
			return allowedOrigin
		},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	return c.Handler(handler)
}

// setupGlobalMiddlewareLocal allows all origins and sets up the handler
func setupGlobalMiddlewareLocal(handler http.Handler) http.Handler {

	log.Debug("Allowing all origins")
	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		//AllowOriginFunc:  func(origin string) bool { return true },
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	return c.Handler(handler)
}

// stringInSlice returns true if the specified string exists in the slice, otherwise returns false
func stringInSlice(a string, list []string) bool {
	if list == nil {
		return false
	}

	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// LoggingResponseWriter is a wrapper around an http.ResponseWriter which captures the
// status code written to the response, so that it can be logged.
type LoggingResponseWriter struct {
	wrapped    http.ResponseWriter
	StatusCode int
	// Response content could also be captured here, but I was only interested in logging the response status code
}

// NewLoggingResponseWriter creates a new logging response writer
func NewLoggingResponseWriter(wrapped http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{wrapped: wrapped}
}

// Header returns the header
func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.wrapped.Header()
}

// Write writes the contents
func (lrw *LoggingResponseWriter) Write(content []byte) (int, error) {
	return lrw.wrapped.Write(content)
}

// WriteHeader writes the header
func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.StatusCode = statusCode
	lrw.wrapped.WriteHeader(statusCode)
}

// responseLoggingMiddleware logs the responses from API endpoints
func responseLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w2 := NewLoggingResponseWriter(w)
		next.ServeHTTP(w2, r)
		fmt.Printf("%s %s, response %d %s\n", r.Method, r.URL.String(), w2.StatusCode, http.StatusText(w2.StatusCode))
	})
}
