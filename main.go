package main

import "bytes"
import "os/exec"
import "os"
import "fmt"
import "bufio"
import "log"
import "flag"
import "sync"

const MAX_ARG_SIZE = 2048

type cmdArgs struct {
	initialLength int
	initialSize   int
	size          int
	args          []string
	maxArgSize    int
}

var wg sync.WaitGroup

func (c *cmdArgs) execute(ch <-chan []string) {

	defer wg.Done()

	for args := range ch {
		cmd := exec.Command(args[0], args[1:]...)
		cmdOutput, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(cmdOutput)

		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
	}

}

func (c *cmdArgs) pushArg(initial bool, arg string, ch chan []string) {

	if initial == true {
		c.initialLength += 1
		c.initialSize += len(arg)
	}

	// result of using --max-chars is going to be different from GNU xargs
	// . GNU Xargs is written in  C - terminating nulls at the end of each
	// argument string  are counted (see man page)

	if initial != true && c.size+len(arg) >= c.maxArgSize {

		argsCopy := make([]string, len(c.args))
		copy(argsCopy, c.args)
		ch <- argsCopy
		c.args = c.args[0:c.initialLength]
		c.size = c.initialSize
	}

	c.size += len(arg)
	c.args = append(c.args, arg)
}

func ScanNullTerminate(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\x00'); i >= 0 {
		// We have a null-terminated.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func main() {

	nullTerminate := flag.Bool("null", false, "items are seperated by a null not whitespace")
	maxProcs := flag.Int("max-procs", 1, "Maximum number of cores to use")
	maxChars := flag.Int("max-chars", MAX_ARG_SIZE, "Use at most max-chars per command line")
	flag.Parse()
	args := flag.Args()

	initialArgs := []string{"/bin/echo"}
	if len(args) > 0 {
		initialArgs = args[0:]
	}
	c := cmdArgs{}
	c.maxArgSize = *maxChars
	ch := make(chan []string)

	for i := 0; i < *maxProcs; i++ {
		wg.Add(1)
		go c.execute(ch)
	}
	for _, a := range initialArgs {
		c.pushArg(true, a, ch)
	}

	splitFunc := bufio.ScanWords

	if *nullTerminate == true {
		splitFunc = ScanNullTerminate
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(splitFunc)

	for scanner.Scan() {
		c.pushArg(false, scanner.Text(), ch)
	}

	ch <- c.args
	close(ch)
	wg.Wait()
}
