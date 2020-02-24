package main

type RNG interface {
  Next() RNG
  H() uint16
  GetSeed() uint32
}

type PokemonRng struct {
  seed uint32
}

func (g *PokemonRng) Next() RNG {
  g.seed *= 0x41C64E6D
  g.seed += 0x00006073
  return g
}

func (g *PokemonRng) H() uint16 {
  return uint16(g.seed >> 0x10)
}

func (g *PokemonRng) GetSeed() uint32 {
  return g.seed
}

func NewPokemonRng(seed uint32) *PokemonRng {
  return &PokemonRng {
    seed: seed,
  }
}

