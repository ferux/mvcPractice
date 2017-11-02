package order

import (
	"time"

	"github.com/google/uuid"
	// "reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestOrder(t *testing.T) {
	db, err := sqlx.Connect("postgres", "postgres://webForum_RW:reachme_123@localhost:5432/simplemvc?sslmode=disable")
	if err != nil {
		t.Fatal("[FATAL] Can't connect to db. Reason: ", err)
	}
	defer db.Close()

	var length int

	if readAllOrder, err := Read(db); err != nil {
		t.Logf("[LOG] Got error while reading whole table, %v", err)
		t.Fail()
	} else {
		t.Logf("[SUCCESS] Read all rows, total number: %d", len(readAllOrder))
		length = len(readAllOrder)
	}

	newOrder := Order{}
	id, err := Create(db, newOrder)
	if err != nil {
		t.Logf("[SUCCESS] Got an error, %v", err)
		newOrder.CartUUID = uuid.New()
		newOrder.ClientUUID = uuid.New()
		newOrder.Date = time.Now()

		newOrder.StatusDate = time.Now()
		id, err = Create(db, newOrder)
	} else {
		t.Error("[ERROR] Should be an error while inserting not full info")
		t.Fail()
		err = nil
	}
	if err != nil {
		t.Fatal("[FATAL] Can't create a row. Reason: ", err)
	} else {
		t.Logf("[SUCCESS] Created new row. UUID: %s", id.String())
	}
	newOrder.UUID = id
	if _, err := ReadUUID(db, id); err != nil {
		t.Log("[ERROR] Can't select inserted line. Reason: ", err)
		t.Fail()
	} else {
		t.Log("[SUCCESS] Inserted line")
		// if !reflect.DeepEqual(newOrder, updatedOrder) {
		// 	t.Errorf("[ERROR] Structures doesn't match.\nOrig:\n%#v\nUpdated:\n%#v", newOrder, updatedOrder)
		// 	t.Fail()
		// }
	}

	newOrder.IsPayed = true
	if err := Update(db, newOrder); err != nil {
		t.Errorf("[ERROR] Can't update line. Reason: %v ", err)
		t.Fail()
	}
	if _, err := ReadUUID(db, id); err != nil {
		t.Errorf("[ERROR] Can't read updated line. Reason %v", err)
		t.Fail()
	} else {
		// if !reflect.DeepEqual(newOrder, updatedOrder) {
		// 	t.Errorf("[ERROR] Structures doesn't match.\nOrig:\n%#v\nUpdated:\n%#v", newOrder, updatedOrder)
		// 	t.Fail()
		// } else {
		t.Logf("[SUCCESS] Row has been updated.")
		// }
	}

	if err := Delete(db, id); err != nil {
		t.Errorf("Got an error deleting row, %v", err)
	}
	t.Log("[SUCCESS] Row has beem deleted")

	if readAllOrder, err := Read(db); err != nil {
		t.Errorf("[ERROR] Can't read whole table, %v", err)
		t.Fail()
	} else {
		t.Logf("[SUCCESS] Read all rows, total number: %d", len(readAllOrder))
		if len(readAllOrder) != length {
			t.Errorf("[ERROR] Number of rows after the tests should be the same.\nExpected: %d", length)
		}
	}
}
