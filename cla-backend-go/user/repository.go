package user

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetUserAndProfilesByLFID(lfidUsername string) (CLAUser, error)
	GetUserProjectIDs(userID string) ([]string, error)
	GetClaManagerCorporateClaIDs(userID string) ([]string, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository {
	return repository{
		db: db,
	}
}

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

func (repo repository) GetUserByGithubID(githubID string) (CLAUser, error) {
	// TODO: Implement when adding authentication to the Corporate Console
	return CLAUser{}, nil
}

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

	return nil, nil
}
