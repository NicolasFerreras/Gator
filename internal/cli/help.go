package cli

import (
	"fmt"
)

func handlerHelp(state *State, cmd Command) error {
	helpText := `
Gator is a command-line RSS feed aggregator.

Available commands:
- login: Log in to your account
- register: Create a new account
- reset: Reset your password
- users: List all users
- agg: Aggregate feeds
- addfeed: Add a new feed
- feeds: List all feeds
- follow: Follow a feed
- following: List all followed feeds
- unfollow: Unfollow a feed
- help: Show this help message`
	fmt.Println(helpText)
	return nil
}
