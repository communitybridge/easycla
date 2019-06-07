package user

type CLAUser struct {
	UserID string
	Name   string

	Emails []string

	LfidProvider   UserProvider
	GithubProvider UserProvider

	ProjectIDs []string
	ClaIDs     []string
}

type UserProvider struct {
	ProviderUserID string
}
