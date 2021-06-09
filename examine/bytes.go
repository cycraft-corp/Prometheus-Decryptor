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


