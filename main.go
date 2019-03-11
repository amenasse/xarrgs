package main
import "os/exec"
import "os"
import "fmt"
import "bufio"
import "log"
import "flag"
import "strings"

func main () {
	flag.Parse()
	args := flag.Args()

	cmdName := "/bin/echo"
	var initialArgs []string

	if len(args) > 0 {
		cmdName = args[0]
	}
	if len(args) > 1 {
		initialArgs = args[1:]
	}
	lineScanner := bufio.NewScanner(os.Stdin)
	for lineScanner.Scan() {
		args := append(initialArgs, strings.Fields(lineScanner.Text())...)
		cmd := exec.Command(cmdName, args...)

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
}
