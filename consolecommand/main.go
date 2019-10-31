package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//https://hackernoon.com/today-i-learned-making-a-simple-interactive-shell-application-in-golang-aa83adcb266a
func main() {

	var NumOfReplica int
	var appName string

	flag.StringVar(&appName, "appName", "default value", "enter name of the app")
	flag.IntVar(&NumOfReplica, "NumOfReplica", 2, "enter number of replica:")

	flag.Parse()

	fmt.Println(appName)
	fmt.Println(NumOfReplica)

	reader := bufio.NewReader(os.Stdin)
	cmdString, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	cmdString = strings.TrimSuffix(cmdString, "\n")
	cmd := exec.Command(cmdString)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()

	for {
		fmt.Print("$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		cmdString = strings.TrimSuffix(cmdString, "\n")
		cmd := exec.Command(cmdString)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	cmdString = strings.TrimSuffix(cmdString, "\n")
	arrCommandStr := strings.Fields(cmdString)
	cmdArray := exec.Command(arrCommandStr[0], arrCommandStr[1:]...)

	str := "Hello World    Beautiful World"
	arrString := strings.Fields(str)
	fmt.Println(arrString)
}
