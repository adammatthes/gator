package config

import (
	"os"
	"encoding/json"
	"reflect"
	"fmt"
	"errors"
	"time"
	"context"
	"github.com/google/uuid"
	"github.com/adammatthes/gator/internal/database"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbUrl string `json:"db_url"`
	Username string `json:"current_user_name"`
}

type State struct {
	Cfg *Config
	Db *database.Queries
}

type Command struct {
	Name string
	Arguments []string
}

type Commands struct {
	CLI map[string]func(*State,Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	err := c.CLI[cmd.Name](s, cmd)
	return err
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CLI[name] = f
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("Number of arguments does not match expected arguments")
	}

	_, err := s.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		return err
	}

	s.Cfg.SetUser(cmd.Arguments[0])
	
	fmt.Printf("Username set to %s\n", cmd.Arguments[0])

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("You must provide a username to register")
	}

	params := database.CreateUserParams{ID:uuid.New(),
					    CreatedAt:time.Now(),
					    UpdatedAt:time.Now(),
					    Name:cmd.Arguments[0]}

	_, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("User %s created\n", cmd.Arguments[0])
	HandlerLogin(s, cmd)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Users Table successfully reset")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	results, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return nil
	}
	
	for n := 0; n < len(results); n++ {
		fmt.Printf("* %v", results[n])
		if results[n] == s.Cfg.Username {
			fmt.Printf(" (current)")
		}
		fmt.Printf("\n")
	}

	return nil
}

func (c *Config) SetUser(newName string) error {
	c.DbUrl = os.Getenv("DATABASE_URL")
	


	c.Username = newName

	write(*c)
	return nil
}

func (c Config) Display() {
	rVal := reflect.ValueOf(c)
	numF := rVal.NumField()
	for n := 0; n < numF; n++ {
		fmt.Printf("%v\n", rVal.Field(n))
	}
}

func write(c Config) error {
	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	HOME, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	err = os.WriteFile(HOME + configFileName, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Read () (Config, error) {
	HOME, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	content, err := os.ReadFile(HOME + configFileName)
	if err != nil {
		return Config{}, err
	}

	result := Config{}

	err = json.Unmarshal(content, &result)
	
	return result, err
}
