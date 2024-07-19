# Quake Cli

## Introduction

This project is a Quake 3 Arena log parser and game simulator. It parses Quake 3 Arena server logs, extracts game data, and generates a report. It also comes with a very (very) simple animation for each game.

The simulation part was done with aid of claude-opus-3.5.


## Requirements

- Go 1.16 or higher

## Usage

1. Install our quake-cli

`go install ./cmd/quake-cli`

2. Parse log and generate report

`quake-cli -file logs.txt`

3. Visualize animations for each game

`quake-cli -file logs.txt -simulate`