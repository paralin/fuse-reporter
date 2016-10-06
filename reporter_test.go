package reporter

import (
	"github.com/fuserobotics/statestream"
	"io/ioutil"
	"testing"
)

func TestSimple(t *testing.T) {
	tempPath, err := CreateTemporaryFile()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("Using temporary db: %s", tempPath)
	{
		r, err := NewReporter(tempPath)
		if err != nil {
			t.Fatal(err.Error())
		}
		cmp, err := r.CreateComponentIfNotExists("test")
		if err != nil {
			t.Fatal(err.Error())
		}
		_, err = cmp.GetState("state")
		if err == nil {
			t.Fatal("Non-existent state should return error.")
		}
		state, err := cmp.CreateStateIfNotExists("state", stream.DefaultStreamConfig())
		if err != nil {
			t.Fatal(err.Error())
		}
		cmpb, err := r.GetComponent("test")
		if err != nil {
			t.Fatal(err.Error())
		}
		if cmpb != cmp {
			t.Fatal("GetComponent should retreive cached instance from the map.")
		}
		stateb, err := cmpb.GetState("state")
		if err != nil {
			t.Fatal(err.Error())
		}
		if stateb != state {
			t.Fatal("GetState should retreive cached instance from the map.")
		}
		r.Close()
	}
	{
		r, err := NewReporter(tempPath)
		if err != nil {
			t.Fatal(err.Error())
		}
		cmp, err := r.GetComponent("test")
		if err != nil {
			t.Fatal(err.Error())
		}
		_, err = cmp.GetState("state")
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func CreateTemporaryFile() (string, error) {
	tf, err := ioutil.TempFile("", "reporter-test")
	if err != nil {
		return "", err
	}
	defer tf.Close()
	return tf.Name(), nil
}
