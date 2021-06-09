// +build windows

package winsup

import(
  "syscall"
  "log"
)

func GetTickCount() int {
  kernel32, err := syscall.LoadLibrary("kernel32.dll")
  if err != nil {
    log.Fatal(err)
  }
  defer syscall.FreeLibrary(kernel32)

  kernel32GetTickCount, err := syscall.GetProcAddress(kernel32, "GetTickCount")
  if err != nil {
    log.Fatal(err)
  }

  tick, _, errno := syscall.Syscall(uintptr(kernel32GetTickCount), 0, 0, 0, 0)
  if errno != 0 {
    log.Fatal(err)
  }

  return int(tick)
}
