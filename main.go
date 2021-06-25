// +build !gui

/*
MIT License

Copyright (c) 2021 CyCraft Technology

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import(
  "flag"
  "log"
  "fmt"
  "strings"
  "strconv"
)

var ctrLogger *log.Logger

type ctrWritter struct{
  lastVal   *int
}

func (w ctrWritter) Write(p []byte) (n int, err error) {
  strp := strings.TrimRight(string(p), "\n\r")
  strp = strings.TrimLeft(strp, "\n\r")
  stri, _ := strconv.Atoi(strp)
  if stri / 1000 != *w.lastVal {
    *w.lastVal = stri / 1000
    fmt.Printf("\r%d000...", stri/1000)
  }
  return len(p), nil
}

func main(){
  inputFile := flag.String("i", "", "Input encrypted file.")
  outputFile := flag.String("o", "", "Output decrypted file.")
  startTick := flag.Int("t", 0, "Start tickcount. (default 0)")
  reversed := flag.Bool("r", false, "Reversed tickcount.")
  useCurTick := flag.Bool("c", false, "Use current tickcount. (only support in Windows)")
  key := flag.String("k", "", "Decrypt with this key.")
  threadCount := flag.Int("p", 1, "Use n thread.")
  format := flag.String("e", "", "Search file extension.")
  customSearch := flag.String("s", "", "Custom search with regular expression.")
  bytesFormat := flag.String("b", "", "Custom search with byte value. (i.e. \\xde\\xad\\xbe\\xef -> deadbeef)\nPlease use ?? to match any byte (i.e. de??beef)")
  flag.Parse()

  ctrLogger = log.New(ctrWritter{new(int)}, "", 0)

  defer func(){
    // abandon panic to prevent process exit
    recover()
  }()

  thanosDecrypt(decOption{
    inputFile:    *inputFile,
    outputFile:   *outputFile,
    startTick:    *startTick,
    reversed:     *reversed,
    useCurTick:   *useCurTick,
    key:          *key,
    threadCount:  *threadCount,
    format:       *format,
    customSearch: *customSearch,
    bytesFormat:  *bytesFormat,
  })
}
