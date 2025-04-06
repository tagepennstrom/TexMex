package main

import (
	"math/rand"
	"sync"
	"time"
)

func SimulateLetter() {
	var wg sync.WaitGroup
	db1 := CRDTEmpty()
	db2 := CRDTEmpty()

	db1 = Insertion("x", 0, db1, 1)
	db1 = Insertion("x", 0, db1, 1)

	db2 = Insertion("x", 0, db2, 1)
	db2 = Insertion("x", 0, db2, 1)

	// Simulate client 1
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 := rand.New(rand.NewSource(time.Now().UnixNano()))

		for i := 0; i < 10; i++ {
			db1 = Insertion("a", 0, db1, 1)

			db2 = Insertion("a", 0, db2, 1)
			time.Sleep(time.Duration(r1.Intn(6)) * time.Millisecond)

		}
	}()

	// Simulate client 2
	wg.Add(1)
	go func() {
		defer wg.Done()
		r2 := rand.New(rand.NewSource(time.Now().UnixNano()))

		for i := 1; i < 11; i++ {
			db2 = Insertion("b", i, db2, 2)

			db1 = Insertion("b", i, db1, 2)
			time.Sleep(time.Duration(r2.Intn(5)) * time.Millisecond)

		}
	}()

	wg.Wait()

	dataToString(db1)
	dataToString(db2)

	//printArrayStructure(db1)
	//printArrayStructure(db2)

}

func SimulateDoc() {
	doc := NewDocument()

	doc.Insert("a", 1) // +a1
	doc.Insert("a", 1) // +a2
	doc.Insert("a", 1) // +a3
	doc.Insert("b", 1) // +b4
	doc.Insert("b", 1) // +b5
	doc.MoveCursor(3)
	doc.Delete() // -a3
	doc.Delete() // -a2
	doc.Delete() // -a1
	doc.Delete() // inget händer för den är i början av dokumentet

	doc.DebugDocContents()
}
