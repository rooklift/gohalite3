package main

import (
	"flag"
	"fmt"
	"time"

	ai "./ai"
	hal "./core"
)

const (
	NAME = "SimTest"
	VERSION = "16.b"				// hash is ??
)

func main() {

	config := new(ai.Config)

	flag.BoolVar(&config.Crash, "crash", false, "randomly crash")
	flag.Parse()

	game := hal.NewGame()

	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("%v", p)
			game.Log("Quitting: %v", p)
			game.Log("Last known hash: %s", game.Hash())
			game.StopLog()
			game.StopFlog()
		}
	}()

	// -------------------------------------------------------------------------------

	game.PrePreParse()				// Reads very early data (including PID, needed for log names)
	true_pid := game.Pid()

	// Both of these fail harmlessly if the directory isn't there:
	game.StartLog(fmt.Sprintf("logs/log-%v.txt", true_pid))
	game.StartFlog(fmt.Sprintf("flogs/flog-%v-%v.json", game.Constants.GameSeed, true_pid))

	game.PreParse()					// Reads the map data.
	game.Init()						// Set game to a valid turn 0 state.

	game.LogWithoutTurn("--------------------------------------------------------------------------------")
	game.LogWithoutTurn("%s %s starting up at %s", NAME, VERSION, time.Now().Format("2006-01-02 15:04:05"))

	var overminds []*ai.Overmind

	for pid := 0; pid < game.Players(); pid++ {
		overminds = append(overminds, ai.NewOvermind(game, config, pid))
	}

	// fmt.Printf("%s %s\n", NAME, VERSION)

	for turn := 0; turn < 500; turn++ {

		for _, o := range overminds {
			o.Step(game)
		}

		if turn < 10 {
			game.Log(game.Hash())
		}

		game = game.SimGen()
	}
}
