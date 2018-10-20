package command

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/scnewma/auditor/lpass"
	"github.com/scnewma/auditor/password"
	"github.com/scnewma/auditor/pwned"
)

const (
	pwnedAPI  = "https://api.pwnedpasswords.com"
	userAgent = "lpass-auditor"
)

type Audit struct{}

func (c Audit) Execute(args []string) error {
	csvFile, err := getFile(args)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	entries, err := lpass.ParseCSV(csvFile)
	if err != nil {
		return fmt.Errorf("could not parse csv: %v", err)
	}

	pwnedEntries, err := determinePwnedEntries(entries)
	if err != nil {
		return err
	}

	printPwnedEntries(pwnedEntries)

	return nil
}

func getFile(args []string) (*os.File, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("path to csv not provided")
	}

	csvFile, err := os.Open(os.Args[1])
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}

	return csvFile, nil
}

func determinePwnedEntries(entries []lpass.Entry) ([]pwnedEntry, error) {
	u, err := url.Parse(pwnedAPI)
	if err != nil {
		return nil, fmt.Errorf("could not parse pwnedAPI url: %v", err)
	}
	pwnedClient := pwned.NewClient(http.DefaultClient, u, userAgent)

	var pwnedEntries []pwnedEntry

	for _, entry := range entries {
		fmt.Printf("checking %s...", entry.Name)
		pwnedCount, err := password.PwnedCount(pwnedClient, entry.Password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to determine if password is pwned for %s: %v", entry.Name, err)
		}

		if pwnedCount > 0 {
			pwnedEntries = append(pwnedEntries, pwnedEntry{
				name:  entry.Name,
				count: pwnedCount,
			})
			fmt.Print("not ok")
		} else {
			fmt.Print("ok")
		}
		fmt.Println()
	}
	return pwnedEntries, nil
}

func printPwnedEntries(pwnedEntries []pwnedEntry) {
	if len(pwnedEntries) == 0 {
		fmt.Println("\nNo pwned passwords found! Keep up the good work!")
	} else {
		fmt.Println("\nYou should consider changing the passwords for the following entries as their passwords have been pwned:")

		for _, pe := range pwnedEntries {
			fmt.Printf("Name: %s - Count: %d\n", pe.name, pe.count)
		}
	}
}

type pwnedEntry struct {
	name  string
	count int
}
