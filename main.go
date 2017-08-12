package main

import (
	"strings"
	"github.com/abiosoft/ishell"
	"github.com/atotto/clipboard"
	"os"
	"strconv"
	"errors"
)


var shell *ishell.Shell

func main(){

    // create new shell.
    // by default, new shell includes 'exit', 'help' and 'clear' commands.
    shell = ishell.New()

    // display welcome info.
    shell.Println("EasyCmd GO")

    shell.AddCmd(&ishell.Cmd{
        Name: "greet",
        Help: "greet user",
        Func: func(c *ishell.Context) {
            c.Println("Hello", strings.Join(c.Args, " "))
        },
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "add",
        Help: "add a new account",
        Func: addAccount,
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "pass",
        Help: "copy pass",
        Func: copyPass,
    })

    findCmd := &ishell.Cmd{
        Name: "find",
        Help: "list accounts; type strings to filter",
        Func: findAny,
    }
    for _, field := range([]string{"name", "pseudo", "email", "notes"}) {
		findCmd.AddCmd(&ishell.Cmd{
			Name: field,
			Help: "find in " + field + " only",
			Func: func(c *ishell.Context) { find(c, field) },
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
    listCmd.AddCmd(&ishell.Cmd{
        Name: "matches",
        Help: "list matches only",
        Func: listMatches,
    })
    listEmptyCmd := &ishell.Cmd{
        Name: "empty",
        Help: "list accounts with empty property",
    }
	for _, field := range([]string{"name", "email", "notes", "pseudo", "pass"}) {
		listEmptyCmd.AddCmd(&ishell.Cmd{
			Name: field,
			Help: "list accounts with an empty " + field,
			Func: func(c *ishell.Context) { listEmpty(c, field) },
		})
	}
	listCmd.AddCmd(listEmptyCmd)
    shell.AddCmd(listCmd)


    shell.NotFound(notFound)

    load()
    // run shell
    shell.Run()
}


func load(){
	if password == "" {
		shell.Print("Password: ")
		password = shell.ReadPassword()
		if password == "" {
			shell.Println("Empty password.")
			os.Exit(0)
		}
	}

	if err := LoadAccounts(sessionPath, password); err != nil {
		shell.Println(err)
		os.Exit(0)
	}
}


func copyPass(c *ishell.Context) {
	acc, err := accountFromHint(c.Args)

	if err != nil {
		c.Println(err)
	}else if (*acc).Password == "" {
		c.Println("empty password.")
	}else{
		if err = clipboard.WriteAll((*acc).Password); err != nil {
			c.Println(err)
		}else{
			c.Println("Copied password from '", (*acc).Name, "' to clipboard.")
		}
	}
	
}

func addAccount(c *ishell.Context) {
	var acc Account
    c.Print("  Name: ")
    acc.Name = c.ReadLine()
    c.Print("  Pseudo: ")
    acc.Name = c.ReadLine()
    c.Print("  Password: ")
    acc.Password = c.ReadPassword()
    c.Println("Notes -- Input multiple lines and end with semicolon ';'.")
    acc.Notes = c.ReadMultiLines(";")
    c.Println()
    c.Println(acc)
}


func list(c *ishell.Context){
	c.Println("accounts: ")
	matches.Fill()
	matches.Print()
}

func listMatches(c *ishell.Context){
	c.Println("matches: ")
	matches.Print()
}

func listEmpty(c *ishell.Context, field string){
	accounts.ListEmpty(field)
	matches.Print()
}

func find(c *ishell.Context, field string){
	if len(c.Args) == 0 {
		list(c)

	}else {
		search := strings.Join(c.Args, " ")
		accounts.FindIn(field, search)
		matches.Print()
	}
}

func findAny(c *ishell.Context){
	find(c, "")
}

func showDetails(c *ishell.Context) {
	acc, err := accountFromHint(c.Args)

	if err != nil {
		c.Println(err)
	}else{
		c.Println("  Name:   ", (*acc).Name)
		c.Println("  Pseudo: ", (*acc).Pseudo)
		c.Println("  Email:  ", (*acc).Email)
		c.Println("  Notes:  ", (*acc).Notes)
	}
}

func notFound(c *ishell.Context) {
	c.Println("not command. Assuming find...")
	findAny(c)
}

func showPass(c *ishell.Context) {
	acc, err := accountFromHint(c.Args)

	if err != nil {
		c.Println(err)
	}else if (*acc).Password == "" {
		c.Println("empty password.")
	}else{
		c.Print("Pass: ", (*acc).Password)
		c.ReadLine()
		// put cursor one line up
		c.Print("\033[1A")
		// clear from the cursor to the end of the screen
		c.Print("\033[0J")
		c.Println()
	}
}


func accountFromHint(args []string) (acc *Account, err error) {
	if len(args) == 0 {
		if matches.Length() == 1 {
			acc, err = matches.AccountAt(0)
		}else{
			err = errors.New("missing account info")
		}
		return
	}

	var idx int

	if idx, err = strconv.Atoi(args[0]); err == nil {
	   acc, err = matches.AccountAt(idx)
	}else if idx, err = accounts.FindOne(strings.Join(args, " ")); err == nil {
		matches.Clear()
		matches.Append(idx)
		acc = &accounts[idx]
	}

	return 

}

