package config

import (
	"os"
	"encoding/json"
	"reflect"
	"fmt"
	"errors"
	"time"
	"context"
	"os/signal"
	"database/sql"
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
	Public bool
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

func Middleware(handler func(*State, Command, database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.Username)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
	
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

func ScrapeFeeds(s *State) error {
	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	nt := sql.NullTime{Time:time.Now(), Valid:true}

	params := database.MarkFeedFetchedParams{ID:nextFeed[0].ID,
						LastFetchedAt: nt}

	_, err = s.Db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		return err
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(),nextFeed[0].Url)
	if err != nil {
		return err
	}

	feedContents, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	
		
	fmt.Printf("%v\n", feedContents.Channel.Title)
	

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("Usage: agg <time between requests>")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		sig := <-sigChan
		time.Sleep(2 * time.Second)
		fmt.Printf("Signal: %v caught. Exiting aggregation.\n", sig)
		os.Exit(0)
	}()

	time_between_reqs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		ScrapeFeeds(s)
	}

	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 2 {
		return errors.New("Usage: addfeed <feedname> <url>")
	}

	params := database.AddFeedParams{ID:uuid.New(),
					 CreatedAt:time.Now(),
					 UpdatedAt:time.Now(),
					 Name:cmd.Arguments[0],
				    	 Url:cmd.Arguments[1],
				    	 UserID:user.ID}
	
	_, err := s.Db.AddFeed(context.Background(), params)
	if err != nil {
		return err
	}

	followParams := database.CreateFeedFollowParams{ID: uuid.New(),
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
							UserID: user.ID,
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

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("Usage: follow <feed url>")
	}

	currFeed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Arguments[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{ID:uuid.New(),
					      CreatedAt:time.Now(),
					      UpdatedAt:time.Now(),
					      UserID:user.ID,
					      FeedID:currFeed.ID}
	
	followReturn, err := s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("%v %v\n", s.Cfg.Username, followReturn)

	return nil

}

func HandlerFollowing(s *State, cmd Command, user database.User) error {


	feedNames, err := s.Db.GetFeedNameById(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", feedNames)

	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("Usage: unfollow <feed url>")
	}

	fid, err := s.Db.GetFeedByUrl(context.Background(), cmd.Arguments[0])
	if err != nil {
		return err
	}

	toRemove := database.RemoveFeedParams{FeedID: fid.ID, UserID: user.ID}

	_, err = s.Db.RemoveFeed(context.Background(), toRemove)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully removed %v from user %v", cmd.Arguments[0], user.Name)

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
