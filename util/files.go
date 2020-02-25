package util

import (
  "log"
  "os"
  "path/filepath"
  "regexp"
)

func GetFiles(dir string) <-chan string {
  ch := make(chan string, 1)

  go func(d string, fch chan<- string) {
    walk := func (path string, info os.FileInfo, err error) error {
      if err != nil {
        return err
      }

      fch <- path

      return nil
    }

    _ = filepath.Walk(d, walk)
    close(fch)
  }(dir, ch)

  return ch
}

func GetFileStat(p string, pattern *regexp.Regexp, size int64) (info os.FileInfo, isDir bool) {
  info, err := os.Stat(p)
  if os.IsNotExist(err) {
    log.Printf("File: %s does not exist%s", p, LineBreak)
    return
  }

  if os.IsPermission(err) {
    log.Printf("File: %s cannot be accessed%s", p, LineBreak)
    return
  }

  if !pattern.MatchString(p) {
    return
  }

  if size > -1 && info.Size() != size {
    return
  }

  isDir = info.IsDir()
  return
}