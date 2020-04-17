package utils

import (
<<<<<<< HEAD
	"fmt"

=======
>>>>>>> e3b809a8a87884d8ca0ade02f84d8237db213c4d
	"github.com/LF-Engineering/lfx-kit/auth"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
)

// SetAuthUserProperties adds username and email to auth user
func SetAuthUserProperties(authUser *auth.User, xUserName *string, xEmail *string) {

	if xUserName != nil {
		authUser.UserName = *xUserName
	}
	if xEmail != nil {
		authUser.Email = *xEmail
	}
<<<<<<< HEAD
	log.Infof(fmt.Sprintf("authuser x-username:%s and x-email:%s", authUser.UserName, authUser.Email))
=======
	log.Debugf("authuser x-username:%s and x-email:%s", authUser.UserName, authUser.Email)
>>>>>>> e3b809a8a87884d8ca0ade02f84d8237db213c4d
}
