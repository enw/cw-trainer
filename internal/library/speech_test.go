package library

import "testing"

func TestGeneratedSpeechAndOverride(t *testing.T) {
	line := Line{CW: "CQ DE KE2ISH HELLO <AR>"}
	if got := Spoken(line); got != "see cue D E K E two I S H hello end of message" {
		t.Fatal(got)
	}
	line.Spoken = "my call"
	line.CW = "KE2ISH"
	if Spoken(line) != "my call" {
		t.Fatal("CW replaced override")
	}
	line.Spoken = ""
	if Spoken(line) != "K E two I S H" {
		t.Fatal("clearing override did not restore generated speech")
	}
}

func TestAdjacentProsignsPreserveWordPronunciation(t *testing.T) {
	for text, want := range map[string]string{
		"KE2ISH<AR>":      "K E two I S H end of message",
		"CQ<BT>HELLO<SK>": "see cue break hello end of contact",
	} {
		if got := Spoken(Line{CW: text}); got != want {
			t.Fatalf("%s: %s", text, got)
		}
	}
}
