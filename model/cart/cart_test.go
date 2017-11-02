package cart

//TODO: Rework test to look like item_test.go
import (
	"github.com/google/uuid"
	// "reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB
var id uuid.UUID

func TestCart(t *testing.T) {
	var err error
	db, err = sqlx.Connect("postgres", "postgres://webForum_RW:reachme_123@localhost:5432/simplemvc?sslmode=disable")
	if err != nil {
		t.Fatal("[FATAL] Can't connect to db. Reason: ", err)
	}
	t.Run("Create", TestCreate)
	t.Run("Read", TestRead)
	t.Run("ReadUUID", TestReadUUID)
	t.Run("Update", TestUpdate)
	t.Run("Delete", TestDelete)
}

func TestCreate(t *testing.T) {
	newCart := Cart{}
	var err error
	id, err = Create(db, newCart)
	if err != nil {
		t.Logf("[SUCCESS] Got simulated error: %v", err)
		newCart.ItemUUID = uuid.New()
		newCart.CartUUID = uuid.New()
		id, err = Create(db, newCart)
	} else {
		t.Error("[ERROR] Should be an error while inserting not full info")
		t.Fail()
	}
	if err != nil {
		t.Fatal("[FATAL] Can't create a row. Reason: ", err)
	}

	t.Logf("[SUCCESS] Created new row. UUID: %s", id.String())
}

func TestRead(t *testing.T) {
	if readAllCart, err := Read(db); err != nil {
		t.Logf("[ERROR] Can't read whole table, %v", err)
		t.Fail()
	} else {
		t.Logf("[SUCCESS] Read all rows, total number: %d", len(readAllCart))
	}
}

func TestReadUUID(t *testing.T) {
	if _, err := ReadUUID(db, id); err != nil {
		t.Log("[ERROR] Can't select line. Reason: ", err)
		t.Fail()
	}
	t.Log("[SUCCESS] Line insert")
}

func TestUpdate(t *testing.T) {
	var updatedCart Cart
	updatedCart.UUID = id
	updatedCart.CartUUID = uuid.New()
	updatedCart.ItemUUID = uuid.New()
	if err := Update(db, updatedCart); err != nil {
		t.Log("[ERROR] Can't update line, ", err)
		t.Fail()
	}
	t.Log("[SUCCESS] Line update")
}

func TestDelete(t *testing.T) {
	if err := Delete(db, id); err != nil {
		t.Errorf("Got an error deleting row, %v", err)
	}
	t.Log("[SUCCESS] Row has beem deleted")
}

// func TestCart(t *testing.T) {
// 	db, err := sqlx.Connect("postgres", "postgres://webForum_RW:reachme_123@localhost:5432/simplemvc?sslmode=disable")
// 	if err != nil {
// 		t.Fatal("[FATAL] Can't connect to db. Reason: ", err)
// 	}
// 	defer db.Close()

// 	newCart := Cart{}
// 	id, err := Create(db, newCart)
// 	if err != nil {
// 		t.Logf("[SUCCESS] Got simulated an error, %v", err)
// 		newCart.ItemUUID = uuid.New()
// 		newCart.CartUUID = uuid.New()
// 		id, err = Create(db, newCart)

// 	} else {
// 		t.Error("[ERROR] Should be an error while inserting not full info")
// 		t.Fail()
// 	}
// 	if err != nil {
// 		t.Fatal("[FATAL] Can't create a row. Reason: ", err)
// 	}

// 	t.Logf("[SUCCESS] Created new row. UUID: %s", id.String())

// 	newCart.UUID = id

// 	if updatedCart, err := ReadUUID(db, id); err != nil {
// 		t.Log("[ERROR] Can't select line. Reasion: ", err)
// 		t.Fail()
// 	} else {
// 		t.Logf("[SUCCESS] Inserted line: %#v", updatedCart)
// 		if !reflect.DeepEqual(newCart, updatedCart) {
// 			t.Logf("[ERROR] Structures doesn't match.\nOrig: %#v\nUpdated:%#v", newCart, updatedCart)
// 			t.Fail()
// 		}
// 	}

// 	newCart.ItemUUID = uuid.New()
// 	if err := Update(db, newCart); err != nil {
// 		t.Log("[ERROR] Can't update line, ", err)
// 		t.Fail()
// 	}
// 	if updatedCart, err := ReadUUID(db, id); err != nil {
// 		t.Errorf("[ERROR] Can't read updated line, %v", err)
// 		t.Fail()
// 	} else {
// 		if updatedCart.ItemUUID != newCart.ItemUUID {
// 			t.Errorf("[ERROR] Expected %s, got: %s",newCart.ItemUUID.String(), updatedCart.ItemUUID.String())
// 			t.Fail()
// 		}
// 		t.Logf("[SUCCESS] Row has been updated.")
// 	}

// 	if readAllCart, err := Read(db); err != nil {
// 		t.Logf("[ERROR] Can't read whole table, %v", err)
// 		t.Fail()
// 	} else {
// 		t.Logf("[SUCCESS] Read all rows, total number: %d", len(readAllCart))
// 	}

// 	if err := Delete(db, id); err != nil {
// 		t.Errorf("Got an error deleting row, %v", err)
// 	}
// 	t.Log("[SUCCESS] Row has beem deleted")
// }
