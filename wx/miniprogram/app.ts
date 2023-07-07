import { IAppOption } from "./appoption"
import { auth } from "./service/proto_gen/auth/auth_pb"
//import { rental } from "./service/proto_gen/rental/rental_pb"
//import { Rentcar } from "./service/request"

// app.ts
App<IAppOption>({
  globalData: {},
  onLaunch() {
    // 展示本地存储能力
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)

   // Rentcar.login()

    //登录
    wx.login({
      success: res => {
        wx.request({
          url: 'http://localhost:8080/v1/auth/login',
          method: 'POST',
          data: {
            code: res.code,
          } as auth.v1.ILoginRequest,
          success: res => {
            console.log(res)
            // const Res: { access_token?: string } = res.data as { access_token?: string };
            // //只有当 typeof res.data 的结果为 'object' 并且 "access_token" 存在于 res.data 中时
            // //，整个条件表达式的结果才为真。
            // console.log(Res.access_token)
            var accessToken = ""
            if (typeof res.data === 'object' && "access_token" in res.data) {
              const access_token = res.data["access_token"];
              console.log(access_token);
              accessToken = access_token
            }


          },
          fail: console.error,
        })


      },
    })

    
  },
})

