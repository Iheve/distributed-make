package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Parse(filename string) (hosts []string) {
	hosts = make([]string, 0)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Cannot open " + filename)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if len(scanner.Text()) == 0 || strings.HasPrefix(scanner.Text(), "#") {
			//Skip empty lines and comments
			continue
		}

		hosts = append(hosts, readLine(scanner.Text()))
	}

	return hosts
}

func readLine(l string) (host string) {
	defaultPort := 4242
	var hostname string
	var port int

	// If there is an empty hostname, we raise error
	c := strings.SplitN(l, ":", 2)

	hostname = c[0]
	if hostname == "" {
		log.Fatal("Invalid hostname : empty hostname")
	}

	// If there is an empty port number, we use the default one (== 4242)
	if !strings.Contains(l, ":") {
		port = defaultPort
	} else {
		port, _ = strconv.Atoi(c[1])
		if port == 0 {
			port = defaultPort
		}
	}

	host = hostname + ":" + strconv.Itoa(port)

	return host
}

func main() {
	// Check if there is an argument
	var path string

	if len(os.Args) != 2 {
		path = "hostfile.cfg"
	} else {
		path = os.Args[1]
	}

	hosts := Parse(path)

	// Print the list of hosts
	for i := range hosts {
		fmt.Println(hosts[i])
	}
}
