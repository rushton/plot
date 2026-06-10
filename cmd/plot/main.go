package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if hasFlag("--help") || hasFlag("-h") {
		fmt.Println(`Usage: plot [--headers] [--date] [--human-numbers] [--title <title>] [--line-color <color>] [--output|-o <file>] [--help|-h]
plot reads input from stdin in space-delimited format and generates a line graph representing the data.

--headers specifies the input has a header row which should be used as the labels
--date specifies the x axis is a date
--human-numbers makes numbers in the y-axis human readable (e.g. 2k, 2M, 2B, etc)
--title <title> sets a title for the plot
--line-color <color> sets the line color for a value column; repeat for each column (e.g. --line-color red --line-color '#00aaff')
--output|-o <file> writes a PNG image to the given file instead of displaying interactively
--help|-h displays this help message`)
		os.Exit(0)
	}
	numParts, fname, err := readInput(bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}
	if numParts < 2 {
		log.Fatalf("Can't plot with fewer than 2 columns.")
	}

	fp, err := os.Create("/tmp/tmp.plot")
	if err != nil {
		log.Fatalf("failed to create script file: %v", err)
	}

	outputFile := flagValue("--output", "-o")

	if hasFlag("--date") {
		setTimeColumn(fp)
	}
	if hasFlag("--human-numbers") {
		approximateNumberFormat(fp)
	}
	if outputFile != "" {
		setPNGOutput(fp, outputFile)
	}
	if title := flagValue("--title"); title != "" {
		setTitle(fp, title)
	}
	plot(fp, numParts, hasFlag("--headers"), fname, flagValues("--line-color"))
	if outputFile == "" {
		keyPressReload(fp)
	}

	cmd := exec.Command("gnuplot", "-c", "/tmp/tmp.plot")
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd failed: %v", err)
	}

	if outputFile != "" {
		exec.Command("open", outputFile).Run()
	}
}

// hasFlag checks whether a target flag is present in the command arguments.
func hasFlag(target string) bool {
	for _, arg := range os.Args {
		if arg == target {
			return true
		}
	}
	return false
}

// flagValue returns the value following any of the given flag names, or empty string if not found.
func flagValue(flags ...string) string {
	for i, arg := range os.Args {
		for _, flag := range flags {
			if arg == flag && i+1 < len(os.Args) {
				return os.Args[i+1]
			}
		}
	}
	return ""
}

// flagValues returns all values following occurrences of flag in the arguments.
func flagValues(flag string) []string {
	var vals []string
	for i, arg := range os.Args {
		if arg == flag && i+1 < len(os.Args) {
			vals = append(vals, os.Args[i+1])
		}
	}
	return vals
}

func setTitle(w io.Writer, title string) {
	w.Write([]byte(fmt.Sprintf("set title '%s' font 'Helvetica Bold,20'\n", title)))
}

func setPNGOutput(w io.Writer, outputFile string) {
	w.Write([]byte(fmt.Sprintf("set terminal pngcairo noenhanced size 2400,1600 font 'Helvetica,16'\nset output '%s'\n", outputFile)))
}

// writeScript generates a gnuplot script file using the specified numColumns. hasHeaderColumn specifies if a header row is present in the data. dataFileName represents the target data file to use as input. The function returns an filepath to the script and an error if present.
func plot(w io.Writer, numColumns int, hasHeaderColumn bool, dataFileName string, lineColors []string) {
	if hasHeaderColumn {
		w.Write([]byte("set key autotitle columnhead\n"))
	}

	w.Write([]byte(`
plot`))
	for i := 1; i < numColumns; i++ {
		if i > 1 {
			w.Write([]byte(", "))
		}
		color := ""
		if i-1 < len(lineColors) {
			color = fmt.Sprintf(" lc rgb '%s'", lineColors[i-1])
		}
		w.Write([]byte(fmt.Sprintf(" '%s' using 1:%d with lines lw 2%s", dataFileName, i+1, color)))
	}
	// TODO: should catch Write errors
}

func keyPressReload(w io.Writer) {
	w.Write([]byte(`
pause mouse keypress
if (MOUSE_KEY == 27) exit 0
reread`))
}

func setTimeColumn(w io.Writer) {
	w.Write([]byte(`
set timefmt '%Y-%m-%dT%H:%M:%SZ'
set xdata time
set format x '%Y-%m-%d'
`))
}

func approximateNumberFormat(w io.Writer) {
	w.Write([]byte(`
set format y '%.s%c'
`))
}

// readInput reads input from the inp and returns the
// number of columns present, the filepath of the temporary
// data location, and any error if present.
func readInput(inp io.Reader) (int, string, error) {
	lines := bufio.NewReader(inp)
	line, err := lines.ReadString('\n')
	if err == io.EOF {
		return 0, "", nil
	} else if err != nil {
		return 0, "", err
	}
	parts := strings.Fields(line)

	fp, err := os.Create("/tmp/tmp.data")
	if err != nil {
		return 0, "", err
	}

	go func() {
		defer fp.Close()
		fp.WriteString(line)
		for {
			line, err := lines.ReadString('\n')
			if err == io.EOF {
				break
			}
			fp.WriteString(line)
		}
	}()
	return len(parts), "/tmp/tmp.data", nil
}
