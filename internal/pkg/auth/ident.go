package auth

type Ident struct {
	uid uint64
}

func NewIdent(uid uint64) *Ident {
	return &Ident{
		uid: uid,
	}
}

func (i *Ident) GetUid() uint64 {
	return i.uid
}
