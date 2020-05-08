# Quantum genetic optimization

This is an implementation of genetic algorithm with quantum mechanics. See [this](https://pointlander.github.io/posts/quantum-genetic-optimization/) post for background information. This implementation allows for a non-unitary matrix, so it is not exactly quantum mechanics that is being simulated. For all numbers composed of two 4 bit prime numbers, the number is always factored. On average it takes 7 population generations to factor a number.

# Running the simulator

```bash
go build
./qgo
```
