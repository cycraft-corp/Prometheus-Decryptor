package main

import(
  "flag"
)

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
