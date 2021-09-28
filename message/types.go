package message

type Subscribe struct {
	Token   []byte
	Caption string
	Details string
}

type About struct {
	Details string
}

type CreateAudience struct {
	Token       []byte
	Description string
}

type JoinAudience struct {
	Audience []byte
	Expire   uint64
}

type Follower struct {
	Token  []byte
	Secret []byte
}

type ChangeAudience struct {
	Audience []byte
	Add      []Follower
	Remove   []Follower
	Details  string
}

type AdvertisingOffer struct {
	Token             []byte
	Audience          []byte
	ContentType       string
	ContentData       []byte
	AdvertisingFee    uint64
	AdvertisingWallet []byte
	Expire            uint64
	Signature         []byte
}

type Content struct {
	Audience         []byte
	ContentSecret    []byte
	ContentType      string
	ContentData      []byte
	AdvertisingToken []byte
}

type GrantPowerOfAttorney struct {
	Token  []byte
	Expire uint64
}

type RevokePowerOfAttorney struct {
	Token []byte
}
