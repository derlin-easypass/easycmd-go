package main

import (
	"fmt"
	"encoding/json"
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

    shell.AddCmd(&ishell.Cmd{
        Name: "find",
        Help: "list accounts; type strings to filter",
        Func: find,
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "show",
        Help: "show details about an account",
        Func: showDetails,
    })

    shell.AddCmd(&ishell.Cmd{
        Name: "showpass",
        Help: "show pass",
        Func: showPass,
    })

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
    acc.Name = c.ReadLineWithDefault("lucy")
    c.Print("  Pseudo: ")
    acc.Name = c.ReadLine()
    c.Print("  Password: ")
    acc.Password = c.ReadPassword()
    c.Println("Notes -- Input multiple lines and end with semicolon ';'.")
    acc.Notes = c.ReadMultiLines(";")
    c.Println()
    c.Println(acc)
}

func find(c *ishell.Context){
	if len(c.Args) == 0 {
		c.Println("show all")
		matches.Fill()
		matches.Print()
	}else {
		search := strings.Join(c.Args, " ")
		c.Println("searching " + search)
		accounts.FindFunc(Account.FindAny, search)
		matches.Print()
	}
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
		err = errors.New("missing account info")
		return
	}

	var idx int

	if idx, err = strconv.Atoi(args[0]); err == nil {
	   acc, err = matches.AccountAt(idx)
	}else if idx, err = accounts.FindOne(strings.Join(args, " ")); err == nil {
		acc = &accounts[idx]
	}

	return 

}

func lala_main() {
	// the content has been generated with:
	//  openssl enc -aes-128-cbc -pass pass:essai -salt -base6
	o := NewOpenSSL()
	result, err := o.DecryptString(`essai`, `U2FsdGVkX1/ktngEqYYCvmt7hpSu0x95/ZikdCf8684Bnt3xDulKlri4xAaNjwd5fNYvmyoayPIl
xT+u8ISTBmCwf9vqsTsl9aSZT2ewACqNQA8IUwb8ZK6B9iQ+CVjxVU21Td6aFRrGO7BNEbv+nChI
T1JxX952iamiSUjwUwMsBezDz/zEAgGTgY1fhzzhLmUMpM1gqdbkQksWBdaTTfHYSnIPjHWhPHyb
WPjDp/o9y5WXZzf79LClu2nzg80JkRAmuou0PG6Uaka2jqjOKeab3qLKeIP4iztZ4/CYv1K9gxSV
zEW/fSnFwZJVhu6tdQSESAiO77hJ+tPfddUk7qVDGzkCmvVJmUgCkQGjz/5ZuQ3iK7218oYbH26J
3821vu0UNkikuNA0/J4ACg9SH+un+aoRVxs4RT2ZH9Jt4NYlIC86zrP9ePgLaP7WKQ8sdSNAJ4Bd
EQyDI39dhpwYd6Bx6G+d6rkZKeNyJCsXy3Y3u0U8VYihSr5pqZf7yJZFgemotEGxlht4RTU3He6C
6VcFE7c8t7+xLrgyPQK2j49yKjAGAEJq2sInxz6j85rTLOc85qhws3rf9C84kbrE9YJbsNE3l4Uq
9UsxcJhvdmqOMO8hQcbjyFV/uCmkA8P6YbILm61VnvCz4P9XnkiDxTAi0SFlQKLxKOqBnpuz3ZXZ
llKkwamQwdyL3vn+hRJn9QICLoc7gAbnykD6107WxxyvMLYWXE+u/FCV0FxHy0mBRSYtkXKTnxaO
/h6GZ8SqmuZib6Mqw+78vQI/mwc2//UJyw/Ju26xVbUWp6roJmRPRdDPAOfO/Rtwcs4SUk1br1g6
h6cmr22ZiLDo4khRNBMaVL5qxKWqGYZ4OCS9FkfZfACDg4nc0xD6k/BaX9SUb3kUKbT4W26ZCuy7
cPnqeS4gUPIFYx2Fn7zIjt+xcd89ZK8dneT44jCeckcrsRO/n94dt/48ht0M3DAWwUJD7zh7kuHh
f4K4T18022eAbn7acY8dHi3npVN+H2EvHzcJrfcqt2Ga0Y67Ed5nN+CGfX9OiIAUwP6W5OaaJrw+
PXeolZXyaMo/T/Q49XlUY3nyV9dsVqPCRfxNa8XCWYKLSFLrhmaNybK3bcxi7ChSff4ig+38LXZu
tXodufiQMc4OMNH7Mh7Pjm/CJxuFJIp0glzL0lE0F8QWs9mKaPqfRSbEkGmRZlHkboVSfnE4zcki
m6rZn/t0sTdH6hekwg3lGNQyhP9WMHjldOkbJon+N/HrFqKVAEHYnJ+At9Bq82mfujsTlY7+3n8a
nsybSUlvymU87/TpS1Rw+d18wgfkW6GnAB55aY6+iizVzXZ16oOFRB3+gIpEY4pZyok2iVfsTGuL
FORVXN0p/5WO76GkLKhjCOphzIebPIPM1Mpkmq6ffrmW3gvRbaLb4Kj7Xknv97iu0rmAeowpgsG/
DnP0gEA5ZbUab8mxukbLfCHvUwCVBfAPlOaeXF/vQ4EJXlcWny1vWhiLkqFsFV2p5jodnDmYezI8
CjdbeJ3hUkCfJFAlmDCiG40Niy/kKc0j1s2AXOq2f6fo8DWCdQSqbqu2Uhtk1Lv04lbz2JqQIYlY
9D4a7dQbGOlWDI5FPMA2CRXjxg==`)
	fmt.Println(err)

	var accounts []Account
	if err := json.Unmarshal(result, &accounts); err != nil {
		fmt.Println("error")
		fmt.Printf("%v", err)
	}
	fmt.Printf("accounts are %v\n", accounts)
	fmt.Printf("Decrypted string is: %s", result)

}
