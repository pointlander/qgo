// Copyright 2020 The QGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/cmplx"
	"math/rand"
)

// Mul multiples a matrix by a vector returning a vector
func Mul(a, b []complex64) []complex64 {
	length, stride := len(a), len(b)
	if length%stride != 0 {
		panic("invalid sizes")
	}
	c := make([]complex64, 0, length/stride)
	for i := 0; i < length; i += stride {
		var sum complex64
		for j, bb := range b {
			sum += a[i+j] * bb
		}
		c = append(c, sum)
	}
	return c
}

// Genome is a genome
type Genome struct {
	Fitness int
	Index   int
	Genome  []complex64
}

// NewGenome creates a new genome
func NewGenome(rnd *rand.Rand) Genome {
	length := Size + Size*Size
	genome := make([]complex64, 0, length)
	for i := 0; i < length; i++ {
		genome = append(genome, complex(float32(rnd.NormFloat64()), float32(rnd.NormFloat64())))
	}
	return Genome{Genome: genome}
}

// ComputeFitness computes the fitness of the genome
func (g *Genome) ComputeFitness(n int, rnd *rand.Rand) {
	a, b := g.Genome[Size:], g.Genome[:Size]
	result := Mul(a, b)
	sum, index := 0.0, 0
	for _, value := range result {
		sum += cmplx.Abs(complex128(value))
	}
	accumulation, r := 0.0, rnd.Float64()
	for i, value := range result {
		accumulation += cmplx.Abs(complex128(value))
		if r < accumulation/sum {
			index = i
			break
		}
	}
	x, y := index&0xF, (index>>4)&0xF
	fit := n - x*y
	if y == 1 || x == 1 {
		fit = n
	} else if fit < 0 {
		fit = -fit
	}
	g.Fitness = fit
	g.Index = index
}

// Copy makes a copy of the genome
func (g Genome) Copy() Genome {
	genome := make([]complex64, len(g.Genome))
	copy(genome, g.Genome)
	return Genome{
		Fitness: g.Fitness,
		Index:   g.Index,
		Genome:  genome,
	}
}
