package pokemon

type PokemonType struct {
	Name string
}

type Pokemon struct {
	ID             int
	Name           string
	Types          []PokemonType
	Habitat        string
	Color          string
	EvolutionStage int
	Height         int
	Generation     int
	SpriteURL      string
}

type GuessResult struct {
	Pokemon          Pokemon
	TypeResult       MatchResult
	HabitatResult    MatchResult
	ColorResult      MatchResult
	EvolutionResult  DirectionalResult
	HeightResult     DirectionalResult
	GenerationResult DirectionalResult
}

type MatchResult string

const (
	MatchCorrect MatchResult = "correct"
	MatchPartial MatchResult = "partial"
	MatchWrong   MatchResult = "wrong"
)

type DirectionalResult string

const (
	DirectionCorrect DirectionalResult = "correct"
	DirectionHigher  DirectionalResult = "higher"
	DirectionLower   DirectionalResult = "lower"
)
