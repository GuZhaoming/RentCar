package wechat

import (
	"fmt"

	"github.com/jasonchance/weapp"
)

type Service struct {
	AppID     string
	AppSecret string
}

func (s *Service) Resolve(code string) (string, error) {
	resp, err := weapp.Login(s.AppID,s.AppSecret,code)
	if err != nil{
		return "",fmt.Errorf("weapp login err :%v",err)
	}

	return resp.OpenID , nil
}