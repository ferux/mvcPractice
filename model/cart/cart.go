package cart

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

//Cart is a model of Cart in the store.
type Cart struct {
	UUID     uuid.UUID `db:"uuid" schema:"uuid"`
	ItemUUID uuid.UUID `db:"item_uuid" schema:"item_uuid"`
	CartUUID uuid.UUID `db:"cart_uuid" schema:"cart_uuid"`
}

//Create insert Cart to db
func Create(db *sqlx.DB, c Cart) (uuid.UUID, error) {
	c.UUID = uuid.New()
	var compare uuid.UUID
	if compare.String() == c.ItemUUID.String() || compare.String() == c.CartUUID.String() {
		return compare, errors.New("You must specify item_uuid or cart_uuid")
	}
	queryStr := "INSERT INTO cart (uuid, item_uuid, cart_uuid)" +
		" VALUES (:uuid, :item_uuid, :cart_uuid)"
	db.NamedExec(queryStr, c)
	return c.UUID, nil
}

//Read returns all Carts.
func Read(db *sqlx.DB) ([]Cart, error) {
	Cart := []Cart{}

	queryStr := "SELECT * FROM cart"
	stmt, _ := db.Preparex(queryStr)
	err := stmt.Select(&Cart)

	return Cart, err
}

//ReadUUID returns specified Cart.
func ReadUUID(db *sqlx.DB, u uuid.UUID) (Cart, error) {
	Cart := Cart{UUID: u}
	queryStr := "SELECT * FROM cart WHERE uuid = $1"
	err := db.Get(&Cart, queryStr, Cart.UUID)
	return Cart, err
}

//Update updates the Cart
func Update(db *sqlx.DB, c Cart) error {
	queryStr := "UPDATE cart SET item_uuid = :item_uuid, cart_uuid = :cart_uuid" +
		" WHERE uuid = :uuid"
	_, err := db.NamedExec(queryStr, c)
	return err
}

//Delete deletes the Cart
func Delete(db *sqlx.DB, id uuid.UUID) error {
	queryStr := "DELETE FROM cart WHERE uuid = $1"
	_, err := db.Exec(queryStr, id)
	return err
}

//DeleteBundle deletes a bunch of rows from table. If length of passed array equals 1 it executes Delete function.
func DeleteBundle(db *sqlx.DB, id []uuid.UUID) error {
	queryStr := "DELETE FROM Client WHERE uuid = $1"
	if len(id) == 0 {
		return errors.New("Array is empty")
	}
	if len(id) == 1 {
		return Delete(db, id[0])
	}
	stmt, err := db.Beginx()
	for _, item := range id {

		if err != nil {
			return errors.New("Can't delete bundle of items. Reason: " + err.Error())
		}
		stmt.Exec(queryStr, item)
	}
	err = stmt.Commit()
	if err != nil {
		return errors.New("transaction commit has failed. Reason: " + err.Error())
	}
	return nil
}
