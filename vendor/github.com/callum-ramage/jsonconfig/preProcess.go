package jsonconfig

import (
  "io"
)

//Outputs json with //comments removed.
type JsonCommentStripper struct {
  R io.Reader
  b []byte
  pos int
  end int
  err error
  withinString bool
  previousCharacter byte
  more bool
}

//Creates a new comment stripper that can be used as an intermediate layer between
//a JSON decoder and a json source reader.
func NewJsonCommentStripper(reader io.Reader) *JsonCommentStripper {
  commentStripper := JsonCommentStripper{reader, make([]byte, 10000), 0, 0, nil, false, 0, true}
  commentStripper.fillBuffer()
  return &commentStripper
}

//Refills the internal buffer from the internal reader.
func (j *JsonCommentStripper) fillBuffer() {
  end, err := j.R.Read(j.b)
  j.end = end
  j.pos = 0
  if err != nil {
    j.more = false
    j.err = err
  } else {
    j.more = true
  }
}

//Reads data from the internal reader, removing //comments as it goes.
func (j *JsonCommentStripper) Read(p []byte) (n int, err error) {
	// Track strings and // comments
  // A comment can't occur within a string
  // Nothing can happen after a comment
  if j.pos == j.end && (j.more || j.err == nil) {
    j.fillBuffer()
  }

  previousCharacter := j.previousCharacter
  commentFound := false
  start := j.pos
  end := j.pos

  for i := j.pos; i <= j.end && cap(p) >= (i - start) && !commentFound; i++ {
    end = i
    if i != j.end {
      if j.b[i] == '"' {
        j.withinString = !j.withinString
      }
      if !j.withinString &&  j.b[i] == '/' && previousCharacter == '/' {
        commentFound = true
      }
      previousCharacter = j.b[i]
    }
  }

  j.pos = end
  j.previousCharacter = previousCharacter

  if commentFound {
    end--
  }

  copy(p, j.b[start:end])
  n = end - start

  if !j.more {
    err = io.EOF
  }

  // Advance to the end of the comment in preparation for the next read
  if commentFound {
    for i := end; commentFound; {
      if i == j.end {
        if j.more {
          j.fillBuffer()
          i = 0
        } else {
          err = j.err
          j.pos = j.end
          commentFound = false
        }
      } else {
        if j.b[i] == byte(10) {
          commentFound = false
          j.pos = i
        }
        i++
      }
    }
  }

	return
}
