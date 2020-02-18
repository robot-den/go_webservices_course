package main

import (
  "strings"
  "bufio"
  "bytes"
  "testing"
)

var testOkInput = `1
2
2
3`

var testOkOutput = `1
2
3
`

func TestOk(t *testing.T) {
  in := bufio.NewReader(strings.NewReader(testOkInput))
  out := new(bytes.Buffer)
  err := uniq(in, out)
  if err != nil {
    t.Errorf("uniq returns error!")
  }
  if out.String() != testOkOutput {
    t.Errorf("uniq returns incorrect result:\n %v - %v", out.String(), testOkOutput)
  }
}

var testErrorInput = `1
2
1
3`

func TestError(t *testing.T) {
  in := bufio.NewReader(strings.NewReader(testErrorInput))
  out := new(bytes.Buffer)
  err := uniq(in, out)
  if err == nil {
    t.Errorf("uniq must return error!")
  }
}
