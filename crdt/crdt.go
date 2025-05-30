package crdt

import (
	"encoding/json"
	"fmt"
	"os"
)

type Document struct {
	CursorPosition *Item
	Textcontent    LinkedList
	Active         bool
}

type LinkedList struct {
	Head   *Item
	Tail   *Item
	Length int
}

type Item struct {
	Letter   string
	Location CoordT // istället för []int här
	ID       int
	Prev     *Item
	Next     *Item
}

type CoordT struct {
	Coordinate []int `json:"coord"`
	ID         int   `json:"id"`
}

type Change struct {
	FromA int    `json:"fromA"` // Start index original document
	ToA   int    `json:"toA"`   // Slut index original document
	FromB int    `json:"fromB"` // Start index new document
	ToB   int    `json:"toB"`   // Slut index new document
	Text  string `json:"text"`  // Tillagd text
}

type EditDocMessage struct {
	Document    string   `json:"document"`
	Changes     []Change `json:"changes"`
	CursorIndex int      `json:"cursorIndex"`
}

type UpdatedDocMessage struct {
	CursorIndex int            `json:"cursorIndex"`
	CChanges    []CoordChanges `json:"coordChanges"`
}

type CRDTNode struct {
	Letter   string `json:"letter"`
	Location CoordT `json:"location"`
	ID       int    `json:"id"`
}

type CoordChanges struct {
	Coordinate CoordT `json:"coordinate"`
	Operation  string `json:"operation"`
	Letter     string `json:"letter"`
}

// *
// Globala variabler
// *

var uID int = -1
var DocuMain Document

// *
// Snapshot state logik (CRDT -> JSON)
// *

func (ll *LinkedList) Snap() []CRDTNode {
	var out []CRDTNode
	for node := ll.Head; node != nil; node = node.Next {
		out = append(out, CRDTNode{
			Letter:   node.Letter,
			Location: node.Location,
			ID:       node.ID,
		})
	}
	return out
}

func (d *Document) Snapshot() ([]byte, error) {
	snapshot := struct {
		Textcontent []CRDTNode `json:"textcontent"`
	}{
		Textcontent: d.Textcontent.Snap(),
	}
	return json.Marshal(snapshot)
	// skickar som byte-kod så det inte blir långsamt med stora dokument
}

func LoadSnapshot(jsonStr string) string {

	jsonBytes := []byte(jsonStr)

	var toLoad struct {
		Textcontent []CRDTNode
	}

	err := json.Unmarshal(jsonBytes, &toLoad)
	if err != nil {
		println("error when unmarshalling loaded snapshot (LoadSnapshot in crdt.go)")
		panic(err)
	}

	loadedDoc := NewDocument()

	// skippa dummy element som redan skapas från NewDocument() ovan
	for i := 1; i < len(toLoad.Textcontent); i++ {
		node := toLoad.Textcontent[i]
		newItem := Item{Letter: node.Letter, Location: node.Location, ID: node.ID}
		loadedDoc.Textcontent.Append(&newItem)
	}
	loadedDoc.Active = true
	DocuMain = loadedDoc

	docAsStr := DocuMain.ToString()
	DocuMain.Textcontent.Length = len(docAsStr)

	return docAsStr
}

// *
// Ladda in state eller ändringar
// *

func SetUserID(id int) {
	uID = id
	println("User ID set in CRDT as ID: ", uID)
}

func buildCoordChange(crd CoordT, op string, ltr string) CoordChanges {

	return CoordChanges{
		Coordinate: crd,
		Operation:  op,
		Letter:     ltr,
	}
}

func (d *Document) HandleCChange(jsonCChange string) string {

	var cChanges []CoordChanges
	c := []byte(jsonCChange)
	json.Unmarshal(c, &cChanges)

	var iChange []Change

	for _, change := range cChanges {
		coord := change.Coordinate

		switch change.Operation {

		case "delete":
			i := d.DeleteAtCoordinate(coord)

			change := Change{
				FromA: -1,
				ToA:   -1,
				FromB: i,
				ToB:   i + 1,
				Text:  "",
			}

			iChange = append(iChange, change)
			break

		case "insert":
			i := d.InsertAtCoordinate(coord, change.Letter)

			change := Change{
				FromA: -1,
				ToA:   -1,
				FromB: i,
				ToB:   i,
				Text:  change.Letter,
			}
			iChange = append(iChange, change)

			break

		case "paste":

			i := d.PasteInsertion(coord, change.Letter)

			change := Change{
				FromA: -1,
				ToA:   -1,
				FromB: i,
				ToB:   i + len(change.Letter),
				Text:  change.Letter,
			}
			iChange = append(iChange, change)
		}

	}

	bytearray, _ := json.Marshal(iChange)

	return string(bytearray)

}

func AddRaise(c CoordT, raise int, raiseGrade int) CoordT {

	ogLen := len(c.Coordinate)

	new := make([]int, ogLen+raiseGrade+1)

	copy(new, c.Coordinate)

	for i := 0; i < raiseGrade; i++ {
		new[i+ogLen] = 0
	}

	new[ogLen+raiseGrade] = raise

	return CoordT{
		Coordinate: new,
		ID:         c.ID,
	}
}

func (d *Document) PasteInsert(paste string, from int) CoordT {

	prevItem, caseFour := d.IndexToCoordinate(from)

	var firstPlacement CoordT
	if caseFour {
		firstPlacement = GetAppendCoordinate(prevItem.Location.Coordinate, uID)
	} else {
		nextItem := prevItem.Next
		coord := findIntermediateCoordinate(prevItem.Location, nextItem.Location)

		firstPlacement = CoordT{
			Coordinate: coord,
			ID:         uID,
		}
	}

	d.PasteInsertion(firstPlacement, paste)

	return firstPlacement
}

// hitta hur högt pastade koordinater måste höjas utan att krocka med next
func findRaiseGrade(ins CoordT, comp CoordT) int {

	if len(ins.Coordinate) >= len(comp.Coordinate) {
		return 0
	} else {
		return len(comp.Coordinate) - len(ins.Coordinate)
	}

}

func (d *Document) PasteInsertion(headCoord CoordT, pasteText string) int {

	prevItem, index := findPrevItem(headCoord, d.Textcontent)
	finalNext := prevItem.Next // nil?

	// skapa en ny nod
	headItem := Item{
		Letter:   string(pasteText[0]),
		Location: headCoord,
		ID:       headCoord.ID,
	}

	headItem.Prev = prevItem
	prevItem.Next = &headItem

	var rg int
	if finalNext != nil {
		rg = findRaiseGrade(headCoord, finalNext.Location)
	} else {
		rg = 0
	}

	userID := headCoord.ID
	prevI := &headItem

	for i, ch := range pasteText {
		if i == 0 {
			continue
		}
		c := AddRaise(headCoord, i, rg)

		new := &Item{Letter: string(ch), Location: c, ID: userID}

		new.Prev = prevI
		prevI.Next = new

		prevI = new
	}

	if finalNext != nil {
		finalNext.Prev = prevI
		prevI.Next = finalNext
	} else {
		d.Textcontent.Tail = prevI
	}

	d.Textcontent.Length = d.Textcontent.Length + len(pasteText)

	return index
}

func applyAndRecordOperations(changes []Change, cursorIndex int) []CoordChanges {
	DocuMain.SetCursorAt(cursorIndex)

	var allChanges []CoordChanges

	for _, change := range changes {

		if change.Text == "" {
			// DELETE Operation
			for i := change.ToA; i > change.FromA; i-- {
				crd := DocuMain.DeleteAtIndex(i)

				change := buildCoordChange(crd, "delete", "")
				allChanges = append(allChanges, change)
			}

		} else if (change.FromA == change.ToA) && len(change.Text) > 1 {
			// PASTE Operation
			crd := DocuMain.PasteInsert(change.Text, change.FromB)

			change := buildCoordChange(crd, "paste", change.Text)
			allChanges = append(allChanges, change)

		} else if change.FromA == change.ToA {
			// INSERT Operation
			for i, ch := range change.Text {
				crd := DocuMain.LoadInsert(string(ch), change.FromB+i, uID)

				change := buildCoordChange(crd, "insert", string(ch))
				allChanges = append(allChanges, change)
			}

		} else {
			// SELECT AND REPLACE Operation
			for i := change.ToA; i > change.FromA; i-- {
				crd := DocuMain.DeleteAtIndex(i)

				change := buildCoordChange(crd, "delete", "")
				allChanges = append(allChanges, change)

			}

			i := 0
			for _, ch := range change.Text {
				crd := DocuMain.LoadInsert(string(ch), change.FromB+i, uID)

				change := buildCoordChange(crd, "insert", string(ch))
				allChanges = append(allChanges, change)

				i++
			}
		}

	}

	return allChanges
}

func UpdateDocument(changes []Change, cursorIndex int) UpdatedDocMessage {

	if uID == -1 {
		println("Error: User ID not initialized")
		os.Exit(69)
	}

	allChanges := applyAndRecordOperations(changes, cursorIndex)

	return UpdatedDocMessage{
		CursorIndex: DocuMain.CursorIndex(),
		CChanges:    allChanges,
	}
}

// *
// Vanliga CRDT Funktioner
// *

func PrintDocument(verbose bool) {
	var result string
	for current := DocuMain.Textcontent.Head; current != nil; current = current.Next {
		result += current.Letter
		if verbose {
			fmt.Println(" x ", current.Location.Coordinate, current.Letter)
		}
	}
	println("Result:", result, "(PrintDocument in crdt.go)")
	println("Doc tail:", DocuMain.Textcontent.Tail.Letter)

}

func DocumentFromStr(str string) Document {
	doc := NewDocument()
	doc.Active = true // "activate" project
	for _, ch := range str {
		doc.Insert(string(ch), 0)
	}
	return doc
}

func (doc *Document) InitEmptyDocumentFromString(str string) {
	doc.Active = true // "activate" project
	for _, ch := range str {
		doc.Insert(string(ch), 0)
	}
}

func (doc *Document) ToString() string {
	str := ""
	item := doc.Textcontent.Head.Next
	for item != nil {
		str += item.Letter
		item = item.Next
	}
	return str
}

func (doc *Document) CursorIndex() int {
	i := 0
	item := doc.Textcontent.Head
	for item != nil {
		if item == doc.CursorPosition {
			return i
		}
		item = item.Next
		i++
	}
	return i
}

func (doc *Document) SetCursorAt(index int) {
	i := 0
	item := doc.Textcontent.Head
	for item != nil {
		if i == index {
			doc.CursorPosition = item
		}
		item = item.Next
		i++
	}
}

func (doc *Document) InsertAtCoordinate(c CoordT, l string) int {

	var index int
	doc.Textcontent, index = Insertion(l, c, doc.Textcontent, c.ID)
	return index
}

func (doc *Document) DeleteAtCoordinate(c CoordT) int {
	temp, index := findPrevItem(c, doc.Textcontent) // blir inte prev då det är samma coord
	toDel := temp.Next
	prev := toDel.Prev

	// forward link
	prev.Next = toDel.Next

	if toDel.Next != nil {
		// backward link
		toDel.Next.Prev = toDel.Prev
	} else {
		doc.Textcontent.Tail = toDel.Prev
	}

	doc.Textcontent.Length--
	return index
}

func (ll *LinkedList) Append(newItem *Item) {
	if ll.Tail == nil {
		println("Error: Detta bör aldrig hända")
		ll.Head = newItem
		ll.Tail = newItem
		return
	}

	ll.Tail.Next = newItem
	newItem.Prev = ll.Tail
	ll.Tail = newItem
	newItem.Next = nil
}

// Returns true if c1 is smaller than c2
func CompareIndexes(c1 CoordT, c2 CoordT) bool {

	coord1 := c1.Coordinate
	coord2 := c2.Coordinate
	len1 := len(coord1)
	len2 := len(coord2)

	count := 0

	for count < len1 && count < len2 {

		if coord1[count] > coord2[count] {
			return true //c1 biggest

		} else if coord1[count] < coord2[count] {
			return false //c2 biggest
		}
		count++
	}

	if len1 == len2 {

		if c1.ID < c2.ID {
			return true
		} else if c2.ID < c1.ID {
			return false
		} else {

			//fmt.Println(c1.Coordinate, "+", c1.ID, "vs", c2.Coordinate, "+", c2.ID)
			//os.Exit(1)
			return true
		}

	}
	return len1 > len2
}

func findIntermediateCoordinate(pCoord CoordT, nCoord CoordT) []int {

	prevCord := pCoord.Coordinate
	nextCord := nCoord.Coordinate

	prevLen := len(prevCord)
	nextLen := len(nextCord)
	count := 0
	newCoordinate := []int{}

	for count < prevLen && count < nextLen {
		comparison := nextCord[count] - prevCord[count]

		if comparison == 0 {
			newCoordinate = append(newCoordinate, prevCord[count])

			// Case 1. Plats finns
		} else if comparison > 1 {
			newCoordinate = append(newCoordinate, prevCord[count]+1)

			//Case Z
			if count+1 < prevLen {
				newCoordinate = append(newCoordinate, prevCord[count+1]+1)

			}
			return newCoordinate

			// Case 2 and Case X. Ingen plats finns
		} else if comparison == 1 {
			newCoordinate = append(newCoordinate, prevCord[count])

			// Case 2.
			if count+1 == prevLen {
				newCoordinate = append(newCoordinate, 1) // Raised

				return newCoordinate

				// Case X.
			} else {
				newCoordinate = append(newCoordinate, prevCord[count+1]+1) // Inkrementera raise

				return newCoordinate

			}

		}
		count++

	}
	if count < nextLen {
		// Case 3. Ingen raised plats
		if nextCord[count] == 1 {
			newCoordinate = append(newCoordinate, 0) // Lower bound

			// varning: case saknas. om next också har nollor måste vi loopa?
			// WHILE LOOP.

			newCoordinate = append(newCoordinate, 1) // Lower bound

			// Case Y. Next längre än prev, men raised plats
		} else if nextCord[count] == 0 {
			// Case W

			for i := nextLen - count; i > 0; i-- {
				newCoordinate = append(newCoordinate, 0)

			}
			newCoordinate = append(newCoordinate, 1)

		} else {
			newCoordinate = append(newCoordinate, 1)

		}
	}

	return newCoordinate
}

func findPrevItem(insertionCoord CoordT, db LinkedList) (*Item, int) {

	prev := db.Head
	index := 0
	for prev.Next != nil {
		if CompareIndexes(prev.Next.Location, insertionCoord) {
			break
		} else {
			index++
			prev = prev.Next
		}
	}
	return prev, index
}

func Insertion(letter string, coordinate CoordT, db LinkedList, uID int) (LinkedList, int) {

	prevItem, i := findPrevItem(coordinate, db)

	newItem := Item{Letter: letter, Location: coordinate, ID: uID} //prev och next
	db.Length++

	// Case 4
	nextItem := prevItem.Next
	if nextItem == nil {
		db.Append(&newItem)
		return db, i
	}

	prevItem.Next = &newItem
	newItem.Prev = prevItem

	nextItem.Prev = &newItem
	newItem.Next = nextItem

	return db, i
}

func Deletion(coordinate CoordT, db LinkedList) LinkedList {
	prevItem, _ := findPrevItem(coordinate, db)

	itemToRemove := prevItem.Next

	nextItem := itemToRemove.Next

	prevItem.Next = nextItem

	if nextItem != nil {
		nextItem.Prev = prevItem
	} else {
		db.Tail = prevItem
	}

	db.Length--

	return db
}

func GetAppendCoordinate(prevCoord []int, uID int) CoordT {
	insertCoord := []int{prevCoord[0] + 1}
	var newLocation CoordT
	newLocation.Coordinate = insertCoord
	newLocation.ID = uID

	return newLocation
}

func (d *Document) findInsertCoord(uID int) CoordT {
	cursorPosCoordinate := d.CursorPosition.Location // TODO REWRITE, det här är oläsbart

	// Case 4
	if d.CursorPosition.Next == nil {
		return GetAppendCoordinate(cursorPosCoordinate.Coordinate, uID)
	}

	cursorPosNextCoord := d.CursorPosition.Next.Location
	insertCoord := findIntermediateCoordinate(cursorPosCoordinate, cursorPosNextCoord)
	return CoordT{
		Coordinate: insertCoord,
		ID:         uID,
	}
}

func (d *Document) Insert(letter string, uID int) {
	location := d.findInsertCoord(uID)
	d.Textcontent, _ = Insertion(letter, location, d.Textcontent, uID)
	d.CursorForward()
}

func (d *Document) CursorForward() {
	if d.CursorPosition.Next != nil {
		d.CursorPosition = d.CursorPosition.Next

	} else {
		println("Can't move cursor further.")
	}
}

func (d *Document) CursorBackwards() {
	// BOF Placeholder har ID 0
	if d.CursorPosition.ID != 0 {
		d.CursorPosition = d.CursorPosition.Prev

	} else {
		println("Error: Can't move cursor further back")
	}
}

func (d *Document) IndexToCoordinate(index int) (Item, bool) {

	docLength := d.Textcontent.Length
	var newPosition Item
	var atEnd bool = false

	if index >= docLength {
		index = docLength
		atEnd = true
		return *d.Textcontent.Tail, true
	}

	if index < 0 {
		println("Error. Can't move cursor out of bounds (IndexToCoordinate)")
		os.Exit(1)

	} else {

		current := d.Textcontent.Head
		for i := 0; i < index; i++ {
			current = current.Next
		}

		newPosition = *current
	}
	return newPosition, atEnd
}

func coordEqual(a CoordT, b CoordT) bool {
	if a.ID != b.ID {
		return false
	}
	if len(a.Coordinate) != len(b.Coordinate) {
		return false
	}

	for i := range a.Coordinate {
		if a.Coordinate[i] != b.Coordinate[i] {
			return false
		}
	}
	return true
}

func (d *Document) CoordinateToIndex(coordJson string) int {
	var coord CoordT
	json.Unmarshal([]byte(coordJson), &coord)
	i := 0
	for current := DocuMain.Textcontent.Head; current != nil; current = current.Next {
		if coordEqual(current.Location, coord) || current.Location.Coordinate[0] > coord.Coordinate[0] {
			return i
		}
		if current.Letter != "" {
			i++
		}
	}
	println("coordinate doesn't exist")
	return i
}

func (d *Document) CoordinateToIndex2(coordJson string) int {
	var coordSearching CoordT
	json.Unmarshal([]byte(coordJson), &coordSearching)
	i := 0

	cur := d.Textcontent.Head

	for cur != nil {
		if CompareIndexes(cur.Location, coordSearching) {
			return i
		}
		cur = cur.Next
		i++
	}

	return i
}

func (d *Document) LoadInsert(letter string, index int, uID int) CoordT {
	prevItem, caseFour := d.IndexToCoordinate(index)

	var location CoordT
	if caseFour {
		location = GetAppendCoordinate(prevItem.Location.Coordinate, uID)
	} else {
		nextItem := prevItem.Next
		coord := findIntermediateCoordinate(prevItem.Location, nextItem.Location)

		location = CoordT{
			Coordinate: coord,
			ID:         uID,
		}
	}

	d.Textcontent, _ = Insertion(letter, location, d.Textcontent, uID)
	if d.CursorIndex() == index {
		d.CursorForward()
	}

	return location
}

func (d *Document) MoveCursor(index int) {
	docLength := d.Textcontent.Length

	if index > docLength || index < 0 {
		println("Error. Can't move cursor out of bounds")
	} else {

		var newPosition Item
		current := d.Textcontent.Head
		for i := 0; i < index; i++ {
			current = current.Next
		}
		newPosition = *current
		d.CursorPosition = &newPosition
	}
}

func (d *Document) DeleteAtIndex(index int) CoordT {
	cursorIndex := d.CursorIndex()

	d.SetCursorAt(index)
	deletedCoord := d.Delete()

	if cursorIndex >= index {
		d.SetCursorAt(cursorIndex - 1)
	} else {
		d.SetCursorAt(cursorIndex)
	}

	return deletedCoord
}

// OBS använder oss bara av current cursor position för deletion just nu
func (d *Document) Delete() CoordT {
	if d.CursorPosition.Prev != nil {
		savedCursor := d.CursorPosition

		deletedCoordinate := savedCursor.Location

		// Link the previous node to the next node
		savedCursor.Prev.Next = savedCursor.Next

		if savedCursor.Next != nil {

			savedCursor.Next.Prev = savedCursor.Prev
		} else {
			// Om det är tailen
			d.Textcontent.Tail = savedCursor.Prev
		}

		d.Textcontent.Length--

		return deletedCoordinate

	} else {
		println("Error: det fanns inte en else förut men jag behövde nån return. Antar att detta aldrig händer ( Delete() i crdt.go )")
		os.Exit(69)
		return d.CursorPosition.Location
	}
}

func NewDocument() Document {

	var location CoordT = CoordT{
		Coordinate: []int{0},
		ID:         0,
	}

	// BOD = Beginning Of File
	BOF := Item{
		Letter:   "",
		Location: location,
		ID:       0,
		Next:     nil,
		Prev:     nil,
	}

	textContent := LinkedList{
		Head:   &BOF,
		Tail:   &BOF,
		Length: 0,
	}

	d := Document{
		Textcontent:    textContent,
		CursorPosition: &BOF,
	}

	return d
}

func (d *Document) CordReset() {
	length := d.Textcontent.Length
	if length < 1 {
		return
	}

	current := d.Textcontent.Head.Next
	for i := 1; i <= length; i++ {
		current.Location.Coordinate = []int{i}
		current = current.Next
	}
}
