package order

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

//Order is a model of Order in the store.
type Order struct {
	UUID       uuid.UUID `db:"uuid" schema:"uuid" json:"uuid"`
	ClientUUID uuid.UUID `db:"client_uuid" schema:"client_uuid" json:"client_uuid"`
	CartUUID   uuid.UUID `db:"cart_uuid" schema:"cart_uuid" json:"cart_uuid"`
	Date       time.Time `db:"date" schema:"date" json:"date"`
	IsPayed    bool      `db:"is_payed" schema:"is_payed" json:"is_payed"`
	Status     int       `db:"status" schema:"status" json:"status"`
	StatusDate time.Time `db:"status_date" schema:"status_date" json:"status_date"`
}

//Currency enums.
const (
	NEW = iota
	PROCESSING
	SENT
	DELIVERED
	CLOSED
)

//Create insert Order to db
func Create(db *sqlx.DB, c Order) (uuid.UUID, error) {
	c.UUID = uuid.New()
	queryStr := "INSERT INTO orders (uuid, client_uuid, cart_uuid, date, is_payed, status, status_date) " +
		"VALUES (:uuid, :client_uuid, :cart_uuid, :date, :is_payed, :status, :status_date)"
	var emptyUUID uuid.UUID
	sEmptyUUID := emptyUUID.String()
	if c.CartUUID.String() == sEmptyUUID || c.ClientUUID.String() == sEmptyUUID {
		return emptyUUID, errors.New("CartUUID and ClientUUID shouldn't be blank")
	}
	db.NamedExec(queryStr, c)

	return c.UUID, nil
}

//Read returns all Orders.
func Read(db *sqlx.DB) ([]Order, error) {
	Order := []Order{}

	queryStr := "SELECT * FROM orders"
	stmt, _ := db.Preparex(queryStr)
	err := stmt.Select(&Order)

	return Order, err
}

//ReadUUID returns specified Order.
func ReadUUID(db *sqlx.DB, u uuid.UUID) (Order, error) {
	Order := Order{UUID: u}
	queryStr := "SELECT * FROM orders WHERE uuid = $1"
	err := db.Get(&Order, queryStr, Order.UUID)
	return Order, err
}

//Update updates the Order
func Update(db *sqlx.DB, c Order) error {
	queryStr := "UPDATE orders SET " +
		"client_uuid = :client_uuid, cart_uuid = :cart_uuid, date = :date, is_payed = :is_payed, status = :status, status_date = :status_date" +
		" WHERE uuid = :uuid"
	_, err := db.NamedExec(queryStr, c)
	return err
}

//Delete deletes the Order
func Delete(db *sqlx.DB, id uuid.UUID) error {
	queryStr := "DELETE FROM orders WHERE uuid = $1"
	_, err := db.Exec(queryStr, id)
	return err
}

//DeleteBundle deletes a bunch of rows from table. If length of passed array equals 1 it executes Delete function.
func DeleteBundle(db *sqlx.DB, id []uuid.UUID) error {
	queryStr := "DELETE FROM Orders WHERE uuid = $1"
	if len(id) == 0 {
		return errors.New("Array is empty")
	}
	if len(id) == 1 {
		return Delete(db, id[0])
	}
	stmt, err := db.Beginx()
	for _, item := range id {
		_, err := stmt.Exec(queryStr, item)
		if err != nil {
			return errors.New("Can't update bundle of items. Reason: " + err.Error())
		}
	}
	err = stmt.Commit()
	if err != nil {
		return errors.New("transaction commit has failed. Reason: " + err.Error())
	}
	return nil
}

//CreateBundle creates a bunch of items.
func CreateBundle(db *sqlx.DB, objects []Order) error {
	queryStr := "INSERT INTO orders (uuid, client_uuid, cart_uuid, date, is_payed, status, status_date) " +
		"VALUES (:uuid, :client_uuid, :cart_uuid, :date, :is_payed, :status, :status_date)"

	if len(objects) == 0 {
		return errors.New("Array is empty")
	}
	if len(objects) == 1 {
		_, err := Create(db, objects[0])
		return err
	}
	stmt, err := db.Beginx()
	for _, item := range objects {
		item.UUID = uuid.New()
		_, err := stmt.NamedExec(queryStr, item)
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
func UpdateBundle(db *sqlx.DB, objects []Order) error {
	queryStr := "UPDATE orders SET " +
		"client_uuid = :client_uuid, cart_uuid = :cart_uuid, date = :date, is_payed = :is_payed, status = :status, status_date = :status_date" +
		" WHERE uuid = :uuid"
	if len(objects) == 0 {
		return errors.New("Array is empty")
	}
	if len(objects) == 1 {
		err := Update(db, objects[0])
		return err
	}
	stmt, err := db.Beginx()
	for _, item := range objects {
		_, err := stmt.NamedExec(queryStr, item)
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
