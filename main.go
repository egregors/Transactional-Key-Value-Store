//nolint:revive // it's ok
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/egregors/Transactional-Key-Value-Store/store"
)

func main() {
	fmt.Println("TKVS :3")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	s := store.NewStore()
	go runStoreAndCLI(s)

	<-stop
	fmt.Println("shutting down...")
}

func runStoreAndCLI(s store.TransactionalKVStorer) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		rawCmd, _ := reader.ReadString('\n')
		rawCmd = strings.TrimSuffix(rawCmd, "\n")
		cmds := strings.Split(rawCmd, " ")

		switch cmds[0] {
		case "SET":
			s.Set(cmds[1], cmds[2])
		case "GET":
			if v, ok := s.Get(cmds[1]); ok {
				fmt.Printf("%s\n", v)
			} else {
				fmt.Printf("key not set\n")
			}
		case "DELETE":
			s.Delete(cmds[1])
		case "COUNT":
			fmt.Printf("%d\n", s.Count(cmds[1]))
		case "BEGIN":
			s.Begin()
		case "COMMIT":
			err := s.Commit()
			if err != nil {
				fmt.Println(err.Error())
			}
		case "ROLLBACK":
			err := s.Rollback()
			if err != nil {
				fmt.Println(err.Error())
			}
		default:
			fmt.Printf(
				"wrong command, expecting: SET, GET, DELETE, COUNT, BEGIN, COMMIT, ROLLBACK; got: %s\n",
				rawCmd,
			)
		}
	}
}
