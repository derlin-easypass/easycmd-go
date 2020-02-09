package main

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/derlin/ishell" // TODO: keep track of changes in "github.com/abiosoft/ishell"
)

// ========= listing

func list(c *ishell.Context) {
	matches.Fill()
	matches.Print()
}

func listMatches(c *ishell.Context) {
	c.Println("matches: ")
	matches.Print()
}

func listEmpty(c *ishell.Context, field string) {
	accounts.ListEmpty(field)
	matches.Print()
}

// ========= finding

func find(c *ishell.Context, field string) {
	if len(c.Args) == 0 {
		list(c)

	} else {
		search := strings.Join(c.Args, " ")
		if field != "" {
			c.Println("searching for '" + search + "' in field " + field)
		} else {
			c.Println("searching for '" + search + "' in all fields")
		}
		accounts.FindIn(field, search)
		matches.Print()
	}
}

func findAny(c *ishell.Context) {
	find(c, "")
}

// ========= showing

func showDetails(c *ishell.Context) {
	acc, _, err := accountFromHint(c.Args)

	if err != nil {
		c.Println(err)
	} else {
		c.Println("  Name:   ", (*acc).Name)
		c.Println("  Pseudo: ", (*acc).Pseudo)
		c.Println("  Email:  ", (*acc).Email)
		c.Println("  Notes:  ", (*acc).Notes)
	}
}

func showPass(c *ishell.Context) {
	acc, _, err := accountFromHint(c.Args)

	if err != nil {
		c.Println(err)
	} else if (*acc).Password == "" {
		c.Println("empty password.")
	} else {
		c.Print("Pass: ", (*acc).Password)
		c.ReadLine()
		// put cursor one line up
		c.Print("\033[1A")
		// clear from the cursor to the end of the screen
		c.Print("\033[0J")
		c.Println()
	}
}

// ========= dumpimg

func dumpClear(c *ishell.Context) {
	var dumpPath string
	if len(c.Args) > 0 {
		dumpPath = c.Args[0]
	} else {
		c.Print("Output File: ")
		dumpPath = c.ReadLine()
	}

	if dumpPath == "" {
		c.Println("Missing mandatory output path.")
		return
	}

	if _, err := os.Stat(dumpPath); !os.IsNotExist(err) {
		c.Printf("File '%s' already exists. Override ? [y|n] ", dumpPath)
		if ok := c.ReadLine(); ok != "y" {
			c.Println("canceled.")
			return
		}
	}

	if err := accounts.DumpAccounts(dumpPath); err != nil {
		c.Println(err)
	} else {
		c.Println("dumped to ", dumpPath)
	}

}

// ========= manipulating one account: add edit delete

func addAccount(c *ishell.Context) {
	var acc Account
	c.Print("  Name: ")
	acc.Name = c.ReadLine()
	if name := strings.TrimSpace(acc.Name); name == "" {
		c.Println("empty name is not allowed")
		return
	}
	c.Print("  Pseudo: ")
	acc.Pseudo = c.ReadLine()
	c.Print("  Email: ")
	acc.Email = c.ReadLine()
	c.Print("  Password: ")
	acc.Password = c.ReadPassword()
	c.Print("  Notes: ")
	acc.Notes = c.ReadLine()

	c.Print("Saving ? [y|n]")
	if ok := c.ReadLine(); ok == "y" {
		(&acc).Sanitize()
		accounts = append(accounts, acc)
		err := accounts.SaveAccounts(sessionPath, password)
		if err == nil {
			matches.Clear()
			matches.Append(len(accounts) - 1)
			c.Println("saved.")
		} else {
			c.Println(err)
		}
	} else {
		c.Println("canceled.")
	}
}

func addJson(c *ishell.Context) {
	var loadPath string
	if len(c.Args) > 0 {
		loadPath = c.Args[0]
	} else {
		c.Print("Json File: ")
		loadPath = c.ReadLine()
	}

	if loadPath == "" {
		c.Println("Missing mandatory input path.")
		return
	}

	var err error
	var cnt int
	if cnt, err = accounts.LoadJson(loadPath); err != nil {
		shell.Println(err)
	} else {
		// save with the json data
		accounts.SaveAccounts(sessionPath, password)
		shell.Printf("loaded %d accounts.\n", cnt)
	}

}

func editAccount(c *ishell.Context) {
	acc, idx, err := accountFromHint(c.Args)
	if err != nil {
		c.Println(err)
		return
	}
	newAcc := Account{}
	shell.SetPrompt("  Name>")
	newAcc.Name = c.ReadLineWithDefault((*acc).Name)
	if name := strings.TrimSpace(newAcc.Name); name == "" {
		c.Println("empty name is not allowed")
		return
	}
	shell.SetPrompt("  Pseudo>")
	newAcc.Pseudo = c.ReadLineWithDefault((*acc).Pseudo)
	shell.SetPrompt("  Email>")
	newAcc.Email = c.ReadLineWithDefault((*acc).Email)
	shell.Print("  Password>")
	newAcc.Password = c.ReadPassword()
	shell.SetPrompt("  Notes>")
	newAcc.Notes = c.ReadLineWithDefault((*acc).Notes)
	shell.SetPrompt(defaultPrompt)

	c.Print("Saving ? [y|n]")
	if ok := c.ReadLine(); ok == "y" {
		(*acc).Name = newAcc.Name
		(*acc).Pseudo = newAcc.Pseudo
		(*acc).Email = newAcc.Email
		if newAcc.Password != "" {
			(*acc).Password = newAcc.Password
		}
		(*acc).Notes = newAcc.Notes
		(*acc).Sanitize()

		err := accounts.SaveAccounts(sessionPath, password)
		if err == nil {
			matches.Clear()
			matches.Append(idx)
			c.Println("saved")
		} else {
			c.Println(err)
		}
	} else {
		c.Println("canceled.")
	}
}

func deleteAccount(c *ishell.Context) {
	acc, idx, err := accountFromHint(c.Args)
	if err != nil {
		c.Println(err)
		return
	}

	c.Print("delete account '", (*acc).Name, "' at index ", idx, "? [y|n] ")
	if ok := c.ReadLine(); ok == "y" {
		accounts = append(accounts[:idx], accounts[idx+1:]...)
		err := accounts.SaveAccounts(sessionPath, password)
		if err == nil {
			matches.RemoveIdx(idx)
			c.Println("saved.")
		} else {
			c.Println(err)
		}
	} else {
		c.Println("canceled.")
	}
}

// ========= copying

func copyProp(c *ishell.Context, field string) {
	acc, _, err := accountFromHint(c.Args)
	if err != nil {
		c.Println(err)
		return
	}

	value, err := acc.GetProp(field)

	if err != nil {
		c.Println(err)
	} else if value == "" {
		c.Printf("empty %s.\n", field)
	} else {
		if err = clipboard.WriteAll(value); err != nil {
			c.Println(err)
		} else {
			c.Printf("Copied %s from '%s' to clipboard.\n", field, (*acc).Name)
		}
	}
}

// ========= utils

func notFound(c *ishell.Context) {
	c.Println("not command. Assuming find...")
	findAny(c)
}

func accountFromHint(args []string) (acc *Account, idx int, err error) {
	if len(args) == 0 {
		if matches.Length() == 1 {
			acc, idx, err = matches.AccountAt(0)
		} else {
			err = errors.New("missing account info")
		}
		return
	}

	if idx, err = strconv.Atoi(args[0]); err == nil {
		acc, idx, err = matches.AccountAt(idx)
	} else if idx, err = accounts.FindOne(strings.Join(args, " ")); err == nil {
		matches.Clear()
		matches.Append(idx)
		acc = &accounts[idx]
	}

	return

}
