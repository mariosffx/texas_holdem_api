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
// Highest card when no other combination exists
// Returns the 5 highest cards
// Example: AKQ84
func getHighCard(cards Cards) Cards {
	// Return the 5 highest cards
	if len(cards) >= 5 {
		return cards[len(cards)-5:]
	}
	return cards
}

// 9. One Pair
// Two of the same card
// Example: AAKQJ
func getPair(cards Cards) Cards {
	// Count rank occurrences
	rankCounts := make(map[string]int)
	rankIndices := make(map[string][]int)
	for i, card := range cards {
		rank := string(card[1:])
		rankCounts[rank]++
		rankIndices[rank] = append(rankIndices[rank], i)
	}

	// Find the highest pair
	var pairIndices []int
	for i := len(cards) - 1; i >= 0; i-- {
		card := cards[i]
		rank := string(card[1:])
		if rankCounts[rank] >= 2 && len(pairIndices) == 0 {
			pairIndices = rankIndices[rank][:2]
			break
		}
	}

	if len(pairIndices) == 0 {
		return Cards{}
	}

	// Collect pair cards
	result := Cards{cards[pairIndices[0]], cards[pairIndices[1]]}

	// Add 3 highest kickers (cards not in the pair)
	usedIndices := make(map[int]bool)
	for _, idx := range pairIndices {
		usedIndices[idx] = true
	}
	usedIndices[pairIndices[0]] = true
	usedIndices[pairIndices[1]] = true

	kickerCount := 0
	for i := len(cards) - 1; i >= 0 && kickerCount < 3; i-- {
		if !usedIndices[i] {
			result = append(result, cards[i])
			usedIndices[i] = true
			kickerCount++
		}
	}

	return result
}

// 8. Two Pair
// Two different pairs of cards
// Example: AAQQJ
func getTwoPair(cards Cards) Cards {
	// Count rank occurrences
	rankCounts := make(map[string]int)
	rankIndices := make(map[string][]int)
	for i, card := range cards {
		rank := string(card[1:])
		rankCounts[rank]++
		rankIndices[rank] = append(rankIndices[rank], i)
	}

	// Find all pairs
	var pairs []string
	for rank, count := range rankCounts {
		if count >= 2 {
			pairs = append(pairs, rank)
		}
	}

	if len(pairs) < 2 {
		return Cards{}
	}

	// Sort pairs by their value (highest first)
	slices.SortFunc(pairs, func(a, b string) int {
		return convertCardToValue(b) - convertCardToValue(a)
	})

	// Take the two highest pairs
	result := Cards{}
	for i := 0; i < 2; i++ {
		pair := pairs[i]
		indices := rankIndices[pair][:2]
		result = append(result, cards[indices[0]], cards[indices[1]])
	}

	// Add 1 highest kicker
	usedIndices := make(map[int]bool)
	for _, idx := range rankIndices[pairs[0]][:2] {
		usedIndices[idx] = true
	}
	for _, idx := range rankIndices[pairs[1]][:2] {
		usedIndices[idx] = true
	}

	for i := len(cards) - 1; i >= 0; i-- {
		if !usedIndices[i] {
			result = append(result, cards[i])
			break
		}
	}

	return result
}

// 7. Three of a Kind
// Three of the same card
// Example: AAAKQ
func getThreeOfAKind(cards Cards) Cards {
	// Count rank occurrences
	rankCounts := make(map[string]int)
	rankIndices := make(map[string][]int)
	for i, card := range cards {
		rank := string(card[1:])
		rankCounts[rank]++
		rankIndices[rank] = append(rankIndices[rank], i)
	}

	// Find the highest three of a kind
	var threeIndices []int
	for i := len(cards) - 1; i >= 0; i-- {
		card := cards[i]
		rank := string(card[1:])
		if rankCounts[rank] >= 3 && len(threeIndices) == 0 {
			threeIndices = rankIndices[rank][:3]
			break
		}
	}

	if len(threeIndices) == 0 {
		return Cards{}
	}

	// Collect three of a kind cards
	result := Cards{cards[threeIndices[0]], cards[threeIndices[1]], cards[threeIndices[2]]}

	// Add 2 highest kickers
	usedIndices := make(map[int]bool)
	for _, idx := range threeIndices {
		usedIndices[idx] = true
	}

	kickerCount := 0
	for i := len(cards) - 1; i >= 0 && kickerCount < 2; i-- {
		if !usedIndices[i] {
			result = append(result, cards[i])
			usedIndices[i] = true
			kickerCount++
		}
	}

	return result
}

// 6. Straight
// Five cards in a sequence, but not of the same suit
// Examples:
// 56789
// A2345
func getStraight(cards Cards) Cards {
	// Extract unique ranks
	rankMap := make(map[int]Card)
	for _, card := range cards {
		rank := string(card[1:])
		rankValue := convertCardToValue(rank)
		if _, exists := rankMap[rankValue]; !exists {
			rankMap[rankValue] = card
		}
	}

	// Get sorted unique rank values
	var uniqueRanks []int
	for rank := range rankMap {
		uniqueRanks = append(uniqueRanks, rank)
	}
	slices.Sort(uniqueRanks)

	// Find the highest straight (5 consecutive cards)
	var straightIndices []int
	for i := len(uniqueRanks) - 1; i >= 4; i-- {
		isConsecutive := true
		for j := 0; j < 4; j++ {
			if uniqueRanks[i-j]-uniqueRanks[i-j-1] != 1 {
				isConsecutive = false
				break
			}
		}
		if isConsecutive {
			for j := 0; j < 5; j++ {
				straightIndices = append(straightIndices, i-j)
			}
			break
		}
	}

	// Check for Ace-low straight (A2345)
	if len(straightIndices) == 0 && len(uniqueRanks) >= 5 {
		// Check if we have A, 2, 3, 4, 5
		hasAce := false
		hasTwo := false
		hasThree := false
		hasFour := false
		hasFive := false

		for _, rank := range uniqueRanks {
			if rank == 14 {
				hasAce = true
			} else if rank == 2 {
				hasTwo = true
			} else if rank == 3 {
				hasThree = true
			} else if rank == 4 {
				hasFour = true
			} else if rank == 5 {
				hasFive = true
			}
		}

		if hasAce && hasTwo && hasThree && hasFour && hasFive {
			// Return A2345 (Ace is used as 1)
			return Cards{rankMap[5], rankMap[4], rankMap[3], rankMap[2], rankMap[14]}
		}
	}

	if len(straightIndices) < 5 {
		return Cards{}
	}

	// Build result from straight indices
	result := Cards{}
	for _, idx := range straightIndices {
		result = append(result, rankMap[uniqueRanks[idx]])
	}

	return result
}

// 5. Flush
// Five cards of the same suit, but not in a sequence
// Example: HA H10 H7 H4 H2
func getFlush(cards Cards) Cards {
	// Count suits
	suitCards := make(map[string]Cards)
	for _, card := range cards {
		suit := string(card[0])
		suitCards[suit] = append(suitCards[suit], card)
	}

	// Find flush (5+ cards of same suit)
	for _, suitHand := range suitCards {
		if len(suitHand) >= 5 {
			// Return the 5 highest cards of this suit
			return suitHand[len(suitHand)-5:]
		}
	}

	return Cards{}
}

// 4. Full House
// Three of a kind and a pair
// Example: AAABB
func getFullHouse(cards Cards) Cards {
	// Count rank occurrences
	rankCounts := make(map[string]int)
	rankIndices := make(map[string][]int)
	for i, card := range cards {
		rank := string(card[1:])
		rankCounts[rank]++
		rankIndices[rank] = append(rankIndices[rank], i)
	}

	// Find all three of a kinds and pairs
	var threes []string
	var pairs []string

	for rank, count := range rankCounts {
		if count >= 3 {
			threes = append(threes, rank)
		}
		if count >= 2 {
			pairs = append(pairs, rank)
		}
	}

	if len(threes) == 0 || len(pairs) == 0 {
		return Cards{}
	}

	// If we only have one three of a kind, we need another pair
	if len(threes) == 1 && len(pairs) == 1 {
		return Cards{}
	}

	// Sort threes by value (highest first)
	slices.SortFunc(threes, func(a, b string) int {
		return convertCardToValue(b) - convertCardToValue(a)
	})

	// Sort pairs by value (highest first)
	slices.SortFunc(pairs, func(a, b string) int {
		return convertCardToValue(b) - convertCardToValue(a)
	})

	// Get the highest three of a kind
	threeRank := threes[0]
	var pairRank string

	// Get the highest pair that's not the same as the three of a kind
	for _, p := range pairs {
		if p != threeRank {
			pairRank = p
			break
		}
	}

	if pairRank == "" {
		return Cards{}
	}

	// Build full house: 3 of a kind + pair
	result := Cards{}
	for _, idx := range rankIndices[threeRank][:3] {
		result = append(result, cards[idx])
	}
	for _, idx := range rankIndices[pairRank][:2] {
		result = append(result, cards[idx])
	}

	return result
}

// 3. Four of a Kind
// Four of the same card
// Example: AAAAC
func getFourOfAKind(cards Cards) Cards {
	// Count rank occurrences
	rankCounts := make(map[string]int)
	rankIndices := make(map[string][]int)
	for i, card := range cards {
		rank := string(card[1:])
		rankCounts[rank]++
		rankIndices[rank] = append(rankIndices[rank], i)
	}

	// Find the highest four of a kind
	var fourIndices []int
	for i := len(cards) - 1; i >= 0; i-- {
		card := cards[i]
		rank := string(card[1:])
		if rankCounts[rank] >= 4 && len(fourIndices) == 0 {
			fourIndices = rankIndices[rank][:4]
			break
		}
	}

	if len(fourIndices) == 0 {
		return Cards{}
	}

	// Collect four of a kind cards
	result := Cards{}
	for _, idx := range fourIndices {
		result = append(result, cards[idx])
	}

	// Add 1 highest kicker
	usedIndices := make(map[int]bool)
	for _, idx := range fourIndices {
		usedIndices[idx] = true
	}

	for i := len(cards) - 1; i >= 0; i-- {
		if !usedIndices[i] {
			result = append(result, cards[i])
			break
		}
	}

	return result
}

// 2. Straight Flush
// Five cards in a sequence and of the same suit
// Examples:
// H5 H6 H7 H8 H9
// HA H2 H3 H4 H5
func getStraightFlush(cards Cards) Cards {
	// Group cards by suit
	suitCards := make(map[string]Cards)
	for _, card := range cards {
		suit := string(card[0])
		suitCards[suit] = append(suitCards[suit], card)
	}

	// For each suit, find straights
	for _, suitHand := range suitCards {
		if len(suitHand) < 5 {
			continue
		}

		// Extract unique ranks for this suit
		rankMap := make(map[int]Card)
		for _, card := range suitHand {
			rank := string(card[1:])
			rankValue := convertCardToValue(rank)
			if _, exists := rankMap[rankValue]; !exists {
				rankMap[rankValue] = card
			}
		}

		// Get sorted unique rank values
		var uniqueRanks []int
		for rank := range rankMap {
			uniqueRanks = append(uniqueRanks, rank)
		}
		slices.Sort(uniqueRanks)

		// Find the highest straight flush (5 consecutive cards of same suit)
		for i := len(uniqueRanks) - 1; i >= 4; i-- {
			isConsecutive := true
			for j := 0; j < 4; j++ {
				if uniqueRanks[i-j]-uniqueRanks[i-j-1] != 1 {
					isConsecutive = false
					break
				}
			}
			if isConsecutive {
				result := Cards{}
				for j := 0; j < 5; j++ {
					result = append(result, rankMap[uniqueRanks[i-j]])
				}
				return result
			}
		}

		// Check for Ace-low straight flush (A2345)
		if len(uniqueRanks) >= 5 {
			hasAce := false
			hasTwo := false
			hasThree := false
			hasFour := false
			hasFive := false

			for _, rank := range uniqueRanks {
				if rank == 14 {
					hasAce = true
				} else if rank == 2 {
					hasTwo = true
				} else if rank == 3 {
					hasThree = true
				} else if rank == 4 {
					hasFour = true
				} else if rank == 5 {
					hasFive = true
				}
			}

			if hasAce && hasTwo && hasThree && hasFour && hasFive {
				return Cards{rankMap[5], rankMap[4], rankMap[3], rankMap[2], rankMap[14]}
			}
		}
	}

	return Cards{}
}

// 1. Royal Flush
// Ten, Jack, Queen, King, Ace, all of the same suit
// Example: H10 HJ HQ HK HA
func getRoyalFlush(cards Cards) Cards {
	// Group cards by suit
	suitCards := make(map[string]Cards)
	for _, card := range cards {
		suit := string(card[0])
		suitCards[suit] = append(suitCards[suit], card)
	}

	// For each suit, check for royal flush (10, J, Q, K, A)
	royalRanks := []string{"10", "J", "Q", "K", "A"}

	for _, suitHand := range suitCards {
		if len(suitHand) < 5 {
			continue
		}

		// Check if this suit has all royal cards
		rankSet := make(map[string]Card)
		for _, card := range suitHand {
			rank := string(card[1:])
			rankSet[rank] = card
		}

		isRoyal := true
		var royalFlush Cards
		for _, rank := range royalRanks {
			if card, exists := rankSet[rank]; exists {
				royalFlush = append(royalFlush, card)
			} else {
				isRoyal = false
				break
			}
		}

		if isRoyal && len(royalFlush) == 5 {
			return royalFlush
		}
	}

	return Cards{}
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
	// Combine all cards and sort them
	allCards := append(hand, communityCards...)
	sortedAllCards := sortHand(allCards)

	return Combinations{
		HighCard:      getHighCard(sortedAllCards),
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

// getStrongestCombination finds the best poker hand from all possible combinations
// Returns the combination name and the cards that make up that combination
func getStrongestCombination(combinations Combinations) (string, Cards) {
	if len(combinations.RoyalFlush) > 0 {
		return "Royal Flush", combinations.RoyalFlush
	}
	if len(combinations.StraightFlush) > 0 {
		return "Straight Flush", combinations.StraightFlush
	}
	if len(combinations.FourOfAKind) > 0 {
		return "Four of a Kind", combinations.FourOfAKind
	}
	if len(combinations.FullHouse) > 0 {
		return "Full House", combinations.FullHouse
	}
	if len(combinations.Flush) > 0 {
		return "Flush", combinations.Flush
	}
	if len(combinations.Straight) > 0 {
		return "Straight", combinations.Straight
	}
	if len(combinations.ThreeOfAKind) > 0 {
		return "Three of a Kind", combinations.ThreeOfAKind
	}
	if len(combinations.TwoPair) > 0 {
		return "Two Pair", combinations.TwoPair
	}
	if len(combinations.Pair) > 0 {
		return "Pair", combinations.Pair
	}
	return "High Card", combinations.HighCard
}

// getHandRank returns a numeric rank for a hand name (higher is better)
func getHandRank(handName string) int {
	switch handName {
	case "Royal Flush":
		return 10
	case "Straight Flush":
		return 9
	case "Four of a Kind":
		return 8
	case "Full House":
		return 7
	case "Flush":
		return 6
	case "Straight":
		return 5
	case "Three of a Kind":
		return 4
	case "Two Pair":
		return 3
	case "Pair":
		return 2
	case "High Card":
		return 1
	default:
		return 0
	}
}

// compareHands compares two hands and returns:
// 1 if player1 wins
// 2 if player2 wins
// 0 if it's a tie
func compareHands(hand1Name string, hand1Cards Cards, hand2Name string, hand2Cards Cards) int {
	rank1 := getHandRank(hand1Name)
	rank2 := getHandRank(hand2Name)

	if rank1 > rank2 {
		return 1
	}
	if rank2 > rank1 {
		return 2
	}

	// Same hand rank, compare card values
	for i := 0; i < len(hand1Cards) && i < len(hand2Cards); i++ {
		card1 := hand1Cards[i]
		card2 := hand2Cards[i]

		rank1Value := convertCardToValue(string(card1[1:]))
		rank2Value := convertCardToValue(string(card2[1:]))

		if rank1Value > rank2Value {
			return 1
		}
		if rank2Value > rank1Value {
			return 2
		}
	}

	// Complete tie
	return 0
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

	// Get the strongest combinations for each player
	player1CombinationName, player1StrongestHand := getStrongestCombination(player1Combinations)
	player2CombinationName, player2StrongestHand := getStrongestCombination(player2Combinations)

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

	// Determine the winner
	winner := "Tie"
	winnerResult := compareHands(player1CombinationName, player1StrongestHand, player2CombinationName, player2StrongestHand)

	if winnerResult == 1 {
		winner = "Player 1"
	} else if winnerResult == 2 {
		winner = "Player 2"
	}

	result := Result{
		Winner:                 winner,
		PlayerCombination1:     player1CombinationName,
		PlayerCombination2:     player2CombinationName,
		PlayerHandCombination1: player1StrongestHand,
		PlayerHandCombination2: player2StrongestHand,
	}

	return TexasHoldemResponse{
		CommunityCards: communityCards,
		Player1:        player1,
		Player2:        player2,
		Result:         result,
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
