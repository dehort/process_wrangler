package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	numberOfNodesFlag := flag.Int("number_of_nodes", 1, "number of receptor nodes to start")
	flag.Parse()

	numberOfNodes := *numberOfNodesFlag
	//var nodeProcesses [numberOfNodes]*exec.Cmd

	nodeProcesses := make([]*exec.Cmd, numberOfNodes)

	fmt.Println("numberOfNodes:", numberOfNodes)
	port := 6000
	for i := 0; i < numberOfNodes; i++ {
		fmt.Println("HERE")
		cmd := exec.Command("sh", "RUN_NODE.txt")
		additionalEnv1 := fmt.Sprintf("LISTEN_PORT=%d", port)
		additionalEnv2 := fmt.Sprintf("NODE_ID=node-%d", i)
		fmt.Println(additionalEnv2)
		newEnv := append(os.Environ(), additionalEnv1)
		newEnv = append(newEnv, additionalEnv2)
		newEnv = append(newEnv, "ACCOUNT=0000001")
		cmd.Env = newEnv
		fmt.Println("Starting process...")
		err := cmd.Start()
		fmt.Println("Process started...")
		if err != nil {
			fmt.Printf("cmd.Run() failed with %s\n", err)
		}

		nodeProcesses[i] = cmd

		port++

		time.Sleep(5 * time.Second)
	}

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Waiting for shutdown signal...")
	<-signalChan
	fmt.Println("Shutting down nodes...")
	for i := 0; i < numberOfNodes; i++ {
		nodeProcesses[i].Process.Kill()
		time.Sleep(5 * time.Second)
	}
}
