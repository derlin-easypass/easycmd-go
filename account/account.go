package main

import (
	"encoding/json"
	"strings"
	"fmt"
	"io/ioutil"
	"bytes"
	"sort"
	"github.com/derlin-easypass/easycmd-go/crypto"
)

const (
	Name = "name"
	Pseudo = "pseudo"
	Email = "email"
	Notes = "notes"
	Password = "password"
)

type Account struct {
	Name     string `json:"name"`
	Pseudo   string `json:"pseudo"`
	Email 	 string `json:"email"`
	Password string `json:"password"`
	Notes    string `json:"notes"`
}

type Accounts []*Account


func (a Accounts) Len() int   		{ return len(a) }
func (a Accounts) Swap(i, j int)    { a[i], a[j] = a[j], a[i] }
func (a Accounts) Less(i, j int) bool  { return a[i].Name < a[j].Name }


/* =======load/save */

func LoadAccounts(path string, pass string) (Accounts, error) {
	var accounts Accounts
	content, err := crypto.DecryptFile(pass, path); 
	if err != nil {
		return nil, err 
	}

	if err := json.Unmarshal(content, &accounts); err != nil {
		return nil, err
	}
	sort.Sort(accounts)
	return accounts, nil	
}

func (accs Accounts) Save(path string, password string) error {
	bs, err := accs.toJson()
	if err == nil {
		return crypto.EncryptFile(password, string(bs), path)
	}
	return err
}

func Import(jsonPath string) (Accounts, error) {
	file, err := ioutil.ReadFile(jsonPath)
    if err != nil {
        return nil, err
    }
    var accounts Accounts
    if err := json.Unmarshal(file, &accounts); err != nil {
    	return nil, err
    }
    return accounts, nil
}

func (acc Accounts) Export(path string) error {
	b, err := acc.toJson()
	if err != nil {
		return err
	}
	var bindent bytes.Buffer 
	if err := json.Indent(&bindent, b, "", "  "); err != nil {
		return err
	}
	return ioutil.WriteFile(path, bindent.Bytes(),  0644)
}

func (accs Accounts) toJson() ([]byte, error){
	// remove empty accounts
	toMarshall := make(Accounts, 0)
	for _, v := range accs {
		if v != nil {
			toMarshall = append(toMarshall, v)
		}
	}
	// serialize
	return json.Marshal(toMarshall)
}


/* ======= find in accounts */

func (accs Accounts) Find(s string) []int {
	return accs.FindIn("", s)
}

func (accs Accounts) FindIn(field string, s string) []int {
	matches := make([]int, 0)
	for idx, acc := range(accs) {
		if acc != nil && acc.FindIn(field, s) {
			matches = append(matches, idx)
		}
	}
	return matches
}


func (accs Accounts) FindOne(s string) (int, error) {
	result := -1
	err := fmt.Errorf("Ambiguous account.")

	for idx, acc := range(accs) {
		if acc != nil && acc.FindIn(Name, s) {
			if result < 0 {
				result = idx
			}else{
				// more than one match. Throw an error
				return -1, err
			}
		}
	}
	if result >= 0 {
		// we found only one match
		return result, nil
	}
	// nothing found
	return -1, err
}

func (accs Accounts) FindEmpty(field string) []int {
	var results = make([]int, 0)
	for idx, acc := range(accs) {
		if acc != nil && acc.IsEmpty(field) {
			results = append(results, idx)
		}
	}
	return results
}

/* ======= get prop */

func (acc *Account) GetProp(field string) (string, error) {
	switch field {
		case Name: return acc.Name, nil
		case Email: return acc.Email, nil
		case Pseudo: return acc.Pseudo, nil
		case Password: return acc.Password, nil
		case "pass": return acc.Password, nil
		case Notes: return acc.Notes, nil
		default: return "", fmt.Errorf("unknown field '%s'", field) 
	}
}


/* ======= find in account */

func (acc *Account) Find(s string) bool {
	return acc.FindIn(Name, s) || acc.FindIn(Pseudo, s) || acc.FindIn(Email, s) || acc.FindIn(Notes, s)
}

func (acc *Account) FindIn(field string, search string) bool {

	if field == "" { // fall back
		return acc.Find(search)
	}
	
	if value, err := acc.GetProp(field); err == nil {
		return strings.Contains(strings.ToLower(value), strings.ToLower(search))
	}
	return false
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