package main

import (
	"log"
	"net"
	"strconv"
)

const HOST = "scanme.nmap.org"
const PORT = 80

func main() {
	_, err := net.Dial("tcp", HOST+":"+strconv.Itoa(PORT))
	if err != nil {
		log.Printf("Port %d is closed!", PORT)
		return
	}

	log.Printf("Port %d is open!", PORT)
}
