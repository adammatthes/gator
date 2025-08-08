package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"os"
	"database/sql"
	"github.com/adammatthes/gator/internal/database"
	"github.com/adammatthes/gator/internal/config"
)

func main() {
	myArgs := os.Args
	if len(myArgs) < 2 {
		fmt.Printf("Usage: %s <command name> ", myArgs[0])
		os.Exit(1)
	}

	myCommand := myArgs[1]
	myArguments := myArgs[2:]

	CMD := config.Command{Name : myCommand, Arguments : myArguments}

	myConfig, err := config.Read()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", myConfig.DbUrl)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	myDb := database.New(db)
	
	myState := config.State{Cfg: &myConfig, Db: myDb}

	myCommands := config.Commands{}
	myCommands.CLI = make(map[string]func(*config.State,config.Command) error)
	myCommands.Register("login", config.HandlerLogin)
	myCommands.Register("register", config.HandlerRegister)


	err = myCommands.Run(&myState, CMD)
	if err != nil {
		fmt.Printf("Command failed: %v ", err)
		os.Exit(1)
	}

	

}
