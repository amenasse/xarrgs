package main

import "bytes"
import "os/exec"
import "os"
import "fmt"
import "bufio"
import "log"
import "flag"

const MAX_ARG_SIZE = 2048

type cmdArgs struct {
	initialLength int
	initialSize   int
	size          int
	args          []string
}

func (c *cmdArgs) execute() {

	cmd := exec.Command(c.args[0], c.args[1:]...)

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

	cmd.Wait()
}

func (c *cmdArgs) pushArg(initial bool, arg string) {

	if initial == true {
		c.initialLength += 1
		c.initialSize += len(arg)
	}

	if initial != true && c.size+len(arg) > MAX_ARG_SIZE {
		c.execute()
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
	flag.Parse()
	args := flag.Args()


	initialArgs := []string{"/bin/echo"}

	if len(args) > 0 {
		initialArgs[0] = args[0]
	}
	if len(args) > 1 {
		initialArgs = args[1:]
	}

	c := cmdArgs{}

	for _, a := range initialArgs {
		c.pushArg(true, a)
	}

	splitFunc := bufio.ScanWords

	if *nullTerminate == true {
		splitFunc = ScanNullTerminate
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(splitFunc)

	for scanner.Scan() {
		c.pushArg(false, scanner.Text())
	}

	c.execute()

}
