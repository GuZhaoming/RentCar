import { routing } from "../../utils/routing"

const centPerSec = 0.7

Page({
  timer: undefined as number | undefined,
  data: {
    location: {
      latitude: 32.92,
      longitude: 118.46,
    },
    setting: {
      skew: 0,
      rotate: 0,
      showLocation: true,
      showScale: true,
      subKey: '',
      layerStyle: 1,
      enableZoom: true,
      enableScroll: true,
      enableRotate: false,
      showCompass: false,
      enable3D: true,
      enableOverlooking: false,
      enableSatellite: false,
      enableTraffic: false,
    },
    scale: 10,
    elapsed: '00:00:00',
    cost: '0.00',
  },
  onLoad(opt:Record<'trip_id',string>) {
    const o : routing.DrivingOpts = opt
    console.log("current trip",o.trip_id)
    this.setupLocationUpdator()
    this.setupTimer()
  },
  onUnload() {
    wx.stopLocationUpdate()       //关闭监听
    if (this.timer){
         clearInterval(this.timer)
    }
  },
  setupLocationUpdator() {
    // 使用该接口时需要到app.json声明该接口名称
    //开启位置实时更新功能
    wx.startLocationUpdate({
      fail: console.error,
      success: () => {
        //使用该接口时需要到app.json声明该接口名称 
        //监听位置变化事件
        wx.onLocationChange(res => {
          console.log('location', res)
          this.setData({
            location: {
              latitude: res.latitude,
              longitude: res.longitude
            }
          })
        })
      },
    })
  },
  setupTimer() {
    let elapsedSec = 0
    let cents = 0
    this.timer = setInterval(() => {
      elapsedSec++
      cents += centPerSec
      this.setData({
        elapsed: formatDuration(elapsedSec),
        cost: formatCost(cents)
      })
    }, 1000);
  },
  onEndTrips(){
     wx.redirectTo({
        url:routing.mytrips()
     })
  }
})

function formatDuration(sec: number) {
  const numString = (n: number) => n < 10 ? '0' + n.toFixed(0) : n.toFixed(0)
  const h = Math.floor(sec / 3600)
  sec -= 3600 * h
  const m = Math.floor(sec / 60)
  sec -= 60 * m
  const s = Math.floor(sec)
  return `${numString(h)}:${numString(m)}:${numString(s)}`
}

function formatCost(cents: number) {
  return (cents / 100).toFixed(2)
}

