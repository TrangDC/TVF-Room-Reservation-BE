package models

type Token struct {
	Audience           string `json:"aud"`
	Issuer             string `json:"iss"`
	IssuedAt           int64  `json:"iat"`
	NotBefore          int64  `json:"nbf"`
	Expiration         int64  `json:"exp"`
	AuthenticationInfo string `json:"aio"`
	AuthorizedParty    string `json:"azp"`
	AuthorizedPartyCR  string `json:"azpacr"`
	Name               string `json:"name"`
	ObjectID           string `json:"oid"`
	PreferredUsername  string `json:"preferred_username"`
	RefreshTokenHash   string `json:"rh"`
	Scopes             string `json:"scp"`
	Subject            string `json:"sub"`
	TenantID           string `json:"tid"`
	UniqueTokenID      string `json:"uti"`
	Version            string `json:"ver"`
}
