package library

import (
	"cw-looper/internal/morse"
	"fmt"
)

// Starter is copied only into a missing library or on explicit import. It is
// never merged into existing user data by version number.
func Starter() Library {
	data := Library{SchemaVersion: 1, SeedVersion: 2}
	items := func(id, name string, words [][2]string) List {
		list := List{ID: id, Name: name, Settings: Defaults()}
		for i, word := range words {
			entryID := fmt.Sprintf("%s-%d", id, i)
			list.Entries = append(list.Entries, Entry{ID: entryID, Kind: "item", Title: word[0], Lines: []Line{{ID: entryID + "-line", Station: "learner", CW: word[0], Spoken: word[1]}}})
		}
		return list
	}
	data.Lists = append(data.Lists, items("callsign", "My callsign · KE2ISH", [][2]string{{"KE2ISH", ""}, {"DE KE2ISH", ""}, {"CQ CQ DE KE2ISH K", ""}}))
	data.Lists = append(data.Lists, items("terms", "Common CW terms", [][2]string{
		{"CQ", "see cue"}, {"DE", "D E, from"}, {"K", "K, invitation to transmit"}, {"R", "R, received"}, {"TU", "T U, thank you"}, {"73", "seventy three, best regards"},
		{"QTH", "Q T H, location"}, {"RST", "R S T, readability strength tone"}, {"QRS", "Q R S, send more slowly"}, {"AGN", "again"}, {"PSE", "please"}, {"TNX", "thanks"}, {"UR", "your"}, {"ES", "and"}, {"NAME", "name"}, {"HW", "H W, how copy"}, {"QSO", "Q S O, contact"}, {"QRM", "Q R M, interference"}, {"QRN", "Q R N, static"}, {"QSB", "Q S B, fading"}, {"QRL", "Q R L, busy"}, {"QRZ", "Q R Z, who is calling"},
	}))
	var prosigns [][2]string
	for _, p := range morse.Prosigns {
		prosigns = append(prosigns, [2]string{p.Token, p.Meaning})
	}
	data.Lists = append(data.Lists, items("prosigns", "Continuous prosigns", prosigns))
	exchanges := List{ID: "exchanges", Name: "Two-station exchanges · practice examples", Settings: Defaults()}
	exchange := func(id, title string, firstOther bool, texts ...string) {
		entry := Entry{ID: id, Kind: "exchange", Title: title}
		for i, text := range texts {
			role := "learner"
			if (i%2 == 0) == firstOther {
				role = "other"
			}
			entry.Lines = append(entry.Lines, Line{ID: fmt.Sprintf("%s-line-%d", id, i), Station: role, CW: text})
		}
		exchanges.Entries = append(exchanges.Entries, entry)
	}
	exchange("calling", "Calling and answering CQ", true,
		"CQ CQ DE W1AW W1AW K",
		"W1AW DE KE2ISH KE2ISH <AR>",
		"KE2ISH DE W1AW R TNX FER CALL K")
	exchange("reports", "Exchanging signal reports", false,
		"W1AW DE KE2ISH UR RST 579 579 K",
		"KE2ISH DE W1AW R UR RST 589 589 K")
	exchange("name-location", "Names and locations", false,
		"W1AW DE KE2ISH NAME ELLIOT <BT> QTH NEW YORK K",
		"KE2ISH DE W1AW R NAME ALEX <BT> QTH CONNECTICUT K")
	exchange("repeat-slower", "Asking for repeats and slower sending", false,
		"W1AW DE KE2ISH PSE QRS <BT> PSE AGN NAME K",
		"KE2ISH DE W1AW R NAME ALEX ALEX K",
		"W1AW DE KE2ISH R TU K")
	exchange("signoff", "Signing off", false,
		"TNX FER QSO 73 <SK> W1AW DE KE2ISH",
		"TU 73 <SK> KE2ISH DE W1AW")
	exchange("full-contact", "A complete short contact", true,
		"CQ CQ DE W1AW W1AW K",
		"W1AW DE KE2ISH KE2ISH <AR>",
		"KE2ISH DE W1AW R UR RST 579 <BT> NAME ALEX <BT> QTH CONNECTICUT K",
		"W1AW DE KE2ISH R UR RST 589 <BT> NAME ELLIOT <BT> QTH NEW YORK K",
		"KE2ISH DE W1AW PSE AGN NAME K",
		"W1AW DE KE2ISH NAME ELLIOT ELLIOT K",
		"R TNX FER QSO 73 <SK> KE2ISH DE W1AW",
		"TU 73 <SK> W1AW DE KE2ISH")
	data.Lists = append(data.Lists, exchanges)
	return data
}
