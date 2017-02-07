package main

import (
	"fmt"
	"os"

	"github.com/pascalw/go-alfred"
)

type searchOutputer interface {
	Error(err error)
	Result(r searchResult)
	Flush()
}

type terminalOutputer struct {
}

func (t *terminalOutputer) Error(err error) {
	fmt.Println(err.Error())
}

func (t *terminalOutputer) Result(r searchResult) {
	fmt.Printf("%s\n",
		r.Cmd,
	)
}

func (t *terminalOutputer) Flush() {
}

type alfredOutputer struct {
	*alfred.AlfredResponse
}

func newAlfredOutputer() searchOutputer {
	return &alfredOutputer{
		AlfredResponse: alfred.NewResponse(),
	}
}

func (a *alfredOutputer) Error(err error) {
	a.AlfredResponse.AddItem(&alfred.AlfredResponseItem{
		Valid:    false,
		Uid:      "error",
		Title:    "Error Occurred",
		Subtitle: err.Error(),
		Icon:     "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/AlertStopIcon.icns",
	})
}

func (a *alfredOutputer) Result(r searchResult) {
	a.AddItem(&alfred.AlfredResponseItem{
		Valid:    true,
		Title:    fmt.Sprintf("%s, %s", r.Dir, r.Cmd),
		Subtitle: r.Date,
		Arg:      fmt.Sprintf("%s, %s", r.Dir, r.Cmd),
	})
}

func (a *alfredOutputer) Flush() {
	a.AlfredResponse.Print()
}

type searchResult struct {
	Date string
	Dir  string
	Cmd  string
}

func searchTimezones(queryTerms []string) ([]searchResult, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}

	alfred.InitTerms(queryTerms)

	rows, err := db.Query("SELECT * FROM history")
	if err != nil {
		return nil, err
	}

	results := []searchResult{}
	for rows.Next() {
		var row searchResult
		err = rows.Scan(&row.Date, &row.Dir, &row.Cmd)
		if err != nil {
			return results, err
		}

		if alfred.MatchesTerms(queryTerms, row.Cmd) {
			results = append(results, row)
			// log.Printf("HIT %+v", row)
		}
	}

	return results, nil
}

func searchCommand(s searchOutputer, queryTerms []string) {
	defer s.Flush()

	results, err := searchTimezones(queryTerms)
	if err != nil {
		s.Error(err)
		os.Exit(1)
	}
	// log.Printf("Found %d matches", len(results))

	for _, r := range results {
		s.Result(r)
	}
}
