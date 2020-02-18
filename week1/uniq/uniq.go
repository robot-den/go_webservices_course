package main

import (
  "os"
  "fmt"
  "bufio"
  "io"
)

func main() {
  err := uniq(os.Stdin, os.Stdout)
  if err != nil {
    panic(err.Error())
  }
}

func uniq(in io.Reader, out io.Writer) error {
  scanner := bufio.NewScanner(in)
  var prev_str string

  for scanner.Scan() {
    text := scanner.Text()

    if text == prev_str {
      continue
    }

    if text < prev_str {
      return fmt.Errorf("input is not sorted!")
    }

    prev_str = text

    fmt.Fprintln(out, text)
  }

  return nil
}
