// +build windows

package winsup

import(
  "syscall"
  "log"
  "unsafe"
)

func GetTickCount() int {
  kernel32, err := syscall.LoadLibrary("kernel32.dll")
  if err != nil {
    log.Fatal(err)
  }
  defer syscall.FreeLibrary(kernel32)

  kernel32GetTickCount, _ := syscall.GetProcAddress(kernel32, "GetTickCount")
  if err != nil {
    log.Fatal(err)
  }

  tick, _, err := syscall.Syscall(uintptr(kernel32GetTickCount), 0, 0, 0, 0)
  if err != nil {
    log.Fatal(err)
  }
  return int(*(*uint64)(unsafe.Pointer(tick)))
}
