package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// JWT is a JWT validation service.
type JWT struct {
	Request  JWTRequest  `json:"request"`
	Response JWTResponse `json:"response"`
}

// JWTRequest is an JWT validation request.
type JWTRequest struct {
	Token         string `json:"token"`
	Key           string `json:"key"`
	SigningMethod string `json:"signingMethod"`
	Issuer        string `json:"iss"`
	Subject       string `json:"sub"`
	Audience      string `json:"aud"`
}

// JWTResponse is a parsed JWT response.
type JWTResponse struct {
	Valid             bool        `json:"valid"`
	Token             ParsedToken `json:"token"`
	ValidationMessage string      `json:"validationMessage"`
	Error             bool        `json:"error"`
	ErrorMessage      string      `json:"errorMessage"`
}

// ParsedToken is a parsed JWT token.
type ParsedToken struct {
	Claims        jwt.MapClaims          `json:"claims"`
	Signature     string                 `json:"signature"`
	SigningMethod string                 `json:"signingMethod"`
	Header        map[string]interface{} `json:"header"`
}

// Execute invokes this JWT service.
func (j *JWT) Execute() error {
	j.Response = JWTResponse{}
	token, err := jwt.Parse(j.Request.Token, func(token *jwt.Token) (interface{}, error) {
		// Make sure signing alg matches what we expect
		switch strings.ToLower(j.Request.SigningMethod) {
		case "hmac":
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
		case "ecdsa":
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
		case "rsa":
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
		case "rsapss":
			if _, ok := token.Method.(*jwt.SigningMethodRSAPSS); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
		case "":
			// Just continue
		default:
			return nil, fmt.Errorf("Unknown signing method expected: %v", j.Request.SigningMethod)
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if j.Request.Issuer != "" && !claims.VerifyIssuer(j.Request.Issuer, true) {
				return nil, jwt.NewValidationError("iss claims do not match", jwt.ValidationErrorIssuer)
			}
			if j.Request.Audience != "" && !claims.VerifyAudience(j.Request.Audience, true) {
				return nil, jwt.NewValidationError("aud claims do not match", jwt.ValidationErrorAudience)
			}
			subClaim, sok := claims["sub"].(string)
			if j.Request.Subject != "" && (!sok || strings.Compare(j.Request.Subject, subClaim) != 0) {
				return nil, jwt.NewValidationError("sub claims do not match", jwt.ValidationErrorClaimsInvalid)
			}
		} else {
			return nil, jwt.NewValidationError("unable to parse claims", jwt.ValidationErrorClaimsInvalid)
		}

		return []byte(j.Request.Key), nil
	})
	if token != nil && token.Valid {
		j.Response.Valid = true
		j.Response.Token = ParsedToken{Signature: token.Signature, SigningMethod: token.Method.Alg(), Header: token.Header}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			j.Response.Token.Claims = claims
		}
		return err
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		j.Response.Valid = false
		j.Response.ValidationMessage = ve.Error()
	} else {
		j.Response.Valid = false
		j.Response.Error = true
		j.Response.ValidationMessage = err.Error()
		j.Response.ErrorMessage = err.Error()
	}
	return nil
}

// InitializeJWT initializes a JWT validation service with provided settings.
func InitializeJWT(settings map[string]interface{}) (jwtService *JWT, err error) {
	jwtService = &JWT{}
	request := JWTRequest{}
	jwtService.Request = request
	err = jwtService.setRequestValues(settings)
	return jwtService, err
}

// UpdateRequest updates a JWT validation service with new provided settings.
func (j *JWT) UpdateRequest(values map[string]interface{}) (err error) {
	return j.setRequestValues(values)
}

func (j *JWT) setRequestValues(settings map[string]interface{}) error {
	for k, v := range settings {
		switch k {
		case "token":
			token, ok := v.(string)
			if !ok {
				return errors.New("invalid type for token")
			}
			// Try to scrub any extra noise from the token string
			tokenSplit := strings.Fields(token)
			token = tokenSplit[len(tokenSplit)-1]
			j.Request.Token = token
		case "key":
			key, ok := v.(string)
			if !ok {
				return errors.New("invalid type for key")
			}
			j.Request.Key = key
		case "signingMethod":
			signingMethod, ok := v.(string)
			if !ok {
				return errors.New("invalid type for signingMethod")
			}
			j.Request.SigningMethod = signingMethod
		case "issuer":
			issuer, ok := v.(string)
			if !ok {
				return errors.New("invalid type for issuer")
			}
			j.Request.Issuer = issuer
		case "subject":
			subject, ok := v.(string)
			if !ok {
				return errors.New("invalid type for subject")
			}
			j.Request.Subject = subject
		case "audience":
			audience, ok := v.(string)
			if !ok {
				return errors.New("invalid type for audience")
			}
			j.Request.Audience = audience
		default:
			// ignore and move on.
		}
	}
	return nil
}
