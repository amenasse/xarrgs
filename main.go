package main
import "os/exec"
import "os"
import "fmt"
import "bufio"
import "log"
import "flag"
//import "strings"

const MAX_ARG_SIZE = 2048

type cmdArgs struct {
	initialLength int
	initialSize int
	size int
	args []string
}


func (c *cmdArgs) execute () {

	cmd := exec.Command(c.args[0], c.args[1:]...)

	cmdOutput , err := cmd.StdoutPipe()
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

	if initial != true && c.size + len(arg) > MAX_ARG_SIZE {
		c.execute()
		c.args = c.args[0:c.initialLength]
		c.size = c.initialSize

	}
	c.size += len(arg)
	c.args = append(c.args, arg)
}


func main () {


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

	wordScanner := bufio.NewScanner(os.Stdin)
	wordScanner.Split(bufio.ScanWords)

	for wordScanner.Scan() {
		c.pushArg(false, wordScanner.Text())
	}

	c.execute()


}

