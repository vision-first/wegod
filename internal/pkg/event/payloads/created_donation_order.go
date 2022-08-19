package payloads

type CreatedDonationOrder struct {
	Money uint32
	UserId uint64
	CreatedAt int64
	Id uint64
}