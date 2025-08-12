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
	"github.com/adammatthes/gator/internal/rss"
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

func HandlerAgg(s *State, cmd Command) error {
	result, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", result)

	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Arguments) != 2 {
		return errors.New("Usage: addfeed <feedname> <url>")
	}

	curr := s.Cfg.Username
	userstats, err := s.Db.GetUser(context.Background(), curr)

	params := database.AddFeedParams{ID:uuid.New(),
					 CreatedAt:time.Now(),
					 UpdatedAt:time.Now(),
					 Name:cmd.Arguments[0],
				    	 Url:cmd.Arguments[1],
				    	 UserID:userstats.ID}
	
	_, err = s.Db.AddFeed(context.Background(), params)
	if err != nil {
		return err
	}

	followParams := database.CreateFeedFollowParams{ID: uuid.New(),
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
							UserID: userstats.ID,
							FeedID: params.ID}
	
	_, err = s.Db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return err	
	}

	fmt.Printf("Successfully added new feed: %v\n", params)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	myFeeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range myFeeds {
		fmt.Printf("%v %v %v\n", feed.Name, feed.Url, feed.Username)	
	}

	return nil
}

func HandlerFollow(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("Usage: follow <feed url>")
	}

	curr, err := s.Db.GetUser(context.Background(), s.Cfg.Username)
	if err != nil {
		return err
	}
	currFeed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Arguments[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{ID:uuid.New(),
					      CreatedAt:time.Now(),
					      UpdatedAt:time.Now(),
					      UserID:curr.ID,
					      FeedID:currFeed.ID}
	
	followReturn, err := s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("%v %v\n", s.Cfg.Username, followReturn)

	return nil

}

func HandlerFollowing(s *State, cmd Command) error {

	currUser, err := s.Db.GetUser(context.Background(), s.Cfg.Username)
	if err != nil {
		return err
	}

	feedNames, err := s.Db.GetFeedNameById(context.Background(), currUser.ID)

	fmt.Printf("%v\n", feedNames)

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
