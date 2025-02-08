package events

type Event interface {
	GetAggregateId() int
	GetType() string
}

type ProductCreatedEvent struct {
	Id       int
	Name     string
	ImageUrl string
	Price    int
}

func (pce ProductCreatedEvent) GetAggregateId() int { return pce.Id }
func (pce ProductCreatedEvent) GetType() string     { return "ProductCreatedEvent" }

type ProductNameUpdatedEvent struct {
	Id   int
	Name string
}

func (pce ProductNameUpdatedEvent) GetAggregateId() int { return pce.Id }
func (pce ProductNameUpdatedEvent) GetType() string     { return "ProductNameUpdatedEvent" }

type ProductDeletedEvent struct {
	Id int
}

func (pce ProductDeletedEvent) GetAggregateId() int { return pce.Id }
func (pce ProductDeletedEvent) GetType() string     { return "ProductDeletedEvent" }
