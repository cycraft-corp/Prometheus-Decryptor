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
  "golang.org/x/crypto/salsa20"
  "github.com/h2non/filetype"
  "thanos_decrypt/csharp_random"
  "thanos_decrypt/examine"
  "thanos_decrypt/winsup"
  "fmt"
  "io/ioutil"
  "log"
  "math"
  "flag"
  "path/filepath"
)


func genKey(seed int32) [32]byte {
  var key [32]byte
  cr := csharp_random.Random(seed)
  for i:=0; i<32; i++ {
    v := cr.Next(33, 127) % 255
    if v == 34 || v == 92 {
      i -= 1
    } else {
      key[i] = byte(v)
    }
  }
  return key
}

func writeFile(data []byte, path string, seed int32, key string) error {
  dir, file := filepath.Split(path)
  writePath := fmt.Sprintf("%s%d_%s", dir, seed, file)
  err := ioutil.WriteFile(writePath, data, 0644)
  log.Printf("\rDecrypt file with seed %d, key: %s, path: %s\n", seed, key, writePath)
  return err
}

func decRoutine(jobs chan int32, result chan bool, file []byte, output string, exam *examine.TypeExam) {
  plain := make([]byte, len(file))
  for seed := range jobs {
    fmt.Printf("\r%d", seed)
    key := genKey(seed)
    salsa20.XORKeyStream(plain, file, []byte{1, 2, 3, 4, 5, 6, 7, 8}, &key)
    if exam.Match(plain) {
      err := writeFile(plain, output, seed, string(key[:]))
      if err != nil {
        log.Println(err)
      }
      result<-true
    } else {
      result<-false
    }
  }
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

  if *inputFile == "" || *outputFile == "" {
    log.Fatal("Please provide input file path and output file path")
  }

  if *key != "" {      // decrypt file with the key
    file, err := ioutil.ReadFile(*inputFile)
    if err != nil {
      log.Fatal(err)
    }
    plain := make([]byte, len(file))
    var key_b [32]byte
    copy(key_b[:], []byte(*key)[:32])
    salsa20.XORKeyStream(plain, file, []byte{1, 2, 3, 4, 5, 6, 7, 8}, &key_b)
    err = ioutil.WriteFile(*outputFile, plain, 0644)
    if err != nil {
      log.Fatal(err)
    }
  } else {            // guess key
    if *threadCount <= 0 {
      log.Fatal("Please provide a positive integer.")
    } else if *format == "" && *customSearch == "" && *bytesFormat == "" {
      log.Fatal("Please provide a possible file extension or custom search string.")
    } else if *customSearch == "" && *bytesFormat == "" && !filetype.IsSupported(*format) {
      log.Fatal("Unsupported format. Please provide a custom search regular expression with -s.")
    } else if len(*bytesFormat) % 2 == 1 {
      log.Fatal("Lemgth of bytes format should be a multiple of 2.")
    }

    if *startTick < 0 {
      *startTick = - *startTick
    }
    if *startTick > math.MaxInt32 {
      log.Fatal("Tick count should between -2147483648 and 2147483648.")
    }

    if *useCurTick {
      *startTick = winsup.GetTickCount()
    }

    // build examine
    exam := examine.Init(*format, *customSearch, *bytesFormat)

    // Read input file
    file, err := ioutil.ReadFile(*inputFile)
    if err != nil {
      log.Fatal(err)
    }

    // start worker
    jobs := make(chan int32, *threadCount)
    result := make(chan bool, *threadCount)
    for i:=0; i<*threadCount; i++ {
      go decRoutine(jobs, result, file, *outputFile, exam)
    }

    // send job (seed)
    go func(){
      for i:=*startTick;; {
        jobs<-int32(i)

        if *reversed {
          i--
          if i < 0 {
            break
          }
        } else {
          i++
          if i > math.MaxInt32 {
            break
          }
        }
      }
      close(jobs)
    }()

    // wait for job
    for i:=*startTick;; {
      <-result

      if *reversed {
        i--
        if i < 0 {
          break
        }
      } else {
        i++
        if i > math.MaxInt32 {
          break
        }
      }
    }
  }
}
