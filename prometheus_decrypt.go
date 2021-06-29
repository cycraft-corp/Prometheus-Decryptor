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
  "prometheus_decrypt/csharp_random"
  "prometheus_decrypt/examine"
  "prometheus_decrypt/winsup"
  "fmt"
  "io/ioutil"
  "log"
  "math"
  "path/filepath"
)

type decOption struct {
  inputFile     string
  outputFile    string
  startTick     int
  reversed      bool
  useCurTick    bool
  key           string
  threadCount   int
  format        string
  customSearch  string
  bytesFormat   string
}

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
    go ctrLogger.Printf("\r%d", seed)
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

func prometheusDecrypt(opt decOption){
  if opt.inputFile == "" || opt.outputFile == "" {
    log.Panic("Please provide input file path and output file path")
  }

  if opt.key != "" {      // decrypt file with the key
    file, err := ioutil.ReadFile(opt.inputFile)
    if err != nil {
      log.Panic(err)
    }
    plain := make([]byte, len(file))
    var key_b [32]byte
    copy(key_b[:], []byte(opt.key)[:32])
    salsa20.XORKeyStream(plain, file, []byte{1, 2, 3, 4, 5, 6, 7, 8}, &key_b)
    err = ioutil.WriteFile(opt.outputFile, plain, 0644)
    if err != nil {
      log.Panic(err)
    }
  } else {            // guess key
    if opt.threadCount <= 0 {
      log.Panic("Please provide a positive integer.")
    } else if opt.format == "" && opt.customSearch == "" && opt.bytesFormat == "" {
      log.Panic("Please provide a possible file extension or custom search string.")
    } else if opt.customSearch == "" && opt.bytesFormat == "" && !filetype.IsSupported(opt.format) {
      log.Panic("Unsupported format. Please provide a custom search regular expression with -s.")
    } else if len(opt.bytesFormat) % 2 == 1 {
      log.Panic("Lemgth of bytes format should be a multiple of 2.")
    }

    if opt.startTick < 0 {
      opt.startTick = - opt.startTick
    }
    if opt.startTick > math.MaxInt32 {
      log.Panic("Tick count should between -2147483648 and 2147483648.")
    }

    if opt.useCurTick {
      opt.startTick = winsup.GetTickCount()
    }

    // build examine
    exam := examine.Init(opt.format, opt.customSearch, opt.bytesFormat)

    // Read input file
    file, err := ioutil.ReadFile(opt.inputFile)
    if err != nil {
      log.Panic(err)
    }

    // start worker
    jobs := make(chan int32, opt.threadCount)
    result := make(chan bool, opt.threadCount)
    for i:=0; i<opt.threadCount; i++ {
      go decRoutine(jobs, result, file, opt.outputFile, exam)
    }

    // send job (seed)
    go func(){
      for i:=opt.startTick;; {
        jobs<-int32(i)

        if opt.reversed {
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
    for i:=opt.startTick;; {
      <-result

      if opt.reversed {
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


