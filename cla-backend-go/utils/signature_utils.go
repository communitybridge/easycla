// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package utils

import (
	"github.com/LF-Engineering/lfx-kit/auth"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// CurrentUserInACL is a helper function to determine if the current logged in user is in the specified CLA Manager list
func CurrentUserInACL(authUser *auth.User, managers []v1Models.User) bool {
	log.Debugf("checking if user: %+v is in the Signature ACL: %+v", authUser, managers)
	var inACL = false
	for _, manager := range managers {
		if manager.LfUsername == authUser.UserName {
			inACL = true
			break
		}
	}

	return inACL
}
