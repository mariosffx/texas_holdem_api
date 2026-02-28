package main

type Card string
type Cards []Card

type FullHouse struct {
	ThreeOfAKind Cards `json:"threeOfAKind,omitempty"`
	Pair         Cards `json:"pair,omitempty"`
}

type Combinations struct {
	HighCard      Cards `json:"highCard"`
	Pair          Cards `json:"pair,omitempty"`
	TwoPair       Cards `json:"twoPair,omitempty"`
	ThreeOfAKind  Cards `json:"threeOfAKind,omitempty"`
	Straight      Cards `json:"straight,omitempty"`
	Flush         Cards `json:"flush,omitempty"`
	FullHouse     Cards `json:"fullHouse,omitempty"`
	FourOfAKind   Cards `json:"fourOfAKind,omitempty"`
	StraightFlush Cards `json:"straightFlush,omitempty"`
	RoyalFlush    Cards `json:"royalFlush,omitempty"`
}

type Player struct {
	Hand         Cards        `json:"hand"`
	Combination  Cards        `json:"combination"`
	Combinations Combinations `json:"combinations"`
	// BestCombination string       `json:"bestCombination"`
}

type Result struct {
	Winner                 string `json:"winner"`                 // "Player 1", "Player 2", or "Tie"
	PlayerHandCombination1 Cards  `json:"playerHandCombination1"` // Player 1 Hand combination (e.g., "Straight", "Flush")
	PlayerHandCombination2 Cards  `json:"playerHandCombination2"` // Player 2 Hand combination (e.g., "Straight", "Flush")
	PlayerCombination1     string `json:"playerCombinationName1"` // Player 1 Hand Combination Name
	PlayerCombination2     string `json:"playerCombinationName2"` // Player 2 Hand Combination Name
}

type TexasHoldemResponse struct {
	CommunityCards Cards  `json:"communityCards"`
	Player1        Player `json:"player1"`
	Player2        Player `json:"player2"`
	Result         Result `json:"result"`
}

// ErrorResponse represents an error message.
type ErrorResponse struct {
	Error string `json:"error"`
}
