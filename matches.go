package main

import (
	"bytes"
	"fmt"
	"errors"
)

type LastMatches []int

var matches *LastMatches

func init(){
	ar := make(LastMatches, 0)
	matches = &ar
}

/* ======= manage cache */

func (lm *LastMatches) Clear(){
	*lm = (*lm)[:0]
}

func (lm *LastMatches) Fill(){
	lm.Clear()
	for idx, _ := range accounts { lm.Append(idx) }
}

func (lm *LastMatches) Append(idx int){
	*lm = append(*lm, idx)
}

func (lm *LastMatches) Length() int {
	return len(*lm)
}

func (lm *LastMatches) AccountAt(idx int) (*Account, int, error) {
	if idx >= 0 && idx < matches.Length() {
		return &accounts[(*lm)[idx]], (*lm)[idx], nil 
	}
	return nil, -1, errors.New("index out of range.")
}

func (lm *LastMatches) Print() {

	var buffer bytes.Buffer

	for idx, accIdx := range *matches {
		 buffer.WriteString(fmt.Sprintf("  [%d]: %s\n", idx, accounts[accIdx].Name))
	}

	if matches.Length() > 30 {
		shell.ShowPaged(buffer.String())
	}else{
		shell.Println(buffer.String())
	}
}

