// Copyright 2020 The QGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
)

const (
	// Bits is the number of bits for a number
	Bits = 4
	// Size of the genome
	Size = 1 << (2 * Bits)
	// Mask masks of the bits
	Mask = (1 << 4) - 1
	// Population the population size
	Population = 10
)

// Verbose enables verbose mode
var Verbose = flag.Bool("v", false, "verbose mode")

func main() {
	flag.Parse()

	done := make(chan int, 8)
	factor := func(n int) {
		done <- Factor(n)
	}
	primes, count, sum := Primes(), 0, 0
	for _, a := range primes {
		for _, b := range primes {
			go factor(a * b)
			count++
		}
	}
	for i := 0; i < count; i++ {
		sum += <-done
		if *Verbose {
			fmt.Println("done")
		}
	}
	fmt.Println(float64(sum) / float64(count))
}

// Primes returns the list of primes
func Primes() []int {
	primes := []int{2}
Search:
	for i := 3; i < (1 << Bits); i++ {
		for _, prime := range primes {
			if (i % prime) == 0 {
				continue Search
			}
		}
		primes = append(primes, i)
	}
	return primes
}

// Factor factors a number
func Factor(n int) int {
	rnd := rand.New(rand.NewSource(2))
	if *Verbose {
		fmt.Println("creating population...")
	}
	genomes := make([]Genome, 0, Population)
	for i := 0; i < Population; i++ {
		if *Verbose {
			fmt.Println("new", i)
		}
		genomes = append(genomes, NewGenome(rnd))
	}
	if *Verbose {
		fmt.Println("searching...")
	}
	generation := 0
	for {
		if *Verbose {
			fmt.Println("generation", generation)
		}
		length := len(genomes)
		for i := 0; i < length; i++ {
			// breed qubits
			if .1 > rnd.Float64() {
				a := genomes[rnd.Intn(length/2)].Copy()
				b := genomes[rnd.Intn(length/2)].Copy()
				x, y := rnd.Intn(Size), rnd.Intn(Size)
				a.Genome[x], b.Genome[y] = b.Genome[y], a.Genome[x]
				genomes = append(genomes, a, b)
			}
			// mutate qubits
			if .1 > rnd.Float64() {
				a := genomes[i].Copy()
				x := rnd.Intn(Size)
				a.Genome[x] += complex(float32(rnd.NormFloat64()), float32(rnd.NormFloat64())) / 2
				genomes = append(genomes, a)
			}
			// breed quantum algorithm
			if .1 > rnd.Float64() {
				a := genomes[rnd.Intn(length/2)].Copy()
				b := genomes[rnd.Intn(length/2)].Copy()
				x, y := Size+Size*a.Index+rnd.Intn(Size), Size+Size*b.Index+rnd.Intn(Size)
				a.Genome[x], b.Genome[y] = b.Genome[y], a.Genome[x]
				genomes = append(genomes, a, b)
			}
			// mutate quantum algorithm
			if .1 > rnd.Float64() {
				a := genomes[i].Copy()
				x := rnd.Intn(Size)
				a.Genome[Size+Size*a.Index+x] += complex(float32(rnd.NormFloat64()), float32(rnd.NormFloat64())) / 2
				genomes = append(genomes, a)
			}
		}
		for i := range genomes {
			genomes[i].ComputeFitness(n, rnd)
		}
		sort.Slice(genomes, func(i, j int) bool {
			return genomes[i].Fitness < genomes[j].Fitness
		})
		genomes = genomes[:Population]
		if *Verbose {
			fmt.Println(genomes[0].Fitness)
		}
		if genomes[0].Fitness == 0 {
			break
		}
		generation++
	}
	if *Verbose {
		x, y := genomes[0].GetValues()
		fmt.Println(x, y)
	}
	return generation
}
