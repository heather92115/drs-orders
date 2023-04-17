package main

import (
	"drs-orders/orders"
	"drs-orders/orm"
	"drs-orders/service"
	"fmt"
)

func main() {
	fmt.Println("Starting Doctor's Orders!")

	err := orm.CreatePool()
	if err != nil {
		fmt.Printf("Failed DB connections, %v\n", err)
		return
	}

	err = orm.MigrateTables()
	if err != nil {
		fmt.Printf("Failed to migrate tables, %v\n", err)
		return
	}

	var previous = &orders.Order{}
	err = service.Find("3", previous)
	if err != nil {
		fmt.Printf("Failed to find order due to %v\n", err)
	}

	target := &orders.Order{Id: previous.Id, PharmacyId: "1", PatientId: "1", ClinicianId: "1", StatusId: "3"}

	err = service.SaveOrder(previous, target)
	if err != nil {
		fmt.Printf("Failed to save order due to %v\n", err)
	}

	var previousOrderItem = &orders.OrderItem{}
	err = service.Find("4", previousOrderItem)
	if err != nil {
		fmt.Printf("Failed to find order item due to %v\n", err)
	}

	targetOrderItem := &orders.OrderItem{OrderId: "3", RxId: "1", ProductId: "2", ActivityDefinitionId: "1", StatusId: "2"}
	err = service.SaveOrderItem(previousOrderItem, targetOrderItem)
	if err != nil {
		fmt.Printf("Failed to save order item due to %v\n", err)
	}

	orderIds, err := service.AcquireNewBatch("3", "1", "2")
	if err != nil {
		fmt.Printf("Failed to save order batch due to %v\n", err)
	}

	fmt.Printf("Created batch with these orders %v", orderIds)

}
