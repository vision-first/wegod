package auth

type Ident struct {
	userId uint64
}

func NewIdent(userId uint64) *Ident {
	return &Ident{
		userId: userId,
	}
}

func (i *Ident) GetUserId() uint64 {
	return i.userId
}
