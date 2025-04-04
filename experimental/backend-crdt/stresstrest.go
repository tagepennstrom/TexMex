// File: main.go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func testLargestCoordinate() {

	t1A := []int{1, 1, 2}
	t1B := []int{4} //largest
	println("T1. Expected false:", CompareIndexes(t1A, t1B))

	t2A := []int{3, 1, 2, 2}
	t2B := []int{3, 1, 2, 3} // largest
	println("T2. Expected false:", CompareIndexes(t2A, t2B))

	t3A := []int{2, 1, 4, 4}
	t3B := []int{2, 1, 4, 4, 0} // largest
	println("T3. Expected false:", CompareIndexes(t3A, t3B))

	t4A := []int{2, 1, 4, 5} // largest
	t4B := []int{2, 1, 4, 4, 1}
	println("T4. Expected true:", CompareIndexes(t4A, t4B))

	t5A := []int{1, 2, 3} // largest
	t5B := []int{1, 2, 2}
	println("T5. Expected true:", CompareIndexes(t5A, t5B))

	t6A := []int{1, 2, 3, 4} // largest
	t6B := []int{1, 2, 3}
	println("T6. Expected true:", CompareIndexes(t6A, t6B))

	t7A := []int{5, 0, 0} // largest
	t7B := []int{5}
	println("T7. Expected true:", CompareIndexes(t7A, t7B))

	t8A := []int{7, 2, 5}
	t8B := []int{7, 2, 900} // largest
	println("T8. Expected false:", CompareIndexes(t8A, t8B))
}

func printArrayStructure(data []Item) {
	for i, s := range data {
		fmt.Println(i, "is:", s.Key, s.Values)
	}
}

func (d *Document) DebugDocContents() {
	printArrayStructure(d.Letters)
}

func CRDT() {

	data := []Item{
		{Key: "", Values: []int{0}, ID: 0},
		{Key: "a", Values: []int{1}, ID: 99},
		{Key: "b", Values: []int{2}, ID: 99},
		{Key: "c", Values: []int{4}, ID: 99},
	}
	//	"abc"

	data = Insertion("w", 2, data, 99) // Case 1 - "abWc"
	data = Insertion("x", 1, data, 99) // Case 2 - "aXbwc"
	data = Insertion("y", 0, data, 99) // Case 5 - "Yaxbwc"
	data = Insertion("z", 3, data, 99) // Case X - "yaxZbwc"

	data = Insertion("p", 7, data, 99) // Case 4 - "yaxzbwcP"

	printArrayStructure(data)
}

func CRDTEmpty() []Item {
	data := []Item{
		{Key: "", Values: []int{0}, ID: 0},
	}
	return data
}

func CRDTStressTest() {
	data := CRDTEmpty()

	data = Insertion("a", 0, data, 99)
	data = Insertion("b", 0, data, 99)
	data = Insertion("c", 1, data, 99)
	data = Insertion("c", 3, data, 99)
	data = Insertion("c", 4, data, 99)

	data = Deletion(4, data)

	data = Insertion("x", 3, data, 99)
	data = Insertion("x", 3, data, 99)
	data = Insertion("x", 3, data, 99)
	data = Insertion("z", 4, data, 99)
	data = Insertion("y", 3, data, 99)

	data = Insertion("po", 6, data, 99)
	data = Insertion("lo", 9, data, 99)
	data = Insertion("bo", 10, data, 99)

	data = Deletion(10, data)

	data = Insertion("xx", 9, data, 99)

	printArrayStructure(data)

	dataToString(data)

}

func dataToString(data []Item) {

	for i, s := range data {
		fmt.Print(s.Key)
		if i == 0 {
		}
	}
	fmt.Print("\n")

}

func letterScrambeler(input string, db []Item) []Item {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Create a new random generator with its own seed.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i, letter := range input {
		wg.Add(1)
		go func(i int, letter rune) {
			defer wg.Done()

			// Wait until the current text length equals the intended index (i).
			for {
				mu.Lock()
				currentLen := len(db) - 1
				mu.Unlock()
				if currentLen == i {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}

			// Simulate a random delay to mimic typing.
			time.Sleep(time.Duration(r.Intn(100)) * time.Millisecond)

			// With a 50% chance, simulate a typo:
			if r.Float64() < 0.5 {
				// Generate a wrong letter.
				wrongLetter := string('a' + rune(r.Intn(26)))
				mu.Lock()
				// Since current length equals i, we can only insert at len(db).
				db = Insertion(wrongLetter, len(db), db, 0)
				mu.Unlock()

				// Short delay before deleting the wrong letter.
				time.Sleep(time.Duration(r.Intn(100)) * time.Millisecond)

				mu.Lock()
				// Delete the wrong letter which is the last element.
				db = Deletion(len(db)-1, db)
				mu.Unlock()
			}

			// Ensure that the text length is still exactly i.
			for {
				mu.Lock()
				currentLen := len(db) - 1
				mu.Unlock()
				if currentLen == i {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}

			// Now insert the correct letter at the end.
			mu.Lock()
			db = Insertion(string(letter), len(db), db, 0)
			mu.Unlock()
		}(i, letter)
	}

	wg.Wait()
	return db
}

func typewriterSimulator(db []Item, durationSec int) (string, []Item) {
	// expected will track the correct text state
	expected := ""

	// Create a new random generator with its own seed.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Determine the end time of the simulation.
	endTime := time.Now().Add(time.Duration(durationSec) * time.Second)

	// Loop until the simulation time expires.
	for time.Now().Before(endTime) {
		// Randomly choose an operation.
		// If expected is empty, we only perform insertion.
		var op string
		if len(expected) == 0 {
			op = "insert"
		} else if r.Float64() < 0.3 {
			op = "delete"
		} else {
			op = "insert"
		}

		if op == "insert" {
			// Choose a random letter between 'a' and 'z'
			letter := string('a' + rune(r.Intn(26)))
			// Choose a random position from 0 to len(expected) inclusive.
			pos := 0
			if len(expected) > 0 {
				pos = r.Intn(len(expected) + 1)
			}
			// Update the expected string: insert the letter at the chosen position.
			expected = expected[:pos] + letter + expected[pos:]
			// Insert the letter into the database at the same position.
			db = Insertion(letter, pos, db, 0)
		} else { // deletion
			// Choose a random position from 0 to len(expected)-1.
			pos := r.Intn(len(expected))
			// Update the expected string: remove the letter at the chosen position.
			expected = expected[:pos] + expected[pos+1:]
			// Delete the letter from the database.
			db = Deletion(pos+1, db)
		}

		// Wait a random short duration between operations.
		time.Sleep(time.Duration(r.Intn(5)) * time.Millisecond)
	}

	return expected, db
}

func WriteThis(input string) {
	data := CRDTEmpty()

	expected, db := typewriterSimulator(data, 1)

	println("SIMULATION RESULTS")

	printArrayStructure(db)

	println("\nExpected output vs Actual:")
	println("------")
	println(expected)
	println("------")
	dataToString(db)

	println("------")

	//println("Expected output:", input)
}
