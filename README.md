**INSTALLATION**

go get gitlab.com/dlarssonse/go-logger

**EXAMPLE**

```golang
package main

import glog "gitlab.com/dlarssonse/go-logger"

type User struct {
  Username string
}

func init() {
  glog.OutputConsole = true
  glog.OutputFilename = "mylogfile.log"
  glog.OutputAsReadable = true
}

func main() {
  glog.Warning("This is a warning message.")
  glog.Debug(int16(16))
  glog.Error(User{Username: "Test"})
}
```