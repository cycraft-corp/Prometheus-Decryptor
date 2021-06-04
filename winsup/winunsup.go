// +build !windows

package winsup

import(
  "log"
)

func GetTickCount() int {
  log.Fatal("GetTickCount is only supported in Windows version.")
  return 0
}
