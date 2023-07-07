
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { routing } from "../../utils/routing"

const shareLocationKey = "share_location"

Page({
  carID: '',
  data: {
    userImgUrl: '',
    shareLocation: false,
  },
  onChooseAvatar(e: any) {
    if (e.detail.avatarUrl) {
      this.setData({
        userImgUrl: e.detail.avatarUrl,
      })
      wx.setStorageSync("head", e.detail.avatarUrl)
    }
  },
  onShareLocation(e: any) {
    this.data.shareLocation = e.detail.value
    wx.setStorageSync(shareLocationKey, this.data.shareLocation)
  },
  //以后若是有多个参数，则opt: Record<'car_id'|'is_vip', string>
  onLoad(opt: Record<'car_id', string>) {
    const o: routing.LockOpts = opt
    this.carID = o.car_id
    this.setData({
      shareLocation: wx.getStorageSync(shareLocationKey) || false
    })
  },
  onUnLockIng() {
    wx.getLocation({
      type: 'gcj02',
      success: async loc => {
        console.log('start a trip', {
          location: {
            latitude: loc.latitude,
            longttude: loc.longitude
          },
          userImgUrl: this.data.shareLocation ? this.data.userImgUrl : '',
        })

        if (!this.carID) {
          console.log("no carID specified")
          return
        }

        wx.request({
          url: 'http://localhost:8080/v1/trip',
          method: 'POST',
          data: {
            carId:"car123456",
          } as rental.v1.ICreateTripRequest,
          success: res => {
            console.log(res)
            
          }
        })

        const tripID = '648badb1b4600ccb093e5b29'
  
        wx.showLoading({
          title: "开锁中",
          mask: true,
        })
        setTimeout(() => {
          wx.redirectTo({
            url: routing.driving({
              //trip_id: trip.id!,
              trip_id: tripID
            }),
            complete: () => {
              wx.hideLoading()
            }
          })
        }, 2000);
      },
      fail: () => {
        wx.showToast({
          icon: 'none',
          title: "请前往设置页授权位置信息"
        })
      }
    })
  },
})