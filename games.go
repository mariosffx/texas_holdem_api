package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
)

var standard_deck = Cards{
	"H2", "H3", "H4", "H5", "H6", "H7", "H8", "H9", "H10", "HJ", "HQ", "HK", "HA",
	"D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9", "D10", "DJ", "DQ", "DK", "DA",
	"C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9", "C10", "CJ", "CQ", "CK", "CA",
	"S2", "S3", "S4", "S5", "S6", "S7", "S8", "S9", "S10", "SJ", "SQ", "SK", "SA",
}

func deal(deck Cards, n int) (Cards, Cards) {
	// Create an independent copy so sorting doesn't affect the original deck
	hand := make(Cards, n)
	copy(hand, deck[:n])
	leftOverDeck := deck[n:]
	return hand, leftOverDeck
}

// 10. High Card
// No combinations
// Example: AKQ84
func getHighCard(cards Cards) Card {
	return cards[len(cards)-1]
}

// 9. One Pair
// Two of the same card
// Example: AAKQJ
func getPair(cards Cards) Cards {
	return cards
}

// 8. Two Pair
// Two different pairs of cards
// Example: AAQQJ
func getTwoPair(cards Cards) Cards {
	return cards
}

// 7. Three of a Kind
// Three of the same card
// Example: AAAKQ
func getThreeOfAKind(cards Cards) Cards {
	return cards
}

// 6. Straight
// Five cards in a sequence, but not of the same suit
// Examples:
// 56789
// A2345
func getStraight(cards Cards) Cards {
	return cards
}

// 5. Flush
// Five cards of the same suit, but not in a sequence
// Example: HA H10 H7 H4 H2
func getFlush(cards Cards) Cards {
	return cards

}

// 4. Full House
// Three of a kind and a pair
// Example: AAABB
func getFullHouse(cards Cards) FullHouse {
	return FullHouse{
		ThreeOfAKind: getThreeOfAKind(cards),
		Pair:         getPair(cards),
	}
}

// 3. Four of a Kind
// Four of the same card
// Example: AAAAC
func getFourOfAKind(cards Cards) Cards {
	return cards
}

// 2. Straight Flush
// Five cards in a sequence and of the same suit
// Examples:
// H5 H6 H7 H8 H9
// HA H2 H3 H4 H5
func getStraightFlush(cards Cards) Cards {
	return cards
}

// 1. Royal Flush
// Ten, Jack, Queen, King, Ace, all of the same suit
// Example: H10 HJ HQ HK HA
func getRoyalFlush(cards Cards) Cards {
	return cards
}

func convertCardToValue(numStr string) int {
	switch numStr {
	case "J":
		return 11
	case "Q":
		return 12
	case "K":
		return 13
	case "A":
		return 14
	default:
		num, err := strconv.Atoi(numStr)
		if err != nil {
			panic("Invalid card number")
		}
		return num
	}

}

func sortHand(hand Cards) Cards {

	slices.SortFunc(hand, func(a Card, b Card) int {
		suitA := string(a[0])
		numStrA := string(a[1:])
		suitB := string(b[0])
		numStrB := string(b[1:])

		numA := convertCardToValue(numStrA)
		numB := convertCardToValue(numStrB)

		if numA != numB {
			return numA - numB
		}

		if suitA < suitB {
			return -1
		}

		if suitA > suitB {
			return 1
		}

		return 0
	})

	return hand
}

func getCombinations(hand Cards, communityCards Cards) Combinations {
	sortedHand := sortHand(hand)
	allCards := append(hand, communityCards...)
	sortedAllCards := sortHand(allCards)

	return Combinations{
		HighCard:      getHighCard(sortedHand),
		Pair:          getPair(sortedAllCards),
		TwoPair:       getTwoPair(sortedAllCards),
		ThreeOfAKind:  getThreeOfAKind(sortedAllCards),
		Straight:      getStraight(sortedAllCards),
		Flush:         getFlush(sortedAllCards),
		FullHouse:     getFullHouse(sortedAllCards),
		FourOfAKind:   getFourOfAKind(sortedAllCards),
		StraightFlush: getStraightFlush(sortedAllCards),
		RoyalFlush:    getRoyalFlush(sortedAllCards),
	}
}

func startGame() (TexasHoldemResponse, error) {
	var deck = make(Cards, len(standard_deck))
	copy(deck, standard_deck)

	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	player1Hand, deckAfterPlayer1 := deal(deck, 2)
	player2Hand, deckAfterPlayer2 := deal(deckAfterPlayer1, 2)
	communityCards, _ := deal(deckAfterPlayer2, 5)

	// Create independent copies of combinations to avoid shared array issues
	player1Combination := make(Cards, 0, 7)
	player1Combination = append(player1Combination, player1Hand...)
	player1Combination = append(player1Combination, communityCards...)

	player2Combination := make(Cards, 0, 7)
	player2Combination = append(player2Combination, player2Hand...)
	player2Combination = append(player2Combination, communityCards...)

	player1Combinations := getCombinations(player1Hand, communityCards)
	player2Combinations := getCombinations(player2Hand, communityCards)

	player1 := Player{
		Hand:         player1Hand,
		Combination:  player1Combination,
		Combinations: player1Combinations,
	}

	player2 := Player{
		Hand:         player2Hand,
		Combination:  player2Combination,
		Combinations: player2Combinations,
	}

	return TexasHoldemResponse{
		CommunityCards: communityCards,
		Player1:        player1,
		Player2:        player2,
	}, nil
}

func texasHoldem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Only GET method is allowed"})
		return
	}

	newGame, err := startGame()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to start game"})
		return
	}

	json.NewEncoder(w).Encode(newGame)

}
