package book

import (
	"database/sql"
	"errors"
	"time"

	"github.com/martin0b101/book-rental-api/types"
)



type Store struct {
	database *sql.DB
}

func NewBookStore(database *sql.DB) *Store {
	return &Store{database: database}
}

func (s *Store) GetAvailableBooks() ([]types.Book, error) {

	rows, err := s.database.Query("SELECT id, title, quantity FROM books WHERE quantity > 0")

	if err != nil {
		return nil, err
	}

	var books []types.Book
	for rows.Next() {
		var book types.Book
		if err := rows.Scan(&book.Id, &book.Title, &book.Quantity); err != nil {
			return books, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return books, err
	}
	return books, nil
}


func (s *Store) BorrowBook(bookId int, userId int) (*types.Book, error) {

	// book, err := getBookById(s, bookId)
	var book types.Book
	err := s.database.QueryRow("SELECT id, title, quantity FROM books WHERE id = $1", 
	bookId).Scan(&book.Id,
		&book.Title,
		&book.Quantity,)

	if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("book not found")
			}
			return nil, err
		}

	
	if book.Quantity < 1 {
		return nil, errors.New("book is not available to borrow")
	}

	errBorowing := borrowBook(s, book, userId)

	if errBorowing != nil{
		return nil, errBorowing
	}

	return &book, nil
}

// u return quantity that is before borrowing
func (s *Store) ReturnBook(bookId int, userId int) (*types.Book, error){

	var book types.Book
	err := s.database.QueryRow("SELECT id, title, quantity FROM books WHERE id = $1", 
	bookId).Scan(&book.Id,
		&book.Title,
		&book.Quantity,)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("book not found")
		}
		return nil, err
	}


	errReturn := bookReturn(s, book, userId)

	if errReturn != nil{
		return nil, errReturn
	}

	return &book, nil
}


// func getBookById(s *Store, bookId int) (*types.Book, error){
	
// 	rows, err := s.database.Query("SELECT id, title, quantity FROM books WHERE id = $1", 1)
// 	if err != nil{
// 		return nil, err
// 	}

// 	var book *types.Book
// 	errScan := rows.Scan(
// 		&book.Id,
// 		&book.Title,
// 		&book.Quantity,
// 	)
// 	if errScan != nil{
// 		return nil, err
// 	}

// 	return book, nil
// }

// create transaction for this.
// tx := db.MustBegin()
// 	tx.MustExec("INSERT INTO borrows (user_id, book_id) VALUES ($1, $2)", borrow.UserID, borrow.BookID)
// 	tx.MustExec("UPDATE books SET quantity = quantity - 1 WHERE id = $1", borrow.BookID)
// 	err = tx.Commit()

func borrowBook(s *Store, book types.Book, userId int) error {
	_, err := s.database.Exec("INSERT INTO borrows (book_id, user_id) VALUES ($1, $2)", 
	book.Id, userId)

	if err != nil{
		return err
	}

	var newQuantity = book.Quantity - 1
	_, errUpdate := s.database.Exec("UPDATE books SET quantity = $1 WHERE id = $2", newQuantity, book.Id)

	if errUpdate != nil{
		return errUpdate
	}
	return nil
}

func bookReturn(s *Store, book types.Book, userId int) error {

	var borrowCountButNotReturned int
	errCount := s.database.QueryRow("SELECT COUNT(id) FROM borrows WHERE book_id = $1 AND user_id = $2 AND returned_at IS NULL", book.Id, userId).Scan(&borrowCountButNotReturned)

	if errCount != nil{
		return errCount
	}

	if borrowCountButNotReturned == 0 {
		return errors.New("error: book is not borrowed u can not return it")
	}


	_, err := s.database.Exec("UPDATE borrows SET returned_at = $1 WHERE book_id = $2 AND user_id = $3", 
	time.Now(), 
	book.Id, 
	userId)

	if err != nil{
		return err
	}

	var newQuantity = book.Quantity + 1
	_, errUpdate := s.database.Exec("UPDATE books SET quantity = $1 WHERE id = $2", newQuantity, book.Id)

	if errUpdate != nil{
		return errUpdate
	}

	return nil
}     