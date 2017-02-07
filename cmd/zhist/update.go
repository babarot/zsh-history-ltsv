package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/najeira/ltsv"
)

type History struct {
	Date string
	Dir  string
	Cmd  string
}

func parseRecord(line string) History {
	bytes := bytes.NewBufferString(strings.Trim(line, "\t"))
	reader := ltsv.NewReader(bytes)
	token, err := reader.Read()
	if err != nil {
		// skip
		//println(line)
	}
	return History{
		Date: token["date"],
		Dir:  token["dir"],
		Cmd:  token["cmd"],
	}
}

func ReadCities(f func(r History) error) error {
	var fp *os.File
	var err error

	fp, err = os.Open(os.Getenv("ZSH_HISTORY_FILE"))
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		if err := f(parseRecord(scanner.Text())); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return nil
}

func updateCommand() {
	os.Remove(databasePath)

	db, err := openDB()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	var counter int64

	err = ReadCities(func(h History) error {
		counter++
		_, err = db.Exec(
			`INSERT OR REPLACE INTO history
				(date, dir, cmd)
			VALUES
				(?, ?, ?)`,
			h.Date,
			h.Dir,
			h.Cmd,
		)
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
}
