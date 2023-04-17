package orm

import (
	"drs-orders/orders"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var UnbatchedSql = `select o.id from "orders" o
     left join order_batches ob on ob.id = o.batch_id
     join statuses s on o.status_id = s.id
     where o.status_id = ?
     and o.pharmacy_id = ?
     and (o.batch_id = '' or ob.owned_by = ? or o.id not in
    (select o.id from "orders" o
        join order_batches ob on ob.id = o.batch_id
        where ob.valid_until > ?))
    order by o.due_on, o.highest_priority desc, o.product_id
    limit 30`

func GetUnBatchedOrders(statusId string, pharmacyId string, ownedBy string) (orderIds []string, err error) {

	orderIds = []string{}
	db := GetConnection()
	rows, err := db.Raw(UnbatchedSql, statusId, pharmacyId, ownedBy, time.Now()).Rows()
	if err != nil {
		fmt.Printf("Failed to query orders, %v\n", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var orderId string
		rows.Scan(&orderId)
		if orderId != "" {
			orderIds = append(orderIds, orderId)
		}
	}

	fmt.Printf("Found these orders %+v\n", orderIds)
	return
}

func CreateOrderBatch(target *orders.OrderBatch) (err error) {

	db := GetConnection()
	target.Id = uuid.New().String()
	err = db.Create(target).Error
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create order batch due to %v", err))
	}
	err = db.Model(orders.Order{}).Where("id in ?", target.OrderIds).Update("order_batch_id", target.Id).Error
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to assign batch %s to orders %v  due to %v", target.Id, target.OrderIds, err))
	}

	return
}

func CancelOrderBatch(batchId string) (err error) {

	db := GetConnection()
	err = db.Model(&orders.OrderBatch{Id: batchId}).Update("valid_until", time.Now()).Error
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to update order batch %s due to %v", batchId, err))
	}

	return nil
}
