package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type Account struct {
	Name     string `json:"name"`
	Pseudo   string `json:"pseudo"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Notes    string `json:"notes"`
}

type Accounts []Account

var accounts Accounts

/* =======load/save */

func LoadAccounts(path string, pass string) error {
	// the content has been generated with:
	//  openssl enc -aes-128-cbc -pass pass:essai -salt -base6
	o := NewOpenSSL()
	content, err := o.DecryptFile(pass, path)
	if err != nil {
		return errors.New(fmt.Sprintf("Bad password (%v)", err))
	}

	if err := json.Unmarshal(content, &accounts); err != nil {
		return err
	}

	return nil
}

func NewAccounts(path string, password string) {
	accounts = make(Accounts, 0)
}

func (acc Accounts) SaveAccounts(path string, password string) error {
	// the content has been generated with:
	//  openssl enc -aes-128-cbc -pass pass:essai -salt -base6
	b, _ := json.Marshal(acc)
	o := NewOpenSSL()
	return o.EncryptFile(password, string(b), path)
}

func (acc Accounts) DumpAccounts(path string) error {
	b, err := json.Marshal(acc)
	if err != nil {
		return err
	}

	var bindent bytes.Buffer
	if err := json.Indent(&bindent, b, "", "  "); err != nil {
		return err
	}

	return ioutil.WriteFile(path, bindent.Bytes(), 0644)
}

func (acc Accounts) LoadJson(jsonPath string) (int, error) {
	file, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return 0, err
	}
	var newAccounts Accounts
	if err := json.Unmarshal(file, &newAccounts); err != nil {
		return 0, err
	}
	accounts = append(accounts, newAccounts...)
	return len(newAccounts), nil
}

/* ======= find in accounts */

func (accs Accounts) FindIn(field string, s string) {
	switch field {
	case "name":
		accs.FindFunc((*Account).FindName, s)
	case "email":
		accs.FindFunc((*Account).FindEmail, s)
	case "pseudo":
		accs.FindFunc((*Account).FindPseudo, s)
	case "notes":
		accs.FindFunc((*Account).FindNote, s)
	default:
		accs.FindFunc((*Account).FindAny, s)
	}
}

func (accs Accounts) FindFunc(f func(acc *Account, s string) bool, s string) {
	matches.Clear()
	for idx, acc := range accs {
		if f(&acc, s) {
			matches.Append(idx)
		}
	}
}

func (accs Accounts) FindOne(s string) (int, error) {
	result := -1
	err := errors.New("Ambiguous account.")

	for idx, acc := range accs {
		if (&acc).FindName(s) {
			if result < 0 {
				result = idx
			} else {
				return -1, err
			}
		}
	}

	if result >= 0 {
		return result, nil
	}
	return -1, err
}

func (accs Accounts) ListEmpty(s string) {
	matches.Clear()
	for idx, acc := range accs {
		if (&acc).IsEmpty(s) {
			matches.Append(idx)
		}
	}
}

/* ======= get prop */

func (acc *Account) GetProp(field string) (string, error) {
	switch field {
	case "name":
		return (*acc).Name, nil
	case "email":
		return (*acc).Email, nil
	case "pseudo":
		return (*acc).Pseudo, nil
	case "password":
		return (*acc).Password, nil
	case "pass":
		return (*acc).Password, nil
	case "notes":
		return (*acc).Notes, nil
	default:
		return "", errors.New("unknown field " + field)
	}
}

/* ======= find in account */

func (acc *Account) FindAny(s string) bool {
	return acc.FindName(s) || acc.FindPseudo(s) || acc.FindEmail(s) || acc.FindNote(s)
}

func (acc *Account) FindName(s string) bool {
	return strings.Contains(strings.ToLower(acc.Name), s)
}

func (acc *Account) FindPseudo(s string) bool {
	return strings.Contains(strings.ToLower(acc.Pseudo), s)
}

func (acc *Account) FindEmail(s string) bool {
	return strings.Contains(strings.ToLower(acc.Email), s)
}

func (acc *Account) FindNote(s string) bool {
	return strings.Contains(strings.ToLower(acc.Notes), s)
}

func (acc *Account) IsEmpty(field string) bool {
	value, err := acc.GetProp(field)
	return err == nil && value == ""
}

/* ======= sanitize account */
func (acc *Account) Sanitize() {
	// TODO
	acc.Name = strings.TrimSpace(acc.Name)
	acc.Pseudo = strings.TrimSpace(acc.Pseudo)
	acc.Email = strings.TrimSpace(acc.Email)
	acc.Notes = strings.TrimSpace(acc.Notes)
}
