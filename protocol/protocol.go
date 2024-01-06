package protocol

import (
	"crypto/rand"
	"math"
)

func GenerateChallenge() []byte {
	challenge := make([]byte, 32)
	_, err := rand.Read(challenge)
	if err != nil {
		panic(err)
	}

	return challenge
}

func CalculateNodeScore(goodResponses, badResponses int) float64 {
	totalVotes := goodResponses + badResponses
	if totalVotes == 0 {
		return 0.5
	}

	average := float64(goodResponses) / float64(totalVotes)
	score := average - (average-0.5)*math.Pow(2, -math.Log(float64(totalVotes+1)))

	return score
}
