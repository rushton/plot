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
		fmt.Println(`Usage: plot [--headers] [--help|-h]
plot reads input from stdin in space-delimited format and generates a line graph representing the data.

--headers specifies the input has a header row which should be used as the labels
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

	plot(fp, numParts, hasFlag("--headers"), fname)
	keyPressReload(fp)

	cmd := exec.Command("gnuplot", "-c", "/tmp/tmp.plot")
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd failed: %v", err)
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

// writeScript generates a gnuplot script file using the specified numColumns. hasHeaderColumn specifies if a header row is present in the data. dataFileName represents the target data file to use as input. The function returns an filepath to the script and an error if present.
func plot(w io.Writer, numColumns int, hasHeaderColumn bool, dataFileName string) {
	if hasHeaderColumn {
		w.Write([]byte("set key autotitle columnhead\n"))
	}

	w.Write([]byte("plot"))
	for i := 1; i < numColumns; i++ {
		if i > 1 {
			w.Write([]byte(", "))
		}
		w.Write([]byte(fmt.Sprintf(" '%s' using 1:%d with lines", dataFileName, i+1)))
	}
	// TODO: should catch Write errors
}

func keyPressReload(w io.Writer) {
	w.Write([]byte(`
pause mouse keypress\n
if (MOUSE_KEY == 27) exit 0\n
reread`))
}

func setTimeColumn(w io.Writer) {
	w.Write([]byte(`set timefmt '%Y-%m-%dT%H:%M:%S'
set format x '%Y-%m-%d'`))
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
