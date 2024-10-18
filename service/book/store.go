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


func (s *Store) BorrowBook(bookId int, userId int) (*types.Book, error) {

	// book, err := getBookById(s, bookId)
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

	
	if book.Quantity < 1 {
		return nil, errors.New("book is not available to borrow")
	}

	errBorowing := borrowBook(s, book, userId)

	if errBorowing != nil{
		return nil, errBorowing
	}

	book.Quantity = book.Quantity - 1

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

	book.Quantity = book.Quantity + 1

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


	var borrowed types.Borrow
	errCheck := s.database.QueryRow("SELECT id, borrowed_quantity FROM borrows WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL", 
	userId, book.Id).Scan(&borrowed.Id, &borrowed.BorrowedQuantity)

	if errCheck != nil {
		if errCheck == sql.ErrNoRows {
			_, err := s.database.Exec("INSERT INTO borrows (book_id, user_id, borrowed_quantity) VALUES ($1, $2, $3)", 
	book.Id, userId, 1)

			if err != nil{
				return err
			}
		}else{
			return errCheck
		}
	}

	
	var quantity = borrowed.BorrowedQuantity + 1
	_, errUpdateBor := s.database.Exec("UPDATE borrows SET borrowed_quantity = $1 WHERE id = $2", 
	quantity, borrowed.Id)

	if errUpdateBor != nil{
		return errUpdateBor
	}
	
	var newQuantity = book.Quantity - 1
	_, errUpdate := s.database.Exec("UPDATE books SET quantity = $1 WHERE id = $2", newQuantity, book.Id)

	if errUpdate != nil{
		return errUpdate
	}
	return nil
}

func bookReturn(s *Store, book types.Book, userId int) error {

	var borrowedQuantity int
	errCount := s.database.QueryRow("SELECT borrowed_quantity FROM borrows WHERE book_id = $1 AND user_id = $2 AND returned_at IS NULL", 
	book.Id, userId).Scan(&borrowedQuantity)

	if errCount != nil{
		return errCount
	}

	if borrowedQuantity == 0 {
		return errors.New("error: book is not borrowed u can not return it")
	}


	if(borrowedQuantity == 1){
		_, err := s.database.Exec("UPDATE borrows SET returned_at = $1 AND borrowed_quantity = 0 AND WHERE book_id = $2 AND user_id = $3", 
		time.Now(), 
		book.Id, 
		userId)

		if err != nil{
			return err
		}
	}

	var newQuantity = book.Quantity + 1
	_, errUpdate := s.database.Exec("UPDATE books SET quantity = $1 WHERE id = $2", newQuantity, book.Id)

	if errUpdate != nil{
		return errUpdate
	}

	return nil
}     