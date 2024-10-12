package main

import (
	"fmt"
	"math"
)

// State represents the game state
type State struct {
	board         [9]int
	currentPlayer int
}

// NewState creates a new game state
func NewState() *State {
	return &State{
		board:         [9]int{},
		currentPlayer: 1,
	}
}

// Clone creates a copy of the current state
func (s *State) Clone() *State {
	newState := *s
	return &newState
}

// GetLegalActions returns all legal moves
func (s *State) GetLegalActions() []int {
	actions := []int{}
	for i, v := range s.board {
		if v == 0 {
			actions = append(actions, i)
		}
	}
	return actions
}

// ApplyAction applies the given action to the state
func (s *State) ApplyAction(action int) {
	s.board[action] = s.currentPlayer
	s.currentPlayer = 3 - s.currentPlayer // Switch player (1 -> 2, 2 -> 1)
}

// IsTerminal checks if the game is over
func (s *State) IsTerminal() bool {
	return s.GetWinner() != 0 || len(s.GetLegalActions()) == 0
}

// GetWinner returns the winner (1 or 2) or 0 if no winner
func (s *State) GetWinner() int {
	lines := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}

	for _, line := range lines {
		if s.board[line[0]] != 0 && s.board[line[0]] == s.board[line[1]] && s.board[line[0]] == s.board[line[2]] {
			return s.board[line[0]]
		}
	}
	return 0
}

// Node represents a node in the MCTS tree
type Node struct {
	state    *State
	parent   *Node
	children []*Node
	visits   int
	value    float64
	untried  []int
}

// NewNode creates a new MCTS node
func NewNode(state *State, parent *Node) *Node {
	return &Node{
		state:    state,
		parent:   parent,
		children: []*Node{},
		visits:   0,
		value:    0,
		untried:  state.GetLegalActions(),
	}
}

// UCB1 calculates the UCB1 value for node selection
func UCB1(node *Node, parentVisits int) float64 {
	if node.visits == 0 {
		return math.Inf(1)
	}
	return node.value/float64(node.visits) + math.Sqrt(2*math.Log(float64(parentVisits))/float64(node.visits))
}

// GetScore calculates the score of the current state
func (s *State) GetScore(currentPlayer int) int {
	winner := s.GetWinner()
	if winner == currentPlayer {
		return 1
	} else if winner != 0 {
		return -1
	}
	return 0
}

// miniMaxScore calculates the score for the minimax algorithm
func miniMaxScore(state *State, depth int) int {
	if state.IsTerminal() || depth == 0 {
		return state.GetScore(state.currentPlayer)
	}
	legalActions := state.GetLegalActions()
	if len(legalActions) == 0 {
		return state.GetScore(state.currentPlayer)
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
func miniMaxAction(state *State, depth int) int {
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

func (n *Node) selectChild() *Node {
	bestScore := math.Inf(-1)
	var bestChild *Node
	for _, child := range n.children {
		score := UCB1(child, n.visits)
		if score > bestScore {
			bestScore = score
			bestChild = child
		}
	}
	return bestChild
}

func (n *Node) addChild(state *State, action int) *Node {
	child := NewNode(state, n)
	n.untried = removeInt(n.untried, action)
	n.children = append(n.children, child)
	return child
}

func removeInt(slice []int, value int) []int {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func main() {
	state := NewState()
	for !state.IsTerminal() {
		fmt.Println(state.board)
		var action int
		if state.currentPlayer == 1 {
			action = miniMaxAction(state, 6) // 深さ3でMinMax法を使用
		} else {
			//
			action = miniMaxAction(state, 4)
			//actions := state.GetLegalActions()
			//action = actions[rand.Intn(len(actions))]
		}
		state.ApplyAction(action)
	}
	fmt.Println(state.board)
	winner := state.GetWinner()
	if winner == 0 {
		fmt.Println("It's a draw!")
	} else {
		fmt.Printf("Player %d wins!\n", winner)
	}
}
