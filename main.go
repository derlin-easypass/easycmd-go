package main

import (
	"github.com/abiosoft/ishell"
	"os"
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

    addCmd := &ishell.Cmd{
        Name: "add",
        Help: "add a new account",
        Func: addAccount,
    }
    addCmd.AddCmd(&ishell.Cmd{
    	Name: "json",
    	Help: "add accounts from a json file",
    	Func: addJson,
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
    for _, field := range([]string{"name", "pseudo", "email", "pass", "notes"}) {
    	createFunc := func (field string) func(c *ishell.Context) {
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
        Func: dumpClear,
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
		// create a new session
		NewAccounts(sessionPath, password)
		// if a json file is specified, try to load accounts from it

		if jsonPath != "" {
			var err error
			var cnt int
			if cnt, err = accounts.LoadJson(jsonPath); err != nil {
				shell.Println(err)
			}else{
				// save with the json data
				accounts.SaveAccounts(sessionPath, password)
				shell.Printf("loaded %d accounts.\n", cnt)
			}
		}

	}else{
		if err := LoadAccounts(sessionPath, password); err != nil {
			shell.Println(err)
			os.Exit(0)
		}
	}
}


