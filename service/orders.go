package service

import (
	"drs-orders/orders"
	"drs-orders/orm"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var dueOnLayout = "2006-12-34"

func Find(id string, dest interface{}) (err error) {
	return orm.Find(id, dest)
}

func SaveOrder(from *orders.Order, target *orders.Order) (err error) {

	//  Make sure they differ
	if from != nil && from.Compare(target) {
		// nothing to do
		return nil
	}

	if from != nil && from.StatusId != target.StatusId {
		target.PreviousStatusId = from.StatusId
	}

	if target.Id == "nil" {
		target.CreatedAt = time.Now()
		target.UpdatedAt = time.Now()
		target.Id = uuid.New().String()

		if target.DueOn == "" {
			now := time.Now().UTC()
			target.DueOn = now.Format(dueOnLayout)
		}
	}

	return
}

func SaveOrderItem(from *orders.OrderItem, target *orders.OrderItem) (err error) {

	//  Make sure they differ
	if from != nil && from.Compare(target) {
		// nothing to do
		return nil
	}

	if target.OrderId == "" {
		return errors.New(fmt.Sprintf("Order item %s has no order id", target.Id))
	}

	if target.Id == "" {
		target.CreatedAt = time.Now()
		target.UpdatedAt = time.Now()
		target.Id = uuid.New().String()
	}

	if target.DueOn == "" {
		now := time.Now().UTC()
		target.DueOn = now.Format("dueOnLayout")
	}

	if from != nil && from.StatusId != target.StatusId {
		target.PreviousStatusId = from.StatusId
	}

	// Check to see if we need to update the parent order
	if from == nil ||
		target.DueOn != from.DueOn ||
		target.Priority != from.Priority ||
		target.StatusId != from.StatusId {
		go syncOrderWithOrderItems(target.OrderId)
	}

	return
}

func syncOrderWithOrderItems(orderId string) {
	var orderDb = &orders.Order{}

	orderItems, err := orm.FindOrderItemsByOrderId(orderId)
	if err != nil {
		fmt.Printf("Sync failed, expected to find at least one order item with orderDb %s, %v", orderId, err)
		return
	}

	err = orm.Find(orderId, orderDb)
	if err != nil {
		fmt.Printf("Sync failed, expected to find order %s, %v", orderId, err)
		return
	}

	// Get copy for comparison
	from := orderDb.Clone()

	if orderDb != nil && len(*orderItems) > 0 {
		needsUpdate := false
		for _, orderItem := range *orderItems {
			if orderItem.Priority > orderDb.HighestPriority {
				orderDb.HighestPriority = orderItem.Priority
				needsUpdate = true
			}

			if orderItem.StatusId < orderDb.StatusId {
				orderDb.StatusId = orderItem.StatusId
				needsUpdate = true
			}

			if orderItem.DueOn < orderDb.DueOn {
				orderDb.DueOn = orderItem.DueOn
				needsUpdate = true
			}
		}

		if needsUpdate {
			err = orm.SaveOrder(from, orderDb)
			fmt.Printf("Sync failed for order %s due to %v", orderId, err)

		}
	}
	return
}
