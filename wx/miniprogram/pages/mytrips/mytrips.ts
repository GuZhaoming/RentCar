import { routing } from "../../utils/routing"

interface Trip {
  id: string,
  start: string,
  end: string,
  duration: string,
  cost: string,
  distance: string
}

Page({
  tripsHeight:0,
  onRegisterTap() {
    wx.navigateTo({
      url: routing.register(),
    })
  },
  onChooseAvatar(e: any) {
    if (e.detail.avatarUrl) {
      this.setData({
        userImgUrl: e.detail.avatarUrl,
      })
      wx.setStorageSync("head", e.detail.avatarUrl)
    }
  },
  onLoad() {
    this.populateTrips()
  },
  populateTrips() {
    const trips: Trip[] = []
    for (let i = 0; i < 100; i++) {
      trips.push({
        id: '10000' + i,
        start: '北京',
        end: '深圳',
        duration: '8小时46分钟',
        cost: '100元',
        distance: '1000公里'
      })
    }
    this.setData({
       trips,
    })
  },
  onReady(){
      wx.createSelectorQuery().select('#heading')
      .boundingClientRect(rect=>{
         const height = wx.getSystemInfoSync().windowHeight - rect.height
         this.setData({
            tripsHeight:height,
         })
      }).exec()

  }
})