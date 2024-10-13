package main

import (
	"fmt"
	"github.com/kumackey/patchwork-game-tree/ttt"
	"math"
	"math/rand"
	"time"
)

// Constants
const C = 1.0
const EXPAND_THRESHOLD = 10

// Node represents a node in the MCTS tree
type Node struct {
	state      *ttt.State
	w          float64
	n          float64
	childNodes []*Node
}

// NewNode creates a new MCTS node
func NewNode(state *ttt.State) *Node {
	return &Node{
		state:      state,
		w:          0,
		n:          0,
		childNodes: []*Node{},
	}
}

// Evaluate evaluates the node
func (node *Node) Evaluate() float64 {
	if node.state.IsTerminal() {
		value := 0.5
		if node.state.GetScore() == -1 {
			value = 0
		}
		node.w += value
		node.n++
		return value
	}
	if len(node.childNodes) == 0 {
		stateCopy := node.state.Clone()
		value := playout(stateCopy)
		node.w += value
		node.n++
		if node.n == EXPAND_THRESHOLD {
			node.Expand()
		}
		return value
	} else {
		value := 1.0 - node.NextChildNode().Evaluate()
		node.w += value
		node.n++
		return value
	}
}

// Expand expands the node
func (node *Node) Expand() {
	legalActions := node.state.GetLegalActions()
	node.childNodes = []*Node{}
	for _, action := range legalActions {
		childState := node.state.Clone()
		childState.ApplyAction(action)
		node.childNodes = append(node.childNodes, NewNode(childState))
	}
}

// NextChildNode selects the next child node to evaluate
func (node *Node) NextChildNode() *Node {
	for _, childNode := range node.childNodes {
		if childNode.n == 0 {
			return childNode
		}
	}
	t := 0.0
	for _, childNode := range node.childNodes {
		t += childNode.n
	}
	bestValue := -math.Inf(1)
	var bestChild *Node
	for _, childNode := range node.childNodes {
		wr := 1.0 - childNode.w/childNode.n
		bias := math.Sqrt(2.0 * math.Log(t) / childNode.n)
		ucb1Value := wr + C*bias
		if ucb1Value > bestValue {
			bestValue = ucb1Value
			bestChild = childNode
		}
	}
	return bestChild
}

// playout simulates a random playout from the given state
func playout(state *ttt.State) float64 {
	for !state.IsTerminal() {
		actions := state.GetLegalActions()
		action := actions[rand.Intn(len(actions))]
		state.ApplyAction(action)
	}
	return float64(state.GetScore()+1) / 2
}

// mctsAction determines the best action using MCTS with a specified number of playouts
func mctsAction(state *ttt.State, playoutNumber int) int {
	rootNode := NewNode(state)
	rootNode.Expand()
	for i := 0; i < playoutNumber; i++ {
		rootNode.Evaluate()
	}
	legalActions := state.GetLegalActions()
	bestN := -1.0
	bestAction := -1
	for i, action := range legalActions {
		if rootNode.childNodes[i].n > bestN {
			bestN = rootNode.childNodes[i].n
			bestAction = action
		}
	}
	return bestAction
}

func main() {
	rand.Seed(time.Now().UnixNano())
	state := ttt.NewState()
	for !state.IsTerminal() {
		fmt.Println(state.Board)
		var action int
		if state.CurrentPlayer == ttt.Player1 {
			actions := state.GetLegalActions()
			action = actions[rand.Intn(len(actions))]
		} else {
			action = mctsAction(state, 100000)
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
