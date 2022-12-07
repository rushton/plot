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
	numParts, fname, err := readInput(bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}
	if numParts < 2 {
		log.Fatalf("Can't plot with fewer than 2 columns.")
	}
	scriptFileName, err := writeScript(numParts, fname)
	if err != nil {
		log.Fatalf("Erorr writing script file: %v", err)
	}

	cmd := exec.Command("gnuplot", "-c", scriptFileName)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd failed: %v", err)
	}
}

func writeScript(numColumns int, dataFileName string) (string, error) {
	fp, err := os.Create("/tmp/tmp.plot")
	if err != nil {
		return "", err
	}

	fp.WriteString("plot")
	for i := 1; i < numColumns; i++ {
		if i > 1 {
			fp.WriteString(", ")
		}
		fp.WriteString(fmt.Sprintf(" '%s' using 1:%d with lines", dataFileName, i+1))
	}
	fp.WriteString("\n")
	fp.WriteString("pause mouse keypress\n")
	fp.WriteString("if (MOUSE_KEY == 27) exit 0\n")
	fp.WriteString("reread")
	return "/tmp/tmp.plot", err
}
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
