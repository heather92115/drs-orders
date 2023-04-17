package service

import (
	"drs-orders/orders"
	"drs-orders/orm"
	"fmt"
	"sync"
	"time"
)

var batchTimeInMins = 15
var m = sync.Mutex{}

func AcquireNewBatch(statusId string, pharmacyId string, userId string) (orderIds []string, err error) {

	m.Lock()
	askOrderIds, err := orm.GetUnBatchedOrders(statusId, pharmacyId, userId)
	if err != nil {
		fmt.Printf("Failed to get orders to batch for status %s, pharmacy %s, user %s due to %v\n",
			statusId, pharmacyId, userId, err)
		return nil, err
	} else if len(askOrderIds) == 0 {
		fmt.Printf("No orders left for status %s, pharmacy %s, user %s\n",
			statusId, pharmacyId, userId)
		return nil, err
	}

	orderBatch := orders.OrderBatch{
		OwnedBy:    userId,
		ValidUntil: time.Now().Local().Add(time.Minute * time.Duration(batchTimeInMins)),
		OrderIds:   askOrderIds,
	}

	err = orm.CreateOrderBatch(&orderBatch)
	if err != nil {
		fmt.Printf("Failed to create order batch for status %s, pharmacy %s, user %s due to %v\n",
			statusId, pharmacyId, userId, err)
		return nil, err
	} else if len(askOrderIds) == 0 {
		fmt.Printf("No orders left for status %s, pharmacy %s, user %s\n",
			statusId, pharmacyId, userId)
		return nil, err
	}
	m.Unlock()

	orderIds = orderBatch.OrderIds
	return
}
