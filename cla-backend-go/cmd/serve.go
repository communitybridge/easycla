package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/config"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/docraptor"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/gen/restapi/operations"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/health"
	"github.com/LF-Engineering/cla-monorepo/cla-backend-go/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-openapi/loads"
	_ "github.com/lib/pq"
	"github.com/lytics/logrus"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	awsRegion = "us-east-1"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the backend server",
	Long:  `Run the backend server which listens for http requests over a given port.`,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := os.Hostname()
		if err != nil {
			logrus.Panicln("unable to get Hostname", err)
		}

		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetLevel(logrus.DebugLevel)
		logrus.WithFields(logrus.Fields{
			"BuildTime": BuildStamp,
			"GitHash":   GitHash,
			"Host":      host,
		}).Info("Service Startup")

		awsSession := session.Must(session.NewSession(
			&aws.Config{
				Region:                        aws.String(awsRegion),
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
		))

		configFile, err := config.LoadConfig(configFile, awsSession, viper.GetString("STAGE"))
		if err != nil {
			log.Panicln("Unable to load config", err)
		}

		// db, err := sqlx.Connect("postgres", viper.GetString("POSTGRESQL_CONNECTION"))
		// if err != nil {
		// 	log.Panicln("unable to connect to DB", err)
		// }

		// db.SetMaxOpenConns(viper.GetInt("DB_MAX_CONNECTIONS"))
		// db.SetMaxIdleConns(5)
		// db.SetConnMaxLifetime(15 * time.Minute)
		// db.MapperFunc(snaker.CamelToSnake)

		swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
		if err != nil {
			logrus.Panicln("Invalid swagger file for initializing cla", err)
		}

		api := operations.NewClaAPI(swaggerSpec)
		docraptorClient, err := docraptor.NewDocraptorClient(configFile.Docraptor.APIKey, configFile.Docraptor.TestMode)
		if err != nil {
			logrus.Panic(err)
		}

		// auth0Validator, err := auth.NewAuth0Validator(
		// 	configFile.Auth0.Domain,
		// 	configFile.Auth0.ClientID,
		// 	configFile.Auth0.UsernameClaim,
		// 	configFile.Auth0.Algorithm)
		// if err != nil {
		// 	logrus.Panic(err)
		// }

		var (
			// userRepo          = user.NewRepository(db)
			// projectRepo       = project.NewRepository(db)
			// contractGroupRepo = contractgroup.NewRepository(db)
			templateRepo = template.NewRepository(awsSession, viper.GetString("STAGE"))
		)

		var (
			healthService = health.New(GitHash, BuildStamp)
			// projectService       = project.NewService(projectRepo)
			//contractGroupService = contractgroup.NewService(contractGroupRepo)
			// userService          = user.NewService(userRepo)
			templateService = template.NewService(viper.GetString("STAGE"), templateRepo, docraptorClient, awsSession)
			//authorizer = auth.NewAuthorizer(auth0Validator)
		)

		//api.OauthSecurityAuth = authorizer.SecurityAuth
		health.Configure(api, healthService)
		template.Configure(api, templateService)
		// project.Configure(api, projectService)
		// contractgroup.Configure(api, contractGroupService)

		// flag.Parse()
		apiHandler := setupGlobalMiddleware(api.Serve(setupMiddlewares))

		server := restapi.NewServer(api)
		defer server.Shutdown() // nolint
		server.Port = viper.GetInt("PORT")

		server.SetHandler(apiHandler)

		err = server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("DB_MAX_CONNECTIONS", 1)
	viper.SetDefault("STAGE", "dev")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return responseLoggingMiddleware(handler)
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	return c.Handler(handler)
}

// LoggingResponseWriter is a wrapper around an http.ResponseWriter which captures the
// status code written to the response, so that it can be logged.
type LoggingResponseWriter struct {
	wrapped    http.ResponseWriter
	StatusCode int
	// Response content could also be captured here, but I was only interested in logging the response status code
}

func NewLoggingResponseWriter(wrapped http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{wrapped: wrapped}
}

func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.wrapped.Header()
}

func (lrw *LoggingResponseWriter) Write(content []byte) (int, error) {
	return lrw.wrapped.Write(content)
}

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
