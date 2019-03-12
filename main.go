package main
import "os/exec"
import "os"
import "fmt"
import "bufio"
import "log"
import "flag"
//import "strings"


func execute (cmdName string, args []string) {

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


func main () {

	var extraArgs  []string
	const MAX_ARG_LENGTH = 2048

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
	wordScanner := bufio.NewScanner(os.Stdin)
	wordScanner.Split(bufio.ScanWords)
	



	current_length := len(cmdName) + len(initialArgs)
	for wordScanner.Scan() {

		word := wordScanner.Text()

		if current_length + len(word) > MAX_ARG_LENGTH {
			execute(cmdName, append(initialArgs, extraArgs...))
			extraArgs = nil
			current_length = len(cmdName) + len(initialArgs)
		}

		current_length += len(word)
		extraArgs = append(extraArgs, word)

	}

	execute(cmdName, append(initialArgs, extraArgs...))


}

