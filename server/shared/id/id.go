package id

//定义aid的统一类型
type AccountID string

func (a AccountID) String() string {
	return string(a)
}

//定义aid的统一类型
type TripID string

func (t TripID) String() string {
	return string(t)
}

type IdentityID string

func (i IdentityID) String() string {
	return string(i)
}

type CarID string

func (c CarID) String() string {
	return string(c)
}
