package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"

	"github.com/dgrijalva/jwt-go"
)

type Auth0Validator struct {
	clientID      string
	usernameClaim string
	algorithm     string
	wellKnownURL  string
}

func NewAuth0Validator(domain, clientID, usernameClaim, algorithm string) (Auth0Validator, error) {
	if domain == "" {
		return Auth0Validator{}, errors.New("missing Domain")
	}
	if clientID == "" {
		return Auth0Validator{}, errors.New("missing ClientID")
	}
	if usernameClaim == "" {
		return Auth0Validator{}, errors.New("missing UsernameClaim")
	}
	if algorithm == "" {
		return Auth0Validator{}, errors.New("missing Algorithm")
	}

	validator := Auth0Validator{
		clientID:      clientID,
		usernameClaim: usernameClaim,
		algorithm:     algorithm,
		wellKnownURL:  "https://" + path.Join(domain, ".well-known/jwks.json"),
	}

	return validator, nil
}

func (av Auth0Validator) VerifyToken(token string) (map[string]interface{}, error) {
	// Using jwt.MapClaims because our username field is set dynamically
	// based on environment
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, av.getPemCert)
	if err != nil {
		return nil, err
	}
	if !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	allClaims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("unable to map claims")
	}

	if err = allClaims.Valid(); err != nil {
		return nil, errors.New("claims are not valid")
	}

	return allClaims, nil
}

type jwks struct {
	Keys []jsonWebKeys `json:"keys"`
}

type jsonWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func (av Auth0Validator) getPemCert(token *jwt.Token) (interface{}, error) {
	cert := ""
	resp, err := http.Get(av.wellKnownURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var jwks = jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return "", err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		return "", errors.New("unable to find appropriate key")
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	if err != nil {
		return "", err
	}

	return rsaPublicKey, nil
}
