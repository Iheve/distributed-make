package config

import (
	"strings"
	"bufio"
	"os"
	"fmt"
	"log"
	"strconv"
)

type Host struct {
	Address string
	Port int64
}

func Parse(filename string) []*Host {
	hosts := make([]*Host, 0)

	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if len(scanner.Text()) == 0 || strings.HasPrefix(scanner.Text(), "#") {
			//Skip empty lines and comments
			continue
		}

		if !strings.Contains(scanner.Text(), ":") {
			// Check if the line as the good syntax
			log.Fatal("Invalid line : can't find separator ':' in line : ", scanner.Text())
			continue
		}

		hosts = append(hosts, readLine(scanner.Text()))
	}

	return hosts
}

func readLine(l string) *Host {
	var err error
	var port int64

	c := strings.SplitN(l, ":", 2)

	host := new(Host)
	host.Address = c[0]
	port, err = strconv.ParseInt(c[1], 0, 0)

	if err != nil {
		log.Fatal(err)
	}

	host.Port = port

	return host
}

func main() {
	// Check if there is an argument
	var path string

	if len(os.Args) != 2 {
		path = "fichier.cfg"
	} else {
		path = os.Args[1]
	}

	hosts := Parse(path)

	// Print the list of hosts
	for i := range hosts {
		fmt.Print(hosts[i].Address)
		fmt.Print(":")
		fmt.Println(hosts[i].Port)
	}
}