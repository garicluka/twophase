# Twophase 3x3 cube solver

This solver is mostly rewrite of Herbert Koicemba's [twophase algorithm](https://github.com/hkociemba/RubiksCube-TwophaseSolver) in go.

## Info

Solve from one cube state to another.

You can easily create scramble by just solving from solvedState to randomState.

Other than getting and generating tables (for me 12 seconds generating, 0.5 seconds just getting) solving should be relatively fast (20 move solutions under 15ms on average for me).

Solutions are not always optimal, they get better with time spent solving.

## Install

Get latest version like this:

```bash
go get github.com/garicluka/twophase@latest
```

## Import

Import like this:

```go
import "github.com/garicluka/twophase"
```

## Tables

First get/generate tables by passing directory path to where you want to save them. This is slow, especially when generating tables.
Tables are in total almost 22 million elements in length.

```go
tables, err := twophase.GetAndGenerateTables("./twophase_tables")
```

## Cube State

Cube states are represented in "UUUUUUUUURRRRRRRRRFFFFFFFFFDDDDDDDDDLLLLLLLLLBBBBBBBBB" format.
You can get solved cube state, get random cube state, and verify if its valid cube state.

```go
randomState := twophase.GetRandomCubeState()
solvedState := twophase.SolvedCubeSate 
myCubeState = "UUUUUUUUUBFFRRRRRRFRRFFFFFFDDDDDDDDDLLLLLLLLLRBBBBBBBB"
err := twophase.VerifyCubeState(myCubeState)
```

## Solve

This is most important function, it turns fromCubeState to toCubeState with returned solution.

Max Length will stop solver early if solution of that or lower length is found.

Timeout will stop solver after that duration. If no solution is found yet it will return when first one is found.

Solution is slice of moves. Moves are 1 or 2 len strings that start with U, R, F, D, L or B, and optionally end with ' or 2.

 This function verifies both fromCubeState and toCubeState and returns error if they are not possible.

```go
timeout := time.Second / 10
maxLength := 20
fromCubeState := twophase.SolvedCubeState
toCubeState := twophase.GetRandomCubeState()
solution, err := twophase.Solve(fromCubeState, toCubeState, maxLength, timeout, tables)
```
