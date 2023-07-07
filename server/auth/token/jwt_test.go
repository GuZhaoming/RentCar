package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA1naX+CMcF4s4KYXP+yriZlQGy9wHbweaBaRiBkTfQs9whMOo
m8gxXMS48OumwCkRWKlJY9i36g8H8EipnXoKDry/k73ulR6mhxVJOX2GUG/ds/lK
DTlWq9dKa5RxsvrZ/G5AadAp6Hi90shjYPzJSlS1dhlGltmbwRODk8FBdDKjYuzX
vALSFwdX9sf9vLGQFb1Deze+HuzijUywrhctbEaCiOThAuZmwvAMcoN9AoKLWzUp
qhSmk13y4kKyMNgO7DhgGb0pODdoS8mQ81BbtP/hu39k83FvMrJPC7LymVugUhA8
w/kapxS7+VW2Kc22GwUcwuoaoGgkeXQpcR2a2wIDAQABAoIBAQDDQsTIioboNLxU
qd2bzAbHrhdmApXOFDi6jFknZgt0E9RZPJ9F/rZaxU2xJIz1Bi8h7ze/rbB9bWQH
9NBhbZy8oEM468POh1KNSOcbbdsdX0yWsREoCx2LZX//hO8kStqvx13kOT/+xffm
csZwppdKkueIeCjZ7ZSu7OpuW5URieFi3yoidnRcvdtEedLyPOaKcsUVpgjfxk9R
MfOg8Ef8xyeEllozvaSZi2tkRk8Q6eDRU9no5/5OD4p/wvVLcxtWMnqu+7x6ex5b
YL+aianiC7hX7qWOjEgtiGrzsMlqMyCK2kuj79y1QvkrhNlYEgbtGh4+fsLQREzi
Y11lh7fxAoGBAP0M1v2Fc73xYNJE/uieK0nP8OXTpfc5vbGgCv07rMZKKtVgVEMp
nqD2CQfb0kn7tJPiFmjbli49/GHoM3AS9zlU3eAyJOFv3mPmgzO0+fmyU3tNAz7i
m0i6Ni8DZruq1Dt0CWTnAPB9iz3XWB3bXXcE1salTGss7ORjgDKAOTRZAoGBANj2
mgHvRKaVfXAd8K0apY5oqh47qI5uVBIIwayAWMP/Tebu50LE0gb5XPJWwBS+k7az
eFjiA/Vo0A2Lnad8q+auv99G3lbvZVKSa+40QQYVKp6/faXxUga7RjTXDwRMGmHK
jJ9Qx9GE5H4/GtAOLARqbzWE9iaW2+tlWDWPEnJTAoGAQxSNRWWEGh0LmpH5tPaA
6S87X+FsRI5E7/pKD3krQuFUW34OuEMnLuop2LB4HW6hHva3FBLpy9ZYuieQwyvz
53nM22rPhgdev8LSkvltrriMEsqGirwNiAj85heTuzn8ysnm1525DQdqyvjz/e4x
56Qbv0sAaATfw2dxC3IcN/kCgYEAlkaQO4DPqxZl8M88EZogS7ghBJnL0QOIBYm9
I88uLGtcOPUGh0+uLZhwkYYWuweZZnV+iQnbNhLn8Eng485NfpVRXsRGYi6knoG+
choNY7orcBMwY0z3xKPYJ+dBhndz2oIhzoN0M6H5ZZwt5Se7wz85UfeLSwU4xB3I
8Cft3fMCgYAb3PVReJz1P1gIp1y1A5XRnI14ZHaEy91yE+nWoAt27pqYD6uhPYYQ
bu/TGHALAPGlNbfry+Mws/jKGXiusec6pBoEv6+VUN4IXal513U6IVaJ7pIoaL60
Ay6GGb+TzjmfJEZztTKRqEI5hHJyYf4QC0SAh5acMMkrR2zEItRpew==
-----END RSA PRIVATE KEY-----`

func TestGenerateToken(t *testing.T) {
	//解析私钥
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("cannot parse private key : %v", err)
	}

	//创建 JWTTokenGen 实例
	g := NewJWTTokenGen("rentCar/auth", key)
	g.nowFunc = func() time.Time {
		return time.Unix(1516239022, 0)
	}

	//生成令牌
	tkn, err := g.GenerateToken("6490014062b46a2dea6f3edf", 2*time.Hour)
	if err != nil {
		t.Errorf("cannot generate token :%v", err)
	}
	
	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoicmVudENhci9hdXRoIiwic3ViIjoiNjQ5MDAxNDA2MmI0NmEyZGVhNmYzZWRmIn0.VwVAE-n0AVyTzaQ5nyfSCclEuSIQhT0th9dB5BQfusoF-8bPMlp2gze0BCOx4YpK05aywuvHMYcVlCTCUE9rhT77aiFoSYBO3BHPaCfU1w9yz5yOgujvCHWbFGRBRGKiiig_M2YzyMETPIaMpBkaCB1xY7SHo2gjEWAVVc-tEqaUnQOerSp74dupaq_ht5IXT90NRjluw6NBHztwmm-_770knT0RsZXsDiw7Xt7JPxCKCyx6Wth_Xio0IyhR_KTCb1TaGwAP9qGZaDhHPxUTVtS3mOooOw-xJWR3MRR2MJadfmF9V3-2eVv6sLjUUPcdPXxM9e2RW9vG7T_MaV_GSA"
	if tkn != want {
		t.Errorf("worng token generate. want:%q ,got:%q", want, tkn)
	}
}
