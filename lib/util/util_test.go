/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package util

import (
	"testing"
)

const XML = `<?xml version="1.0"?>

<soap:Envelope
xmlns:soap="http://www.w3.org/2003/05/soap-envelope/"
soap:encodingStyle="http://www.w3.org/2003/05/soap-encoding">

<soap:Body>
  <m:GetPrice xmlns:m="https://www.w3schools.com/prices">
    <m:Item>Apples</m:Item>
  </m:GetPrice>
</soap:Body>

</soap:Envelope>
`

const JSON = `{
 "_body": [
  {
   "_inst": "version=\"1.0\"",
   "_target": "xml",
   "_type": "ProcInst"
  },
  {
   "_body": "\n\n",
   "_type": "CharData"
  },
  {
   "_body": [
    {
     "_body": "\n\n",
     "_type": "CharData"
    },
    {
     "_body": [
      {
       "_body": "\n  ",
       "_type": "CharData"
      },
      {
       "_body": [
        {
         "_body": "\n    ",
         "_type": "CharData"
        },
        {
         "_body": [
          {
           "_body": "Apples",
           "_type": "CharData"
          }
         ],
         "_name": "Item",
         "_space": "m",
         "_type": "Element"
        },
        {
         "_body": "\n  ",
         "_type": "CharData"
        }
       ],
       "_name": "GetPrice",
       "_space": "m",
       "_type": "Element",
       "xmlns___m": "https://www.w3schools.com/prices"
      },
      {
       "_body": "\n",
       "_type": "CharData"
      }
     ],
     "_name": "Body",
     "_space": "soap",
     "_type": "Element"
    },
    {
     "_body": "\n\n",
     "_type": "CharData"
    }
   ],
   "_name": "Envelope",
   "_space": "soap",
   "_type": "Element",
   "soap___encodingStyle": "http://www.w3.org/2003/05/soap-encoding",
   "xmlns___soap": "http://www.w3.org/2003/05/soap-envelope/"
  },
  {
   "_body": "\n",
   "_type": "CharData"
  }
 ]
}`

func TestXMLUnmarshal(t *testing.T) {
	var output map[string]interface{}
	err := XMLUnmarshal([]byte(XML), &output)
	if err != nil {
		t.Fatal(err)
	}
	body, ok := output[XMLKeyBody].([]interface{})
	if !ok {
		t.Fatal("invalid _body element")
	}
	if len(body) != 4 {
		t.Fatal("len of body should be 1 is", len(body))
	}
	element, ok := body[2].(map[string]interface{})
	if !ok {
		t.Fatal("invalid element")
	}
	typ, ok := element[XMLKeyType]
	if !ok {
		t.Fatal("there should be a type")
	}
	ty, ok := typ.(string)
	if !ok || ty != XMLTypeElement {
		t.Fatal("invalid type")
	}
}

func TestXMLMarshal(t *testing.T) {
	var output map[string]interface{}
	err := XMLUnmarshal([]byte(XML), &output)
	if err != nil {
		t.Fatal(err)
	}
	data, err := XMLMarshal(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != len([]byte(XML)) {
		t.Fatal("length of input is not same as the length of output")
	}
}
