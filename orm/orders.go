package orm

import (
	"database/sql"
	"drs-orders/orders"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var findOrderItemsByOrderId = `select id, due_on, priority, status_id from order_items where order_id = ?`

func Find(id string, dest interface{}) (err error) {

	db := GetConnection()
	if err != nil {
		return err
	}

	err = db.First(dest, id).Error
	if err != nil {
		return err
	} else if dest == nil {
		return errors.New(fmt.Sprintf("%T Record with id %s not found\n", dest, id))
	}

	return nil
}

func FindOrderItemsByOrderId(orderId string) (orderItems *[]orders.OrderItem, err error) {

	db := GetConnection()
	if err != nil {
		return nil, err
	}

	orderItems = &[]orders.OrderItem{}
	err = db.Where("order_id = ?", orderId).Find(orderItems).Error
	if err != nil {
		return nil, err
	}

	return
}

func SaveOrder(from *orders.Order, target *orders.Order) error {
	var dbFrom = &orders.Order{}

	fmt.Printf("Doa is saving order %v, previously %v\n", target, from)

	db := GetConnection()

	if from == nil {
		db.Create(target)
	} else {
		if from.Id != target.Id {
			return errors.New(fmt.Sprintf("Id mismatch, previous %s, target %s\n", from.Id, target.Id))
		}

		tx := db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

		err := tx.First(dbFrom, from.Id).Error
		if err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("Failed to find previous %s due to %v\n", from.Id, err))
		} else if dbFrom == nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("Previous id %s not found\n", from.Id))
		}

		if !from.Compare(dbFrom) {
			tx.Rollback()
			return errors.New(fmt.Sprintf("%v and %v doesn't match\n", from, dbFrom))
		}

		err = tx.Model(dbFrom).Updates(target).Error
		if err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("Failed to update %s due to %v\n", target.Id, err))
		}

		err = tx.Commit().Error
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to commit changes, %s\n", target.Id))
		}
	}

	return nil
}

func SaveOrderItem(from *orders.OrderItem, target *orders.OrderItem) error {
	var dbFrom = &orders.OrderItem{}

	fmt.Printf("Doa is saving order item %v, previously %v\n", target, from)

	db := GetConnection()

	if from == nil {
		target.CreatedAt = time.Now()
		target.UpdatedAt = time.Now()
		target.Id = uuid.New().String()
		db.Create(target)
	} else {
		if from.Id != target.Id {
			return errors.New(fmt.Sprintf("Id mismatch, previous %s, target %s\n", from.Id, target.Id))
		}

		tx := db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

		err := tx.First(dbFrom, from.Id).Error
		if err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("Failed to find %s due to %v\n", from.Id, err))
		} else if dbFrom == nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("Id %s not found\n", from.Id))
		}

		if !from.Compare(dbFrom) {
			tx.Rollback()
			return errors.New(fmt.Sprintf("%v and %v doesn't match\n", from, dbFrom))
		}

		err = tx.Model(dbFrom).Updates(target).Error
		if err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("Failed to update %s due to %v\n", target.Id, err))
		}

		err = tx.Commit().Error
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to commit changes for order %s\n", target.Id))
		}
	}

	return nil
}
