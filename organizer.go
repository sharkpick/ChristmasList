package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	GiftTable      = "tGifts"
	RecipientTable = "tRecipients"
)

type Gift struct {
	Id        int
	Name      string
	Price     int
	URL       string
	Purchased bool
	BelongsTo string
}

func (g *Gift) GetPrice() string {
	return fmt.Sprintf("$%0.2f", float64(g.Price)/float64(100))
}

type Recipient struct {
	Id       int
	Name     string
	Gifts    []Gift
	Finished bool
}

func (r *Recipient) TotalCost() string {
	var sum int
	for _, g := range r.Gifts {
		sum += g.Price
	}
	floatSum := float64(sum) / float64(100)
	floatSum += floatSum * SalesTaxRate
	return fmt.Sprintf("$%0.2f", floatSum)
}

type ListOrganizer struct {
	db   *sql.DB
	Name string
}

func NewListOrganizer(db *sql.DB, name string) *ListOrganizer {
	return &ListOrganizer{Name: name, db: db}
}

func (l *ListOrganizer) GetGifts(id int, name string) []Gift {
	giftsQuery := `SELECT id, name, price, url, purchased FROM ` + GiftTable + ` WHERE userid=?`
	stmt, err := l.db.Prepare(giftsQuery)
	if err != nil {
		log.Fatalln("Error preparing GetGifts statement", err)
	}
	gifts := make([]Gift, 0)
	row, err := stmt.Query(id)
	if err != nil {
		log.Fatalln("Error executing GetGifts statement", err)
	}
	defer row.Close()
	for row.Next() {
		var gift Gift
		var purchased int
		err := row.Scan(&gift.Id, &gift.Name, &gift.Price, &gift.URL, &purchased)
		if err != nil {
			log.Fatalln("Error in row.Scan():", err)
		}
		gift.Purchased = purchased == 1
		gift.BelongsTo = name
		gifts = append(gifts, gift)
	}
	return gifts
}

func (l *ListOrganizer) GetRecipient(name string) (Recipient, error) {
	userQuery := `SELECT id, name, finished FROM ` + RecipientTable + ` WHERE name=?`
	stmt, err := l.db.Prepare(userQuery)
	if err != nil {
		log.Fatalln("Error preparing GetRecipient statement", err)
	}
	var recipient Recipient
	var fin int
	err = stmt.QueryRow(name).Scan(&recipient.Id, &recipient.Name, &fin)
	if err != nil {
		if err == sql.ErrNoRows {
			return Recipient{}, fmt.Errorf("recipient not found")
		}
		log.Fatalln("Error querying in GetRecipient", err)
	}
	recipient.Gifts = l.GetGifts(recipient.Id, name)

	recipient.Finished = func() bool {
		if fin == 1 {
			return true
		}
		for _, g := range recipient.Gifts {
			if !g.Purchased {
				return false
			}
		}
		return true
	}()
	return recipient, nil
}

func (l *ListOrganizer) CreateRecipient(name string) {
	insertSQL := `INSERT INTO ` + RecipientTable + `(name, finished) VALUES (?, ?)`
	statement, err := l.db.Prepare(insertSQL)
	if err != nil {
		log.Fatalln("Error preparing SQL statement:", err)
	}
	_, err = statement.Exec(name, false)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
	}
}

func (l *ListOrganizer) DropRecipient(id int, name string) {
	rcpt, err := l.GetRecipient(name)
	if err != nil {
		log.Println("Error in DropRecipient:", err)
	}
	for _, g := range rcpt.Gifts {
		l.DropGift(g.Id)
	}
	deleteSQL := `DELETE FROM ` + RecipientTable + ` WHERE id=?`
	statement, err := l.db.Prepare(deleteSQL)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = statement.Exec(id)
	if err != nil {
		log.Println("error removing", name, "from the database:", err)
	} else {
		log.Println("removed", name, "from database")
	}

}

func (l *ListOrganizer) AddGift(recipientName, giftName string, price int, url string) {
	recipient, err := l.GetRecipient(recipientName)
	if err != nil || recipient.Name == "" {
		l.CreateRecipient(recipientName)
		recipient, err = l.GetRecipient(recipientName)
		if err != nil || recipient.Name == "" {
			log.Fatalln("Error:", err)
		}
	}
	if giftName != "" {
		insertSQL := `INSERT INTO ` + GiftTable + `(name, price, url, purchased, userid) VALUES (?, ?, ?, ?, ?)`
		statement, err := l.db.Prepare(insertSQL)
		if err != nil {
			log.Fatalln("Error preparing SQL statement:", err)
		}
		_, err = statement.Exec(giftName, price, url, false, recipient.Id)
		if err != nil {
			log.Println("Error executing SQL statement:", err)
		}
	}
}

func (l *ListOrganizer) GetAllUsersAndGifts() []Recipient {
	selectSQL := `SELECT id, name, finished FROM ` + RecipientTable
	statement, err := l.db.Prepare(selectSQL)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	recipients := make([]Recipient, 0)
	row, err := statement.Query()
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer row.Close()
	for row.Next() {
		var recipient Recipient
		var fin int
		err := row.Scan(&recipient.Id, &recipient.Name, &fin)
		if err != nil {
			log.Fatalln("Error:", err)
		}
		recipient.Finished = fin == 1
		recipient.Gifts = l.GetGifts(recipient.Id, recipient.Name)
		recipients = append(recipients, recipient)
	}
	return recipients
}

func (l *ListOrganizer) SetPurchased(id int) {
	updateSQL := `UPDATE ` + GiftTable + ` SET purchased=1 WHERE id=?`
	statement, err := l.db.Prepare(updateSQL)
	if err != nil {
		log.Fatalln(err)
	}
	statement.Exec(id)
}

func (l *ListOrganizer) DropGift(id int) {
	deleteSQL := `DELETE FROM ` + GiftTable + ` WHERE id=?`
	statement, err := l.db.Prepare(deleteSQL)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = statement.Exec(id)
	if err != nil {
		log.Println("Error in DropGift:", err)
	} else {
		log.Println("Deleted gift", id)
	}
}
