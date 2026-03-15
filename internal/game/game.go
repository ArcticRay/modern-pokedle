package game

import (
	"github.com/ArcticRay/modern-pokedle/internal/pokemon"
)

func CompareGuess(guess, target pokemon.Pokemon) pokemon.GuessResult {
	result := pokemon.GuessResult{
		Pokemon: guess,
	}

	result.TypeResult = compareTypes(guess.Types, target.Types)

	if guess.Habitat == target.Habitat {
		result.HabitatResult = pokemon.MatchCorrect
	} else {
		result.HabitatResult = pokemon.MatchWrong
	}

	if guess.Color == target.Color {
		result.ColorResult = pokemon.MatchCorrect
	} else {
		result.ColorResult = pokemon.MatchWrong
	}

	result.GenerationResult = compareDirectional(guess.Generation, target.Generation)

	result.HeightResult = compareDirectional(guess.Height, target.Height)

	result.EvolutionResult = compareDirectional(guess.EvolutionStage, target.EvolutionStage)

	return result
}

func compareTypes(guess, target []pokemon.PokemonType) pokemon.MatchResult {
	if len(guess) == 0 && len(target) == 0 {
		return pokemon.MatchCorrect
	}

	targetTypes := make(map[string]bool)
	for _, t := range target {
		targetTypes[t.Name] = true
	}

	matches := 0
	for _, t := range guess {
		if targetTypes[t.Name] {
			matches++
		}
	}

	if matches == len(target) && len(guess) == len(target) {
		return pokemon.MatchCorrect
	}
	if matches > 0 {
		return pokemon.MatchPartial
	}
	return pokemon.MatchWrong
}

func compareDirectional(guess, target int) pokemon.DirectionalResult {
	if guess == target {
		return pokemon.DirectionCorrect
	}
	if guess < target {
		return pokemon.DirectionHigher
	}
	return pokemon.DirectionLower
}
