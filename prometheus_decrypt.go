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
  "os"
  "strings"
  "sync"
)

type decOption struct {
  inputFile     string
  outputFile    string
  startTick     int
  reversed      bool
  useCurTick    bool
  found         int
  backTime      int
  decSize       int
  key           string
  threadCount   int
  format        string
  customSearch  string
  bytesFormat   string
}

const HEADER_MAX_SIZE = 100     // except document

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
  // create all dir
  err := os.MkdirAll(dir, 0755)
  if err != nil {
    return err
  }
  // write file
  err = ioutil.WriteFile(writePath, data, 0644)
  log.Printf("\rDecrypt file with seed %d, key: %s, path: %s\n", seed, key, writePath)
  return err
}

func decRoutine(jobs chan int32, result chan int32, file []byte, output string, exam *examine.TypeExam, wg *sync.WaitGroup, decSize int) {
  defer wg.Done()
  plain := make([]byte, len(file))
  for{
    if seed, ok := <-jobs; ok {
      ctrLogger.Printf("\r%d", seed)
      key := genKey(seed)
      salsa20.XORKeyStream(plain, file[0:decSize], []byte{1, 2, 3, 4, 5, 6, 7, 8}, &key)
      if exam.Match(plain) {
        // the file header matches -> decrypt the whole file now
        salsa20.XORKeyStream(plain, file, []byte{1, 2, 3, 4, 5, 6, 7, 8}, &key)
        err := writeFile(plain, output, seed, string(key[:]))
        if err != nil {
          log.Println(err)
        }
        result<-seed
      }
    } else {
      break
    }
  }
}

func decWithoutKey(opt decOption, quitCh chan bool) int32 {
  // catch local error
  defer func(){
    // do nothing just return
    recover()
  }()

  log.Println("Start decrypt", opt.inputFile)

  // local check
  if opt.customSearch == "" && opt.bytesFormat == "" && !filetype.IsSupported(opt.format) {
    log.Panic("Unsupported format. Please provide a custom search regular expression with -s.")
  }

  // build examine
  exam := examine.Init(opt.format, opt.customSearch, opt.bytesFormat)

  // Read input file
  file, err := ioutil.ReadFile(opt.inputFile)
  if err != nil {
    log.Panic(err)
  }

  // check decrypt size
  switch opt.decSize {
  case -1:                    // max header size according to github.com/h2non/filetype)
    if opt.format == "docx" || opt.format == "xlsx" || opt.format == "pptx" {
      opt.decSize = len(file)
    } else {
      opt.decSize = HEADER_MAX_SIZE
    }
  case 0:                     // full file size
    opt.decSize = len(file)
  default:
    if opt.decSize < 0 {
      opt.decSize = - opt.decSize
    }
  }
  if opt.decSize > len(file) {
    opt.decSize = len(file)
  }

  // start worker
  var wg sync.WaitGroup
  wg.Add(opt.threadCount)
  jobs := make(chan int32, opt.threadCount)
  result := make(chan int32, opt.threadCount)
  for i:=0; i<opt.threadCount; i++ {
    go decRoutine(jobs, result, file, opt.outputFile, exam, &wg, opt.decSize)
  }

  // send job (seed)
  go func(){
    end := false
    for i:=opt.startTick; !end; {
      select {
      case <-quitCh:
        end = true
      default:
        jobs<-int32(i)

        if opt.reversed {
          i--
          if i < 0 {
            i = math.MaxInt32
          }
        } else {
          i++
          if i > math.MaxInt32 {
            i = 0
          }
        }

        if i == opt.startTick {
          end = true
        }
      }
    }
    close(jobs)
  }()

  // wait for job
  var lastTick int32 = -1
  found := 0
  go func(){
    for true {
      if tick, ok := <-result; ok {
        lastTick = tick
        found++
        if found == opt.found {
          quitCh<-true
        }
      } else {
        break
      }
    }
  }()

  wg.Wait()
  close(result)

  return lastTick
}

func decWithKey(opt decOption){
  // catch local error
  defer func(){
    // do nothing just return
    recover()
  }()

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
}

func prometheusDecrypt(opt decOption, quitCh chan bool){
  // catch global error
  defer func(){
    // do nothing just return
    recover()
  }()

  if opt.inputFile == "" || opt.outputFile == "" {
    log.Panic("Please provide input file path and output file path")
  }

  // check file or dir
  files := make([]struct{i string; o string}, 0)
  if fstatI, err := os.Stat(opt.inputFile); err != nil {
    log.Panic("Open input file failed:", err)
  } else if fstatI.IsDir() {
    if fstatO, err := os.Stat(opt.outputFile); err != nil{
      log.Panic("Open output file failed:", err)
    } else if !fstatO.IsDir() {
      log.Panic("Input path and output path type should be the same.")
    } else {
      // generate all files in dir
      err = filepath.Walk(opt.inputFile, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() {
          path, err := filepath.Rel(opt.inputFile, path)
          if err != nil {
            return err
          }
          files = append(files, struct{i string; o string}{
            filepath.Join(opt.inputFile, path),
            filepath.Join(opt.outputFile, path),
          })
        }
        return nil
      })
      if err != nil {
        log.Panic("Parse path failed:", err)
      }
    }
  } else {
    files = append(files, struct{i string; o string}{opt.inputFile, opt.outputFile})
  }

  if opt.key != "" {      // decrypt file with the key
    for _, file := range files {
      tmpOpt := opt
      opt.inputFile = file.i
      opt.outputFile = file.o
      decWithKey(tmpOpt)
    }
  } else {            // guess key
    // global check
    if opt.threadCount <= 0 {
      log.Panic("Please provide a positive integer.")
    } else if len(opt.bytesFormat) % 2 == 1 {
      log.Panic("Length of bytes format should be a multiple of 2.")
    } else if opt.found <= 0 {
      log.Panic("Candidate found should be greater than 0.")
    }

    // set global startTick
    if opt.startTick < 0 {
      opt.startTick = - opt.startTick
    }
    if opt.startTick > math.MaxInt32 {
      log.Panic("Tick count should between -2147483648 and 2147483648.")
    }

    if opt.useCurTick {
      opt.startTick = winsup.GetTickCount()
    }

    // set backtime to millisecond
    opt.backTime *= 1000*60

    // decrypt each files
    var tick int = opt.startTick
    for _, file := range files {
      tmpOpt := opt
      tmpOpt.inputFile = file.i
      tmpOpt.outputFile = file.o
      // file format
      if tmpOpt.format == "" && tmpOpt.customSearch == "" && tmpOpt.bytesFormat == "" {
        tmpOpt.format = filepath.Ext(strings.Split(file.i, ".PROM")[0])
        if len(tmpOpt.format) != 0 {
          tmpOpt.format = tmpOpt.format[1:]
        }
      }
      tmpOpt.format = strings.ToLower(tmpOpt.format)
      // startTick
      if tmpOpt.reversed {
        if tick+tmpOpt.backTime < tmpOpt.startTick {
          tmpOpt.startTick = tick+tmpOpt.backTime
        }
      } else {
        if tick-tmpOpt.backTime > tmpOpt.startTick {
          tmpOpt.startTick = tick-tmpOpt.backTime
        }
      }

      lastTick := int(decWithoutKey(tmpOpt, quitCh))
      if lastTick != -1 {
        tick = lastTick
      }
    }
  }
}


