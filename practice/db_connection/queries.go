package db_connection

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Order struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func (p *DBPool) GetAllUsers() ([]User, error) {
	query := "SELECT id, username, email, created_at, updated_at FROM users ORDER BY id"
	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}

func (p *DBPool) GetUserByID(id int) (*User, error) {
	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?"
	row := p.DB.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return &user, nil
}

func (p *DBPool) GetOrdersByUserID(userID int) ([]Order, error) {
	query := `SELECT id, user_id, product_name, quantity, price, status, created_at, updated_at 
			  FROM orders WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := p.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.UserID, &order.ProductName, &order.Quantity,
			&order.Price, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return orders, nil
}

func (p *DBPool) CreateUser(username, email string) error {
	query := "INSERT INTO users (username, email) VALUES (?, ?)"
	_, err := p.DB.Exec(query, username, email)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
