package main

import (
	"fmt"
	"github.com/kumackey/patchwork-game-tree/ttt"
	"math"
)

// miniMaxScore calculates the score for the minimax algorithm
func miniMaxScore(state *ttt.State, depth int) int {
	if state.IsTerminal() || depth == 0 {
		return state.GetScore()
	}
	legalActions := state.GetLegalActions()
	if len(legalActions) == 0 {
		return state.GetScore()
	}
	bestScore := math.MinInt
	for _, action := range legalActions {
		nextState := state.Clone()
		nextState.ApplyAction(action)
		score := -miniMaxScore(nextState, depth-1)
		if score > bestScore {
			bestScore = score
		}
	}
	return bestScore
}

// miniMaxAction determines the best action using the minimax algorithm with a specified depth
func miniMaxAction(state *ttt.State, depth int) int {
	bestAction := -1
	bestScore := math.MinInt
	for _, action := range state.GetLegalActions() {
		nextState := state.Clone()
		nextState.ApplyAction(action)
		score := -miniMaxScore(nextState, depth)
		if score > bestScore {
			bestAction = action
			bestScore = score
		}
	}
	return bestAction
}

func main() {
	state := ttt.NewState()
	for !state.IsTerminal() {
		fmt.Println(state.Board)
		var action int
		if state.CurrentPlayer == 1 {
			action = miniMaxAction(state, 6)
		} else {
			action = miniMaxAction(state, 3)
		}
		state.ApplyAction(action)
	}
	fmt.Println(state.Board)
	winner := state.GetWinner()
	if winner == 0 {
		fmt.Println("It's a draw!")
	} else {
		fmt.Printf("Player %d wins!\n", winner)
	}
}
