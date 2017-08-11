package main

import (
	"encoding/json"
	"strings"
	"errors"
)


type Account struct {
	Name     string `json:"name"`
	Pseudo   string `json:"pseudo"`
	Email 	 string `json:"email"`
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
	content, err := o.DecryptFile(pass, path); 
	if err != nil {
		return err 
	}

	if err := json.Unmarshal(content, &accounts); err != nil {
		return err
	}

	return nil	
}

func NewAccounts(path string, password string) Accounts {
	accounts := make(Accounts, 0)
	return accounts
}

func (acc Accounts) SaveAccounts(path string, password string) error {
		// the content has been generated with:
	//  openssl enc -aes-128-cbc -pass pass:essai -salt -base6
	b, _ := json.Marshal(acc)
	o := NewOpenSSL()
	return o.EncryptFile(password, path, string(b))
}

/* ======= find in accounts */

func (accs Accounts) FindFunc(f func(acc Account, s string) bool, s string ) {
	matches.Clear()
	for idx, acc := range(accs) {
		if f(acc, s) {
			matches.Append(idx)
		}
	}
}

func (accs Accounts) FindOne(s string) (int, error) {
	result := -1

	for idx, acc := range(accs) {
		if acc.FindName(s) {
			if result < 0 {
				result = idx
			}else{
				return -1, errors.New("Ambiguous account")
			}
		}
	}
	return result, nil
}


/* ======= find in account */

func (acc Account) FindAny(s string) bool {
	return acc.FindName(s) || acc.FindPseudo(s) || acc.FindEmail(s) || acc.FindNotes(s)
}

func (acc Account) FindName(s string) bool {
	return strings.Contains(strings.ToLower(acc.Name), s) 
}

func (acc Account) FindPseudo(s string) bool {
	return strings.Contains(strings.ToLower(acc.Pseudo), s) 
}

func (acc Account) FindEmail(s string) bool {
	return strings.Contains(strings.ToLower(acc.Email), s) 
}

func (acc Account) FindNotes(s string) bool {
	return strings.Contains(strings.ToLower(acc.Notes), s) 
}



