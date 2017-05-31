package split

import "testing"

const shortText = "Bacon ipsum dolor amet kielbasa tenderloin brisket short ribs."
const baconText = "Bacon ipsum dolor amet kielbasa tenderloin brisket short ribs tri-tip venison turkey ground round sausage corned beef hamburger pork chop turducken jerky. Drumstick doner landjaeger, frankfurter flank meatball strip steak pig jowl sirloin cow ham hock tenderloin venison. Turkey capicola salami, bresaola biltong meatloaf andouille shankle cupim. Pastrami ribeye meatball, strip steak pig picanha shankle ham hock drumstick pancetta jowl kielbasa. Turducken andouille doner, salami kielbasa meatloaf strip steak biltong pork belly alcatra."

func TestSplitterWihShortMessage(t *testing.T) {
	parts := Splitter(shortText)
	expectOnly1Part(t, parts)
}

func TestSplitter(t *testing.T) {
	parts := Splitter(baconText)
	expectPartNumber(t, parts, 4)
}

func expectOnly1Part(t *testing.T, parts []string) {
	expectPartNumber(t, parts, 1)
}

func expectPartNumber(t *testing.T, parts []string, expected int) {
	if n := len(parts); n != expected {
		t.Error("Expected part number is", expected, ", got", n)
	}
}

func TestAlternativeSplitter(t *testing.T) {
	parts := RecursiveSplitter(baconText, make([]string, 0))
	expectPartNumber(t, parts, 4)
}
