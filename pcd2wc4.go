package main

import (
  "bufio"
  . "github.com/projectpokemon/PCD2WC4/util"
  "log"
  "os"
  "regexp"
  "sync"
)

const version = "1.1.0"

func main() {
  log.Printf("PCD to WC4 by Sabresite v%s%s", version, LineBreak)

  args := os.Args[1:]

  if len(args) == 0 {
    log.Printf("Usage: pcd2wc4 [file1.pcd, <file2.pcd...>]%s", LineBreak)
    end()
  }

  wait := sync.WaitGroup{}

  filech := make(chan string, 100)

  wait.Add(1)
  go func(ch chan string) {
    for p := range ch {
      log.Printf("Processing %s%s", p, LineBreak)

      fi, isDir := GetFileStat(p, regexp.MustCompile("\\.pcd$"), 856)

      if isDir {
        go func(dir string) {
          for file := range GetFiles(dir) {
            ch <- file
          }
        }(p)
        continue
      }

      if fi == nil { continue }

      f, err := os.Open(p)
      if err != nil {
        log.Printf("Error: Unable to open file %s%s", p, LineBreak)
        continue
      }

      saveData := ConvertWondercard(f)
      _ = f.Close()

      if saveData == nil { return }
      SaveWondercard(saveData, p[0:len(p)-4] + ".wc4")
    }

    // Signal
    defer wait.Done()
  }(filech)

  // Loop through all files/dirs provided
  for _, arg := range args {
    filech <- arg
  }
  close(filech)

  // Wait for processing to finish
  wait.Wait()
  end()
}

func end() {
  log.Print("Press any key to exit...")
  reader := bufio.NewReader(os.Stdin)
  _, _, _ = reader.ReadRune()
  os.Exit(0)
}