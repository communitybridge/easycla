package user

type service struct {
	repo Repository
}

func NewService(repo Repository) service {
	return service{
		repo: repo,
	}
}

func (s service) GetUserAndProfilesByLFID(lfidUsername string) (CLAUser, error) {
	user, err := s.repo.GetUserAndProfilesByLFID(lfidUsername)
	if err != nil {
		return CLAUser{}, err
	}

	return user, nil
}

func (s service) GetUserProjectIDs(userID string) ([]string, error) {
	projectIDs, err := s.repo.GetUserProjectIDs(userID)
	if err != nil {
		return nil, err
	}

	return projectIDs, nil
}

func (s service) GetClaManagerCorporateClaIDs(userID string) ([]string, error) {
	return nil, nil
}
