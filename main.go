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
}

var wg sync.WaitGroup

func (c *cmdArgs) execute(ch <-chan []string) {

	defer wg.Done()

	for args := range ch {
		cmd := exec.Command(args[0], args[1:]...)
		cmdOutput, err := cmd.StdoutPipe()
		cmd.Wait()
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
	}
}

func (c *cmdArgs) pushArg(initial bool, arg string, ch chan []string) {

	if initial == true {
		c.initialLength += 1
		c.initialSize += len(arg)
	}

	c.size += len(arg)
	c.args = append(c.args, arg)
	if initial != true && c.size+len(arg) > MAX_ARG_SIZE {

		argsCopy := make([]string, len(c.args))
		copy(argsCopy, c.args)
		ch <- argsCopy
		c.args = c.args[0:c.initialLength]
		c.size = c.initialSize
	}
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
	flag.Parse()
	args := flag.Args()

	initialArgs := []string{"/bin/echo"}
	if len(args) > 0 {
		initialArgs = args[0:]
	}
	c := cmdArgs{}
	ch := make(chan []string)

	wg.Add(1)
	go c.execute(ch)
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
