// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package project

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Repository defines the interfaces for the project service
type Repository interface {
	GetProjectIDsForUser(ctx context.Context, userID string) ([]string, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new instance of the repository service
func NewRepository(db *sqlx.DB) repository {
	return repository{
		db: db,
	}
}

// GetProjectIDsForUser returns a list of project ID's for the specified user
func (repo repository) GetProjectIDsForUser(ctx context.Context, userID string) ([]string, error) {
	// Getting SFDC ids from the DB
	sql := `
	SELECT
		project_sfdc_id
	FROM
	    cla.project_manager pm
	WHERE
	    pm.user_id = $1;`
	rows, err := repo.db.Queryx(sql, userID)
	if err != nil {
		return nil, err
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
