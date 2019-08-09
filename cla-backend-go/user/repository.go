// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package user

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository interface methods
type Repository interface {
	GetUserAndProfilesByLFID(lfidUsername string) (CLAUser, error)
	GetUserProjectIDs(userID string) ([]string, error)
	GetClaManagerCorporateClaIDs(userID string) ([]string, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new user repository
func NewRepository(db *sqlx.DB) repository {
	return repository{
		db: db,
	}
}

// GetUserAndProfilesByLFID get user profile by LFID
func (repo repository) GetUserAndProfilesByLFID(lfidUsername string) (CLAUser, error) {
	fmt.Println("lfidUsername:", lfidUsername)
	sql := `
		SELECT
			u.user_id,
			u.name,
			lf.provider_user_id AS lfid,
			gh.provider_user_id AS github
		FROM
			cla.user_auth_provider lf
		JOIN
			cla."user" u
		ON
			lf.user_id = u.user_id
		LEFT JOIN LATERAL (
			SELECT
				gh.provider_user_id
			FROM
				cla.user_auth_provider gh
			WHERE
				gh.user_id = u.user_id
			AND
				gh.provider = 'github'
		) gh ON TRUE
		WHERE
			lf.provider_user_id = $1
		AND
			lf.provider = 'lfid';`

	user := CLAUser{}
	err := repo.db.QueryRowx(sql, lfidUsername).Scan(
		&user.UserID,
		&user.Name,
		&user.LfidProvider.ProviderUserID,
		&user.GithubProvider.ProviderUserID,
	)
	if err != nil {
		return CLAUser{}, err
	}

	return user, nil
}

// GetUserByGithubID returns the user details based on the github ID
func (repo repository) GetUserByGithubID(githubID string) (CLAUser, error) {
	// TODO: Implement when adding authentication to the Corporate Console
	return CLAUser{}, nil
}

// GetUserProjectIDs get the user project ID's based on the specified user ID
func (repo repository) GetUserProjectIDs(userID string) ([]string, error) {
	getUserProjectIDsSQL := `
		SELECT
			project_sfdc_id
		FROM
			cla.project_manager
		WHERE
			user_id = $1;`

	rows, err := repo.db.Queryx(getUserProjectIDsSQL, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return []string{}, nil
	}

	projectIDs := []string{}
	for rows.Next() {
		var projectID string
		err := rows.Scan(&projectID)
		if err != nil {
			rows.Close()
			return nil, err
		}

		projectIDs = append(projectIDs, projectID)
	}

	return projectIDs, nil
}

// GetClaManagerCorporateClaIDs returns a list of CLA manager corporate CLAs associated with the specified user
func (repo repository) GetClaManagerCorporateClaIDs(userID string) ([]string, error) {
	getClaManagerCorporateClaIDsSQL := `
		SELECT
			corporate_cla_group_id
		FROM
			cla.cla_manager
		WHERE
			user_id = $1;`

	rows, err := repo.db.Queryx(getClaManagerCorporateClaIDsSQL, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return []string{}, nil
	}

	corporateClaIDs := []string{}
	for rows.Next() {
		var corporateClaID string
		err := rows.Scan(&corporateClaID)
		if err != nil {
			rows.Close()
			return nil, err
		}

		corporateClaIDs = append(corporateClaIDs, corporateClaID)
	}

	return corporateClaIDs, nil
}

// GetUserCompanyIDs returns a list of company IDs based on the user
func (repo repository) GetUserCompanyIDs(userID string) ([]string, error) {
	return []string{}, nil
}
