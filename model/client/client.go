package client

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

//Client is a model of Client in the store.
type Client struct {
	UUID      uuid.UUID `db:"uuid" schema:"uuid" json:"uuid"`
	FirstName string    `db:"first_name" schema:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" schema:"last_name" json:"last_name"`
	Email     string    `db:"email" schema:"email" json:"email"`
	Phone     string    `db:"phone" schema:"phone" json:"phone"`
	Address   string    `db:"address" schema:"address" json:"address"`
}

//Create insert Client to db
func Create(db *sqlx.DB, c Client) (uuid.UUID, error) {
	c.UUID = uuid.New()
	queryStr := "INSERT INTO client (uuid, first_name, last_name, email, phone, address)" +
		" VALUES (:uuid, :first_name, :last_name, :email, :phone, :address)"
	_, err := db.NamedExec(queryStr, c)

	if err != nil {
		return uuid.UUID{}, err
	}
	return c.UUID, nil
}

//Read returns all Clients.
func Read(db *sqlx.DB) ([]Client, error) {
	Client := []Client{}

	queryStr := "SELECT * FROM client"
	stmt, _ := db.Preparex(queryStr)
	err := stmt.Select(&Client)

	return Client, err
}

//ReadUUID returns specified Client.
func ReadUUID(db *sqlx.DB, u uuid.UUID) (Client, error) {
	Client := Client{UUID: u}
	queryStr := "SELECT * FROM client WHERE uuid = $1"
	err := db.Get(&Client, queryStr, Client.UUID)
	return Client, err
}

//Update updates the Client
func Update(db *sqlx.DB, c Client) error {
	queryStr := "UPDATE client SET first_name = :first_name, last_name = :last_name, email = :email, phone = :phone, address = :address" +
		" WHERE uuid = :uuid"
	_, err := db.NamedExec(queryStr, c)
	return err
}

//Delete deletes the Client
func Delete(db *sqlx.DB, id uuid.UUID) error {
	queryStr := "DELETE FROM Client WHERE uuid = $1"
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

		_, err := stmt.Exec(queryStr, item)
		if err != nil {
			return errors.New("Can't delete bundle of items. Reason: " + err.Error())
		}
	}
	err = stmt.Commit()
	if err != nil {
		return errors.New("Transaction commit has been failed. Reason: " + err.Error())
	}
	return nil
}

//CreateBundle creates a bunch of items.
func CreateBundle(db *sqlx.DB, objects []Client) error {
	queryStr := "INSERT INTO client (uuid, first_name, last_name, email, phone, address)" +
		" VALUES (:uuid, :first_name, :last_name, :email, :phone, :address)"

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
func UpdateBundle(db *sqlx.DB, objects []Client) error {
	queryStr := "UPDATE client SET first_name = :first_name, last_name = :last_name, email = :email, phone = :phone, address = :address WHERE uuid = :uuid"
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
