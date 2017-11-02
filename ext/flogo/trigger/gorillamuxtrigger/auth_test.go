/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package gorillamuxtrigger

import (
	"io/ioutil"
	"os"
	"testing"
)

const dataPlain string = `foo:bar
moo:mar`

const dataHashed string = `foo:5VvmQnTXZ10wGZu_Gkjb8umfUPIOQTQ3p1YFadAWTl8=:6267beb3f851b7fee14011f6aa236556f35b186a6791b80b48341e990c367643
`

func TestBasicAuthVerify(t *testing.T) {
	b := basicAuth()

	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(dataPlain)); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}

	path := file.Name()
	loadCreds(path)

	if !b.authenticate(base64encode("foo", "bar")) {
		t.Error("Creds should match")
	}

	if b.authenticate(base64encode("foo", "badpass")) {
		t.Error("Creds should match")
	}

	if b.authenticate(base64encode("foo2", "bar")) {
		t.Error("User doesn't exist")
	}

}

func TestBasicAuthHashVerify(t *testing.T) {
	b := basicAuth()

	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(dataHashed)); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}

	path := file.Name()
	loadCreds(path)

	if !b.authenticate(base64encode("foo", "bar")) {
		t.Error("Creds should match")
	}

	if b.authenticate(base64encode("foo", "badpass")) {
		t.Error("Creds should match")
	}

	if b.authenticate(base64encode("foo2", "bar")) {
		t.Error("User doesn't exist")
	}

}
