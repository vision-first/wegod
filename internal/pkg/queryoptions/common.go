package queryoptions

const (
	SelectColumns = iota
	EqualUserId
	EqualCategoryId
	IsUnexpired
	EqualBuddhaId
	OnShelfStatus
	InIds
	TimeRangeAt
	OrderByPayedMoneyDesc
	GroupByUserId
	LikeName
	EqualStatus
	CreatedAtRange
	PayedAtRange
)