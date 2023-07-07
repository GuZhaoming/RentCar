import { routing } from "../../utils/routing"

//const app = getApp<IAppOption>()

Page({
  isPageShowing: false,
  data: {
    userImgUrl: '',
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
      showCompass: true,
      enable3D: true,
      enableOverlooking: false,
      enableSatellite: false,
      enableTraffic: false,
    },
    location: {
      latitude: 31,   //经度
      longitude: 120,   //纬度
    },
    scale: 10,   //地图缩放级别
    markers: [
      {
        iconPath: "/pages/image/location.jpg",
        id: 0,
        latitude: 23.099994,
        longitude: 113.324520,
        width: 50,
        height: 50,
      },
      {
        iconPath: "/pages/image/location.jpg",
        id: 1,
        latitude: 23.099994,
        longitude: 114.324520,
        width: 50,
        height: 50,
      },
      {
        iconPath: "/pages/image/location.jpg",
        id: 2,
        latitude: 29.756825521115363,
        longitude: 121.87222114786053,
        width: 50,
        height: 50,
      },
    ],
  },
  onLoad() {
    this.setData({
      userImgUrl: wx.getStorageSync("head")
    })
  }
  ,
  onMyLocationTap() {
    wx.getLocation({
      type: 'gcj02',
      success: res => {
        this.setData({
          location: {
            latitude: res.latitude,
            longitude: res.longitude
          },
        })
      },
      fail: () => {
        wx.showToast({
          icon: 'none',
          title: "请前往设置页授权"
        })
      },
    })
  },
  onShow() {
    this.isPageShowing = true;
  },
  onHide() {
    this.isPageShowing = false;
  },
  onMyTrips(){
       wx.navigateTo({
         //url:'/pages/mytrips/mytrips',
         url:routing.mytrips(),
       })
  },
  moveCars() {
    const map = wx.createMapContext("map")
    const dest = {
      latitude: 23.099994,
      longitude: 113.324520,
    }
    const moveCar = () => {
      dest.latitude += 0.1
      dest.longitude += 0.1
      map.translateMarker({
        destination: {
          latitude: dest.latitude,
          longitude: dest.longitude,
        },
        markerId: 0,
        autoRotate: false,
        rotate: 0,
        duration: 5000,
        animationEnd: () => {
          if (this.isPageShowing) {
            moveCar()
          }
        }
      })
    }
    moveCar()
  },
  onScanClicked() {
    wx.scanCode({
      success: () => {
        const carID = 'car123'
        const redirectURL = routing.lock({
            car_id:carID
        })
        wx.navigateTo({
          //url: `/pages/register/register?redirect=${encodeURIComponent(redirectURL)}`
          url:routing.register({
             redirectURL:redirectURL,
          })
        })
      },
      fail: console.error,
    })
  },
})
