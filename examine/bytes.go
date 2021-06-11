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

package examine

import(
  "bytes"
  "strconv"
  "log"
)

func matchBytes(str []byte, sub string) bool {
  if len(str) < len(sub) {
    return false
  }

  n := len(sub) / 2
  i0, err := strconv.ParseUint(sub[0:2], 16, 8)
  if err != nil {
    log.Fatal(err)
  }
  c0 := byte(i0)

  for i, t := 0, len(str)-n+1; i < t; i++ {
    if str[i] != c0 {
      o := bytes.IndexByte(str[i+1:t], c0)
      if o < 0 {
        return false
      }
      i += o + 1
    }
    if equalBytes(str[i:i+n], sub) {
      return true
    }
  }

  return false
}

func equalBytes(str []byte, sub string) bool {
  for i:=0; i<len(sub)/2; i++ {
    c := sub[i*2: (i+1)*2]
    if c == "??" {
      continue
    }

    b, err := strconv.ParseUint(c, 16, 8)
    if err != nil {
      log.Fatal(err)
    }

    if byte(b) != str[i] {
      return false
    }
  }

  return true
}


