package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/Jangari-nTK/go-vautomation/vsphere"
	"golang.org/x/term"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Hostname: ")
	scanner.Scan()
	vc_hostname := scanner.Text()

	c, err := vsphere.NewClient("https://"+vc_hostname, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Username: ")
	scanner.Scan()
	username := scanner.Text()
	fmt.Printf("Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = c.CreateSession(ctx, username, string(bytePassword))
	if err != nil {
		panic(err)
	}

	tlsInfo, err := c.GetVcenterTls(ctx)
	if err != nil {
		panic(err)
	}

	jsonBytes, err := json.Marshal(tlsInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
}
