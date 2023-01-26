package files

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

type testCase struct {
	path    string
	storage *LocalStorage
	t       *testing.T
}

func newTestCase(t *testing.T) *testCase {
	p := "test"

	storage := NewLocalStorage(p)
	return &testCase{
		storage: storage,
		path:    p,
	}
}

// run creates test env, run the test and clean env
func (tc *testCase) run(f func()) {
	tc.up()
	f()
	tc.down()
}

// up creates base dir and notes.txt
func (tc *testCase) up() {
	// Delete storage dir if exits
	os.RemoveAll(tc.path)

	// Create base dir
	err := os.Mkdir(tc.path, 0775)
	if err != nil {
		tc.t.Fatal(err)
	}

	// Create notes.txt
	f := "notes.txt"
	content := []byte{0x21, 0xaa, 0xf8}
	b := bytes.NewBuffer(content)

	err = tc.storage.Save(f, b)
	if err != nil {
		tc.t.Errorf("Unable to save file (%s)", err.Error())
	}
}

// down delete all files and base directory
func (tc *testCase) down() {
	err := os.RemoveAll(tc.path)
	if err != nil {
		tc.t.Fatal(err)
	}
}

func TestSave(t *testing.T) {
	tc := newTestCase(t)

	t.Run("Save", func(t *testing.T) {
		tc.run(func() {
			// Save file
			f := "file.txt"
			content := []byte{0x21, 0xaa, 0xf8}
			b := bytes.NewBuffer(content)

			err := tc.storage.Save(f, b)
			if err != nil {
				t.Errorf("Unable to save file (%s)", err.Error())
			}

			// Read file and compare content
			fpath := fmt.Sprintf("%s/%s", tc.path, f)
			d, err := os.ReadFile(fpath)
			if err != nil {
				t.Errorf("Unable to read new file (%s)", err.Error())
			}

			if bytes.Compare(content, d) != 0 {
				t.Errorf("Wrong saving")
			}
		})
	})

	t.Run("Save overwriting", func(t *testing.T) {
		tc.run(func() {
			f := "notes.txt"

			// Read existing file data
			fpath := fmt.Sprintf("%s/%s", tc.path, f)
			oldData, err := os.ReadFile(fpath)
			if err != nil {
				t.Errorf("Unable to read existing file (%s)", err.Error())
			}

			// Overwrite file
			content := []byte{0x33, 0x49, 0xc2}
			b := bytes.NewBuffer(content)

			err = tc.storage.Save(f, b)
			if err != nil {
				t.Errorf("Unable to save file (%s)", err.Error())
			}

			// Read file content after overwrite
			data, err := os.ReadFile(fpath)
			if err != nil {
				t.Errorf("Unable to read new file (%s)", err.Error())
			}

			t.Log(data, oldData)

			if bytes.Compare(oldData, data) == 0 {
				t.Errorf("Wrong overwriting")
			}

		})
	})
}

func TestWrite(t *testing.T) {
	tc := newTestCase(t)

	t.Run("Write", func(t *testing.T) {
		tc.run(func() {
			b := bytes.NewBuffer([]byte{})

			fpath := fmt.Sprintf("%s/%s", tc.path, "notes.txt")
			sf, err := tc.storage.Write(fpath, b)
			if err != nil {
				t.Errorf("Unable to write file (%s)", err.Error())
			}

			if sf.Size != 3 {
				t.Errorf("Wrong writing, invalid size")
			}
		})
	})

	t.Run("Write non existing file", func(t *testing.T) {
		tc.run(func() {
			b := bytes.NewBuffer([]byte{})

			_, err := tc.storage.Write("films.txt", b)
			if err == nil {
				t.Errorf("Unexpected reading")
			}
		})
	})
}

func TestGetAllInfo(t *testing.T) {
	tc := newTestCase(t)

	t.Run("GetAllInfo", func(t *testing.T) {
		tc.run(func() {
			sf, err := tc.storage.GetAllInfo()
			if err != nil {
				t.Errorf("Unable to get information of files (%s)", err.Error())
			}

			if len(*sf) != 1 {
				t.Errorf("Wrong number of files")
			}
		})
	})
}
