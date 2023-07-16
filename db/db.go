package db

const (
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
	DBURI      = "mongodb://localhost:27017"
)

type PaginationFilter struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
