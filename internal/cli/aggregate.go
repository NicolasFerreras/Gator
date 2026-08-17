package cli

import (
	"fmt"
	"time"
)

func handlerAggregate(state *State, cmd Command) error {
	if len(cmd.args) < 1 {
		return ErrNoArgument
	}

	reqTime, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v", reqTime)

	ticker := time.NewTicker(reqTime)
	//i := 0
	for ; ; <-ticker.C {
		scrapeFeeds(state)
		// i = i + 1
		// fmt.Printf("Fetch number %v", i) simple test para saber si esta haciendo el fetch correctamente

	}
}
