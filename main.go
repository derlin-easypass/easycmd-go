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

    shell.AddCmd(&ishell.Cmd{
        Name: "edit",
        Help: "edit account",
        Func: editAccount,
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "add",
        Help: "add a new account",
        Func: addAccount,
    })

     shell.AddCmd(&ishell.Cmd{
        Name: "delete",
        Help: "delete an account",
        Func: deleteAccount,
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
    	createFunc := func (field string) func(c *ishell.Context) {
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
		createFunc := func (field string) func(c *ishell.Context) {
			return func(c *ishell.Context) { listEmpty(c, field) }
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
	if len(remainingArgs) > 1 && remainingArgs[0] == "exit" {
	    shell.Process(remainingArgs[1:]...)
	} else {	
	    // start shell
	    shell.SetPrompt("\x1b[33measypass>\x1b[0m ") // yellow prompt
	    shell.Run()
	}
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

	if isNewSession {
		NewAccounts(sessionPath, password)
	}else{
		if err := LoadAccounts(sessionPath, password); err != nil {
			shell.Println(err)
			os.Exit(0)
		}
	}
}


func copyPass(c *ishell.Context) {
	acc, _, err := accountFromHint(c.Args)

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
	    	matches.Clear()
    		c.Println("saved.")
    	}else{
    		c.Println(err)
    	}
    }else{
    	c.Println("canceled.")
    }

}


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
    		matches.Append(len(accounts) -1)
    		c.Println("saved.")
    	}else{
    		c.Println(err)
    	}
    }else{
    	c.Println("canceled.")
    }
}

func editAccount(c *ishell.Context) {
	acc, idx, err := accountFromHint(c.Args)
	if err != nil {
		c.Println(err)
		return
	}
	newAcc := Account{}
    c.Println("  Name: ")
    newAcc.Name = c.ReadLineWithDefault((*acc).Name)
    if name := strings.TrimSpace(newAcc.Name); name == "" {
    	c.Println("empty name is not allowed")
    	return
    }
    c.Println("  Pseudo: ")
    newAcc.Pseudo = c.ReadLineWithDefault((*acc).Pseudo)
    c.Println("  Email: ")
    newAcc.Email = c.ReadLineWithDefault((*acc).Email)
    c.Print("  Password: ")
    newAcc.Password = c.ReadPassword()  
    c.Println("Notes")
    newAcc.Notes = c.ReadLineWithDefault((*acc).Notes)

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
    	}else{
    		c.Println(err)
    	}
    }else{
    	c.Println("canceled.")
    }
}

func list(c *ishell.Context){
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
		c.Println("searching for '" + search + "' in field " + field)
		accounts.FindIn(field, search)
		matches.Print()
	}
}

func findAny(c *ishell.Context){
	find(c, "")
}

func showDetails(c *ishell.Context) {
	acc, _, err := accountFromHint(c.Args)

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
	acc, _, err := accountFromHint(c.Args)

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


func accountFromHint(args []string) (acc *Account, idx int, err error) {
	if len(args) == 0 {
		if matches.Length() == 1 {
			acc, idx, err = matches.AccountAt(0)
		}else{
			err = errors.New("missing account info")
		}
		return
	}

	if idx, err = strconv.Atoi(args[0]); err == nil {
		acc, idx, err = matches.AccountAt(idx)
	}else if idx, err = accounts.FindOne(strings.Join(args, " ")); err == nil {
		matches.Clear()
		matches.Append(idx)
		acc = &accounts[idx]
	}

	return 

}

