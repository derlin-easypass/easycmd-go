package account

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func createAccounts() Accounts {
	accs := make(Accounts, 4)
	accs[0] = &Account{Name: "test1", Pseudo: "", Email: "lala@lulu.test"}
	accs[1] = &Account{Name: "test2", Pseudo: "test", Email: "", Password: "pass"}
	accs[2] = &Account{Name: "test3", Pseudo: "", Email: "lala@lala.com", Notes: "this is a note."}
	accs[3] = nil // ensure that nil doesn't throw an error
	return accs
}

func createDummyAccount() *Account {
	return &Account{Name: Name, Pseudo: Pseudo, Password: Password, Email: Email, Notes: Notes}
}

func TestLoad(t *testing.T) {
	creds := &Creds{"test.json.enc", "test"}
	accounts, err := LoadAccounts(creds)
	assert.Nil(t, err, "LoadAccounts failed with an error")
	assert.Equal(t, 3, len(accounts), "not 3 accounts")
	assert.Equal(t, "easypass", accounts[0].Pseudo, "sort not ok.")
	assert.Equal(t, "stackoverflow", accounts[2].Name, "sort not ok.")
}

func TestFindEmpty(t *testing.T) {
	accs := createAccounts()

	var results []int
	results = accs.FindEmpty(Name)
	assert.Equal(t, 0, len(results), "list empty failed.")
	results = accs.FindEmpty(Pseudo)
	assert.Equal(t, 2, len(results), "list empty failed.")
	results = accs.FindEmpty(Email)
	assert.Equal(t, 1, len(results), "list empty failed.")
}

func TestGetProp(t *testing.T) {
	acc := createDummyAccount()

	for _, f := range []string{Name, Pseudo, Password, Email, Notes} {
		v, _ := acc.GetProp(f)
		assert.Equal(t, f, v, fmt.Sprintf("get %s failed", f))
	}

	_, err := acc.GetProp("lala")
	assert.NotNil(t, err, "get prop lala returned a result")
}

func TestFindIn(t *testing.T) {
	accs := createAccounts()

	var results []int
	// find everywhere
	results = accs.Find("test")
	assert.Equal(t, 3, len(results), "find test everywhere failed.")
	// find in name
	results = accs.FindIn(Name, "dummy")
	assert.Equal(t, 0, len(results), "find dummy in name failed.")
	results = accs.FindIn(Name, "test")
	assert.Equal(t, 3, len(results), "find test in name failed.")
	// find in email
	results = accs.FindIn(Email, "@")
	assert.Equal(t, 2, len(results), "find @ in email failed.")
	// find in pseudo
	results = accs.FindIn(Pseudo, "t")
	assert.Equal(t, 1, len(results), "find t in pseudo failed.")
	// find in password
	results = accs.FindIn(Password, "p")
	assert.Equal(t, 1, len(results), "find p in password failed.")
	// find in notes
	results = accs.FindIn(Notes, "a note")
	assert.Equal(t, 1, len(results), "find 'a note' in notes failed.")
	results = accs.FindIn(Notes, "the note")
	assert.Equal(t, 0, len(results), "find 'the note' in notes failed.")

	// special cases
	results = accs.FindIn(Name, "")
	assert.Equal(t, 3, len(results), "find '' in name failed.")
	results = accs.Find("TEST")
	assert.Equal(t, 3, len(results), "uppercase should be ignored")
	results = accs.Find("TeST")
	assert.Equal(t, 3, len(results), "camelcase should be ignored")
	results = accs.Find(" TeST ")
	assert.Equal(t, 0, len(results), "spaces should make the find fail")
}

func TestImportExport(t *testing.T) {
	// create temp file
	tmpfile, err := ioutil.TempFile("", "accs-test")
	if err != nil {
		t.Error(err)
	}
	t.Logf("tmpfile is: %s\n", tmpfile.Name())
	defer os.Remove(tmpfile.Name()) // clean up

	// export
	accs := createAccounts()

	err = accs.Export(tmpfile.Name())
	assert.Nil(t, err, "export failed.")

	// import
	newAccs, err := Import(tmpfile.Name())
	assert.Nil(t, err, "import threw error.")
	assert.Equal(t, 3, len(newAccs), "import failed.") // null should not be serialized, so len = 3 and not 4

	for i, a := range newAccs {
		assert.Equal(t, accs[i].Name, a.Name, "import: values not equal")
	}
}

func TestSaveLoad(t *testing.T) {
	// create temp file
	tmpfile, err := ioutil.TempFile("", "accs-test")
	if err != nil {
		t.Error(err)
	}
	t.Logf("tmpfile is: %s\n", tmpfile.Name())
	defer os.Remove(tmpfile.Name()) // clean up

	creds := &Creds{Path: tmpfile.Name(), Password: "test-lala"}
	// save
	accs := createAccounts()

	err = accs.Save(creds)
	assert.Nil(t, err, "save failed.")

	// import
	newAccs, err := LoadAccounts(creds)
	assert.Nil(t, err, "import threw error.")
	assert.Equal(t, 3, len(newAccs), "load failed.") // null should not be serialized, so len = 3 and not 4

}
