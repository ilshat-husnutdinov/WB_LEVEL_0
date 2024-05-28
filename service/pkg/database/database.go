package database

import (
	"encoding/json"
	"fmt"
	"log"
	"service/service/internal/lib/model"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type databaseConfig interface {
	GetDSN() string // get data source name, connStr := "postgres://admin:admin@localhost:5432/wb-service"
}

// ConnectDB создаёт подключение к БД
func ConnectDB(cfg databaseConfig) (*sqlx.DB, error) {
	const op = "db.db.ConnectDB"

	db, err := sqlx.Connect("pgx", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return db, nil
}

// InsertOrder выполняет запрос к БД для размещения полученного из Nats-Streaming заказа(order)
func InsertOrder(db *sqlx.DB, orderUID string, order *model.Order) error {
	const op = "db.db.InsertOrder"

	jsonByte, err := json.Marshal(order)
	if err != nil {
		log.Printf("[%v]: failed marshal to json: %v", op, err)
	}
	query := "INSERT INTO orders (order_uid, data) VALUES ($1,$2)"

	_, err = db.Exec(query, orderUID, jsonByte)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetOrder выполняет запрос к БД для получения заказа(order)
func GetOrder(db *sqlx.DB, orderUID string) (string, error) {
	const op = "db.db.getOrder"

	var jsonData string

	err := db.Get(&jsonData, "SELECT data FROM orders WHERE order_uid=$1", orderUID)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return jsonData, nil
}

// GetAllOrders выполняет запрос к БД для получения ВСЕХ заказов(orders)
func GetAllOrders(db *sqlx.DB) (map[string]*model.Order, error) {
	const op = "db.db.GetAllOrders"

	orders := make(map[string]*model.Order)

	rows, err := db.Queryx("SELECT order_uid, data FROM orders")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var data string

		if err := rows.Scan(&id, &data); err != nil {
			return nil, err
		}

		var order model.Order
		err := json.Unmarshal([]byte(data), &order)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		orders[id] = &order
	}
	return orders, nil
}
