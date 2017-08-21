package main

import (
	"bytes"
	"fmt"
	"github.com/derlin-easypass/easycmd-go/account"
	"github.com/derlin/ishell" //"github.com/abiosoft/ishell"
	"os"
	"strconv"
	"strings"
)

var accounts account.Accounts

var selectedAccount int = -1
var lastMatch int = -1

const defaultPrompt = "\x1b[33measypass>\x1b[0m "

var shell *ishell.Shell

func main() {

	ReadCmd()
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell = ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "edit",
		Help: "edit account",
		Func: editAccount,
	})

	addCmd := &ishell.Cmd{
		Name: "add",
		Help: "add a new account",
		Func: addAccount,
	}
	addCmd.AddCmd(&ishell.Cmd{
		Name: "import",
		Help: "add accounts from a json file",
		Func: importAccounts,
	})
	shell.AddCmd(addCmd)

	shell.AddCmd(&ishell.Cmd{
		Name: "delete",
		Help: "delete an account",
		Func: deleteAccount,
	})

	copyCmd := &ishell.Cmd{
		Name: "copy",
		Help: "copy a property to the clipboard",
	}
	for _, field := range []string{"name", "pseudo", "email", "pass", "notes"} {
		createFunc := func(field string) func(c *ishell.Context) {
			return func(c *ishell.Context) { copyProp(c, field) }
		}
		copyCmd.AddCmd(&ishell.Cmd{
			Name: field,
			Help: "copy " + field + " to clipboard",
			Func: createFunc(field),
		})
	}
	shell.AddCmd(copyCmd)

	shell.AddCmd(&ishell.Cmd{
		Name: "pass",
		Help: "copy pass",
		Func: func(c *ishell.Context) { copyProp(c, "password") },
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "dump",
		Help: "dump to file. ",
		Func: exportAccounts,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "select",
		Help: "select an account given an index. ",
		Func: selectAccount,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "unselect",
		Help: "unselect the account. ",
		Func: func(*ishell.Context) { unselect() },
	})

	findCmd := &ishell.Cmd{
		Name: "find",
		Help: "list accounts; type strings to filter",
		Func: findAny,
	}
	for _, field := range []string{"name", "pseudo", "email", "notes"} {
		createFunc := func(field string) func(c *ishell.Context) {
			return func(c *ishell.Context) { find(c, field) }
		}
		findCmd.AddCmd(&ishell.Cmd{
			Name: field,
			Help: "find in " + field + " only",
			Func: createFunc(field),
		})
	}
	shell.AddCmd(findCmd)

	showCmd := &ishell.Cmd{
		Name: "show",
		Help: "show details about an account",
		Func: showDetails,
	}
	showCmd.AddCmd(&ishell.Cmd{
		Name: "pass",
		Help: "show pass",
		Func: showPass,
	})
	shell.AddCmd(showCmd)

	listCmd := &ishell.Cmd{
		Name: "list",
		Help: "list accounts",
		Func: list,
	}
	listEmptyCmd := &ishell.Cmd{
		Name: "empty",
		Help: "list accounts with empty property",
	}
	for _, field := range []string{"name", "email", "notes", "pseudo", "pass"} {
		createFunc := func(field string) func(c *ishell.Context) {
			return func(c *ishell.Context) { findEmpty(c, field) }
		}
		listEmptyCmd.AddCmd(&ishell.Cmd{
			Name: field,
			Help: "list accounts with an empty " + field,
			Func: createFunc(field),
		})
	}
	listCmd.AddCmd(listEmptyCmd)
	shell.AddCmd(listCmd)

	shell.NotFound(notFound)

	load()
	// run shell
	// when started with "exit" as first argument, assume non-interactive execution
	if commands != "" {
		shell.Process(commands)
	} else {
		// start shell
		shell.SetPrompt(defaultPrompt) // yellow prompt
		shell.Run()
	}
}

func load() {
	if creds.Password == "" {
		shell.Print("Password: ")
		creds.Password = shell.ReadPassword()
		if creds.Password == "" {
			shell.Println("Empty password.")
			os.Exit(0)
		}
	}

	if isNewSession {
		// if a json file is specified, try to load accounts from it
		if jsonPath != "" {
			var err error
			var cnt int
			if accounts, err = account.Import(jsonPath); err != nil {
				shell.Println(err)
			} else {
				// save with the json data
				accounts.Save(creds)
				shell.Printf("loaded %d accounts.\n", cnt)
			}
		} else {
			// create a new session
			accounts = make(account.Accounts, 0)
		}

	} else {
		var err error
		if accounts, err = account.LoadAccounts(creds); err != nil {
			shell.Println(err)
			os.Exit(0)
		}
	}
}

// ========= utils

func printAll() {
	shell.Println(len(accounts), " accounts.")
	var buffer bytes.Buffer
	for idx, acc := range accounts {
		if acc != nil {
			buffer.WriteString(fmt.Sprintf("  [%d]: %s\n", idx, acc.Name))
		}
	}
	// use paged functionality if too much accounts
	if len(accounts) > 30 {
		shell.ShowPaged(buffer.String())
	} else {
		shell.Println(buffer.String())
	}
}

func printMatches(matches []int) {
	shell.Println(len(matches), " match(es).")
	var buffer bytes.Buffer
	for _, idx := range matches {
		buffer.WriteString(fmt.Sprintf("  [%d]: %s\n", idx, accounts[idx].Name))
	}
	// use paged functionality if too much accounts
	if len(matches) > 30 {
		shell.ShowPaged(buffer.String())
	} else {
		shell.Println(buffer.String())
	}
}

func accountFromHint(args []string) (acc *account.Account, idx int, err error) {
	if len(args) == 0 {
		if selectedAccount > -1 {
			acc = accounts[selectedAccount]
			idx = selectedAccount
		} else {
			err = fmt.Errorf("missing account info")
		}
		return
	}

	if idx, err = strconv.Atoi(args[0]); err == nil {
		if idx < len(accounts) && accounts[idx] != nil {
			acc = accounts[idx]
			selecct(idx)
		} else {
			err = fmt.Errorf("%d: wrong index.", idx)
		}
	} else if idx, err = accounts.FindOne(strings.Join(args, " ")); err == nil {
		acc = accounts[idx]
		selecct(idx)
	}

	return
}

func selecct(accountIdx int) {
	selectedAccount = accountIdx
	shell.SetPrompt("\x1b[33measypass [" + accounts[selectedAccount].Name + "]>\x1b[0m ")
}

func unselect() {
	selectedAccount = -1
	shell.SetPrompt(defaultPrompt)
}
