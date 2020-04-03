package utils

import (
	"fmt"

	"github.com/LF-Engineering/lfx-kit/auth"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

func SetAuthUserProperties(authUser *auth.User, xUserName *string, xEmail *string) {

	if xUserName != nil {
		authUser.UserName = *xUserName
	}
	if xEmail != nil {
		authUser.Email = *xEmail
	}
	log.Infof(fmt.Sprintf("authuser x-username:%s and x-email:%s", authUser.UserName, authUser.Email))
}
