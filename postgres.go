package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"

	data "github.com/fops9311/wbl0_231124/data"
	_ "github.com/lib/pq"
)

func initTable() error {
	query := `    
	CREATE TABLE IF NOT EXISTS public.orders2
	(
		id text NOT NULL,
		data bytea,
		PRIMARY KEY (id)
	);
	
	ALTER TABLE IF EXISTS public.orders2
		OWNER to postgres;`

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDBNAME)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func insertOrder(key string, data string) error {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDBNAME)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return err
	}

	defer db.Close()

	insertDynStmt := `insert into "orders2"("id", "data") values($1, $2)`
	_, err = db.Exec(insertDynStmt, key, data)
	if err != nil {
		return err
	}
	ordersWriteDb.Inc()
	return nil
}

type Record struct {
	Key  string
	Data []byte
}

func selectOrders(c chan OrderWithKey) error {

	if !isOrderTypeRegistered {
		gob.Register(data.RawOrderData{})
		isOrderTypeRegistered = true
	}

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDBNAME)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return err
	}

	defer db.Close()

	rows, err := db.Query("select * from \"orders2\"")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		r := Record{}
		err := rows.Scan(&r.Key, &r.Data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		buff := bytes.NewBuffer(r.Data)
		dec := gob.NewDecoder(buff)
		order := data.RawOrderData{}
		if err := dec.Decode(&order); err != nil {
			continue
		}
		ordersReadDb.Inc()
		orderwk := OrderWithKey{Val: order}
		orderwk.Key.Set(r.Key)
		c <- orderwk
	}
	return nil
}
