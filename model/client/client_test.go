package client

//TODO: Rework item to look like item_test.go
import (
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestClient(t *testing.T) {
	db, err := sqlx.Connect("postgres", "postgres://webForum_RW:reachme_123@localhost:5432/simplemvc?sslmode=disable")
	if err != nil {
		t.Fatal("Can't connect to db. Reason: ", err)
	}
	defer db.Close()

	newClient := Client{
		FirstName: "Alex",
		Email:     "atrushkin@outlook.com",
		Phone:     "972.989.4341",
		Address:   "5520 Somerset drive apt #202",
	}
	id, err := Create(db, newClient)
	if err != nil {
		t.Logf("[SUCCESS] Got an error, %v", err)
		newClient.LastName = "True"
		id, err = Create(db, newClient)
	} else {
		t.Error("Should be an error while inserting not full info")
		t.Fail()
	}
	if err != nil {
		t.Fatal("Can't create a row. Reason: ", err)
	}

	t.Logf("[SUCCESS] Created new row. UUID: %s", id.String())

	newClient.UUID = id
	if updatedClient, err := ReadUUID(db, id); err != nil {
		t.Log("Error selecting inserted line. Reasion: ", err)
		t.Fail()
	} else {
		t.Logf("[SUCCESS] Inserted line: %#v", updatedClient)
		if !reflect.DeepEqual(newClient, updatedClient) {
			t.Logf("Structures doesn't match.\nOrig: %#v\nUpdated:%#v", newClient, updatedClient)
			t.Fail()
		}
	}

	newClient.FirstName = "Kate"
	if err := Update(db, newClient); err != nil {
		t.Log("Error updating line, ", err)
		t.Fail()
	}
	if updatedClient, err := ReadUUID(db, id); err != nil {
		t.Errorf("Error read updated line, %v", err)
		t.Fail()
	} else {
		if updatedClient.FirstName != "Kate" {
			t.Errorf("Expected Kate, got: %s", updatedClient.FirstName)
			t.Fail()
		}
		t.Logf("[SUCCESS] Row has been updated.")
	}

	if readAllClient, err := Read(db); err != nil {
		t.Logf("Got error while reading whole table, %v", err)
		t.Fail()
	} else {
		t.Logf("[SUCCESS] Read all rows, total number: %d", len(readAllClient))
	}

	if err := Delete(db, id); err != nil {
		t.Errorf("Got an error deleting row, %v", err)
	}
	t.Log("[SUCCESS] Row has beem deleted")
}
