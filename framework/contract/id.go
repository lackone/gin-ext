package contract

const IdKey = "ext:id"

type IDService interface {
	NewID() string
}
