package ttt

type Player int

const (
	Player1 Player = 1
	Player2 Player = 2
)

// State represents the game state
type State struct {
	Board         [9]int
	CurrentPlayer Player
}

// NewState creates a new game state
func NewState() *State {
	return &State{
		Board:         [9]int{},
		CurrentPlayer: Player1,
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
	for i, v := range s.Board {
		if v == 0 {
			actions = append(actions, i)
		}
	}
	return actions
}

func switchPlayer(player Player) Player {
	return 3 - player
}

// ApplyAction applies the given action to the state
func (s *State) ApplyAction(action int) {
	s.Board[action] = int(s.CurrentPlayer)
	s.CurrentPlayer = switchPlayer(s.CurrentPlayer)
}

// IsTerminal checks if the game is over
func (s *State) IsTerminal() bool {
	return s.GetWinner() != 0 || len(s.GetLegalActions()) == 0
}

// GetWinner returns the winner (1 or 2) or 0 if no winner
func (s *State) GetWinner() Player {
	lines := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}

	for _, line := range lines {
		if s.Board[line[0]] != 0 && s.Board[line[0]] == s.Board[line[1]] && s.Board[line[0]] == s.Board[line[2]] {
			return Player(s.Board[line[0]])
		}
	}
	return 0
}

// GetScore calculates the score of the current state
func (s *State) GetScore() int {
	winner := s.GetWinner()
	if winner == s.CurrentPlayer {
		return 1
	} else if winner == switchPlayer(s.CurrentPlayer) {
		return -1
	}
	return 0
}
