package main

/*
*

	#
	#
	#
	LÄGGER FILER HÄR OM JAG KANSKE BEHÖVER DOM IGEN
	MEDANS JAG REFACTORAR ALLT
	#
	#
	#
*/
func findInsertionCoordinate(insertionCoord []int, db LinkedList) []int {

	prev := db.Head
	for prev.Next != nil {
		if CompareIndexes(prev.Coordinate, insertionCoord) {
			break
		} else {
			prev = prev.Next
		}
	}

	prevCoordinate := prev.Coordinate

	if prev.Next == nil {
		// Case 4: Insertion at end. [3, 1] EOF -> [4]
		return []int{prevCoordinate[0] + 1}
		// TODO: Move tail
	}

	nextCoordinate := prev.Next.Coordinate // edge case om den inte existerar?

	return getCoordinateProperties(prevCoordinate, nextCoordinate)

}
