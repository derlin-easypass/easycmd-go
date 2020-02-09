# Easycmd-go

Yet a new version of easypass commandline, this time written in Go !

## Installation

Clone the repository and (assuming you have a valid go setup), simply run `go build`:
```bash
git clone git@github.com:derlin-easypass/easycmd-go.git
cd easycmd-go
go build -o easycmd
```

## Getting started

**vocabulary**

* `account`: information about one "account", identified by its name and with properties such as a password, a username, an email address and a note.
* `session`: a group of one or more accounts, which is serialized into an encrypted JSON file

**usage**
```text
Usage of easycmd:
  -json string
        Combined with '-new', load session from a json file, which will be encrypted to 'path'.
  -new
        create a new session saved to 'path'
  -pass string
        password (not recommended !)
  -path string
        path to the session file
```

**Create a new session** 
use the `-path` argument to specify the new session path and `-new` to let the program knows it needs to create it:
```bash
easycmd -path <path to a non-existing file> -new
```

Here is an example:
```bash
easycmd  -path /tmp/example.ec -new 
Password: 
easypass> add
  Name: github.com
  Pseudo: derlin
  Email: derlin@example.com
  Password: 
  Notes: here I can describe what this account is about, or leave keywords
Saving ? [y|n]y
saved.
easypass> exit
```

**Open a session**
simply give the path to the session to easycmd:
```bash
easycmd -path <path to a session file>
```
Here is an example:
```bash
easycmd -path /tmp/example.ec     
Password: 
easypass> list
1 accounts.
  [0]: github.com
```

## Prompt

The easycmd prompt is a full-featured interactive prompt:
* arrows support
* history (limited to current prompt)
* built-in help (type `help` or `<command> help`)
* autocompletion with tab

How it works:
* each time you run `list` or `find`, it stores the list of results along with an index, you can then query the information of a match using the index
* in case `list`|`find` has only one result, the index may be omitted
* if a search term returns only one result, it can be used instead of an index
* if you don't remember what is the current result list, type `list matches`
* if what you type in the prompt is not a command, it is assumed to be search keywords: `[search terms]` is thus short for `find [search terms]`

## Prompt commands

The easiest way to discover it is to type `help` in the prompt:
```text
easypass> help

Commands:
  add         add a new account
  clear       clear the screen
  copy        copy a property to the clipboard
  delete      delete an account
  dump        dump to file. 
  edit        edit account
  exit        exit the program
  find        list accounts; type strings to filter
  help        display help
  list        list accounts
  pass        copy pass
  show        show details about an account
```

The most typical example usage is:
```text
easycmd -path example.ec
Password: 
easypass> github
not command. Assuming find...
searching for 'github' in all fields
3 accounts.
  [0]: docker.com
  [1]: github.com
  [2]: github.com master

easypass> pass 1
Copied password from 'github.com' to clipboard.
easypass> exit
```

## FAQ

**How can I backup my sessions?**

Well, it is up to you, since the sessions are *in fine* just regular (json encrypted) files. 
I personally save them to Dropbox.

**I don't have the program installed, how can I get my data back ?**

The Easypass suite is using the standard aes-128-cbc encryption. Hence, you can also decode your session file using tools 
such as OpenSSL:
```bash
openssl enc -aes-128-cbc -d -a -pass pass:YOUR_PASSWORD -in PATH/TO/YOUR_SESSION_FILE
```

**Can I import data from other providers ?**

Well, as long as you have data in a JSON-compliant format (or some programming skills to code a little converter), yes.
Use the `-json` argument in conjonction with `-new` to tell the program to load data from a JSON file, or directly 
create the encrypted file using OpenSSL:
```bash
openssl enc -aes-128-cbc -pass pass:YOUR_PASSWORD -salt -base64 -in PATH/TO/YOUR_JSON_FILE
```

**Can I export data to a clear JSON ?**

Yes, simply use the `dump [path/to/exported.json]` command in the easypass prompt.

**This project is silly, why reinvent the wheel ?**

Well, because programming is fun !

My seriously, I am of course using other services to store my passwords, but I love knowing that I also have one backup that I completely control.
Before my programming days, I was using a basic Excel file. When I started CS, I figured there should be a better way. Hence EasyPass !