package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/imdario/mergo"
)

// JWT is a JWT validation service.
type JWT struct {
	Request JWTRequest `json:"request"`
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
func (j *JWT) Execute(requestValues map[string]interface{}) (Response, error) {
	response := JWTResponse{}
	request, err := j.createRequest(requestValues)
	if err != nil {
		return response, err
	}
	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		// Make sure signing alg matches what we expect
		switch strings.ToLower(request.SigningMethod) {
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
			return nil, fmt.Errorf("Unknown signing method expected: %v", request.SigningMethod)
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if request.Issuer != "" && !claims.VerifyIssuer(request.Issuer, true) {
				return nil, jwt.NewValidationError("iss claims do not match", jwt.ValidationErrorIssuer)
			}
			if request.Audience != "" && !claims.VerifyAudience(request.Audience, true) {
				return nil, jwt.NewValidationError("aud claims do not match", jwt.ValidationErrorAudience)
			}
			subClaim, sok := claims["sub"].(string)
			if request.Subject != "" && (!sok || strings.Compare(request.Subject, subClaim) != 0) {
				return nil, jwt.NewValidationError("sub claims do not match", jwt.ValidationErrorClaimsInvalid)
			}
		} else {
			return nil, jwt.NewValidationError("unable to parse claims", jwt.ValidationErrorClaimsInvalid)
		}

		return []byte(request.Key), nil
	})
	if token != nil && token.Valid {
		response.Valid = true
		response.Token = ParsedToken{Signature: token.Signature, SigningMethod: token.Method.Alg(), Header: token.Header}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			response.Token.Claims = claims
		}
		return response, err
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		response.Valid = false
		response.ValidationMessage = ve.Error()
	} else {
		response.Valid = false
		response.Error = true
		response.ValidationMessage = err.Error()
		response.ErrorMessage = err.Error()
	}
	return response, nil
}

// InitializeJWT initializes a JWT validation service with provided settings.
func InitializeJWT(settings map[string]interface{}) (jwtService *JWT, err error) {
	jwtService = &JWT{}
	jwtService.Request, err = jwtService.createRequest(settings)
	return jwtService, err
}

func (j *JWT) createRequest(settings map[string]interface{}) (JWTRequest, error) {
	request := JWTRequest{}
	for k, v := range settings {
		switch k {
		case "token":
			token, ok := v.(string)
			if !ok {
				return request, errors.New("invalid type for token")
			}
			// Try to scrub any extra noise from the token string
			tokenSplit := strings.Fields(token)
			token = tokenSplit[len(tokenSplit)-1]
			request.Token = token
		case "key":
			key, ok := v.(string)
			if !ok {
				return request, errors.New("invalid type for key")
			}
			request.Key = key
		case "signingMethod":
			signingMethod, ok := v.(string)
			if !ok {
				return request, errors.New("invalid type for signingMethod")
			}
			request.SigningMethod = signingMethod
		case "issuer":
			issuer, ok := v.(string)
			if !ok {
				return request, errors.New("invalid type for issuer")
			}
			request.Issuer = issuer
		case "subject":
			subject, ok := v.(string)
			if !ok {
				return request, errors.New("invalid type for subject")
			}
			request.Subject = subject
		case "audience":
			audience, ok := v.(string)
			if !ok {
				return request, errors.New("invalid type for audience")
			}
			request.Audience = audience
		default:
			// ignore and move on.
		}
	}
	if err := mergo.Merge(&request, j.Request); err != nil {
		return request, errors.New("unable to merge request values")
	}
	return request, nil
}
