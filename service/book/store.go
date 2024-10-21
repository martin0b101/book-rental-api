package book

import (
	"database/sql"
	"errors"
	"github.com/martin0b101/book-rental-api/types"
)


type Store struct {
	database *sql.DB
}

func NewStore(database *sql.DB) *Store {
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


func (s *Store) BorrowBook(bookId int, userId int) error {

	book, err := getBookById(s, bookId)

	if err != nil{
		return err
	}

	if book.Quantity < 1 {
		return errors.New("book is not available to borrow")
	}

	errBorowing := borrowBook(s, *book, userId)

	if errBorowing != nil{
		return errBorowing
	}

	return nil
}


func (s *Store) ReturnBook(bookId int, userId int) error{

	book, err := getBookById(s, bookId)

	if err != nil{
		return err
	}

	errReturn := bookReturn(s, *book, userId) 

	if errReturn != nil{
		return errReturn
	}

	//book.Quantity = book.Quantity + 1

	return nil
}


func borrowBook(s *Store, book types.Book, userId int) error {

	var borrowed types.Borrow
	errCheck := s.database.QueryRow("SELECT id, borrowed_quantity FROM borrows WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL", 
	userId, book.Id).Scan(&borrowed.Id, &borrowed.BorrowedQuantity)

	// Start a transaction
	tx, err := s.database.Begin()

	if err != nil {
		return err
	}

	if errCheck != nil {
		if errCheck == sql.ErrNoRows {

			_, err := tx.Exec("INSERT INTO borrows (book_id, user_id, borrowed_quantity) VALUES ($1, $2, $3)", 
			book.Id, userId, 1)
			
			if err != nil{
				tx.Rollback()
				return err
			}
		}else{
			return errCheck
		}
	}

	var quantity = borrowed.BorrowedQuantity + 1
 
	_, errUpdateBor := tx.Exec("UPDATE borrows SET borrowed_quantity = $1 WHERE id = $2", 
	quantity, borrowed.Id)

	if errUpdateBor != nil{
		tx.Rollback()
		return errUpdateBor
	}
	
	var newQuantity = book.Quantity - 1

	_, errUpdate := tx.Exec("UPDATE books SET quantity = $1 WHERE id = $2", newQuantity, book.Id)

	if errUpdate != nil{
		tx.Rollback()
		return errUpdate
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func bookReturn(s *Store, book types.Book, userId int) error {

	var borrowedQuantity int
	errCount := s.database.QueryRow("SELECT borrowed_quantity FROM borrows WHERE book_id = $1 AND user_id = $2 AND returned_at IS NULL", 
	book.Id, userId).Scan(&borrowedQuantity)


	if errCount != nil{
		if errCount == sql.ErrNoRows {
			return errors.New("error: book is not borrowed u can not return it")
		}else{
			return errCount
		}
	}


	tx, err := s.database.Begin()
	if err != nil {
		return err
	}

	if borrowedQuantity == 1{
		_, err := tx.Exec("DELETE FROM borrows WHERE book_id = $1 AND user_id = $2", 
		book.Id, 
		userId)

		if err != nil{
			tx.Rollback()
			return err
		}
	}

	var newBorrowedQuant = borrowedQuantity - 1;

	_, errUpdateBor := tx.Exec("UPDATE borrows SET borrowed_quantity = $1 WHERE book_id = $2 AND user_id = $3", 
	newBorrowedQuant, book.Id, userId)
	
	if errUpdateBor != nil{
		tx.Rollback()
		return errUpdateBor
	}

	var newQuantity = book.Quantity + 1
	_, errUpdate := tx.Exec("UPDATE books SET quantity = $1 WHERE id = $2", newQuantity, book.Id)

	if errUpdate != nil{
		tx.Rollback()
		return errUpdate
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}     


func getBookById(s *Store, bookId int) (*types.Book, error){
	var book types.Book
	err := s.database.QueryRow("SELECT id, title, quantity FROM books WHERE id = $1", 
	bookId).Scan(&book.Id,
		&book.Title,
		&book.Quantity,)

	if err != nil {
			if err == sql.ErrNoRows {
				return nil, types.NotFoundError
			}
			return nil, err
		}
	return &book, nil
}