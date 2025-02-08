package events

type ProductCreatedEvent struct {
	Id       string
	Name     string
	ImageUrl string
	Price    int
}

type ProductNameUpdatedEvent struct {
	Id   string
	Name string
}

type ProductDeletedEvent struct {
	Id string
}
