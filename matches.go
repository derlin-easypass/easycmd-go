package main

import (
	"bytes"
	"errors"
	"fmt"
)

const notAssigned = -1

type LastMatches []int

var matches *LastMatches

func init() {
	ar := make(LastMatches, 0)
	matches = &ar
}

/* ======= manage cache */

func (lm *LastMatches) Clear() {
	*lm = (*lm)[:0]
}

func (lm *LastMatches) Fill() {
	lm.Clear()
	for idx, _ := range accounts {
		lm.Append(idx)
	}
}

func (lm *LastMatches) Append(idx int) {
	*lm = append(*lm, idx)
}

func (lm *LastMatches) Length() int {
	return len(*lm)
}

func (lm *LastMatches) AccountAt(idx int) (*Account, int, error) {
	if idx >= 0 && idx < matches.Length() {
		accountIdx := (*lm)[idx]
		if accountIdx >= 0 && accountIdx < len(accounts) {
			return &accounts[(*lm)[idx]], (*lm)[idx], nil
		}
	}
	return nil, -1, errors.New("index out of range.")
}

func (lm *LastMatches) RemoveIdx(idx int) {
	for i, val := range *lm {
		if val == idx {
			shell.Printf("making %d unassigned.\n", i)
			(*lm)[i] = notAssigned
			return
		}
	}
}

func (lm *LastMatches) Print() {

	shell.Println(lm.Length(), "accounts.")

	var buffer bytes.Buffer
	for idx, accIdx := range *lm {
		if accIdx != notAssigned {
			buffer.WriteString(fmt.Sprintf("  [%d]: %s\n", idx, accounts[accIdx].Name))
		}
	}

	if lm.Length() > 30 {
		shell.ShowPaged(buffer.String())
	} else {
		shell.Println(buffer.String())
	}
}
