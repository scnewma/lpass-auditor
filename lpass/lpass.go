package lpass

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

// Entry represents a set of values in LastPass (website, secure note, etc)
type Entry struct {
	URL      url.URL
	Username string
	Password string
	Extra    string
	Name     string
	Fav      int
}

// ParseCSV will convert a LastPass exported CSV to []Entry
func ParseCSV(r io.Reader) ([]Entry, error) {
	reader := csv.NewReader(r)

	// skip header
	reader.Read()

	var entries []Entry
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("encountered error while parsing csv: %v", err)
		}

		entries = append(entries, entryFor(line))
	}

	return entries, nil
}

func entryFor(line []string) Entry {
	u, _ := url.Parse(line[0])

	fav, _ := strconv.Atoi(line[5])
	return Entry{
		URL:      *u,
		Username: line[1],
		Password: line[2],
		Extra:    line[3],
		Name:     line[4],
		Fav:      fav,
	}
}
