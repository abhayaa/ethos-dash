package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserId   string
	EthosKey string
}

type BotKey struct {
	BotKey string
	Bot    string
	UserId string
}

const (
	username = "root"
	password = "Asb5535070"
	hostname = "127.0.0.1:3306"
	dbname   = "ethos"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func connectDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	//shouldnt ever have to do this, but could be useful in the future
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return nil, err
	}

	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return nil, err
	}
	log.Printf("Rows affected %d\n", no)

	db.Close()
	db, err = sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
	return db, nil
}

func InsertUser(u User) error {

	db, err := connectDb()
	if err != nil {
		log.Printf("Error %s when getting db connection", err)
		return nil
	}

	defer db.Close()
	log.Printf("Successfully connected to the database")

	query := "INSERT INTO users(ethosKey, userId, expiration) VALUES (?, ?, DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 30 DAY))"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Errors %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, u.EthosKey, u.UserId)
	if err != nil {
		log.Printf("Error %s when inserting row into users table", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}

	log.Printf("%d rows affected", rows)
	log.Printf("User %s inserted into table", u.UserId)
	return nil
}

func GetEthosKey(user string) string {
	db, err := connectDb()
	if err != nil {
		log.Printf("Error %s when getting db connection", err)
	}

	defer db.Close()
	log.Printf("Successfully connected to database")

	query := "SELECT * FROM users where userId=?;"

	var ethosKey string
	var userId string
	var expiration string
	row := db.QueryRow(query, user)
	switch err := row.Scan(&ethosKey, &userId, &expiration); err {
	case sql.ErrNoRows:
		log.Printf("no rows returned while querying for ethoskey for user %s", user)
	case nil:
		log.Printf("Successfully retreived ethoskey for user %s : %s", userId, ethosKey)
	default:
		panic(err)
	}

	return ethosKey
}

func InsertBot(b BotKey) error {
	db, err := connectDb()
	if err != nil {
		log.Printf("Error %s when getting db connection", err)
		return nil
	}

	ethosKey := GetEthosKey(b.UserId)

	defer db.Close()
	log.Printf("Successfully connected to database")

	query := "INSERT INTO " + b.Bot + " (botKey, ethosKey) VALUES (?, ?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Errors %s when preparing SQL statement", err)
		return err
	}

	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, b.BotKey, ethosKey)
	if err != nil {
		log.Printf("Error %s when inserting into %s table", err, b.Bot)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected on table %s", err, b.Bot)
		return err
	}

	log.Printf("linked new %s key for ethos key : %s", b.Bot, ethosKey)
	log.Printf("%d rows affected", rows)
	return nil
}

func CreateBotTable(bot string) error {
	db, err := connectDb()
	if err != nil {
		log.Printf("Error %s when getting db connection", err)
		return nil
	}

	defer db.Close()
	log.Printf("Successfully connected to db")

	query := "CREATE TABLE " + bot +
		`(
			botKey varchar(40) NOT NULL UNIQUE,
			ethosKey varchar(30) NOT NULL UNIQUE,
			PRIMARY KEY (botKey),
			FOREIGN KEY (ethosKey) REFERENCES users(ethosKey)
		);`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Errors %s when preparing SQL statement for table creation for %s", err, bot)
		return err
	}

	defer stmt.Close()
	res, err := stmt.ExecContext(ctx)
	if err != nil {
		log.Printf("Errors %s when creating table for %s", err, bot)
		return nil
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected on create table %s", err, bot)
		return err
	}

	log.Printf("New table created for %s ", bot)
	log.Printf("%d rows affected, should be 0 most of the time", rows)
	return nil
}

func UpdateMembership(userId string) error {
	db, err := connectDb()
	if err != nil {
		log.Printf("error %s when getting db connection", err)
		return nil
	}

	defer db.Close()
	log.Printf("Successfully connected to the db")

	query := `UPDATE users 
			SET expiration=DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 30 DAY) 
			WHERE userId = ?;`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Errors %s when preparing SQL statement for membership update", err)
		return err
	}

	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, userId)
	if err != nil {
		log.Printf("Error  %s when updating expiration for user %s ", err, userId)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("%d rows affected", rows)
	log.Printf("Membership updated for %s ", userId)

	return nil
}
