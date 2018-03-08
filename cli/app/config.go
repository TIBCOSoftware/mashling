/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

type ComponentType int

const (
	LINK ComponentType = 1 + iota
	TRIGGER
	HANDLER
	ALL
)

var ctStr = [...]string{
	"all",
	"link",
	"trigger",
	"handler",
}

type Component struct {
	Name string
	Type ComponentType
	Ref  string
}

func (m ComponentType) String() string { return ctStr[m] }
