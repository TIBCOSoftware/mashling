/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package gorillamuxtrigger

import (
	"testing"
)

func TestBasicAuthVerify(t *testing.T) {
	b := basicAuth("foo", "bar")

	if !b.authenticate(base64encode("foo", "bar")) {
		t.Error("Creds should match")
	}

}

func TestBasicAuthNegativeVerify(t *testing.T) {
	b := basicAuth("foo", "bar")

	if b.authenticate(base64encode("foo", "bar2")) {
		t.Error("Creds does not match and so should not authenticate")
	}

}
