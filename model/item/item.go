package item

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

//Item is a model of Item in the store.
type Item struct {
	UUID        uuid.UUID `db:"uuid" schema:"uuid" json:"uuid,omitempty"`
	ItemName    string    `db:"item_name" schema:"item_name" json:"item_name,omitempty"`
	DisplayName string    `db:"display_name" schema:"display_name" json:"display_name,omitempty"`
	Price       float32   `db:"price" schema:"price" json:"price,omitempty"`
	Currency    int       `db:"currency" schema:"currency" json:"currency,omitempty"`
	ItemTypeID  int       `db:"item_type_id" schema:"item_type_id" json:"item_type_id,omitempty"`
	Available   int       `db:"available" schema:"available" json:"available,omitempty"`
	Description string    `db:"description" schema:"description" json:"description,omitempty"`
	ImagePath   string    `db:"image_path" schema:"image_path" json:"image_path,omitempty"`
}

//Type is a model for Item Type table
type Type struct {
	ID      int    `db:"id" schema:"id" json:"id"`
	Display string `db:"display"  schema:"display" json:"display"`
}

//Currency enums.
const (
	USD = iota + 1
	RUB
	EUR
)

//Create insert Item to db
func Create(db *sqlx.DB, c Item) (uuid.UUID, error) {
	c.UUID = uuid.New()
	queryStr := "INSERT INTO item (uuid, item_name, display_name, price, currency, " +
		"item_type_id, available, description, image_path)" +
		" VALUES(:uuid, :item_name, :display_name, :price, :currency, " +
		":item_type_id, :available, :description, :image_path)"
	_, err := db.NamedExec(queryStr, c)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Can't insert into db. %v", err)
	}
	return c.UUID, err
}

//Read returns all Items.
func Read(db *sqlx.DB) ([]Item, error) {
	Item := []Item{}
	queryStr := "SELECT * FROM item"
	stmt, _ := db.Preparex(queryStr)
	err := stmt.Select(&Item)
	return Item, err
}

//ReadUUID returns specified Item.
func ReadUUID(db *sqlx.DB, u uuid.UUID) (Item, error) {
	Item := Item{UUID: u}
	queryStr := `SELECT * FROM item i WHERE uuid = $1`
	err := db.Get(&Item, queryStr, Item.UUID)
	return Item, err
}

//Update updates the Item
func Update(db *sqlx.DB, c Item) error {
	queryStr := "UPDATE item SET " +
		"uuid = :uuid, item_name = :item_name, display_name = :display_name, " +
		" price = :price, currency = :currency, item_type_id = :item_type_id, " +
		"available = :available, description = :description, image_path = :image_path" +
		" WHERE uuid = :uuid"
	_, err := db.NamedExec(queryStr, c)
	return err
}

//Delete deletes the Item
func Delete(db *sqlx.DB, id uuid.UUID) error {
	queryStr := "DELETE FROM item WHERE uuid = $1"
	_, err := db.Exec(queryStr, id)
	return err
}

//GetItemTypes returns item types
func GetItemTypes(db *sqlx.DB) ([]Type, error) {
	var it []Type
	const queryStr = "SELECT id, display FROM item_type;"
	stmt, _ := db.Preparex(queryStr)
	err := stmt.Select(&it)
	return it, err
}

//GetItemTypeID returns ID of the related display row.
func GetItemTypeID(db *sqlx.DB, display string, dst interface{}) error {
	const queryStr = "SELECT id FROM public.item_type WHERE display = $1;"
	row := db.QueryRow(queryStr, display)
	err := row.Scan(dst)
	if err != nil {
		log.Printf("Can't read row, %v", err)
		return err
	}
	return err
}

//DeleteBundle deletes a bunch of rows from table. If length of passed array equals 1 it executes Delete function.
func DeleteBundle(db *sqlx.DB, id []uuid.UUID) error {
	queryStr := "DELETE FROM Item WHERE uuid = $1"
	if len(id) == 0 {
		return errors.New("Array is empty")
	}
	if len(id) == 1 {
		return Delete(db, id[0])
	}
	stmt, err := db.Beginx()
	if err != nil {
		return err
	}
	stmtPrep, err := stmt.Prepare(queryStr)
	if err != nil {
		return err
	}
	for _, item := range id {

		_, err := stmtPrep.Exec(item)
		if err != nil {
			return errors.New("Can't delete bundle of items. Reason: " + err.Error())
		}
	}
	err = stmt.Commit()
	if err != nil {
		return errors.New("transaction commit has failed. Reason: " + err.Error())
	}
	return nil
}

//CreateBundle creates a bunch of items.
func CreateBundle(db *sqlx.DB, objects []Item) error {
	queryStr := "INSERT INTO item (uuid, item_name, display_name, price, currency, " +
		"item_type_id, available, description, image_path)" +
		" VALUES(:uuid, :item_name, :display_name, :price, :currency, " +
		":item_type_id, :available, :description, :image_path)"

	if len(objects) == 0 {
		return errors.New("Array is empty")
	}
	if len(objects) == 1 {
		_, err := Create(db, objects[0])
		return err
	}
	stmt, err := db.Beginx()
	if err != nil {
		return err
	}
	stmtPrep, err := stmt.PrepareNamed(queryStr)
	if err != nil {
		return err
	}
	for _, item := range objects {
		item.UUID = uuid.New()
		_, err := stmtPrep.Exec(item)
		if err != nil {
			return errors.New("Can't insert bundle of items. Reason: " + err.Error())
		}
	}
	err = stmt.Commit()
	if err != nil {
		return errors.New("Transaction commit has been failed. Reason: " + err.Error())
	}
	return nil
}

//UpdateBundle updates a bunch of items
func UpdateBundle(db *sqlx.DB, objects []Item) error {
	queryStr := "UPDATE item SET " +
		"uuid = :uuid, item_name = :item_name, display_name = :display_name, " +
		" price = :price, currency = :currency, item_type_id = :item_type_id, " +
		"available = :available, description = :description, image_path = :image_path" +
		" WHERE uuid = :uuid"
	if len(objects) == 0 {
		return errors.New("Array is empty")
	}
	if len(objects) == 1 {
		err := Update(db, objects[0])
		return err
	}
	stmt, err := db.Beginx()
	if err != nil {
		return err
	}
	stmtPrep, err := stmt.Prepare(queryStr)
	if err != nil {
		return err
	}
	for _, item := range objects {
		_, err := stmtPrep.Exec(item)
		if err != nil {
			return errors.New("Can't update bundle of items. Reason: " + err.Error())
		}
	}
	err = stmt.Commit()
	if err != nil {
		if errRollback := stmt.Rollback(); errRollback != nil {
			return errors.New("Can't commit transaction and can't proceed a rollback. Reasons:" + errRollback.Error() + err.Error())
		}
		return errors.New("Transaction commit has been failed. Reason: " + err.Error())
	}
	return nil
}
