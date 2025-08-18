# Plot

A fast and flexible command-line plotting tool that creates visualizations from data streams using gnuplot.

## Quick Start

```bash
# Basic line plot
seq 100 | awk '{print $1, $1*2, $1*4}' | plot

# With headers
echo "x,y1,y2" && seq 100 | awk '{print $1, $1*2, $1*4}' | plot --headers

# Date-based time series
echo "date,value" && for i in {1..30}; do echo "2023-01-$i,$((RANDOM % 100))"; done | plot --date
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

# Comma-separated
1,2,4
2,4,8
3,6,12

# With headers
x,y1,y2
1,2,4
2,4,8
3,6,12
```

## Examples

### Basic Line Plot
```bash
# Generate simple line plot
seq 100 | awk '{print $1, $1*2, $1*4}' | plot
```

### Time Series Data
```bash
# Create time series with dates
echo "date,value" && for i in {1..30}; do echo "2023-01-$i,$((RANDOM % 100))"; done | plot --date
```

### Multiple Data Series
```bash
# Compare multiple metrics
echo "x,metric1,metric2,metric3" && seq 100 | awk '{print $1, $1*2, $1*3, $1*1.5}' | plot --headers
```

### System Monitoring
```bash
# Monitor CPU usage over time
while true; do echo "$(date +%H:%M:%S),$(top -l 1 | grep "CPU usage" | awk '{print $3}' | sed 's/%//')"; sleep 5; done | plot --date
```

### Data Analysis
```bash
# Analyze log file data
grep "ERROR" app.log | awk '{print $1, $2}' | plot --date
```

### Export for Reports
```bash
# Save as PNG for presentations
seq 100 | awk '{print $1, $1*2}' | plot --output png > report_plot.png
```

## Advanced Usage

### Working with CSV Files
```bash
# Plot CSV data
cat data.csv | plot --headers

# Filter and plot specific columns
cut -d',' -f1,3,5 data.csv | plot --headers
```

### Real-time Data Visualization
```bash
# Monitor network traffic
while true; do echo "$(date +%H:%M:%S),$(netstat -i | grep en0 | awk '{print $7}')"; sleep 1; done | plot --date
```

### Custom Data Processing
```bash
# Process and visualize data
cat sensor_data.txt | awk '{print $1, $2*1000, $3/100}' | plot
```
