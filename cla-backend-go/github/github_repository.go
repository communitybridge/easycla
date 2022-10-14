// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package github

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/communitybridge/easycla/cla-backend-go/logging"
	"github.com/google/go-github/v37/github"
)

// errors
var (
	ErrGithubRepositoryNotFound = errors.New("github repository not found")
)

const (
	helpLink = "https://help.github.com/en/github/committing-changes-to-your-project/why-are-my-commits-linked-to-the-wrong-user"
	unknown  = "Unknown"
)

func GetGitHubRepository(ctx context.Context, installationID, githubRepositoryID int64) (*github.Repository, error) {
	f := logrus.Fields{
		"functionName":       "github.github_repository.GetGitHubRepository",
		"githubRepositoryID": githubRepositoryID,
	}
	client, clientErr := NewGithubAppClient(installationID)
	if clientErr != nil {
		return nil, clientErr
	}

	repository, _, repoErr := client.Repositories.GetByID(ctx, githubRepositoryID)
	if repoErr != nil {
		log.WithFields(f).WithError(repoErr).Warnf("unable to fetch repository by ID: %d", githubRepositoryID)
		return nil, repoErr
	}

	return repository, nil
}

func GetPullRequest(ctx context.Context, pullRequestID int, owner, repo string, client *github.Client) (*github.PullRequest, error) {
	f := logrus.Fields{
		"functionName":  "github.github_repository.GetPullRequest",
		"pullRequestID": pullRequestID,
		"owner":         owner,
		"repo":          repo,
	}

	pullRequest, _, err := client.PullRequests.Get(ctx, owner, repo, pullRequestID)
	if err != nil {
		logging.WithFields(f).WithError(err).Warn("unable to get pull request")
		return nil, err
	}

	return pullRequest, nil
}

type UserCommitSummary struct {
	SHA          string
	CommitAuthor *github.CommitAuthor
	Affiliated   bool
	Authorized   bool
}

func (u UserCommitSummary) IsValid() bool {
	return (*u.CommitAuthor.Login != "" && *u.CommitAuthor.Name != "")
}

func (u UserCommitSummary) GetDisplayText(tagUser bool) string {
	if !u.IsValid() {
		return "Invalid author details.\n"
	}
	if u.Affiliated && u.Authorized {
		return fmt.Sprintf("%s is authorized.\n ", u.getUserInfo(tagUser))
	}
	if u.Affiliated {
		return fmt.Sprintf("%s is associated with a company, but not an approval list.\n", u.getUserInfo(tagUser))
	} else {
		return fmt.Sprintf("%s is not associated with a company.\n", u.getUserInfo(tagUser))
	}
}

func (u UserCommitSummary) getUserInfo(tagUser bool) string {
	userInfo := ""
	tagValue := ""
	var sb strings.Builder
	sb.WriteString(userInfo)

	if tagUser {
		tagValue = "@"
	}
	if u.CommitAuthor.Login != nil {
		sb.WriteString(fmt.Sprintf("login: %s%s / ", tagValue, *u.CommitAuthor.Login))
	}
	if u.CommitAuthor.Name != nil {
		sb.WriteString(fmt.Sprintf("%sname: %s / ", userInfo, *u.CommitAuthor.Name))
	}
	return strings.Replace(sb.String(), "/ $", "", -1)
}

// GetPullRequestCommitAuthors gets UserSummary and latesSHA value
func GetPullRequestCommitAuthors(ctx context.Context, installationID int64, pullRequestID int, owner, repo string) ([]*UserCommitSummary, *string, error) {
	f := logrus.Fields{
		"functionName":  "github.github_repository.GetPullRequestCommitAuthors",
		"pullRequestID": pullRequestID,
	}
	var userCommitSummary []*UserCommitSummary

	client, err := NewGithubAppClient(installationID)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("unable to create Github client")
		return nil, nil, err
	}

	commits, _, comErr := client.PullRequests.ListCommits(ctx, owner, repo, pullRequestID, &github.ListOptions{})
	if comErr != nil {
		log.WithFields(f).WithError(comErr).Warn("unable to get commits")
		return nil, nil, comErr
	}

	for _, commit := range commits {
		userCommitSummary = append(userCommitSummary, &UserCommitSummary{
			SHA:          *commit.SHA,
			CommitAuthor: commit.Commit.Author,
		})
	}

	// get latest commit SHA
	latestCommitSHA := commits[len(commits)-1].SHA

	return userCommitSummary, latestCommitSHA, nil

}

func UpdatePullRequest(ctx context.Context, apiBaseURL string, installationID int64, pullRequestID int, owner, repo, repositoryID, repositoryType string, signed []*UserCommitSummary, missing []*UserCommitSummary, CLALandingPage, latestSHA string) error {
	f := logrus.Fields{
		"functionName":   "github.github_repository.UpdatePullRequest",
		"installationID": installationID,
		"owner":          owner,
		"repo":           repo,
		"pullRequestID":  pullRequestID,
		"repositoryID":   repositoryID,
	}

	client, err := NewGithubAppClient(installationID)
	if err != nil {
		log.WithFields(f).WithError(err).Warn("unable to create Github client")
		return err
	}

	if len(missing) > 0 {
		helpURL := ""
		text := ""
		var sb strings.Builder
		sb.WriteString(text)
		for _, userSummary := range missing {
			if !userSummary.IsValid() {
				helpURL = helpLink
			} else {
				helpURL = getFullSignURL(installationID, apiBaseURL, repositoryType, repositoryID, strconv.Itoa(pullRequestID))
			}
			if userSummary.SHA != latestSHA {
				continue
			}
			sb.WriteString(userSummary.GetDisplayText(true))
		}
		log.WithFields(f).Debugf("creating check run for PR:%d on repo: %s , sha: %s ", pullRequestID, repo, latestSHA)
		text = sb.String()
		status := "completed"
		conclusion := "action_required"
		title := "EasyCLA: Signed CLA not found"
		summaryList := "One or more committers are authorized under a signed CLA"
		checkRunOptions := github.CreateCheckRunOptions{
			Name:       "CLA check",
			HeadSHA:    latestSHA,
			Status:     &status,
			Conclusion: &conclusion,
			DetailsURL: &helpURL,
			Output: &github.CheckRunOutput{
				Title:   &title,
				Summary: &summaryList,
				Text:    &text,
			},
		}
		_, _, checkRunErr := client.Checks.CreateCheckRun(ctx, owner, repo, checkRunOptions)
		if checkRunErr != nil {
			log.WithFields(f).WithError(checkRunErr).Warnf("unable to create check run update for repo: %s", repo)
			return checkRunErr
		}
	}

	// Update the comment
	body := assembleClaComment(apiBaseURL, repositoryType, installationID, repositoryID, pullRequestID, signed, missing, CLALandingPage)
	hasPreviouslyFailed, comment, failedErr := hasCheckPreviouslyFailed(ctx, client, owner, repo, pullRequestID)
	if failedErr != nil {
		log.WithFields(f).WithError(failedErr).Warnf("unable to check previously failed status PR: %d", pullRequestID)
		return failedErr
	}

	if len(missing) == 0 {
		// After Issue #167 was in place, they decided via Issue #289 that we
		// DO want to update the comment, but only after we've previously failed
		if hasPreviouslyFailed {
			log.WithFields(f).Debugf("Found previously failed checks - updating CLA comment in PR : %d", pullRequestID)
			comment.Body = &body
			_, _, editErr := client.Issues.EditComment(ctx, owner, repo, *comment.ID, comment)
			if editErr != nil {
				log.WithFields(f).WithError(editErr).Warnf("unable to edit comment: %+v", comment)
				return editErr
			}
		}
		log.WithFields(f).Debugf("EasyCLA App checks pass for PR: %d with authors : %+v", pullRequestID, signed)
	} else {
		// Per Issue #167, only add a comment if check fails
		// update_cla_comment(pull_request, body)
		if hasPreviouslyFailed {
			log.WithFields(f).Debugf("Found previously failed checks - updating CLA comment in PR : %d", pullRequestID)
			comment.Body = &body
			_, _, editErr := client.Issues.EditComment(ctx, owner, repo, *comment.ID, comment)
			if editErr != nil {
				log.WithFields(f).WithError(editErr).Warnf("unable to edit comment: %+v", comment)
				return editErr
			}
		} else {
			_, _, createErr := client.Issues.Create(ctx, owner, repo, &github.IssueRequest{Body: &body})
			if createErr != nil {
				log.WithFields(f).WithError(createErr).Debugf("unable to create comment for content : %s", body)
				return createErr
			}
		}
		log.WithFields(f).Debugf("EasyCLA App checks fail for PR: %d. CLA signatures with signed authors: %+v and missing authors: %+v", pullRequestID, signed, missing)
	}

	// Update the status
	context_name := "communitybridge/cla"

	var state string
	var ClASigned bool
	const failure = "failure"
	const success = "success"

	if len(missing) > 0 {
		state = failure
		ClASigned = false
	} else if len(signed) > 0 {
		state = success
		ClASigned = true
	} else {
		state = failure
	}

	log.WithFields(f).Debugf("Creating status for PR: %d", pullRequestID)
	context, body := assembleCLAStatus(context_name, ClASigned)

	signURL := getFullSignURL(installationID, apiBaseURL, repositoryType, repositoryID, strconv.Itoa(pullRequestID))
	_, resp, commitStatusErr := createCommitStatus(ctx, client, owner, repo, latestSHA, StatusRequest{
		State:       &state,
		TargetURL:   &signURL,
		Description: &body,
		Context:     &context,
	})

	if commitStatusErr != nil {
		log.WithFields(f).Debugf("Could not POST status on PR: %d commit: %s Response code: %d", pullRequestID, latestSHA, resp.StatusCode)
		return commitStatusErr
	}

	return nil
}

func getFullSignURL(installationID int64, apiBaseURL, repositoryType, repositoryID, pullRequestID string) string {
	return fmt.Sprintf("%s/v2/repository-provider/%s/sign/%d/%s/%s/#/?version=2", apiBaseURL, repositoryType, installationID, repositoryID, pullRequestID)
}

func assembleClaComment(apiBaseURL, repositoryType string, installationID int64, repositoryID string, pullRequestID int, signed, missing []*UserCommitSummary, CLALandingPage string) string {
	var allSigned bool
	var noUserID bool
	var isApproved = false
	signURL := getFullSignURL(installationID, apiBaseURL, repositoryType, repositoryID, strconv.Itoa(pullRequestID))
	commentBody := getCommentBody(signURL, signed, missing)
	if len(missing) == 0 {
		allSigned = true
	}
	commentBadge := getCommentBadge(signURL, allSigned, noUserID, isApproved, apiBaseURL, CLALandingPage)
	return fmt.Sprintf("%s<br />%s", commentBadge, commentBody)
}

func assembleCLAStatus(authorName string, signed bool) (string, string) {
	author := authorName
	if authorName == "" {
		author = unknown
	}
	if signed {
		return author, "EasyCLA check passed. You are authorized to contribute."
	}
	return author, "Missing CLA Authorization."
}

func getAggregatedAuthorSummary(summaryList []*UserCommitSummary) map[string][]*UserCommitSummary {
	// keep track of authors commit aggregated
	aggregated := make(map[string][]*UserCommitSummary, 0)
	for _, committer := range summaryList {
		authorInfo := committer.getUserInfo(false)
		if _, found := aggregated[authorInfo]; !found {
			aggregated[authorInfo] = []*UserCommitSummary{
				committer,
			}
		} else {
			aggregated[authorInfo] = append(aggregated[authorInfo], committer)
		}
	}
	return aggregated
}

func getCommentBody(signURL string, signed, missing []*UserCommitSummary) string {
	failed := ":x:"
	success := ":white_check_mark:"

	var committersComment strings.Builder
	var text string

	// start of the HTML to render the list of committers
	if len(signed) > 0 || len(missing) > 0 {
		committersComment.WriteString("<ul>")
	}

	if len(signed) > 0 {
		summaryList := getAggregatedAuthorSummary(signed)
		for _, committer := range signed {
			authorInfo := committer.getUserInfo(false)
			var aggregated []string
			for _, v := range summaryList[authorInfo] {
				aggregated = append(aggregated, v.SHA)
			}
			committersComment.WriteString(fmt.Sprintf("<li>%s %s (%s)", success, authorInfo, strings.Join(aggregated, ", ")))
		}
	}
	if len(missing) > 0 {
		supportURL := "https://jira.linuxfoundation.org/servicedesk/customer/portal/4"
		summaryList := getAggregatedAuthorSummary(missing)
		helpURL := helpLink
		for authorInfo, v := range summaryList {
			if authorInfo == unknown {
				var shas []string
				for _, item := range v {
					shas = append(shas, item.SHA)
				}
				committersComment.WriteString(fmt.Sprintf(`<li>%s The commit (%s).
				This user is missing the Users ID, preventing the EasyCLA check, 
				<a href='%s' target='_blank'>Consult GitHub Help<a> to resolve.
				For further assistance with EasyCLA, 
				<a href='%s' target='_blank'>please submit a support request ticket</a>.</li>`, failed, strings.Join(shas, ", "), helpURL, supportURL))
			} else {
				var missingAffiliations []*UserCommitSummary
				var shas []string
				for _, item := range v {
					if !item.Affiliated {
						missingAffiliations = append(missingAffiliations, item)
					}
					shas = append(shas, item.SHA)
				}
				if len(missingAffiliations) > 0 {
					committersComment.WriteString(fmt.Sprintf(`<li>%s %s (%s).
					This user is authorized, but they must confirm their affiliation with their company.
					Start the authorization process
					<a href='%s' target='_blank'> by clicking here</a>, click \"Corporate\", 
					select the appropriate company from the list, then confirm
					your affiliation on the page that appears. 
					For further assistance with EasyCLA,
					<a href='%s' target='_blank'>please submit a support request ticket</a>.</li>
					`, failed, authorInfo, strings.Join(shas, ", "), signURL, supportURL))
				} else {
					var commitShas []string
					for _, item := range v {
						commitShas = append(shas, item.SHA)
					}
					committersComment.WriteString(fmt.Sprintf(`<li>
					<a href='%s' target='_blank'>%s</a> -
					%s. The commit %s is not authorized under a signed CLA.
					<a href='%s' target='_blank'>Please click here to be authorized</a>
					For further assistance with EasyCLA,
					<a href='%s' target='_blank'> please submit a support request ticket</a>.
					</li>
					`, signURL, failed, authorInfo, strings.Join(commitShas, ", "), signURL, supportURL))
				}
			}
		}

	}

	if len(signed) > 0 || len(missing) > 0 {
		committersComment.WriteString("</ul>")
	}

	if len(signed) > 0 && len(missing) == 0 {
		text = "The committers listed above are authorized under a signed CLA."
	}

	return fmt.Sprintf("%s%s", text, committersComment.String())
}

func getCommentBadge(signURL string, allSigned, missingUserID, isApproved bool, ClaV1ApiURL string, CLALandingPage string) string {
	var alt string
	badgeHyperLink := signURL
	var badgeURL string

	if allSigned {
		badgeURL = fmt.Sprintf("%s/cla-signed.svg", ClaV1ApiURL)
		badgeHyperLink = fmt.Sprintf("%s#/?version=2", CLALandingPage)
		alt = "CLA Signed"
	} else if missingUserID {
		badgeURL = fmt.Sprintf("%s/cla-missing-id.svg", CLALandingPage)
		alt = "CLA Missing ID"
	} else if isApproved {
		badgeURL = fmt.Sprintf("%s/cla-confirmation-needed.svg", CLALandingPage)
		alt = "CLA Confirmation Needed"
	} else {
		badgeURL = fmt.Sprintf("%s/cla-not-signed.svg", CLALandingPage)
		alt = "CLA Not Signed"
	}

	return fmt.Sprintf("<a href=%s><img src=\"%s\" alt=\"%s\" align=left height=\"28\" width=\"328\"></a><br/>", badgeHyperLink, badgeURL, alt)
}

func hasCheckPreviouslyFailed(ctx context.Context, client *github.Client, owner, repo string, pullRequestID int) (bool, *github.IssueComment, error) {
	f := logrus.Fields{
		"functionName": "github.github_repository.GetPullRequest",
	}
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, pullRequestID, &github.IssueListCommentsOptions{})
	// comments, _, err := client.Issues.ListComments(ctx, owner, repo, pullRequestID, &github.IssueListCommentsOptions{})
	if err != nil {
		log.WithFields(f).WithError(err).Warnf("unable to get fetch comments for repo: %s, pr: %d", repo, pullRequestID)
		return false, nil, err
	}
	for _, comment := range comments {
		if strings.Contains(*comment.Body, "is not authorized under a signed CLA") {
			return true, comment, nil
		}
		if strings.Contains(*comment.Body, "they must confirm their affiliation") {
			return true, comment, nil
		}
		if strings.Contains(*comment.Body, "is missing the User") {
			return true, comment, nil
		}
	}
	return false, nil, nil
}

// GetRepositoryByExternalID finds github repository by github repository id
func GetRepositoryByExternalID(ctx context.Context, installationID, id int64) (*github.Repository, error) {
	client, err := NewGithubAppClient(installationID)
	if err != nil {
		return nil, err
	}
	org, resp, err := client.Repositories.GetByID(ctx, id)
	if err != nil {
		logging.Warnf("GitHubGetRepository %v failed. error = %s", id, err.Error())
		if resp.StatusCode == 404 {
			return nil, ErrGithubRepositoryNotFound
		}
		return nil, err
	}
	return org, nil
}

func createCommitStatus(ctx context.Context, client *github.Client, owner, repo, latestSHA string, statusRequest StatusRequest) (*Status, *github.Response, error) {

	endpoint := fmt.Sprintf("/repos/%s/%s/statuses/%s", owner, repo, latestSHA)
	req, err := client.NewRequest("POST", endpoint, statusRequest)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	status := new(Status)

	resp, err := client.Do(ctx, req, status)
	if err != nil {
		return nil, resp, err
	}
	return status, resp, nil
}

// GetRepositories gets github repositories by organization
func GetRepositories(ctx context.Context, organizationName string) ([]*github.Repository, error) {
	f := logrus.Fields{
		"functionName":     "GetRepositories",
		utils.XREQUESTID:   ctx.Value(utils.XREQUESTID),
		"organizationName": organizationName,
	}

	// Get the client with token
	client := NewGithubOauthClient()

	var responseRepoList []*github.Repository
	var nextPage = 1
	for {
		// API https://docs.github.com/en/free-pro-team@latest/rest/reference/repos
		// API Pagination: https://docs.github.com/en/free-pro-team@latest/rest/guides/traversing-with-pagination
		repoList, resp, err := client.Repositories.ListByOrg(ctx, organizationName, &github.RepositoryListByOrgOptions{
			Type:      "public",
			Sort:      "full_name",
			Direction: "asc",
			ListOptions: github.ListOptions{
				Page:    nextPage,
				PerPage: 100,
			},
		})
		if err != nil {
			log.WithFields(f).WithError(err).Warn("unable to list repositories for organization")
			if resp != nil && resp.StatusCode == 404 {
				return nil, ErrGithubOrganizationNotFound
			}
			return nil, err
		}

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			msg := fmt.Sprintf("GetRepositories %s failed with no success response code %d. error = %s", organizationName, resp.StatusCode, err.Error())
			log.WithFields(f).Warnf(msg)
			return nil, errors.New(msg)
		}

		// Append our results to the response...
		responseRepoList = append(responseRepoList, repoList...)
		// if no more pages...
		if resp.NextPage == 0 {
			break
		}

		// update our next page value
		nextPage = resp.NextPage
	}

	return responseRepoList, nil
}
