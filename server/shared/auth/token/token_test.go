package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1naX+CMcF4s4KYXP+yri
ZlQGy9wHbweaBaRiBkTfQs9whMOom8gxXMS48OumwCkRWKlJY9i36g8H8EipnXoK
Dry/k73ulR6mhxVJOX2GUG/ds/lKDTlWq9dKa5RxsvrZ/G5AadAp6Hi90shjYPzJ
SlS1dhlGltmbwRODk8FBdDKjYuzXvALSFwdX9sf9vLGQFb1Deze+HuzijUywrhct
bEaCiOThAuZmwvAMcoN9AoKLWzUpqhSmk13y4kKyMNgO7DhgGb0pODdoS8mQ81Bb
tP/hu39k83FvMrJPC7LymVugUhA8w/kapxS7+VW2Kc22GwUcwuoaoGgkeXQpcR2a
2wIDAQAB
-----END PUBLIC KEY-----`

//用于测试 JWTTokenVerifier 的 Verify 方法
func TestVerify(t *testing.T) {

	//将公钥的 PEM 编码形式解析为 RSA 公钥
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("cannot parse public key : %v", err)
	}

	v := &JWTTokenVerifier{
		PublicKey: pubKey,
	}

	//定义了一个测试案例的切片，每个测试案例包含名称、令牌、当前时间、期望的主题和是否期望发生错误等信息。
	cases := []struct {
		name    string
		tkn     string
		now     time.Time
		want    string
		wantErr bool
	}{
		{
			name:    "valid_token",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoicmVudENhci9hdXRoIiwic3ViIjoiNjQ5MDAxNDA2MmI0NmEyZGVhNmYzZWRmIn0.VwVAE-n0AVyTzaQ5nyfSCclEuSIQhT0th9dB5BQfusoF-8bPMlp2gze0BCOx4YpK05aywuvHMYcVlCTCUE9rhT77aiFoSYBO3BHPaCfU1w9yz5yOgujvCHWbFGRBRGKiiig_M2YzyMETPIaMpBkaCB1xY7SHo2gjEWAVVc-tEqaUnQOerSp74dupaq_ht5IXT90NRjluw6NBHztwmm-_770knT0RsZXsDiw7Xt7JPxCKCyx6Wth_Xio0IyhR_KTCb1TaGwAP9qGZaDhHPxUTVtS3mOooOw-xJWR3MRR2MJadfmF9V3-2eVv6sLjUUPcdPXxM9e2RW9vG7T_MaV_GSA",
			now:     time.Unix(1516239122, 0),
			want:    "6490014062b46a2dea6f3edf",
			wantErr: false,
		},
		{
			name:    "token_expired",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoicmVudENhci9hdXRoIiwic3ViIjoiNjQ5MDAxNDA2MmI0NmEyZGVhNmYzZWRmIn0.VwVAE-n0AVyTzaQ5nyfSCclEuSIQhT0th9dB5BQfusoF-8bPMlp2gze0BCOx4YpK05aywuvHMYcVlCTCUE9rhT77aiFoSYBO3BHPaCfU1w9yz5yOgujvCHWbFGRBRGKiiig_M2YzyMETPIaMpBkaCB1xY7SHo2gjEWAVVc-tEqaUnQOerSp74dupaq_ht5IXT90NRjluw6NBHztwmm-_770knT0RsZXsDiw7Xt7JPxCKCyx6Wth_Xio0IyhR_KTCb1TaGwAP9qGZaDhHPxUTVtS3mOooOw-xJWR3MRR2MJadfmF9V3-2eVv6sLjUUPcdPXxM9e2RW9vG7T_MaV_GSA",
			now:     time.Unix(1526239122, 0),
			wantErr: true,
		},
		{
			name:    "bad_token",
			tkn:     "bad_token",
			now:     time.Unix(1516239122, 0),
			wantErr: true,
		},
		{
			name:    "wrong_signature",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoicmVudENhci9hdXRoIiwic3ViIjoiNjQ5MDAxNDA2MmI0NmEyZGVhNmYzZWRmNCJ9.VwVAE-n0AVyTzaQ5nyfSCclEuSIQhT0th9dB5BQfusoF-8bPMlp2gze0BCOx4YpK05aywuvHMYcVlCTCUE9rhT77aiFoSYBO3BHPaCfU1w9yz5yOgujvCHWbFGRBRGKiiig_M2YzyMETPIaMpBkaCB1xY7SHo2gjEWAVVc-tEqaUnQOerSp74dupaq_ht5IXT90NRjluw6NBHztwmm-_770knT0RsZXsDiw7Xt7JPxCKCyx6Wth_Xio0IyhR_KTCb1TaGwAP9qGZaDhHPxUTVtS3mOooOw-xJWR3MRR2MJadfmF9V3-2eVv6sLjUUPcdPXxM9e2RW9vG7T_MaV_GSA",
			now:     time.Unix(1516239122, 0),
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return c.now
			}

			//调用 v.Verify 方法验证令牌，并检查是否出现了错误。
			accountID, err := v.Verify(c.tkn)
			if !c.wantErr && err != nil {
				t.Errorf("verification failed : %v", err)
			}

			if c.wantErr && err == nil {
				t.Errorf("want error, got no error : %v", err)
			}

			if accountID != c.want {
				t.Errorf("wrong accountID id . want :%q ,got :%q", c.want, accountID)
			}
		})
	}

	// tkn := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoicmVudENhci9hdXRoIiwic3ViIjoiNjQ5MDAxNDA2MmI0NmEyZGVhNmYzZWRmIn0.VwVAE-n0AVyTzaQ5nyfSCclEuSIQhT0th9dB5BQfusoF-8bPMlp2gze0BCOx4YpK05aywuvHMYcVlCTCUE9rhT77aiFoSYBO3BHPaCfU1w9yz5yOgujvCHWbFGRBRGKiiig_M2YzyMETPIaMpBkaCB1xY7SHo2gjEWAVVc-tEqaUnQOerSp74dupaq_ht5IXT90NRjluw6NBHztwmm-_770knT0RsZXsDiw7Xt7JPxCKCyx6Wth_Xio0IyhR_KTCb1TaGwAP9qGZaDhHPxUTVtS3mOooOw-xJWR3MRR2MJadfmF9V3-2eVv6sLjUUPcdPXxM9e2RW9vG7T_MaV_GSA"
	// jwt.TimeFunc = func() time.Time {
	// 	return time.Unix(1516239122, 0)
	// }

	// accountID, err := v.Verify(tkn)
	// if err != nil{
	// 	t.Errorf("verification failed : %v",err)
	// }

	// want := "6490014062b46a2dea6f3edf"
	// if accountID != want{
	// 	t.Errorf("wrong accountID id . want :%q ,got :%q",want,accountID)
	// }

}
