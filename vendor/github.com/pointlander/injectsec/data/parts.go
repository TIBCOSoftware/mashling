// Copyright 2018 The InjectSec Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var (
	// ErrorNotSupported means the part type is not supported
	ErrorNotSupported = fmt.Errorf("part type is not supported")
)

// PartType is a type of a part
type PartType int

const (
	// PartTypeLiteral is a literal part type
	PartTypeLiteral PartType = iota
	// PartTypeNumber is a number
	PartTypeNumber
	// PartTypeName is a name
	PartTypeName
	// PartTypeOr is a or part type with spaces
	PartTypeOr
	// PartTypeHexOr is a or part type with hex spaces
	PartTypeHexOr
	// PartTypeAnd is a and part type with spaces
	PartTypeAnd
	// PartTypeSpaces represents spaces
	PartTypeSpaces
	// PartTypeSpacesOptional represents spaces or nothing
	PartTypeSpacesOptional
	// PartTypeHexSpaces represents hex spaces
	PartTypeHexSpaces
	// PartTypeHexSpaces represents hex spaces or nothing
	PartTypeHexSpacesOptional
	// PartTypeComment represents a comment
	PartTypeComment
	// PartTypeObfuscated is an obfuscated string
	PartTypeObfuscated
	// PartTypeObfuscatedWithComments is an comment obfuscated string
	PartTypeObfuscatedWithComments
	// PartTypeHex is a hex string
	PartTypeHex
	// PartTypeNumberList is a list of numbers
	PartTypeNumberList
	// PartTypeScientificNumber is a sciencetific number
	PartTypeScientificNumber
	// PartTypeSQL is a sql part type
	PartTypeSQL
)

// Part is part of a regex
type Part struct {
	PartType
	Variable int
	Literal  string
	Max      int
	Parts    *Parts
}

// Parts is a bunch of Part
type Parts struct {
	Parts []Part
}

// NewParts creates a new set of parts
func NewParts() *Parts {
	return &Parts{
		Parts: make([]Part, 0, 16),
	}
}

// AddType adds a part with type to the parts
func (p *Parts) AddType(partType PartType) {
	part := Part{
		PartType: partType,
	}
	p.Parts = append(p.Parts, part)
	return
}

// AddParts adds parts
func (p *Parts) AddParts(partType PartType, adder func(p *Parts)) {
	part := Part{
		PartType: partType,
		Parts:    NewParts(),
	}
	adder(part.Parts)
	p.Parts = append(p.Parts, part)
}

// AddLiteral adds a literal to the parts
func (p *Parts) AddLiteral(literal string) {
	part := Part{
		PartType: PartTypeLiteral,
		Literal:  literal,
	}
	p.Parts = append(p.Parts, part)
	return
}

// AddNumber adds a literal to the parts
func (p *Parts) AddNumber(variable, max int) {
	part := Part{
		PartType: PartTypeNumber,
		Variable: variable,
		Max:      max,
	}
	p.Parts = append(p.Parts, part)
	return
}

// AddName adss a PartTypeName
func (p *Parts) AddName(variable int) {
	part := Part{
		PartType: PartTypeName,
		Variable: variable,
	}
	p.Parts = append(p.Parts, part)
	return
}

// AddOr adds a part type or
func (p *Parts) AddOr() {
	p.AddType(PartTypeOr)
}

// AddHexOr adds a part type hex or
func (p *Parts) AddHexOr() {
	p.AddType(PartTypeHexOr)
}

// AddAnd adds a part type and
func (p *Parts) AddAnd() {
	p.AddType(PartTypeAnd)
}

// AddSpaces adds a part type spaces
func (p *Parts) AddSpaces() {
	p.AddType(PartTypeSpaces)
}

// AddSpacesOptional adds a part type spaces optional
func (p *Parts) AddSpacesOptional() {
	p.AddType(PartTypeSpacesOptional)
}

// AddHexSpaces adds a part type hex spaces
func (p *Parts) AddHexSpaces() {
	p.AddType(PartTypeHexSpaces)
}

// AddHexSpaces adds a part type hex spaces optional
func (p *Parts) AddHexSpacesOptional() {
	p.AddType(PartTypeHexSpacesOptional)
}

// AddComment adds a part type comment
func (p *Parts) AddComment() {
	p.AddType(PartTypeComment)
}

// AddHex adds a hex type
func (p *Parts) AddHex(max int) {
	part := Part{
		PartType: PartTypeHex,
		Max:      max,
	}
	p.Parts = append(p.Parts, part)
}

// AddNumberList adds a list of numbers
func (p *Parts) AddNumberList(max int) {
	part := Part{
		PartType: PartTypeNumberList,
		Max:      max,
	}
	p.Parts = append(p.Parts, part)
}

// AddBenchmark add a SQL benchmark statement
func (p *Parts) AddBenchmark() {
	p.AddLiteral("benchmark(")
	p.AddSpacesOptional()
	p.AddNumber(1024, 10000000)
	p.AddSpacesOptional()
	p.AddLiteral(",MD5(")
	p.AddSpacesOptional()
	p.AddNumber(1025, 10000000)
	p.AddSpacesOptional()
	p.AddLiteral("))#")
}

// AddWaitfor adds a waitfor statement
func (p *Parts) AddWaitfor() {
	p.AddLiteral(";waitfor")
	p.AddSpaces()
	p.AddLiteral("delay")
	p.AddSpaces()
	p.AddLiteral("'")
	p.AddNumber(1024, 24)
	p.AddLiteral(":")
	p.AddNumber(1025, 60)
	p.AddLiteral(":")
	p.AddNumber(1026, 60)
	p.AddLiteral("'--")
}

// AddSQL adds a part type SQL
func (p *Parts) AddSQL() {
	p.AddType(PartTypeSQL)
}

// RegexFragment is part of a regex
func (p *Parts) RegexFragment() (string, error) {
	last, regex := len(p.Parts)-1, ""
	for i, part := range p.Parts {
		switch part.PartType {
		case PartTypeLiteral:
			regex += regexp.QuoteMeta(strings.ToLower(part.Literal))
		case PartTypeNumber:
			regex += "-?[[:digit:]]+([[:space:]]*[+\\-*/][[:space:]]*-?[[:digit:]]+)*"
		case PartTypeName:
			regex += "[\\p{L}_\\p{Cc}][\\p{L}\\p{N}_\\p{Cc}]*"
		case PartTypeOr:
			a := ""
			if i == 0 {
				a += "[[:space:]]*"
			} else {
				a += "[[:space:]]+"
			}
			a += "or"
			if i == last {
				a += "[[:space:]]*"
			} else {
				a += "[[:space:]]+"
			}
			b := "[[:space:]]*" + regexp.QuoteMeta("||") + "[[:space:]]*"
			regex += "((" + a + ")|(" + b + "))"
		case PartTypeHexOr:
			hex := "(" + regexp.QuoteMeta("%20") + ")"
			a := ""
			if i == 0 {
				a += hex + "*"
			} else {
				a += hex + "+"
			}
			a += "or"
			if i == last {
				a += hex + "*"
			} else {
				a += hex + "+"
			}
			b := hex + "*" + regexp.QuoteMeta("||") + hex + "*"
			regex += "((" + a + ")|(" + b + "))"
		case PartTypeAnd:
			a := ""
			if i == 0 {
				a += "[[:space:]]*"
			} else {
				a += "[[:space:]]+"
			}
			a += "and"
			if i == last {
				a += "[[:space:]]*"
			} else {
				a += "[[:space:]]+"
			}
			b := "[[:space:]]*" + regexp.QuoteMeta("&&") + "[[:space:]]*"
			regex += "((" + a + ")|(" + b + "))"
		case PartTypeSpaces:
			regex += "[[:space:]]+"
		case PartTypeSpacesOptional:
			regex += "[[:space:]]*"
		case PartTypeHexSpaces:
			regex += "(" + regexp.QuoteMeta("%20") + ")+"
		case PartTypeHexSpacesOptional:
			regex += "(" + regexp.QuoteMeta("%20") + ")*"
		case PartTypeComment:
			regex += regexp.QuoteMeta("/*") + "[[:alnum:][:space:]]*" + regexp.QuoteMeta("*/")
		case PartTypeObfuscated:
			regex += "'[\\p{L}\\p{N}_\\p{Cc}[:space:]]*'([[:space:]]*([|]{2}|[+])[[:space:]]*'[\\p{L}\\p{N}_\\p{Cc}[:space:]]*')*"
		case PartTypeObfuscatedWithComments:
			regex += "([\\p{L}\\p{N}_\\p{Cc}[:space:]]+|(/[*][\\p{L}\\p{N}_\\p{Cc}[:space:]]*[*]/))+"
		case PartTypeHex:
			regex += "0x[[:xdigit:]]+"
		case PartTypeNumberList:
			regex += "([[:digit:]]*[[:space:]]*,[[:space:]]*)*[[:digit:]]+"
		case PartTypeScientificNumber:
			regex += "[+-]?[[:digit:]]+" + regexp.QuoteMeta(".") + "?[[:digit:]]*(e[+-]?[[:digit:]]+)?"
		case PartTypeSQL:
			regex += "select([[:space:]]+[\\p{L}\\p{N}_\\p{Cc}]+[[:space:]]*,)*([[:space:]]+[\\p{L}\\p{N}_\\p{Cc}]+)" +
				"[[:space:]]+from([[:space:]]+[\\p{L}\\p{N}_\\p{Cc}]+[[:space:]]*,)*([[:space:]]+[\\p{L}\\p{N}_\\p{Cc}]+)" +
				"[[:space:]]+where[[:space:]]+[\\p{L}\\p{N}_\\p{Cc}]+[[:space:]]*[=><][[:space:]]*[\\p{L}\\p{N}_\\p{Cc}]+"
		}
	}
	return regex, nil
}

// Regex generates a regex from the parts
func (p *Parts) Regex() (string, error) {
	regex, err := p.RegexFragment()
	if err != nil {
		return "", err
	}
	return "^" + regex + "$", nil
}

// Sample samples from the parts
func (p *Parts) Sample(rnd *rand.Rand) (string, error) {
	last, sample, state := len(p.Parts)-1, "", make(map[int]string)
	for i, part := range p.Parts {
		switch part.PartType {
		case PartTypeLiteral:
			sample += part.Literal
		case PartTypeNumber:
			if value, ok := state[part.Variable]; ok {
				sample += value
				break
			}
			s := strconv.Itoa(rand.Intn(part.Max))
			state[part.Variable] = s
			sample += s
		case PartTypeName:
			if value, ok := state[part.Variable]; ok {
				sample += value
				break
			}
			s, count := "", rand.Intn(8)+1
			for i := 0; i < count; i++ {
				s += string(rune(int('a') + rnd.Intn(int('z'-'a'))))
			}
			state[part.Variable] = s
			sample += s
		case PartTypeOr:
			if rnd.Intn(2) == 0 {
				if i == 0 {
					count := rnd.Intn(8)
					for i := 0; i < count; i++ {
						sample += " "
					}
				} else {
					count := rnd.Intn(8) + 1
					for i := 0; i < count; i++ {
						sample += " "
					}
				}
				sample += "or"
				if i == last {
					count := rnd.Intn(8)
					for i := 0; i < count; i++ {
						sample += " "
					}
				} else {
					count := rnd.Intn(8) + 1
					for i := 0; i < count; i++ {
						sample += " "
					}
				}
			} else {
				count := rnd.Intn(8)
				for i := 0; i < count; i++ {
					sample += " "
				}
				sample += "||"
				count = rnd.Intn(8)
				for i := 0; i < count; i++ {
					sample += " "
				}
			}
		case PartTypeHexOr:
			if rnd.Intn(2) == 0 {
				if i == 0 {
					count := rnd.Intn(8)
					for i := 0; i < count; i++ {
						sample += "%20"
					}
				} else {
					count := rnd.Intn(8) + 1
					for i := 0; i < count; i++ {
						sample += "%20"
					}
				}
				sample += "or"
				if i == last {
					count := rnd.Intn(8)
					for i := 0; i < count; i++ {
						sample += "%20"
					}
				} else {
					count := rnd.Intn(8) + 1
					for i := 0; i < count; i++ {
						sample += "%20"
					}
				}
			} else {
				count := rnd.Intn(8)
				for i := 0; i < count; i++ {
					sample += "%20"
				}
				sample += "||"
				count = rnd.Intn(8)
				for i := 0; i < count; i++ {
					sample += "%20"
				}
			}
		case PartTypeAnd:
			if rnd.Intn(2) == 0 {
				if i == 0 {
					count := rnd.Intn(8)
					for i := 0; i < count; i++ {
						sample += " "
					}
				} else {
					count := rnd.Intn(8) + 1
					for i := 0; i < count; i++ {
						sample += " "
					}
				}
				sample += "and"
				if i == last {
					count := rnd.Intn(8)
					for i := 0; i < count; i++ {
						sample += " "
					}
				} else {
					count := rnd.Intn(8) + 1
					for i := 0; i < count; i++ {
						sample += " "
					}
				}
			} else {
				count := rnd.Intn(8)
				for i := 0; i < count; i++ {
					sample += " "
				}
				sample += "&&"
				count = rnd.Intn(8)
				for i := 0; i < count; i++ {
					sample += " "
				}
			}
		case PartTypeSpaces:
			count := rnd.Intn(8) + 1
			for i := 0; i < count; i++ {
				sample += " "
			}
		case PartTypeSpacesOptional:
			count := rnd.Intn(8)
			for i := 0; i < count; i++ {
				sample += " "
			}
		case PartTypeHexSpaces:
			count := rnd.Intn(8) + 1
			for i := 0; i < count; i++ {
				sample += "%20"
			}
		case PartTypeHexSpacesOptional:
			count := rnd.Intn(8)
			for i := 0; i < count; i++ {
				sample += "%20"
			}
		case PartTypeComment:
			sample += "/*"
			count := rand.Intn(8) + 1
			for i := 0; i < count; i++ {
				sample += string(rune(int('a') + rnd.Intn(int('z'-'a'))))
			}
			sample += "*/"
		case PartTypeObfuscated:
			s, err := part.Parts.Sample(rnd)
			if err != nil {
				return "", err
			}
			sample += "'"
			for _, v := range s {
				sample += string(v)
				if rnd.Intn(3) == 0 {
					sample += "'"
					if rnd.Intn(2) == 0 {
						sample += "+"
					} else {
						sample += "||"
					}
					sample += "'"
				}
			}
			sample += "'"
		case PartTypeObfuscatedWithComments:
			s, err := part.Parts.Sample(rnd)
			if err != nil {
				return "", err
			}
			for _, v := range s {
				sample += string(v)
				if rnd.Intn(3) == 0 {
					sample += "/**/"
				}
			}
		case PartTypeHex:
			sample += fmt.Sprintf("%#x", rnd.Intn(part.Max))
		case PartTypeNumberList:
			for i := 0; i < 7; i++ {
				sample += strconv.Itoa(rand.Intn(part.Max))
				sample += ","
			}
			sample += strconv.Itoa(rand.Intn(part.Max))
		case PartTypeScientificNumber:
			const factor = 1337 * 1337
			sample += fmt.Sprintf("%E", rnd.Float64()*factor-factor/2)
		case PartTypeSQL:
			a, count := "", rand.Intn(8)+1
			for i := 0; i < count; i++ {
				a += string(rune(int('a') + rnd.Intn(int('z'-'a'))))
			}
			b, count := "", rand.Intn(8)+1
			for i := 0; i < count; i++ {
				b += string(rune(int('a') + rnd.Intn(int('z'-'a'))))
			}
			n := strconv.Itoa(rand.Intn(1337))

			sample += "select " + a + " from " + b + " where " + n + "=" + n
		}
	}
	return sample, nil
}
