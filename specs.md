# Commands

* list
* find
* empty <field>
* duplicate <field>
* show
* copy
* add/edit/delete/undo
* (import/export)


## Immutable indexes

On load, the accounts are sorted by name and attributed an index (increasing). 
Then, this index is printed for commands with more than one output. 

``` 
easypass> find

	[0] alibaba
	[1] github.com
	[2] google.com
	...
``` 

It is possible to select one account using this number (`select`).

When an account is selected, it is possible to type something like `edit` or `copy`
and it will be applied directly to that account.

In the previous example, to show the entry `github.com` then copy the password, we can do:

```
easypass> show 1
easypass [github.com]> copy
``` 

or 

```
easypass> select 1
easypass [github.com]> show
easypass [github.com]> copy 
```

In case of deletion, the index becomes invalid, but there is no "shift"; the other accounts keep the
same indices for the whole session. Upon reload, the names will be sorted once again, so the indices 
will change !

### Impact of commands on the selection

* list, find, empty, duplicate: doesn't change the selection, but if it returns only one result,
	it will be stored in memory. It is then possible to type `select` without any argument to select
	the account.
* show, copy, edit, delete: without argument, it acts on the current selection. If an argument is
	given, it must be either an index or a regex with a unique output, which becomes the current
	selection.
* add: replaces the current selection by the new account.
* (import/export): does nothing on the selection.



