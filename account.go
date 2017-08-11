package main

import (
	"encoding/json"
)


type Account struct {
	Name     string `json:"name"`
	Pseudo   string `json:"pseudo"`
	Password string `json:"password"`
	Notes    string `json:"notes"`
}

type Accounts struct {
	Pass string
	Path string
	Accounts []Account
}

func LoadAccounts(path string, pass string) (*Accounts, error){
		// the content has been generated with:
	//  openssl enc -aes-128-cbc -pass pass:essai -salt -base6
	o := NewOpenSSL()
	content, err := o.DecryptFile(pass, path); 
	if err != nil {
		return nil, err 
	}

	var accounts []Account
	if err := json.Unmarshal(content, &accounts); err != nil {
		return nil, err
	}

	return &Accounts{ pass, path, accounts }, nil	
}

func NewAccounts(path string, password string) *Accounts {
	return &Accounts{ password, path, make([]Account, 0) }
}

func (acc *Accounts) SaveAccounts(path string, password string) error {
		// the content has been generated with:
	//  openssl enc -aes-128-cbc -pass pass:essai -salt -base6
	b, _ := json.Marshal(acc)
	o := NewOpenSSL()
	return o.EncryptFile(acc.Pass, acc.Path, string(b))
}