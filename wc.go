package main

import (
  "bufio"
  "encoding/binary"
  . "github.com/projectpokemon/PCD2WC4/util"
  "io"
  "log"
  "os"
)

var blockPositions = [][]int{
  {0, 0, 0, 0, 0, 0, 1, 1, 2, 3, 2, 3, 1, 1, 2, 3, 2, 3, 1, 1, 2, 3, 2, 3},
  {1, 1, 2, 3, 2, 3, 0, 0, 0, 0, 0, 0, 2, 3, 1, 1, 3, 2, 2, 3, 1, 1, 3, 2},
  {2, 3, 1, 1, 3, 2, 2, 3, 1, 1, 3, 2, 0, 0, 0, 0, 0, 0, 3, 2, 3, 2, 1, 1},
  {3, 2, 3, 2, 1, 1, 3, 2, 3, 2, 1, 1, 3, 2, 3, 2, 1, 1, 0, 0, 0, 0, 0, 0},
}

func ConvertWondercard(f *os.File) []byte {
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

func SaveWondercard(saveData []byte, newPath string) {
  // Save the data
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