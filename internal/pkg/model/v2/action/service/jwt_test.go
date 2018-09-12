package service

import (
	"testing"

	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

func TestJWT(t *testing.T) {
	service := types.Service{
		Type: "jwt",
		Settings: map[string]interface{}{
			"signingMethod": "HMAC",
			"key":           "qwertyuiopasdfghjklzxcvbnm789101",
			"aud":           "www.mashling.io",
			"iss":           "Mashling",
		},
	}
	instance, err := Initialize(service)
	if err != nil {
		t.Fatal(err)
	}
	err = instance.UpdateRequest(map[string]interface{}{
		"token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJNYXNobGluZyIsImlhdCI6MTUzNjE4MTE4OSwiZXhwIjo0MTIzODYxMTg5LCJhdWQiOiJ3d3cubWFzaGxpbmcuaW8iLCJzdWIiOiJqcm9ja2V0QGV4YW1wbGUuY29tIn0.Zl4l68Z9VcuFXEFQt8kCH7fcaiMmRRGtrC28lSWvJWw",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = instance.Execute()
	if err != nil {
		t.Fatal(err)
	}

	if !instance.(*JWT).Response.Valid {
		t.Fatal("JWT token should be valid")
	}
}
