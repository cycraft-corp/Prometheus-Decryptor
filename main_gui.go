// +build gui

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
  "log"
  "time"
  "math"
  "strings"
  "strconv"

  "github.com/lxn/walk"
  . "github.com/lxn/walk/declarative"
)

var ctrLogger *log.Logger

type logWritter struct {
  results   *walk.TextEdit
}

func (w logWritter) Write(p []byte) (n int, err error) {
  w.results.AppendText(string(p) + "\r\n")
  return len(p), nil
}

type ctrWritter struct {
  ctrText   *walk.TextLabel
  lastVal   *int
}

func (w ctrWritter) Write(p []byte) (n int, err error) {
  strp := strings.TrimLeft(string(p), "\r")
  strp = strings.TrimRight(strp, "\n")
  intp, err := strconv.Atoi(strp)
  if (intp / 1000) != *w.lastVal {
    *w.lastVal = intp / 1000
    w.ctrText.SetText(strconv.Itoa(intp/1000) + "000...")
  }
  return len(p), nil
}

type mainWindow struct {
  *walk.MainWindow
  inputFile     *walk.LineEdit
  outputFile    *walk.LineEdit
  results       *walk.TextEdit
  counter       *walk.TextLabel
  // tickcount
  useCurTick    *walk.CheckBox
  startTick     *walk.CheckBox
  startTickNum  *walk.NumberEdit
  useRevTick    *walk.CheckBox
  // key
  useKey        *walk.CheckBox
  decKey        *walk.LineEdit
  // thread
  useThread     *walk.CheckBox
  threadCount   *walk.NumberEdit
  // search target
  useExt        *walk.CheckBox
  searchExt     *walk.LineEdit
  useStr        *walk.CheckBox
  searchStr     *walk.LineEdit
  useBytes      *walk.CheckBox
  searchBytes   *walk.LineEdit
  // decrypt
  opt           decOption
}

func (mw *mainWindow) selectInputFile(){
  dlg := &walk.FileDialog{}
  dlg.Title = "Select input file"
  dlg.Filter = "*.*"

  if _, err := dlg.ShowOpen(mw); err != nil {
      log.Println(err)
      return
  }
  mw.inputFile.SetText(dlg.FilePath)
}

func (mw *mainWindow) selectOutputFile(){
  dlg := &walk.FileDialog{}
  dlg.Title = "Select output file"
  dlg.Filter = "*.*"

  if _, err := dlg.ShowSave(mw); err != nil {
      log.Println(err)
      return
  }
  mw.outputFile.SetText(dlg.FilePath)
}

func (mw *mainWindow) selectUseCurTick(){
  if mw.useCurTick.CheckState() == walk.CheckChecked {
    mw.startTick.SetCheckState(walk.CheckUnchecked)
  }
}

func (mw *mainWindow) selectStartTick(){
  if mw.startTick.CheckState() == walk.CheckChecked {
    mw.useCurTick.SetCheckState(walk.CheckUnchecked)
  }
}

func (mw *mainWindow) decrypt(){
  mw.opt.inputFile = mw.inputFile.Text()
  if mw.opt.inputFile == "Input file" {
    mw.opt.inputFile = ""
  }
  mw.opt.outputFile = mw.outputFile.Text()
  if mw.opt.outputFile == "Output file" {
    mw.opt.outputFile = ""
  }
  if mw.startTick.CheckState() == walk.CheckChecked {
    mw.opt.startTick = int(mw.startTickNum.Value())
  }
  if mw.useRevTick.CheckState() == walk.CheckChecked {
    mw.opt.reversed = true
  }
  if mw.useCurTick.CheckState() == walk.CheckChecked {
    mw.opt.useCurTick = true
  }
  if mw.useKey.CheckState() == walk.CheckChecked {
    mw.opt.key = mw.decKey.Text()
  }
  if mw.useThread.CheckState() == walk.CheckChecked {
    mw.opt.threadCount = int(mw.threadCount.Value())
  }
  if mw.useExt.CheckState() == walk.CheckChecked {
    mw.opt.format = mw.searchExt.Text()
  }
  if mw.useStr.CheckState() == walk.CheckChecked {
    mw.opt.customSearch = mw.searchStr.Text()
  }
  if mw.useBytes.CheckState() == walk.CheckChecked {
    mw.opt.bytesFormat = mw.searchBytes.Text()
  }

  go func(){
    defer func(){
      // abandon panic to prevent process exit
      recover()
    }()
    prometheusDecrypt(mw.opt)
  }()
}


func main(){
  mw := &mainWindow{opt: decOption{
    inputFile:    "",
    outputFile:   "",
    startTick:    0,
    reversed:     false,
    useCurTick:   false,
    key:          "",
    threadCount:  1,
    format:       "",
    customSearch: "",
    bytesFormat:  "",
  }}

  // log to results (set after run)
  go func(){
    time.Sleep(3 * time.Second)
    log.SetOutput(logWritter{mw.results})
    ctrLogger = log.New(ctrWritter{mw.counter, new(int)}, "", 0)
  }()

  // mainWindow
  if _, err := (MainWindow{
    AssignTo: &mw.MainWindow,
    Title:    "Prometheus Decrypt",
    MinSize:  Size{600, 400},
    Layout:   VBox{},
    Children: []Widget{
      // input & output
      GroupBox{
        Title: "Select Input / Output File",
        Layout: VBox{},
        Children: []Widget{
          Composite{
            Layout: HBox{},
            Children: []Widget{
              LineEdit{
                Text: "Input file",
                AssignTo: &mw.inputFile,
                ReadOnly: true,
              },
              PushButton{
                Text: "select",
                OnClicked: mw.selectInputFile,
              },
            },
          },
          Composite{
            Layout: HBox{},
            Children: []Widget{
              LineEdit{
                Text: "Output file",
                AssignTo: &mw.outputFile,
                ReadOnly: true,
              },
              PushButton{
                MaxSize: Size{100, 20},
                Text: "select",
                OnClicked: mw.selectOutputFile,
              },
            },
          },
        },
      },
      // Option
      GroupBox{
        Title: "Options",
        Layout: VBox{},
        Children: []Widget{
          Composite{
            Layout: HBox{},
            Children: []Widget{
              GroupBox{
                Title: "Search strategy",
                MinSize: Size{350, 100},
                Layout: VBox{Alignment: AlignHNearVCenter},
                Children: []Widget{
                  CheckBox {
                    AssignTo: &mw.useCurTick,
                    Text: "Use current tickcount",
                    Checked: false,
                    OnCheckedChanged: mw.selectUseCurTick,
                  },
                  Composite{
                    Layout:HBox{Alignment: AlignHNearVCenter, MarginsZero: true},
                    Children: []Widget{
                      CheckBox {
                        AssignTo: &mw.startTick,
                        Text: "Start tickcount (default: 0)",
                        Checked: false,
                        OnCheckedChanged: mw.selectStartTick,
                      },
                      NumberEdit {
                        MinSize: Size{Width: 150},
                        MaxValue: math.MaxInt32,
                        MinValue: 0,
                        AssignTo: &mw.startTickNum,
                      },
                    },
                  },
                  CheckBox {
                    AssignTo: &mw.useRevTick,
                    Text: "Reversed tickcount",
                    Checked: false,
                    OnCheckedChanged: func(){},
                  },
                },
              },
              Composite{
                Layout: VBox{MarginsZero: true, SpacingZero: true},
                Children: []Widget{
                  GroupBox{
                    Title: "Key",
                    Layout: VBox{Alignment: AlignHNearVCenter},
                    Children: []Widget{
                      CheckBox {
                        AssignTo: &mw.useKey,
                        Text: "Key (use this key to decrypt it directly)",
                        Checked: false,
                        OnCheckedChanged: mw.selectStartTick,
                      },
                      LineEdit {
                        MaxSize: Size{Width: 150},
                        AssignTo: &mw.decKey,
                      },
                    },
                  },
                  GroupBox{
                    Title: "Thread",
                    Layout: VBox{Alignment: AlignHNearVCenter},
                    Children: []Widget{
                      CheckBox {
                        AssignTo: &mw.useThread,
                        Text: "Use Thread (please input amount of thread, max: 256)",
                        Checked: false,
                        OnCheckedChanged: mw.selectStartTick,
                      },
                      NumberEdit {
                        MaxSize: Size{Width: 150},
                        MaxValue: 256,
                        MinValue: 1,
                        AssignTo: &mw.threadCount,
                      },
                    },
                  },
                },
              },
            },
          },
          GroupBox{
            Title: "Search Target",
            Layout: VBox{},
            Children: []Widget{
              Composite{
                Layout:HBox{Alignment: AlignHNearVCenter, MarginsZero: true},
                Children: []Widget{
                  CheckBox {
                    AssignTo: &mw.useExt,
                    Text: "Search extension",
                    Checked: true,
                    OnCheckedChanged: func(){},
                  },
                  LineEdit {
                    Alignment: AlignHFarVCenter,
                    MaxSize: Size{Width: 300},
                    AssignTo: &mw.searchExt,
                  },
                },
              },
              Composite{
                Layout:HBox{Alignment: AlignHNearVCenter, MarginsZero: true},
                Children: []Widget{
                  CheckBox {
                    AssignTo: &mw.useStr,
                    Text: "Search string",
                    Checked: false,
                    OnCheckedChanged: func(){},
                  },
                  LineEdit {
                    Alignment: AlignHFarVCenter,
                    MaxSize: Size{Width: 300},
                    AssignTo: &mw.searchStr,
                  },
                },
              },
              Composite{
                Layout:HBox{Alignment: AlignHNearVCenter, MarginsZero: true},
                Children: []Widget{
                  CheckBox {
                    AssignTo: &mw.useBytes,
                    Text: "Search bytes string",
                    Checked: false,
                    OnCheckedChanged: func(){},
                  },
                  LineEdit {
                    Alignment: AlignHFarVCenter,
                    MaxSize: Size{Width: 300},
                    AssignTo: &mw.searchBytes,
                  },
                },
              },
            },
          },
        },
      },
      // Decrypt
      Composite{
        Layout: HBox{},
        Children: []Widget{
          PushButton{
            Text: "Decrypt",
            OnClicked: mw.decrypt,
          },
          TextLabel{
            AssignTo: &mw.counter,
          },
        },
      },
      // result
      TextEdit{
	      AssignTo: &mw.results,
        ReadOnly: true,
        HScroll: true,
        VScroll: true,
        MinSize: Size{Height: 200},
      },
      TextLabel{
        Alignment: AlignHFarVCenter,
        Text: "powered by CyCraft Technology",
      },
    },
  }.Run()); err != nil {
    log.Fatal(err)
  }
}


