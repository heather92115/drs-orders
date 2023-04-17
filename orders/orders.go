package orders

import (
	"encoding/json"
	"fmt"
	"time"
)

type Order struct {
	Id               string    `gorm:"primary_key;not_null" json:"id"`
	PatientId        string    `json:"patient_id"`
	ClinicianId      string    `json:"clinician_id"`
	PharmacyId       string    `json:"pharmacy_id"`
	OrderBatchId     string    `gorm:"foreignKey:OrderBatch" json:"order_batch_id"`
	StatusId         string    `gorm:"foreignKey:Status" json:"status_id"`
	PreviousStatusId string    `gorm:"foreignKey:Status" json:"previous_status_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DueOn            string    `json:"due_on"`
	HighestPriority  int       `json:"highest_priority"`
	ProductId        string    `json:"product_id"`
}

func (o *Order) Compare(order2 *Order) bool {
	return o.Id == order2.Id &&
		o.PatientId == order2.PatientId &&
		o.PharmacyId == order2.PharmacyId &&
		o.ClinicianId == order2.ClinicianId &&
		o.OrderBatchId == order2.OrderBatchId &&
		o.StatusId == order2.StatusId &&
		o.PreviousStatusId == order2.PreviousStatusId &&
		o.CreatedAt == order2.CreatedAt &&
		o.UpdatedAt == order2.UpdatedAt &&
		o.DueOn == order2.DueOn &&
		o.HighestPriority == order2.HighestPriority &&
		o.ProductId == order2.ProductId
}

func (o *Order) Clone() (cloned *Order) {
	cloned = &Order{}

	cloned.Id = o.Id
	cloned.PatientId = o.PatientId
	cloned.PharmacyId = o.PharmacyId
	cloned.ClinicianId = o.ClinicianId
	cloned.OrderBatchId = o.OrderBatchId
	cloned.StatusId = o.StatusId
	cloned.PreviousStatusId = o.PreviousStatusId
	cloned.CreatedAt = o.CreatedAt
	cloned.UpdatedAt = o.UpdatedAt
	cloned.DueOn = o.DueOn
	cloned.HighestPriority = o.HighestPriority
	cloned.ProductId = o.ProductId

	return
}

func (o *Order) JSON() string {
	b, err := json.Marshal(o)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ""
	}
	return string(b)
}

type OrderItem struct {
	Id                   string    `gorm:"primary_key;not_null" json:"id"`
	OrderId              string    `gorm:"foreignKey:Order;not_null" json:"order_id"`
	RxId                 string    `json:"rx_id"`
	ProductId            string    `json:"product_id;not_null"`
	ActivityDefinitionId string    `json:"activity_definition_id;not_null"`
	PriceInCents         int       `json:"price_in_cents"`
	StatusId             string    `gorm:"foreignKey:Status;not_null" json:"status_id"`
	PreviousStatusId     string    `gorm:"foreignKey:Status" json:"previous_status_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	DueOn                string    `json:"due_on"`
	Priority             int       `json:"priority"`
}

func (oi *OrderItem) Compare(oi2 *OrderItem) bool {
	return oi.Id == oi2.Id &&
		oi.OrderId == oi2.OrderId &&
		oi.RxId == oi2.RxId &&
		oi.ProductId == oi2.ProductId &&
		oi.ActivityDefinitionId == oi2.ActivityDefinitionId &&
		oi.PriceInCents == oi2.PriceInCents &&
		oi.StatusId == oi2.StatusId &&
		oi.PreviousStatusId == oi2.PreviousStatusId &&
		oi.CreatedAt == oi2.CreatedAt &&
		oi.UpdatedAt == oi2.UpdatedAt &&
		oi.DueOn == oi2.DueOn &&
		oi.Priority == oi2.Priority
}

func (oi *OrderItem) JSON() string {
	b, err := json.Marshal(oi)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ""
	}
	return string(b)
}

type OrderBatch struct {
	Id         string    `gorm:"primary_key;not_null" json:"id"`
	OwnedBy    string    `json:"owned_by"`
	ValidUntil time.Time `json:"valid_until"`
	OrderIds   []string  `gorm:"-:all" json:"order_ids"`
}

func (ob *OrderBatch) Compare(ob2 *OrderBatch) bool {
	return ob.Id == ob2.Id &&
		ob.OwnedBy == ob2.OwnedBy &&
		ob.ValidUntil == ob2.ValidUntil
}

func (ob *OrderBatch) JSON() string {
	b, err := json.Marshal(ob)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ""
	}
	return string(b)
}

type Status struct {
	Id     string `gorm:"primary_key;not_null" json:"id"`
	Status string
}

type AuditTrail struct {
	Id        string `gorm:"primary_key;not_null" json:"id"`
	Domain    string
	SubDomain string
	ObjectId  string
	Comments  string
	Diff      string
	Before    string
	After     string
	CreatedBy string
	CreatedAt time.Time
}
