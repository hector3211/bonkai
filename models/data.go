package models

import "time"

type WorkOrders struct {
	Id        int       `json:"id"`
	UserId    string    `json:"userId"`
	Completed bool      `json:"completed"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type ComplaintCategory string

const (
	SOUND     ComplaintCategory = "sound"
	LEAK      ComplaintCategory = "leak"
	APPLIANCE ComplaintCategory = "appliance"
)

type Complaint struct {
	Id             int               `json:"id"`
	UserID         string            `json:"userId"`
	Category       ComplaintCategory `json:"category"` // sound / leak / appliance
	BuildingNumber int               `json:"buldingNumber"`
	Completed      bool              `json:"completed"`
	WorkOrderId    int               `json:"workOrderId"`
	CreatedAt      time.Time         `json:"createdAt"`
}

type User struct {
	Id             string      `json:"id"`
	FirstName      []byte      `json:"firstName"`
	LastName       []byte      `json:"lastname"`
	Email          []byte      `json:"email"`
	EmailVerified  bool        `json:"emailVerified"`
	IsAdmin        bool        `json:"isAdmin"`
	BuildingNumber int         `json:"buldingNumber"`
	Complaints     []Complaint `json:"complaints"`
}
