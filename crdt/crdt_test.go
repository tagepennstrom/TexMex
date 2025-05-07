package crdt

import (
	"testing"
)

func assertDocumentsAreEqual(actual UpdatedDocMessage, expected UpdatedDocMessage, t *testing.T) {
	if actual.Document != expected.Document {
		t.Errorf(`Document = %q, but expected %q`, actual.Document, expected.Document)
	}
	if actual.CursorIndex != expected.CursorIndex {
		t.Errorf(`CursorIndex = %v, but expected %v`, actual.CursorIndex, expected.CursorIndex)
	}
}

func TestInsert(t *testing.T) {
	document := "one three"
	cursorPos := 3
	text := " two"
	changes := []Change{
		{
			FromA: cursorPos,
			ToA:   cursorPos,
			FromB: cursorPos,
			ToB:   cursorPos + len(text),
			Text:  text,
		},
	}
	actual := UpdateDocument(document, changes, cursorPos)

	expected := UpdatedDocMessage{
		Document:    "one two three",
		CursorIndex: 7,
	}
	assertDocumentsAreEqual(actual, expected, t)
}

func TestInsertAtEnd(t *testing.T) {
	document := "one two"
	cursorPos := 7
	text := " three"
	changes := []Change{
		{
			FromA: cursorPos,
			ToA:   cursorPos,
			FromB: cursorPos,
			ToB:   cursorPos + len(text),
			Text:  text,
		},
	}
	actual := UpdateDocument(document, changes, cursorPos)

	expected := UpdatedDocMessage{
		Document:    "one two three",
		CursorIndex: 13,
	}
	assertDocumentsAreEqual(actual, expected, t)
}

func TestDeleteLetter(t *testing.T) {
	document := "abc"
	cursorPos := 1
	changes := []Change{
		{
			FromA: cursorPos,
			ToA:   cursorPos,
			FromB: cursorPos,
			ToB:   cursorPos,
			Text:  "",
		},
	}
	actual := UpdateDocument(document, changes, cursorPos)

	expected := UpdatedDocMessage{
		Document:    "ac",
		CursorIndex: 1,
	}
	assertDocumentsAreEqual(actual, expected, t)
}

func TestDeleteSelection(t *testing.T) {
	document := "one three two"
	cursorPos := 4
	changes := []Change{
		{
			FromA: cursorPos,
			ToA:   cursorPos + 5,
			FromB: cursorPos,
			ToB:   cursorPos,
			Text:  "",
		},
	}
	actual := UpdateDocument(document, changes, cursorPos)

	expected := UpdatedDocMessage{
		Document:    "one two",
		CursorIndex: cursorPos,
	}
	assertDocumentsAreEqual(actual, expected, t)
}

func TestInsertEarlierOnSameLine(t *testing.T) {
	document := "word"
	cursorPos := 4
	text := "a "
	changes := []Change{
		{
			FromA: 0,
			ToA:   0,
			FromB: 0,
			ToB:   len(text),
			Text:  text,
		},
	}
	actual := UpdateDocument(document, changes, cursorPos)

	expected := UpdatedDocMessage{
		Document:    "a word",
		CursorIndex: cursorPos + len(text),
	}
	assertDocumentsAreEqual(actual, expected, t)
}

func TestInsertAfterOnSameLine(t *testing.T) {
	document := "one"
	cursorPos := 0
	text := " two"
	changes := []Change{
		{
			FromA: 3,
			ToA:   3,
			FromB: 3,
			ToB:   3 + len(text),
			Text:  text,
		},
	}
	actual := UpdateDocument(document, changes, cursorPos)

	expected := UpdatedDocMessage{
		Document:    "one two",
		CursorIndex: 0,
	}
	assertDocumentsAreEqual(actual, expected, t)
}

func TestDeleteEarlierOnSameLine(t *testing.T) {
	document := "a word"
	cursorPos := 6
	changes := []Change{
		{
			FromA: 0,
			ToA:   1,
			FromB: 0,
			ToB:   0,
			Text:  "",
		},
	}
	actual := UpdateDocument(document, changes, cursorPos)

	expected := UpdatedDocMessage{
		Document:    "word",
		CursorIndex: cursorPos - 2,
	}
	assertDocumentsAreEqual(actual, expected, t)
}

func TestDeleteAfterOnSameLine(t *testing.T) {
	document := "one two"
	cursorPos := 0
	changes := []Change{
		{
			FromA: 3,
			ToA:   6,
			FromB: 3,
			ToB:   3,
			Text:  "",
		},
	}
	actual := UpdateDocument(document, changes, cursorPos)

	expected := UpdatedDocMessage{
		Document:    "one",
		CursorIndex: cursorPos,
	}
	assertDocumentsAreEqual(actual, expected, t)
}
