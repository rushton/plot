# Plot

A fast and flexible command-line plotting tool that creates visualizations from data streams using gnuplot.

## Quick Start

```bash
# Basic line plot
seq 100 | awk '{print $1, $1*2, $1*4}' | plot

# With headers
{echo "x y1 y2" && seq 100} | awk 'NR==1{print}NR>1{print $1, $1*2, $1*4}' | plot --headers

# Date-based time series
{echo "date value" && for i in {1..30}; do echo "2023-01-$i $((RANDOM % 100))"; done} | plot --date
```

## Installation

### From Source
```bash
git clone https://github.com/rushton/plot.git
cd plot
go build -o plot cmd/plot/main.go
sudo mv plot /usr/local/bin/
```

### Requirements
- Go 1.19 or later
- gnuplot (for plotting)

## Usage

```bash
plot [OPTIONS]
```

### Options

| Option | Description |
|--------|-------------|
| `--help`, `-h` | Show help message |
| `--headers` | Input has header row |
| `--date` | X-axis contains dates |
| `--human-numbers` | Make Y-axis numbers human readable (e.g., 2k, 2M, 2B) |

### Input Format

Plot accepts space-delimited or comma-separated data from stdin:

```
# Space-delimited
1 2 4
2 4 8
3 6 12

# With headers
x y1 y2
1 2 4
2 4 8
3 6 12
```
