// Copyright The Linux Foundation and each contributor to LFX.
// SPDX-License-Identifier: MIT

package sign

import (
	"time"

	log "github.com/linuxfoundation/easycla/cla-backend-go/logging"
	"github.com/linuxfoundation/easycla/cla-backend-go/utils"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

func jwtToken(docusignPrivateKey string) (string, error) {
	f := logrus.Fields{
		"functionName": "v2.sign.jwtToken",
	}

	claims := jwt.MapClaims{
		"iss":   utils.GetProperty("DOCUSIGN_INTEGRATOR_KEY"), // integration key / client_id
		"sub":   utils.GetProperty("DOCUSIGN_USER_ID"),        // user_id, in PROD should be the EasyCLA Admin user account
		"aud":   utils.GetProperty("DOCUSIGN_AUTH_SERVER"),    // account.docusign.com or account-d.docusign.com (for dev)
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(), // one hour appears to be the max, minus 60 seconds
		"scope": "signature impersonation",
	}
	// log.WithFields(f).Debugf("claims: %+v", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token.Header["alg"] = "RS256"
	token.Header["typ"] = "JWT"

	privateKey, privateKeyErr := jwt.ParseRSAPrivateKeyFromPEM([]byte(docusignPrivateKey))
	if privateKeyErr != nil {
		log.WithFields(f).WithError(privateKeyErr).Warnf("problem decoding docusign private key")
		return "", privateKeyErr
	}
	// log.WithFields(f).Debugf("private key: %s", utils.GetProperty("DOCUSIGN_RSA_PRIVATE_KEY"))

	signedToken, signedTokenErr := token.SignedString(privateKey)
	if signedTokenErr != nil {
		log.WithFields(f).WithError(signedTokenErr).Warnf("problem generating the signed token")
	}
	// log.WithFields(f).Debugf("signed token: %s", signedToken)

	return signedToken, signedTokenErr
}
