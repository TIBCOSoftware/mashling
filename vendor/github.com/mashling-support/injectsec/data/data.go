// Copyright 2018 The InjectSec Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"math/rand"
)

// Generator generates training data
type Generator struct {
	Form      string
	Case      string
	SkipTrain bool
	SkipMatch bool
	Regex     func(p *Parts)
}

// TrainingDataGenerator returns a data generator
func TrainingDataGenerator(rnd *rand.Rand) []Generator {
	generators := []Generator{
		// Generic-SQLi.txt
		{
			Form: ")%20or%20('x'='x",
			Regex: func(p *Parts) {
				p.AddLiteral(")")
				p.AddHexOr()
				p.AddLiteral("('")
				p.AddName(0)
				p.AddLiteral("'='")
				p.AddName(0)
			},
		},
		{
			Form: "%20or%201=1",
			Regex: func(p *Parts) {
				p.AddHexOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: "; execute immediate 'sel' || 'ect us' || 'er'",
			Regex: func(p *Parts) {
				p.AddLiteral(";")
				p.AddSpaces()
				p.AddLiteral("execute")
				p.AddSpaces()
				p.AddLiteral("immediate")
				p.AddSpaces()
				p.AddParts(PartTypeObfuscated, func(p *Parts) {
					p.AddLiteral("select")
					p.AddSpaces()
					p.AddName(0)
				})
			},
		},
		{
			Form: "benchmark(10000000,MD5(1))#",
			Regex: func(p *Parts) {
				p.AddBenchmark()
			},
		},
		{
			Form: "update",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("update")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "\";waitfor delay '0:0:__TIME__'--",
			Case: "\";waitfor delay '0:0:24'--",
			Regex: func(p *Parts) {
				p.AddLiteral("\"")
				p.AddWaitfor()
			},
		},
		{
			Form: "1) or pg_sleep(__TIME__)--",
			Case: "1) or pg_sleep(123)--",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddLiteral(")")
				p.AddOr()
				p.AddLiteral("pg_sleep(")
				p.AddNumber(1, 1337)
				p.AddLiteral(")--")
			},
		},
		{
			Form: "||(elt(-3+5,bin(15),ord(10),hex(char(45))))",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("(elt(")
				p.AddNumber(0, 1337)
				p.AddLiteral(",bin(")
				p.AddNumber(1, 1337)
				p.AddLiteral("),ord(")
				p.AddNumber(2, 10)
				p.AddLiteral("),hex(char(")
				p.AddNumber(3, 256)
				p.AddLiteral("))))")
			},
		},
		{
			Form: "\"hi\"\") or (\"\"a\"\"=\"\"a\"",
			Regex: func(p *Parts) {
				p.AddLiteral("\"")
				p.AddName(0)
				p.AddLiteral("\"\")")
				p.AddOr()
				p.AddLiteral("(\"\"")
				p.AddName(1)
				p.AddLiteral("\"\"=\"\"")
				p.AddName(1)
				p.AddLiteral("\"")
			},
		},
		{
			Form: "delete",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("delete")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "like",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("like")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "\" or sleep(__TIME__)#",
			Case: "\" or sleep(123)#",
			Regex: func(p *Parts) {
				p.AddLiteral("\"")
				p.AddOr()
				p.AddLiteral("sleep(")
				p.AddNumber(0, 1337)
				p.AddLiteral(")#")
			},
		},
		{
			Form: "pg_sleep(__TIME__)--",
			Case: "pg_sleep(123)--",
			Regex: func(p *Parts) {
				p.AddLiteral("pg_sleep(")
				p.AddNumber(0, 1337)
				p.AddLiteral(")--")
			},
		},
		{
			Form: "*(|(objectclass=*))",
		},
		{
			Form: "declare @q nvarchar (200) 0x730065006c00650063 ...",
			Case: "declare @q nvarchar (200) 0x730065006c00650063",
			Regex: func(p *Parts) {
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("nvarchar")
				p.AddSpaces()
				p.AddLiteral("(")
				p.AddNumber(1, 1337)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddHex(1337 * 1337)
			},
		},
		{
			Form: " or 0=0 #",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("#")
			},
		},
		{
			Form: "insert",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("insert")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "1) or sleep(__TIME__)#",
			Case: "1) or sleep(567)#",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddLiteral(")")
				p.AddOr()
				p.AddLiteral("sleep(")
				p.AddNumber(1, 1337)
				p.AddLiteral(")#")
			},
		},
		{
			Form: ") or ('a'='a",
			Regex: func(p *Parts) {
				p.AddLiteral(")")
				p.AddOr()
				p.AddLiteral("('")
				p.AddName(0)
				p.AddLiteral("'='")
				p.AddName(0)
			},
		},
		{
			Form: "; exec xp_regread",
			Regex: func(p *Parts) {
				p.AddLiteral(";")
				p.AddSpaces()
				p.AddLiteral("exec")
				p.AddSpaces()
				p.AddLiteral("xp_regread")
			},
		},
		{
			Form: "*|",
		},
		{
			Form: "@var select @var as var into temp end --",
			Regex: func(p *Parts) {
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("as")
				p.AddSpaces()
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("into")
				p.AddSpaces()
				p.AddName(1)
				p.AddSpaces()
				p.AddLiteral("end")
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "1)) or benchmark(10000000,MD5(1))#",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddLiteral("))")
				p.AddOr()
				p.AddBenchmark()
			},
		},
		{
			Form: "asc",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("asc")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "(||6)",
			Regex: func(p *Parts) {
				p.AddLiteral("(||")
				p.AddNumber(0, 1337)
				p.AddLiteral(")")
			},
		},
		{
			Form: "\"a\"\" or 3=3--\"",
			Regex: func(p *Parts) {
				p.AddLiteral("\"")
				p.AddName(0)
				p.AddLiteral("\"\"")
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddLiteral("--\"")
			},
		},
		{
			Form: "\" or benchmark(10000000,MD5(1))#",
			Regex: func(p *Parts) {
				p.AddLiteral("\"")
				p.AddOr()
				p.AddBenchmark()
			},
		},
		{
			Form: "# from wapiti",
			Regex: func(p *Parts) {
				p.AddLiteral("#")
				p.AddSpaces()
				p.AddLiteral("from")
				p.AddSpaces()
				p.AddLiteral("wapiti")
			},
		},
		{
			Form: " or 0=0 --",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "1 waitfor delay '0:0:10'--",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("waitfor")
				p.AddSpaces()
				p.AddLiteral("delay")
				p.AddSpaces()
				p.AddLiteral("'")
				p.AddNumber(1, 24)
				p.AddLiteral(":")
				p.AddNumber(2, 60)
				p.AddLiteral(":")
				p.AddNumber(3, 60)
				p.AddLiteral("'--")
			},
		},
		{
			Form: " or 'a'='a",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(0)
				p.AddLiteral("'='")
				p.AddName(0)
			},
		},
		{
			Form: "hi or 1=1 --\"",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddSpaces()
				p.AddLiteral("--\"")
			},
		},
		{
			Form: "or a = a",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddName(0)
			},
		},
		{
			Form: " UNION ALL SELECT",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("union")
				p.AddSpaces()
				p.AddLiteral("all")
				p.AddSpaces()
				p.AddLiteral("select")
			},
		},
		{
			Form: ") or sleep(__TIME__)='",
			Case: ") or sleep(123)='",
			Regex: func(p *Parts) {
				p.AddLiteral(")")
				p.AddOr()
				p.AddLiteral("sleep(")
				p.AddNumber(0, 1337)
				p.AddLiteral(")='")
			},
		},
		{
			Form: ")) or benchmark(10000000,MD5(1))#",
			Regex: func(p *Parts) {
				p.AddLiteral("))")
				p.AddOr()
				p.AddBenchmark()
			},
		},
		{
			Form: "hi' or 'a'='a",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(1)
				p.AddLiteral("'='")
				p.AddName(1)
			},
		},
		{
			Form:      "0",
			SkipTrain: true,
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: "21 %",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("%")
			},
		},
		{
			Form: "limit",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("limit")
				p.AddSpacesOptional()
			},
		},
		{
			Form: " or 1=1",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: " or 2 > 1",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral(">")
				p.AddSpaces()
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: "\")) or benchmark(10000000,MD5(1))#",
			Regex: func(p *Parts) {
				p.AddLiteral("\"))")
				p.AddOr()
				p.AddBenchmark()
			},
		},
		{
			Form: "PRINT",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("print")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "hi') or ('a'='a",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("')")
				p.AddOr()
				p.AddLiteral("('")
				p.AddName(1)
				p.AddLiteral("'='")
				p.AddName(1)
			},
		},
		{
			Form: " or 3=3",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: "));waitfor delay '0:0:__TIME__'--",
			Case: "));waitfor delay '0:0:42'--",
			Regex: func(p *Parts) {
				p.AddLiteral("))")
				p.AddWaitfor()
			},
		},
		{
			Form: "a' waitfor delay '0:0:10'--",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddSpaces()
				p.AddLiteral("waitfor")
				p.AddSpaces()
				p.AddLiteral("delay")
				p.AddSpaces()
				p.AddLiteral("'")
				p.AddNumber(1, 24)
				p.AddLiteral(":")
				p.AddNumber(2, 60)
				p.AddLiteral(":")
				p.AddNumber(3, 60)
				p.AddLiteral("'--")
			},
		},
		{
			Form: "1;(load_file(char(47,101,116,99,47,112,97,115, ...",
			Case: "1;(load_file(char(47,101,116,99,47,112,97,115)))",
			Regex: func(p *Parts) {
				p.AddNumber(0, 256)
				p.AddLiteral(";(load_file(char(")
				p.AddNumberList(256)
				p.AddLiteral(")))")
			},
		},
		{
			Form: "or%201=1",
			Regex: func(p *Parts) {
				p.AddHexOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: "1 or sleep(__TIME__)#",
			Case: "1 or sleep(123)#",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddOr()
				p.AddLiteral("sleep(")
				p.AddSpacesOptional()
				p.AddNumber(1, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")#")
			},
		},
		{
			Form: "or 1=1",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: " and 1 in (select var from temp)--",
			Regex: func(p *Parts) {
				p.AddAnd()
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("in")
				p.AddSpaces()
				p.AddLiteral("(select")
				p.AddSpaces()
				p.AddName(1)
				p.AddSpaces()
				p.AddLiteral("from")
				p.AddSpaces()
				p.AddName(2)
				p.AddLiteral(")--")
			},
		},
		{
			Form: " or '7659'='7659",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("'")
				p.AddNumber(0, 1337)
				p.AddLiteral("'='")
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: " or 'text' = n'text'",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(0)
				p.AddLiteral("'")
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddLiteral("n'")
				p.AddName(0)
				p.AddLiteral("'")
			},
		},
		{
			Form: " --",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("--")
			},
		},
		{
			Form: " or 1=1 or ''='",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddOr()
				p.AddLiteral("''='")
			},
		},
		{
			Form: "declare @s varchar (200) select @s = 0x73656c6 ...",
			Case: "declare @s varchar (200) select @s = 0x73656c6",
			Regex: func(p *Parts) {
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("varchar")
				p.AddSpaces()
				p.AddLiteral("(")
				p.AddNumber(1, 200)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddHex(1337 * 1337)
			},
		},
		{
			Form: "exec xp",
			Regex: func(p *Parts) {
				p.AddLiteral("exec")
				p.AddSpaces()
				p.AddName(0)
			},
		},
		{
			Form: "; exec master..xp_cmdshell 'ping 172.10.1.255'--",
			Regex: func(p *Parts) {
				p.AddLiteral(";")
				p.AddSpaces()
				p.AddLiteral("exec")
				p.AddSpaces()
				p.AddLiteral("master..xp_cmdshell")
				p.AddSpaces()
				p.AddLiteral("'ping")
				p.AddSpaces()
				p.AddNumber(0, 256)
				p.AddLiteral(".")
				p.AddNumber(1, 256)
				p.AddLiteral(".")
				p.AddNumber(2, 256)
				p.AddLiteral(".")
				p.AddNumber(3, 256)
				p.AddLiteral("'--")
			},
		},
		{
			Form: "3.10E+17",
			Regex: func(p *Parts) {
				p.AddType(PartTypeScientificNumber)
			},
		},
		{
			Form: "\" or pg_sleep(__TIME__)--",
			Case: "\" or pg_sleep(123)--",
			Regex: func(p *Parts) {
				p.AddLiteral("\"")
				p.AddOr()
				p.AddLiteral("pg_sleep(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")--")
			},
		},
		{
			Form: "x' AND email IS NULL; --",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddAnd()
				p.AddName(1)
				p.AddSpaces()
				p.AddLiteral("is")
				p.AddSpaces()
				p.AddLiteral("null;")
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "&",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("&")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "admin' or '",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddOr()
				p.AddLiteral("'")
			},
		},
		{
			Form: " or 'unusual' = 'unusual'",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(0)
				p.AddLiteral("'")
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddLiteral("'")
				p.AddName(0)
				p.AddLiteral("'")
			},
		},
		{
			Form: "//",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("//")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "truncate",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("truncate")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "1) or benchmark(10000000,MD5(1))#",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddLiteral(")")
				p.AddOr()
				p.AddBenchmark()
			},
		},
		{
			Form: "\x27UNION SELECT",
			Regex: func(p *Parts) {
				p.AddLiteral("\x27union")
				p.AddSpaces()
				p.AddLiteral("select")
			},
		},
		{
			Form: "declare @s varchar(200) select @s = 0x77616974 ...",
			Case: "declare @s varchar(200) select @s = 0x77616974",
			Regex: func(p *Parts) {
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("varchar(")
				p.AddNumber(1, 200)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddHex(1337 * 1337)
			},
		},
		{
			Form: "tz_offset",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("tz_offset")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "sqlvuln",
			Case: "select a from b where 1=1",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddSQL()
				p.AddSpacesOptional()
			},
		},
		{
			Form: "\"));waitfor delay '0:0:__TIME__'--",
			Case: "\"));waitfor delay '0:0:23'--",
			Regex: func(p *Parts) {
				p.AddLiteral("\"))")
				p.AddWaitfor()
			},
		},
		{
			Form: "||6",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: "or%201=1 --",
			Regex: func(p *Parts) {
				p.AddHexOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "%2A%28%7C%28objectclass%3D%2A%29%29",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("%2A%28%7C%28objectclass%3D%2A%29%29")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "or a=a",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddName(0)
				p.AddLiteral("=")
				p.AddName(0)
			},
		},
		{
			Form: ") union select * from information_schema.tables;",
			Regex: func(p *Parts) {
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("union")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("*")
				p.AddSpaces()
				p.AddLiteral("from")
				p.AddSpaces()
				p.AddLiteral("information_schema.tables;")
			},
		},
		{
			Form: "PRINT @@variable",
			Regex: func(p *Parts) {
				p.AddLiteral("print")
				p.AddSpaces()
				p.AddLiteral("@@")
				p.AddName(0)
			},
		},
		{
			Form: "or isNULL(1/0) /*",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("isnull(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
				p.AddLiteral("/")
				p.AddSpacesOptional()
				p.AddLiteral("0")
				p.AddSpacesOptional()
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("/*")
			},
		},
		{
			Form: "26 %",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("%")
			},
		},
		{
			Form: "\" or \"a\"=\"a",
			Regex: func(p *Parts) {
				p.AddLiteral("\"")
				p.AddOr()
				p.AddLiteral("\"")
				p.AddName(0)
				p.AddLiteral("\"=\"")
				p.AddName(0)
			},
		},
		{
			Form: "(sqlvuln)",
			Case: "(select a from b where 1=1)",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("(")
				p.AddSQL()
				p.AddLiteral(")")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "x' AND members.email IS NULL; --",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddAnd()
				p.AddLiteral("members.email")
				p.AddSpaces()
				p.AddLiteral("is")
				p.AddSpaces()
				p.AddLiteral("null;")
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: " or 1=1--",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddLiteral("--")
			},
		},
		{
			Form: " and 1=( if((load_file(char(110,46,101,120,11 ...",
			Case: " and 1=( if((load_file(char(110,46,101,120,11)))))",
			Regex: func(p *Parts) {
				p.AddAnd()
				p.AddNumber(0, 1337)
				p.AddLiteral("=(")
				p.AddSpaces()
				p.AddLiteral("if((load_file(char(")
				p.AddNumberList(256)
				p.AddLiteral(")))))")
			},
		},
		{
			Form: "0x770061006900740066006F0072002000640065006C00 ...",
			Case: "0x770061006900740066006F0072002000640065006C00",
			Regex: func(p *Parts) {
				p.AddHex(1337 * 1336)
			},
		},
		{
			Form: "%20'sleep%2050'",
			Regex: func(p *Parts) {
				p.AddHexSpaces()
				p.AddLiteral("'sleep")
				p.AddHexSpaces()
				p.AddNumber(0, 1337)
				p.AddLiteral("'")
			},
		},
		{
			Form: "as",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("as")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "1)) or pg_sleep(__TIME__)--",
			Case: "1)) or pg_sleep(123)--",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddLiteral("))")
				p.AddOr()
				p.AddLiteral("pg_sleep(")
				p.AddSpacesOptional()
				p.AddNumber(1, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")--")
			},
		},
		{
			Form: "/**/or/**/1/**/=/**/1",
			Regex: func(p *Parts) {
				p.AddComment()
				p.AddLiteral("or")
				p.AddComment()
				p.AddNumber(0, 1337)
				p.AddComment()
				p.AddLiteral("=")
				p.AddComment()
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: " union all select @@version--",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("union")
				p.AddSpaces()
				p.AddLiteral("all")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("@@")
				p.AddName(0)
				p.AddLiteral("--")
			},
		},
		{
			Form: ",@variable",
			Regex: func(p *Parts) {
				p.AddLiteral(",@")
				p.AddName(0)
			},
		},
		{
			Form: "(sqlattempt2)",
			Case: "(select a from b where 1=1)",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("(")
				p.AddSQL()
				p.AddLiteral(")")
				p.AddSpacesOptional()
			},
		},
		{
			Form: " or (EXISTS)",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("(exists)")
			},
		},
		{
			Form: "t'exec master..xp_cmdshell 'nslookup www.googl ...",
			Case: "t'exec master..xp_cmdshell 'nslookup www.google.com",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'exec")
				p.AddSpaces()
				p.AddLiteral("master..xp_cmdshell")
				p.AddSpaces()
				p.AddLiteral("'nslookup")
				p.AddSpaces()
				p.AddName(1)
				p.AddLiteral(".")
				p.AddName(2)
				p.AddLiteral(".")
				p.AddName(3)
			},
		},
		{
			Form: "%20$(sleep%2050)",
			Regex: func(p *Parts) {
				p.AddHexSpaces()
				p.AddLiteral("$(sleep")
				p.AddHexSpaces()
				p.AddNumber(0, 1337)
				p.AddLiteral(")")
			},
		},
		{
			Form: "1 or benchmark(10000000,MD5(1))#",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddOr()
				p.AddBenchmark()
			},
		},
		{
			Form: "%20or%20''='",
			Regex: func(p *Parts) {
				p.AddHexOr()
				p.AddLiteral("''='")
			},
		},
		{
			Form: "||UTL_HTTP.REQUEST",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("utl_http.request")
				p.AddSpacesOptional()
			},
		},
		{
			Form: " or pg_sleep(__TIME__)--",
			Case: " or pg_sleep(123)--",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("pg_sleep(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")--")
			},
		},
		{
			Form: "hi' or 'x'='x';",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(1)
				p.AddLiteral("'='")
				p.AddName(1)
				p.AddLiteral("';")
			},
		},
		{
			Form: "\") or sleep(__TIME__)=\"",
			Case: "\") or sleep(857)=\"",
			Regex: func(p *Parts) {
				p.AddLiteral("\")")
				p.AddOr()
				p.AddLiteral("sleep(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")=\"")
			},
		},
		{
			Form: " or 'whatever' in ('whatever')",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(0)
				p.AddLiteral("'")
				p.AddSpaces()
				p.AddLiteral("in")
				p.AddSpaces()
				p.AddLiteral("('")
				p.AddName(0)
				p.AddLiteral("')")
			},
		},
		{
			Form: "; begin declare @var varchar(8000) set @var=' ...",
			Case: "; begin declare @var varchar(8000) set @var='abc'",
			Regex: func(p *Parts) {
				p.AddLiteral(";")
				p.AddSpaces()
				p.AddLiteral("begin")
				p.AddSpaces()
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("varchar(")
				p.AddNumber(1, 8000)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("set")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddLiteral("='")
				p.AddName(2)
				p.AddLiteral("'")
			},
		},
		{
			Form: " union select 1,load_file('/etc/passwd'),1,1,1;",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("union")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddNumber(0, 1337)
				p.AddLiteral(",load_file('/etc/passwd'),1,1,1;")
			},
		},
		{
			Form: "0x77616974666F722064656C61792027303A303A313027 ...",
			Case: "0x77616974666F722064656C61792027303A303A313027",
			Regex: func(p *Parts) {
				p.AddHex(1337 * 1337)
			},
		},
		{
			Form: "exec(@s)",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("exec(@")
				p.AddName(0)
				p.AddLiteral(")")
				p.AddSpacesOptional()
			},
		},
		{
			Form: ") or pg_sleep(__TIME__)--",
			Case: ") or pg_sleep(123)--",
			Regex: func(p *Parts) {
				p.AddLiteral(")")
				p.AddOr()
				p.AddLiteral("pg_sleep(")
				p.AddNumber(0, 1337)
				p.AddLiteral(")--")
			},
		},
		{
			Form: " union select",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("union")
				p.AddSpaces()
				p.AddLiteral("select")
			},
		},
		{
			Form: " or sleep(__TIME__)#",
			Case: " or sleep(123)#",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("sleep(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")#")
			},
		},
		{
			Form: " select * from information_schema.tables--",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("*")
				p.AddSpaces()
				p.AddLiteral("from")
				p.AddSpaces()
				p.AddLiteral("information_schema.tables--")
			},
		},
		{
			Form: "a' or 1=1--",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddLiteral("--")
			},
		},
		{
			Form: "a' or 'a' = 'a",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(1)
				p.AddLiteral("'")
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddLiteral("'")
				p.AddName(1)
			},
		},
		{
			Form: "declare @s varchar(22) select @s =",
			Regex: func(p *Parts) {
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("varchar(")
				p.AddNumber(1, 22)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("=")
			},
		},
		{
			Form: " or 2 between 1 and 3",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("between")
				p.AddSpaces()
				p.AddNumber(1, 1337)
				p.AddSpaces()
				p.AddLiteral("and")
				p.AddSpaces()
				p.AddNumber(2, 1337)
			},
		},
		{
			Form: " or a=a--",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddName(0)
				p.AddLiteral("=")
				p.AddName(0)
				p.AddLiteral("--")
			},
		},
		{
			Form: " or '1'='1",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("'")
				p.AddNumber(0, 1337)
				p.AddLiteral("'='")
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: "|",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("|")
				p.AddSpacesOptional()
			},
		},
		{
			Form: " or sleep(__TIME__)='",
			Case: " or sleep(123)='",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("sleep(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")='")
			},
		},
		{
			Form: " or 1 --'",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("--'")
			},
		},
		{
			Form: "or 0=0 #\"",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("#\"")
			},
		},
		{
			Form: "having",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("having")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "a'",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
			},
		},
		{
			Form: "\" or isNULL(1/0) /*",
			Regex: func(p *Parts) {
				p.AddLiteral("\"")
				p.AddOr()
				p.AddLiteral("isnull(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
				p.AddLiteral("/")
				p.AddSpacesOptional()
				p.AddLiteral("0")
				p.AddSpacesOptional()
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("/*")
			},
		},
		{
			Form: "declare @s varchar (8000) select @s = 0x73656c ...",
			Case: "declare @s varchar (8000) select @s = 0x73656c",
			Regex: func(p *Parts) {
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("varchar")
				p.AddSpaces()
				p.AddLiteral("(")
				p.AddNumber(1, 1337)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddHex(1337 * 1337)
			},
		},
		{
			Form: "â or 1=1 --",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "char%4039%41%2b%40SELECT",
			Regex: func(p *Parts) {
				p.AddLiteral("char%40")
				p.AddNumber(0, 256)
				p.AddLiteral("%41%2b%40select")
			},
		},
		{
			Form: "order by",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("order")
				p.AddSpaces()
				p.AddLiteral("by")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "bfilename",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("bfilename")
				p.AddSpacesOptional()
			},
		},
		{
			Form: " having 1=1--",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("having")
				p.AddSpaces()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddLiteral("--")
			},
		},
		{
			Form: ") or benchmark(10000000,MD5(1))#",
			Regex: func(p *Parts) {
				p.AddLiteral(")")
				p.AddOr()
				p.AddBenchmark()
			},
		},
		{
			Form: " or username like char(37);",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("like")
				p.AddSpaces()
				p.AddLiteral("char(")
				p.AddNumber(1, 256)
				p.AddLiteral(");")
			},
		},
		{
			Form: ";waitfor delay '0:0:__TIME__'--",
			Case: ";waitfor delay '0:0:123'--",
			Regex: func(p *Parts) {
				p.AddWaitfor()
			},
		},
		{
			Form: "\" or 1=1--",
			Regex: func(p *Parts) {
				p.AddLiteral("\"")
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddLiteral("--")
			},
		},
		{
			Form: "x' AND userid IS NULL; --",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddAnd()
				p.AddName(1)
				p.AddSpaces()
				p.AddLiteral("is")
				p.AddSpaces()
				p.AddLiteral("null;")
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "*/*",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("*")
				p.AddSpacesOptional()
				p.AddLiteral("/*")
				p.AddSpacesOptional()
			},
		},
		{
			Form: " or 'text' > 't'",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(0)
				p.AddLiteral("'")
				p.AddSpaces()
				p.AddLiteral(">")
				p.AddSpaces()
				p.AddLiteral("'")
				p.AddName(1)
				p.AddLiteral("'")
			},
		},
		{
			Form: " (select top 1",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("(select")
				p.AddSpaces()
				p.AddLiteral("top")
				p.AddSpaces()
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: " or benchmark(10000000,MD5(1))#",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddBenchmark()
			},
		},
		{
			Form: "\");waitfor delay '0:0:__TIME__'--",
			Case: "\");waitfor delay '0:0:42'--",
			Regex: func(p *Parts) {
				p.AddLiteral("\")")
				p.AddWaitfor()
			},
		},
		{
			Form: "a' or 3=3--",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddLiteral("--")
			},
		},
		{
			Form: " -- &password=",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("--")
				p.AddSpaces()
				p.AddLiteral("&")
				p.AddName(0)
				p.AddLiteral("=")
			},
		},
		{
			Form: " group by userid having 1=1--",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("group")
				p.AddSpaces()
				p.AddLiteral("by")
				p.AddSpaces()
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("having")
				p.AddSpaces()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddLiteral("--")
			},
		},
		{
			Form: " or ''='",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("''='")
			},
		},
		{
			Form: "; exec master..xp_cmdshell",
			Regex: func(p *Parts) {
				p.AddLiteral(";")
				p.AddSpaces()
				p.AddLiteral("exec")
				p.AddSpaces()
				p.AddLiteral("master..xp_cmdshell")
			},
		},
		{
			Form: "%20or%20x=x",
			Regex: func(p *Parts) {
				p.AddHexOr()
				p.AddName(0)
				p.AddLiteral("=")
				p.AddName(0)
			},
		},
		{
			Form: "select",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("select")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "\")) or sleep(__TIME__)=\"",
			Case: "\")) or sleep(123)=\"",
			Regex: func(p *Parts) {
				p.AddLiteral("\"))")
				p.AddOr()
				p.AddLiteral("sleep(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")=\"")
			},
		},
		{
			Form: "0x730065006c0065006300740020004000400076006500 ...",
			Case: "0x730065006c0065006300740020004000400076006500",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddHex(1337 * 1337)
				p.AddSpacesOptional()
			},
		},
		{
			Form: "hi' or 1=1 --",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "\") or pg_sleep(__TIME__)--",
			Case: "\") or pg_sleep(123)--",
			Regex: func(p *Parts) {
				p.AddLiteral("\")")
				p.AddOr()
				p.AddLiteral("pg_sleep(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddLiteral(")--")
			},
		},
		{
			Form: "%20or%20'x'='x",
			Regex: func(p *Parts) {
				p.AddHexOr()
				p.AddLiteral("'")
				p.AddName(0)
				p.AddLiteral("'='")
				p.AddName(0)
			},
		},
		{
			Form: " or 'something' = 'some'+'thing'",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(0)
				p.AddLiteral("'")
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddParts(PartTypeObfuscated, func(p *Parts) {
					p.AddName(0)
				})
			},
		},
		{
			Form: "exec sp",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("exec")
				p.AddSpaces()
				p.AddLiteral("sp")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "29 %",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("%")
			},
		},
		{
			Form: "(",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("(")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "Ã½ or 1=1 --",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "1 or pg_sleep(__TIME__)--",
			Case: "1 or pg_sleep(123)--",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddOr()
				p.AddLiteral("pg_sleep(")
				p.AddSpacesOptional()
				p.AddNumber(1, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")--")
			},
		},
		{
			Form: "0 or 1=1",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
			},
		},
		{
			Form: ") or (a=a",
			Regex: func(p *Parts) {
				p.AddLiteral(")")
				p.AddOr()
				p.AddLiteral("(")
				p.AddName(0)
				p.AddLiteral("=")
				p.AddName(0)
			},
		},
		{
			Form:      "uni/**/on sel/**/ect",
			SkipMatch: true,
			Regex: func(p *Parts) {
				p.AddParts(PartTypeObfuscatedWithComments, func(p *Parts) {
					p.AddLiteral("union")
					p.AddSpaces()
					p.AddLiteral("select")
				})
			},
		},
		{
			Form: "replace",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("replace")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "%27%20or%201=1",
			Regex: func(p *Parts) {
				p.AddLiteral("%27")
				p.AddHexOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: ")) or pg_sleep(__TIME__)--",
			Case: ")) or pg_sleep(343)--",
			Regex: func(p *Parts) {
				p.AddLiteral("))")
				p.AddOr()
				p.AddLiteral("pg_sleep(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")--")
			},
		},
		{
			Form: "%7C",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("%7C")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "x' AND 1=(SELECT COUNT(*) FROM tabname); --",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddAnd()
				p.AddLiteral("1=(select")
				p.AddSpaces()
				p.AddLiteral("count(*)")
				p.AddSpaces()
				p.AddLiteral("from")
				p.AddSpaces()
				p.AddName(1)
				p.AddLiteral(");")
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "&apos;%20OR",
			Regex: func(p *Parts) {
				p.AddLiteral("&apos;")
				p.AddHexOr()
			},
		},
		{
			Form: "; or '1'='1'",
			Regex: func(p *Parts) {
				p.AddLiteral(";")
				p.AddOr()
				p.AddLiteral("'")
				p.AddNumber(0, 1337)
				p.AddLiteral("'='")
				p.AddNumber(0, 1337)
				p.AddLiteral("'")
			},
		},
		{
			Form: "declare @q nvarchar (200) select @q = 0x770061 ...",
			Case: "declare @q nvarchar (200) select @q = 0x770061",
			Regex: func(p *Parts) {
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("nvarchar")
				p.AddSpaces()
				p.AddLiteral("(")
				p.AddNumber(1, 200)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddHex(1337 * 1337)
			},
		},
		{
			Form: "1 or 1=1",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
			},
		},
		{
			Form: "; exec ('sel' + 'ect us' + 'er')",
			Regex: func(p *Parts) {
				p.AddLiteral(";")
				p.AddSpaces()
				p.AddLiteral("exec")
				p.AddSpaces()
				p.AddLiteral("(")
				p.AddParts(PartTypeObfuscated, func(p *Parts) {
					p.AddLiteral("select")
					p.AddSpaces()
					p.AddName(0)
				})
				p.AddLiteral(")")
			},
		},
		{
			Form: "23 OR 1=1",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
			},
		},
		{
			Form: "/",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("/")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "anything' OR 'x'='x",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(1)
				p.AddLiteral("'='")
				p.AddName(1)
			},
		},
		{
			Form: "declare @q nvarchar (4000) select @q =",
			Regex: func(p *Parts) {
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("nvarchar")
				p.AddSpaces()
				p.AddLiteral("(")
				p.AddNumber(1, 4000)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("=")
			},
		},
		{
			Form: "or 0=0 --",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "desc",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("desc")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "||'6",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("||'6")
				p.AddSpacesOptional()
			},
		},
		{
			Form: ")",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral(")")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "1)) or sleep(__TIME__)#",
			Case: "1)) or sleep(123)#",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddLiteral("))")
				p.AddOr()
				p.AddLiteral("sleep(")
				p.AddSpacesOptional()
				p.AddNumber(1, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")#")
			},
		},
		{
			Form: "or 0=0 #",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("#")
			},
		},
		{
			Form: " select name from syscolumns where id = (sele ...",
			Case: " select name from syscolumns where id = (select 3)",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("from")
				p.AddSpaces()
				p.AddName(1)
				p.AddSpaces()
				p.AddLiteral("where")
				p.AddSpaces()
				p.AddName(2)
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddLiteral("(select ")
				p.AddNumber(3, 1337)
				p.AddLiteral(")")
			},
		},
		{
			Form: "hi or a=a",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddOr()
				p.AddName(1)
				p.AddLiteral("=")
				p.AddName(1)
			},
		},
		{
			Form: "*(|(mail=*))",
			Regex: func(p *Parts) {
				p.AddLiteral("*(|(")
				p.AddName(0)
				p.AddLiteral("=*))")
			},
		},
		{
			Form: "password:*/=1--",
			Regex: func(p *Parts) {
				p.AddLiteral("password:*/=")
				p.AddNumber(0, 1337)
				p.AddLiteral("--")
			},
		},
		{
			Form: "distinct",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("distinct")
				p.AddSpacesOptional()
			},
		},
		{
			Form: ");waitfor delay '0:0:__TIME__'--",
			Case: ");waitfor delay '0:0:123'--",
			Regex: func(p *Parts) {
				p.AddLiteral(")")
				p.AddWaitfor()
			},
		},
		{
			Form: "to_timestamp_tz",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("to_timestamp_tz")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "\") or benchmark(10000000,MD5(1))#",
			Regex: func(p *Parts) {
				p.AddLiteral("\")")
				p.AddOr()
				p.AddBenchmark()
			},
		},
		{
			Form: " UNION SELECT",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("union")
				p.AddSpaces()
				p.AddLiteral("select")
			},
		},
		{
			Form: "%2A%28%7C%28mail%3D%2A%29%29",
			Regex: func(p *Parts) {
				p.AddHexSpacesOptional()
				p.AddLiteral("%2A%28%7C%28")
				p.AddName(0)
				p.AddLiteral("%3D%2A%29%29")
				p.AddHexSpacesOptional()
			},
		},
		{
			Form: "+sqlvuln",
			Case: "+select a from b where 1=1",
			Regex: func(p *Parts) {
				p.AddLiteral("+")
				p.AddSQL()
			},
		},
		{
			Form: " or 1=1 /*",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("/*")
			},
		},
		{
			Form: ")) or sleep(__TIME__)='",
			Case: ")) or sleep(123)='",
			Regex: func(p *Parts) {
				p.AddLiteral("))")
				p.AddOr()
				p.AddLiteral("sleep(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")='")
			},
		},
		{
			Form: "or 1=1 or \"\"=",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddOr()
				p.AddLiteral("\"\"=")
			},
		},
		{
			Form: " or 1 in (select @@version)--",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("in")
				p.AddSpaces()
				p.AddLiteral("(select")
				p.AddSpaces()
				p.AddLiteral("@@")
				p.AddName(1)
				p.AddLiteral(")--")
			},
		},
		{
			Form: "sqlvuln;",
			Case: "select a from b where 1=1;",
			Regex: func(p *Parts) {
				p.AddSQL()
				p.AddLiteral(";")
			},
		},
		{
			Form: " union select * from users where login = char ...",
			Case: " union select * from users where login = char 1, 2, 3",
			Regex: func(p *Parts) {
				p.AddSpaces()
				p.AddLiteral("union")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("*")
				p.AddSpaces()
				p.AddLiteral("from")
				p.AddSpaces()
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("where")
				p.AddSpaces()
				p.AddName(1)
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddLiteral("char")
				p.AddSpaces()
				p.AddNumberList(256)
			},
		},
		{
			Form: "x' or 1=1 or 'x'='y",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(2)
				p.AddLiteral("'='")
				p.AddName(2)
			},
		},
		{
			Form: "28 %",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("%")
			},
		},
		{
			Form: "â or 3=3 --",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "@variable",
			Regex: func(p *Parts) {
				p.AddLiteral("@")
				p.AddName(0)
			},
		},
		{
			Form: " or '1'='1'--",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddLiteral("'")
				p.AddNumber(0, 1337)
				p.AddLiteral("'='")
				p.AddNumber(0, 1337)
				p.AddLiteral("'--")
			},
		},
		{
			Form: "\"a\"\" or 1=1--\"",
			Regex: func(p *Parts) {
				p.AddLiteral("\"")
				p.AddName(0)
				p.AddLiteral("\"\"")
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddLiteral("--\"")
			},
		},
		{
			Form: "//*",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("//*")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "%2A%7C",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("%2A%7C")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "\" or 0=0 --",
			Regex: func(p *Parts) {
				p.AddLiteral("\"")
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "\")) or pg_sleep(__TIME__)--",
			Case: "\")) or pg_sleep(123)--",
			Regex: func(p *Parts) {
				p.AddLiteral("\"))")
				p.AddOr()
				p.AddLiteral("pg_sleep(")
				p.AddSpacesOptional()
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
				p.AddLiteral(")--")
			},
		},
		{
			Form: "?",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("?")
				p.AddSpacesOptional()
			},
		},
		{
			Form: " or 1/*",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("/*")
			},
		},
		{
			Form: "!",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("!")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "'",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddSpacesOptional()
			},
		},
		{
			Form: " or a = a",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddName(0)
			},
		},
		{
			Form: "declare @q nvarchar (200) select @q = 0x770061006900740066006F0072002000640065006C00610079002000270030003A0030003A0031003000270000 exec(@q)",
			Regex: func(p *Parts) {
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("nvarchar")
				p.AddSpaces()
				p.AddLiteral("(")
				p.AddNumber(1, 200)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddHex(1337 * 1337)
				p.AddSpaces()
				p.AddLiteral("exec(@")
				p.AddName(0)
				p.AddLiteral(")")
			},
		},
		{
			Form: "declare @s varchar(200) select @s = 0x77616974666F722064656C61792027303A303A31302700 exec(@s) ",
			Regex: func(p *Parts) {
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("varchar(")
				p.AddNumber(1, 200)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddHex(1337 * 1337)
				p.AddSpaces()
				p.AddLiteral("exec(@")
				p.AddName(0)
				p.AddLiteral(")")
				p.AddSpaces()
			},
		},
		{
			Form: "declare @q nvarchar (200) 0x730065006c00650063007400200040004000760065007200730069006f006e00 exec(@q)",
			Regex: func(p *Parts) {
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("nvarchar")
				p.AddSpaces()
				p.AddLiteral("(")
				p.AddNumber(1, 200)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddHex(1337 * 1337)
				p.AddSpaces()
				p.AddLiteral("exec(@")
				p.AddName(0)
				p.AddLiteral(")")
			},
		},
		{
			Form: "declare @s varchar (200) select @s = 0x73656c65637420404076657273696f6e exec(@s)",
			Regex: func(p *Parts) {
				p.AddLiteral("declare")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("varchar")
				p.AddSpaces()
				p.AddLiteral("(")
				p.AddNumber(1, 200)
				p.AddLiteral(")")
				p.AddSpaces()
				p.AddLiteral("select")
				p.AddSpaces()
				p.AddLiteral("@")
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("=")
				p.AddSpaces()
				p.AddHex(1337 * 1337)
				p.AddSpaces()
				p.AddLiteral("exec(@")
				p.AddName(0)
				p.AddLiteral(")")
			},
		},
		{
			Form: "' or 1=1",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: " or 1=1 --",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "x' OR full_name LIKE '%Bob%",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddLiteral("'")
				p.AddOr()
				p.AddName(1)
				p.AddSpaces()
				p.AddLiteral("like")
				p.AddSpaces()
				p.AddLiteral("'%")
				p.AddName(2)
				p.AddLiteral("%")
			},
		},
		{
			Form: "'; exec master..xp_cmdshell 'ping 172.10.1.255'--",
			Regex: func(p *Parts) {
				p.AddLiteral("';")
				p.AddSpaces()
				p.AddLiteral("exec")
				p.AddSpaces()
				p.AddLiteral("master..xp_cmdshell")
				p.AddSpaces()
				p.AddLiteral("'ping")
				p.AddSpaces()
				p.AddNumber(0, 256)
				p.AddLiteral(".")
				p.AddNumber(1, 256)
				p.AddLiteral(".")
				p.AddNumber(2, 256)
				p.AddLiteral(".")
				p.AddNumber(3, 256)
				p.AddLiteral("'--")
			},
		},
		{
			Form: "'%20or%20''='",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddHexOr()
				p.AddLiteral("''='")
			},
		},
		{
			Form: "'%20or%20'x'='x",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddHexOr()
				p.AddLiteral("'")
				p.AddName(0)
				p.AddLiteral("'='")
				p.AddName(0)
			},
		},
		{
			Form: "')%20or%20('x'='x",
			Regex: func(p *Parts) {
				p.AddLiteral("')")
				p.AddHexOr()
				p.AddLiteral("('")
				p.AddName(0)
				p.AddLiteral("'='")
				p.AddName(0)
			},
		},
		{
			Form: "' or 0=0 --",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
		{
			Form: "' or 0=0 #",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("#")
			},
		},
		{
			Form: " or 0=0 #\"",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("#\"")
			},
		},
		{
			Form: "' or 1=1--",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddLiteral("--")
			},
		},
		{
			Form: "' or '1'='1'--",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddLiteral("'")
				p.AddNumber(0, 1337)
				p.AddLiteral("'='")
				p.AddNumber(0, 1337)
				p.AddLiteral("'--")
			},
		},
		{
			Form: "' or 1 --'",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddSpaces()
				p.AddLiteral("--'")
			},
		},
		{
			Form: "or 1=1--",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddLiteral("--")
			},
		},
		{
			Form: "' or 1=1 or ''='",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddOr()
				p.AddLiteral("''='")
			},
		},
		{
			Form: " or 1=1 or \"\"=",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
				p.AddOr()
				p.AddLiteral("\"\"=")
			},
		},
		{
			Form: "' or a=a--",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddName(0)
				p.AddLiteral("=")
				p.AddName(0)
				p.AddLiteral("--")
			},
		},
		{
			Form: " or a=a",
			Regex: func(p *Parts) {
				p.AddOr()
				p.AddName(0)
				p.AddLiteral("=")
				p.AddName(0)
			},
		},
		{
			Form: "') or ('a'='a",
			Regex: func(p *Parts) {
				p.AddLiteral("')")
				p.AddOr()
				p.AddLiteral("('")
				p.AddName(0)
				p.AddLiteral("'='")
				p.AddName(0)
			},
		},
		{
			Form: "'hi' or 'x'='x';",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddName(0)
				p.AddLiteral("'")
				p.AddOr()
				p.AddLiteral("'")
				p.AddName(1)
				p.AddLiteral("'='")
				p.AddName(1)
				p.AddLiteral("';")
			},
		},
		{
			Form: "or",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("or")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "procedure",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("procedure")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "handler",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("handler")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "' or username like '%",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddName(0)
				p.AddSpaces()
				p.AddLiteral("like")
				p.AddSpaces()
				p.AddLiteral("'%")
			},
		},
		{
			Form: "' or uname like '%",
		},
		{
			Form: "' or userid like '%",
		},
		{
			Form: "' or uid like '%",
		},
		{
			Form: "' or user like '%",
		},
		{
			Form: "'; exec master..xp_cmdshell",
			Regex: func(p *Parts) {
				p.AddLiteral("';")
				p.AddSpaces()
				p.AddLiteral("exec")
				p.AddSpaces()
				p.AddLiteral("master..xp_cmdshell")
			},
		},
		{
			Form: "'; exec xp_regread",
			Regex: func(p *Parts) {
				p.AddLiteral("';")
				p.AddSpaces()
				p.AddLiteral("exec")
				p.AddSpaces()
				p.AddLiteral("xp_regread")
			},
		},
		{
			Form: "t'exec master..xp_cmdshell 'nslookup www.google.com'--",
			Regex: func(p *Parts) {
				p.AddLiteral("t'exec")
				p.AddSpaces()
				p.AddLiteral("master..xp_cmdshell")
				p.AddSpaces()
				p.AddLiteral("'nslookup")
				p.AddSpaces()
				p.AddName(0)
				p.AddLiteral(".")
				p.AddName(1)
				p.AddLiteral(".")
				p.AddName(2)
				p.AddLiteral("'--")
			},
		},
		{
			Form: "--sp_password",
			Regex: func(p *Parts) {
				p.AddLiteral("--")
				p.AddSpacesOptional()
				p.AddLiteral("sp_password")
				p.AddSpacesOptional()
			},
		},
		{
			Form: "' UNION SELECT",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddSpaces()
				p.AddLiteral("union")
				p.AddSpaces()
				p.AddLiteral("select")
			},
		},
		{
			Form: "' UNION ALL SELECT",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddSpaces()
				p.AddLiteral("union")
				p.AddSpaces()
				p.AddLiteral("all")
				p.AddSpaces()
				p.AddLiteral("select")
			},
		},
		{
			Form: "' or (EXISTS)",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddLiteral("(exists)")
			},
		},
		{
			Form: "' (select top 1",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddSpaces()
				p.AddLiteral("(select")
				p.AddSpaces()
				p.AddLiteral("top")
				p.AddSpaces()
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: "'||UTL_HTTP.REQUEST",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddLiteral("utl_http.request")
			},
		},
		{
			Form: "1;SELECT%20*",
			Regex: func(p *Parts) {
				p.AddNumber(0, 1337)
				p.AddLiteral(";select")
				p.AddHexSpaces()
				p.AddLiteral("*")
			},
		},
		{
			Form: "<>\"'%;)(&+",
		},
		{
			Form: "'%20or%201=1",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddHexOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: "'sqlattempt1",
			Case: "'select a from b where 1=1",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddSQL()
			},
		},
		{
			Form: "%28",
			Regex: func(p *Parts) {
				p.AddSpacesOptional()
				p.AddLiteral("%")
				p.AddNumber(0, 1337)
				p.AddSpacesOptional()
			},
		},
		{
			Form: "%29",
		},
		{
			Form: "%26",
		},
		{
			Form: "%21",
		},
		{
			Form: "' or ''='",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddLiteral("''='")
			},
		},
		{
			Form: "' or 3=3",
			Regex: func(p *Parts) {
				p.AddLiteral("'")
				p.AddOr()
				p.AddNumber(0, 1337)
				p.AddLiteral("=")
				p.AddNumber(0, 1337)
			},
		},
		{
			Form: " or 3=3 --",
			Regex: func(p *Parts) {
				p.AddName(0)
				p.AddOr()
				p.AddNumber(1, 1337)
				p.AddLiteral("=")
				p.AddNumber(1, 1337)
				p.AddSpaces()
				p.AddLiteral("--")
			},
		},
	}
	return generators
}
