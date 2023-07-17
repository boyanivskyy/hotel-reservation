package db

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
