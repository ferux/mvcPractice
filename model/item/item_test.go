package item

import (
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestItem(t *testing.T) {
	db, err := sqlx.Connect("postgres", "postgres://webForum_RW:reachme_123@localhost:5432/simplemvc?sslmode=disable")
	if err != nil {
		t.Fatal("[FATAL] Can't connect to db. Reason: ", err)
	}
	defer db.Close()

	var length int

	if readAllItem, err := Read(db); err != nil {
		t.Logf("[LOG] Got error while reading whole table, %v", err)
		t.Fail()
	} else {
		t.Logf("[SUCCESS] Read all rows, total number: %d", len(readAllItem))
		length = len(readAllItem)
	}

	newItem := Item{}
	id, err := Create(db, newItem)
	if err != nil {
		t.Logf("[SUCCESS] Got an error, %v", err)
		newItem.ItemName = "True"
		id, err = Create(db, newItem)
	} else {
		t.Error("[ERROR] Should be an error while inserting not full info")
		t.Fail()
	}
	if err != nil {
		t.Fatal("[FATAL] Can't create a row. Reason: ", err)
	} else {
		t.Logf("[SUCCESS] Created new row. UUID: %s", id.String())
	}
	newItem.UUID = id
	if updatedItem, err := ReadUUID(db, id); err != nil {
		t.Log("[ERROR] Can't select inserted line. Reason: ", err)
		t.Fail()
	} else {
		t.Log("[SUCCESS] Row has been inserted")
		if !reflect.DeepEqual(newItem, updatedItem) {
			t.Errorf("[ERROR] Structures doesn't match.\nOrig: %#v\nUpdated:%#v", newItem, updatedItem)
			t.Fail()
		}
	}

	newItem.ItemName = "Kate"
	if err := Update(db, newItem); err != nil {
		t.Errorf("[ERROR] Can't update line. Reason: %v ", err)
		t.Fail()
	}
	if updatedItem, err := ReadUUID(db, id); err != nil {
		t.Errorf("[ERROR] Can't read updated line. Reason %v", err)
		t.Fail()
	} else {
		if !reflect.DeepEqual(newItem, updatedItem) {
			t.Errorf("[ERROR] Structures doesn't match.\nOrig: %#v\nUpdated:%#v", newItem, updatedItem)
			t.Fail()
		} else {
			t.Logf("[SUCCESS] Row has been updated.")
		}
	}

	if err := Delete(db, id); err != nil {
		t.Errorf("Got an error deleting row, %v", err)
	}
	t.Log("[SUCCESS] Row has beem deleted")

	if readAllItem, err := Read(db); err != nil {
		t.Errorf("[ERROR] Can't read whole table, %v", err)
		t.Fail()
	} else {
		t.Logf("[SUCCESS] Read all rows, total number: %d", len(readAllItem))
		if len(readAllItem) != length {
			t.Errorf("[ERROR] Number of rows after the tests should be the same.\nExpected: %d", length)
		}
	}
}
