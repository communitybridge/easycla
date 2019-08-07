// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/communitybridge/easycla/cla-backend-go/cmd"
)

var (
	// version the application version
	version string

	// build/Commit the application build number
	commit string

	// branch the build branch
	branch string

	// build date
	buildDate string
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Branch = branch
	cmd.BuildDate = buildDate

	cmd.Execute()
}
