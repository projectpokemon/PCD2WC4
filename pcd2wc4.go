package main

import (
  "bufio"
  "encoding/binary"
  . "github.com/projectpokemon/PCD2WC4/util"
  "io"
  "log"
  "os"
  "path/filepath"
  "strings"
  "sync"
)

const version = "1.0.0"

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
    for f := range ch {
      log.Printf("Processing file %s%s", f, LineBreak)
      processFile(f, ch)
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

func processFile(p string, ch chan<- string) {
  info, err := os.Stat(p)
  if os.IsNotExist(err) {
    log.Printf("File: %s does not exist%s", p, LineBreak)
    return
  }

  if os.IsPermission(err) {
    log.Printf("File: %s cannot be accessed%s", p, LineBreak)
    return
  }

  if !strings.HasSuffix(p, ".pcd") {
    log.Printf("File: %s is not a PCD file%s", p, LineBreak)
    return
  }

  if info.Size() != 856 {
    log.Printf("File: %s is not the correct size%s", p, LineBreak)
    return
  }

  if info.IsDir() {
    getFiles(info.Name(), ch)
    return
  }

  f, err := os.Open(info.Name())
  if err != nil {
    log.Printf("%v%s", err, LineBreak)
    return
  }

  saveData := convertWondercard(f)
  _ = f.Close()

  if saveData != nil {
    // Save the data
    newPath := p[0:len(p)-4] + ".wc4"
    newFile, err := os.Create(newPath)
    if err != nil {
      log.Printf("File: Error opening new file file %s%s", newPath, LineBreak)
      return
    }

    writer := bufio.NewWriter(newFile)
    if _, err := writer.Write(saveData); err != nil {
      log.Printf("File: Error writing new file %v%s", err, LineBreak)
      return
    }

    if err := writer.Flush(); err != nil {
      log.Printf("File: Error flushing new file %v%s", err, LineBreak)
      return
    }

    _ = newFile.Close()
    log.Printf("Created %s%s", newPath, LineBreak)
  }
}

func convertWondercard(f *os.File) []byte {
  r := bufio.NewReader(f)

  buf := make([]byte, 856)
  _, err := io.ReadFull(r, buf)
  if err != nil {
    log.Printf("Error reading data for %s%s", f.Name(), LineBreak)
    return nil
  }

  pkmBuf := buf[8:244]

  pid := binary.LittleEndian.Uint32(pkmBuf[0:4])
  seed := uint32(binary.LittleEndian.Uint16(pkmBuf[6:8]))

  shiftValue := ((pid & 0x3E000) >> 0xD) % 24

  decryptBuf := make([]byte, 228)

  prng := NewPokemonRng(seed)

  // Decrypt
  for i := 0; i < 228; i += 2 {
    if i == 128 { prng = NewPokemonRng(pid) }
    v := binary.LittleEndian.Uint16(pkmBuf[i+8:i+10])
    binary.LittleEndian.PutUint16(decryptBuf[i:i+2], v ^ prng.Next().H())
  }

  // Unshuffle the blocks
  newBlocks := make([]byte, 228)

  for block := 0; block < 4; block++ {
    pos := 32 * blockPositions[block][shiftValue]
    blockPos := block * 32
    nextBlockPos := (block+1) * 32
    copy(newBlocks[blockPos:nextBlockPos], decryptBuf[pos:32+pos])
  }

  // Copy party stuff
  copy(newBlocks[128:228], decryptBuf[128:228])

  wcBuf := make([]byte, 856)
  copy(wcBuf[0:16], buf[0:16])
  copy(wcBuf[16:244], newBlocks)
  copy(wcBuf[244:], buf[244:])

  return wcBuf
}

func getFiles(dir string, ch chan<- string) {
  walk := func (path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    f := info.Name()

    if strings.HasSuffix(f, ".pcd") {
      ch <- f
    }

    return nil
  }

  err := filepath.Walk(dir, walk)
  if err != nil {
    log.Printf("%v%s", err, LineBreak)
  }
}

func end() {
  log.Print("Press any key to exit...")
  reader := bufio.NewReader(os.Stdin)
  _, _, _ = reader.ReadRune()
  os.Exit(0)
}

var blockPositions = [][]int{
  {0, 0, 0, 0, 0, 0, 1, 1, 2, 3, 2, 3, 1, 1, 2, 3, 2, 3, 1, 1, 2, 3, 2, 3},
  {1, 1, 2, 3, 2, 3, 0, 0, 0, 0, 0, 0, 2, 3, 1, 1, 3, 2, 2, 3, 1, 1, 3, 2},
  {2, 3, 1, 1, 3, 2, 2, 3, 1, 1, 3, 2, 0, 0, 0, 0, 0, 0, 3, 2, 3, 2, 1, 1},
  {3, 2, 3, 2, 1, 1, 3, 2, 3, 2, 1, 1, 3, 2, 3, 2, 1, 1, 0, 0, 0, 0, 0, 0},
}