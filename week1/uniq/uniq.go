package main

import (
  "os"
  "fmt"
  "bufio"
)

func main() {
  scanner := bufio.NewScanner(os.Stdin)
  var prev_str string

  for scanner.Scan() {
    text := scanner.Text()

    if text == prev_str {
      continue
    }

    if text < prev_str {
      panic("input is not sorted!")
    }

    prev_str = text

    fmt.Println(text)
  }
}
