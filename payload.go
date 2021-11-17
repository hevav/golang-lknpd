package lknpd

type DeviceInfo struct {
	AppVersion     string      `json:"appVersion,omitempty"`
	SourceDeviceId string      `json:"sourceDeviceId,omitempty"`
	SourceType     string      `json:"sourceType,omitempty"`
	MetaDetails    MetaDetails `json:"metaDetails"`
}

type MetaDetails struct {
	UserAgent string `json:"userAgent,omitempty"`
}

type RefreshTokenPayload struct {
	DeviceInfo   DeviceInfo `json:"deviceInfo"`
	RefreshToken string     `json:"refreshToken,omitempty"`
}

type AuthPayload struct {
	Username   string     `json:"username,omitempty"`
	Password   string     `json:"password,omitempty"`
	DeviceInfo DeviceInfo `json:"deviceInfo"`
}

type AuthResponse struct {
	Profile       AuthResponseProfile `json:"profile"`
	RefreshToken  string              `json:"refreshToken,omitempty" json:"refresh_token,omitempty"`
	Token         string              `json:"token,omitempty" json:"token,omitempty"`
	TokenExpireIn string              `json:"tokenExpireIn,omitempty" json:"token_expire_in,omitempty"`
}

type AuthResponseProfile struct {
	INN string `json:"inn,omitempty"`
}

type IncomeResponse struct {
	ApprovedReceiptUuid string `json:"approvedReceiptUuid,omitempty"`
}
