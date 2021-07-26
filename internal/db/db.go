package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"ethos-dash/internal/utils"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserId       string
	EthosKey     string
	Email        string
	AccessToken  string
	RefreshToken string
	Username     string
}

type QUser struct {
	EthosKey   string
	UserId     string
	Expiration string
	MemberType string
	Email      string
	PlanType   string
	Username   string
}

type BotKey struct {
	BotKey string
	Bot    string
	UserId string
}

var username string = utils.GetEnvKey("DB_USERNAME")
var password string = utils.GetEnvKey("DB_PASSWORD")
var hostname string = utils.GetEnvKey("DB_HOSTNAME")
var dbname string = utils.GetEnvKey("DB_NAME")

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
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

	query := "INSERT INTO users(ethosKey, userId, expiration, email, accessToken, refreshToken, username) VALUES (?, ?, DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 30 DAY), ?, ?, ?, ?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Errors %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, u.EthosKey, u.UserId, u.Email, u.AccessToken, u.RefreshToken, u.Username)
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

func GetUserById(id string) QUser {

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
	var membertype string
	var email string
	var accesstoken string
	var refreshtoken string
	var planType string
	var username string

	row := db.QueryRow(query, id)
	switch err := row.Scan(&ethosKey, &userId, &expiration, &membertype, &email, &accesstoken, &refreshtoken, &planType, &username); err {
	case sql.ErrNoRows:
		log.Printf("user query did not return for id %s", id)
	case nil:
		log.Printf("Successfully retreived ethoskey for user %s : %s", userId, ethosKey)
	default:
		panic(err)
	}

	var retUser QUser
	retUser.EthosKey = ethosKey
	retUser.UserId = userId
	retUser.Expiration = expiration
	retUser.MemberType = membertype
	retUser.Email = email
	retUser.PlanType = planType
	retUser.Username = username

	return retUser
}

func GetEthosKey(userId string) string {

	user := GetUserById(userId)

	return user.EthosKey
}

//linking a user's key with a bot key -> upgrade ethos key
func UpgradeKey(b BotKey) error {
	db, err := connectDb()
	if err != nil {
		log.Printf("Error %s when getting db connection", err)
		return nil
	}

	ethosKey := GetEthosKey(b.UserId)

	defer db.Close()
	log.Printf("Successfully connected to database")

	query := "INSERT INTO " + b.Bot + "(botKey, ethosKey, expiration) VALUES (?, ?, DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 30 DAY));"
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

func DowngradeKey(b BotKey) error {
	db, err := connectDb()
	if err != nil {
		log.Printf("Error %s when getting db connection", err)
		return nil
	}

	ethosKey := GetEthosKey(b.UserId)

	defer db.Close()
	log.Printf("Successfully connected to database")
	query := "DELETE FROM " + b.Bot + " WHERE ethosKey=\"" + ethosKey + "\";"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Errors %s when preparing SQL statement", err)
		return err
	}

	defer stmt.Close()
	res, err := stmt.ExecContext(ctx)
	if err != nil {
		log.Printf("Error %s when deleting from %s table", err, b.Bot)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected on table %s", err, b.Bot)
		return err
	}

	log.Printf("unlinked %s key from ethos key : %s", b.Bot, ethosKey)
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
			expiration datetime NOT NULL,
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

func ValidateKey(key string, bot string) (string, bool) {
	db, err := connectDb()
	if err != nil {
		log.Printf("error %s when getting db connection", err)
	}

	defer db.Close()
	log.Printf(("Successfully connected to database"))
	initialQuery := "SELECT COUNT(ethosKey) FROM " + bot + " WHERE ethosKey=\"" + key + "\";"
	query := "SELECT ethosKey, expiration FROM " + bot + " WHERE ethosKey=\"" + key + "\";"

	var count int
	initialRow := db.QueryRow(initialQuery)

	switch err := initialRow.Scan(&count); err {
	case sql.ErrNoRows:
		log.Printf("key doesnt exist in table %s", bot)
	case nil:
		log.Printf("Successfully retreived ethoskey for %s", bot)
	default:
		panic(err)
	}

	var ethoskey string
	var expdate string
	row := db.QueryRow(query)

	switch err := row.Scan(&ethoskey, &expdate); err {
	case sql.ErrNoRows:
		log.Printf("key doesnt exist in table %s", bot)
	case nil:
		log.Printf("Successfully retreived ethoskey for %s", bot)
	default:
		panic(err)
	}

	return expdate, count > 0
}

func UserCheck(id string) bool {
	db, err := connectDb()
	if err != nil {
		log.Printf("error %s when getting db connection", err)
	}

	defer db.Close()
	log.Printf(("Successfully connected to database"))
	query := "SELECT COUNT(userId) FROM users WHERE userId=\"" + id + "\";"
	var count int
	initialRow := db.QueryRow(query)

	switch err := initialRow.Scan(&count); err {
	case sql.ErrNoRows:
		log.Print("id doesnt exist in table")
	case nil:
		log.Print("Successfully counted user")
	default:
		panic(err)
	}

	log.Print(count)

	return count > 0
}
