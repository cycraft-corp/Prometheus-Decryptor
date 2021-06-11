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
  "github.com/h2non/filetype"
  "regexp"
)

type TypeExam struct {
  re    *regexp.Regexp
  ext   string
  bytes string
}

func Init(ext, format, bytes string) *TypeExam {
  if format == "" {
    return &TypeExam{re: nil, ext: ext, bytes: bytes}
  }

  return &TypeExam{re: regexp.MustCompile(format), ext: ext, bytes: bytes}
}

func (exam *TypeExam) Match(data []byte) bool {
  extFound, reFound, bytesFound := true, true, true

  if exam.ext != "" {
    extFound = filetype.Is(data, exam.ext)
  }

  if exam.re != nil {
    reFound = exam.re.Match(data)
  }

  if exam.bytes != "" {
    bytesFound = matchBytes(data, exam.bytes)
  }

  return extFound && reFound && bytesFound && exam.matchExt(data)
}

// For examine extension
func (exam *TypeExam) matchExt(data []byte) bool {
  result := true
  allExt := map[string]([](func([]byte)bool)){}

  if exam.ext != "" {
    matcher, ok := allExt[exam.ext]
    if !ok {
      for _, f := range matcher {
        result = result && f(data)
      }
    }
  }

  return result
}


