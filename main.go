package main

import (
	"fmt"
	"net"
	"sync"
)

const (
	Host              = "scanme.nmap.org"
	MaxPort    uint16 = 65535
	Goroutines        = 5000
)

type portScanResult struct {
	port uint16
	open bool
}

func main() {
	results := checkHost(Host)
	for result := range results {
		if result.open {
			fmt.Printf("Port: %d is open!\n", result.port)
		}
	}
}

func checkHost(host string) <-chan portScanResult {
	var waitGroup sync.WaitGroup
	results := make(chan portScanResult, MaxPort) //Buf channel for results

	waitGroup.Add(int(MaxPort)) //We are going to wait for all to be processed
	go func() {
		waitGroup.Wait() //We wait all the results to be sent
		close(results)   //And only then we close the channel
	}()

	workChan := make(chan uint16, MaxPort) //We dispatch all ports to this chan for processing
	for i := 0; i < Goroutines; i++ {      //We are going to use limited amount of goroutine to reduce memory consumption
		go func() {
			for port := range workChan {
				checkPort(host, port, &waitGroup, results)
			}
		}()
	}

	var port uint16 = 1
	for port != 0 { //overflow goes to zero on 65535
		workChan <- port //We send every possible port to this channel
		port++
	}
	close(workChan) //and we are ready to close the channel
	return results
}

func checkPort(host string, port uint16, waitGroup *sync.WaitGroup, results chan<- portScanResult) { // Function to check the ports
	defer waitGroup.Done() //Anyway we release one counter in waitGroup

	hostPort := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", hostPort)
	var result portScanResult
	if err != nil { //If no errors -> the port is reachable
		result = portScanResult{port, false}
	} else {
		_ = conn.Close() //If port is reachable, we have to close the connection
		result = portScanResult{port, true}
	}

	results <- result
}
