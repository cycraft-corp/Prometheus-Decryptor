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


