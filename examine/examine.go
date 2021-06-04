package examine

import(
  "github.com/h2non/filetype"
  "regexp"
)

type TypeExam struct {
  re    *regexp.Regexp
  ext   string
}

func Init(ext, format string) *TypeExam {
  if format == "" {
    return &TypeExam{re: nil, ext: ext}
  }

  return &TypeExam{re: regexp.MustCompile(format), ext: ext}
}

func (exam *TypeExam) Match(data []byte) bool {
  extFound, reFound := true, true

  if exam.ext != "" {
    extFound = filetype.Is(data, exam.ext)
  }

  if exam.re != nil {
    reFound = exam.re.Match(data)
  }

  return extFound && reFound && exam.matchExt(data)
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


