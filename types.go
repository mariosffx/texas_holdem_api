package main

type Card string
type Cards []Card

type FullHouse struct {
	ThreeOfAKind Cards `json:"threeOfAKind,omitempty"`
	Pair         Cards `json:"pair,omitempty"`
}

type Combinations struct {
	HighCard      Card      `json:"highCard"`
	Pair          Cards     `json:"pair,omitempty"`
	TwoPair       Cards     `json:"twoPair,omitempty"`
	ThreeOfAKind  Cards     `json:"threeOfAKind,omitempty"`
	Straight      Cards     `json:"straight,omitempty"`
	Flush         Cards     `json:"flush,omitempty"`
	FullHouse     FullHouse `json:"fullHouse,omitempty"`
	FourOfAKind   Cards     `json:"fourOfAKind,omitempty"`
	StraightFlush Cards     `json:"straightFlush,omitempty"`
	RoyalFlush    Cards     `json:"royalFlush,omitempty"`
}

type Player struct {
	Hand         Cards        `json:"hand"`
	Combination  Cards        `json:"combination"`
	Combinations Combinations `json:"combinations"`
	// BestCombination string       `json:"bestCombination"`
}

type TexasHoldemResponse struct {
	CommunityCards Cards  `json:"communityCards"`
	Player1        Player `json:"player1"`
	Player2        Player `json:"player2"`
}

// ErrorResponse represents an error message.
type ErrorResponse struct {
	Error string `json:"error"`
}
