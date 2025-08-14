# Gator
## A [boot.dev](https://boot.dev) Follow-Along project

### Requirements
You will need Postgres and Go installed to use this repository.

### Installation
`go install` or `go build` will create an executable binary.

### What is Gator?
A CLI tool to follow and view RSS feeds.

### Available Commands

- login
	- Usage: `login <username>`
	- the username must be registered using the _register_ command

- register
	- Usage: `register <username>`
	- Adds a new user to the users table and automatically logs them as the current user

- users
	- Usage: `users`
	- Displays a list of all registered users and labels the currently logged-in user

- reset
	- Usage: `reset`
	- Clears all users from the users table

- addfeed
	- Usage: `addfeed <feed url>`
	- Adds an RSS feed to the feeds table

- feeds
	- Usage: `feeds`
	- Displays all feeds currently in the feeds table

- follow
	- Usage: `follow <feed url>`
	- Adds the current user id and the specified feed to the feed\_follows table

- unfollow
	- Usage: `unfollow <feed url>`
	- Removes the user id and feed id combination from the feed\_follows table

- agg
	- Usage: `agg <\time interval string>`
	- This command will run in a continous loop, updating the contents of the user's followed feeds
	- Use CTRL + C to end this command

- browse
	- Usage: `browse [optional limit number; default is 2]`
	- Displays the contents of the currently followed feeds.

### Config File
In your home directory, you will need to setup a file called `.gatorconfig.json`.

The important field to specify is db\_url, which is the url to connect to your postgres. It will look something like:

`db_url: postgres://username:password@localhost:5432/gator?sslmode=disable`


### Running the program   
Assuming you compiled the binary, usage is:
`gator <command> <arguments>`

Alternatively, from the root of the project:
`go run . <command> <arguments>`
 
