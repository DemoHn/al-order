package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

const (
	tbName  = "orders"
	attrAll = "sequence_id, id, status, distance, location_info"
)

// OrderStatus - defines the enum of all available order statuses.
type OrderStatus = string

const (
	// Unassigned - status: this order is unassigned to anyone
	Unassigned OrderStatus = "UNASSIGNED"
	// Taken - status: this order has been assigned to somebody
	Taken OrderStatus = "TAKEN"
)

// Order - order model
type Order struct {
	ID       string
	Status   string
	Distance float32
	Location LocationInfo
	// internal sequence id
	sequenceID int
}

// LocationInfo - record the start & dest points' latitude and longitude
type LocationInfo struct {
	StartLat float32
	StartLng float32
	EndLat   float32
	EndLng   float32
}

// SaveNewOrder - store new order data to DB
func SaveNewOrder(db *sql.DB, id string, distance float32, location LocationInfo) error {
	strLocation := encodeLocationInfo(location)
	insertStmt := fmt.Sprintf("insert into `%s` (id, status, distance, location_info) values (?, ?, ?, ?)", tbName)

	if _, err := db.Exec(insertStmt, id, Unassigned, distance, strLocation); err != nil {
		return err
	}
	return nil
}

// FindOneOrder - find one order info if exists
// NOTICE: errors will be returned IF AND ONLY IF sql execution failed!
func FindOneOrder(db *sql.DB, id string) (*Order, error) {
	findStmt := fmt.Sprintf("select %s from %s where id = ?", attrAll, tbName)
	// scan variables
	var order Order
	var strLocation string
	if err := db.QueryRow(findStmt, id).Scan(&order.sequenceID, &order.ID, &order.Status, &order.Distance, &strLocation); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	order.Location = decodeLocationInfo(strLocation)
	return &order, nil
}

// UpdateOrderStatus - update current status of one order
// Please ensure id exists! Or the update operation won't executed WITHOUT ANY WARNINGS!
func UpdateOrderStatus(db *sql.DB, id string, newStatus OrderStatus) error {
	updateStmt := fmt.Sprintf("update %s set status = ? where id = ?", tbName)

	if _, err := db.Exec(updateStmt, newStatus, id); err != nil {
		return err
	}
	return nil
}

// ListOrders by limits
func ListOrders(db *sql.DB, limit int, offset int) ([]Order, error) {
	listStmt := fmt.Sprintf("select %s from %s limit ? offset ?", attrAll, tbName)
	orders := make([]Order, 0)

	rows, err := db.Query(listStmt, limit, offset)
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		var strLocation string
		if err := rows.Scan(&order.sequenceID, &order.ID, &order.Status, &order.Distance, &strLocation); err != nil {
			return []Order{}, err
		}
		order.Location = decodeLocationInfo(strLocation)
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return []Order{}, err
	}
	return orders, nil
}

//// private function

// encodeLocationInfo - from LocationInfo to json string
// notice we won't return error!
func encodeLocationInfo(info LocationInfo) string {
	b, _ := json.Marshal(info)
	return string(b)
}

// decodeLocationInfo - from json string
func decodeLocationInfo(data string) LocationInfo {
	var lInfo LocationInfo
	json.Unmarshal([]byte(data), &lInfo)

	return lInfo
}
