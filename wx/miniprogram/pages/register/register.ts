import { routing } from "../../utils/routing"

Page({
  redirectURL:'',
  data: {
    licImgUrl: '',
    genderIndex:0,
    genders:['未知','男','女','其他'],
    birthDate:"1990-01-01",
    licNo:'',
    Name:'',
    state:'UNSUBMIT' as 'UNSUBMIT'|'VERIFYING'|'VERIFIED',
    checkImgUrl:''
  },
  onLoad(opt:Record<'redirect',string>){
    const o : routing.RegisterOpts = opt
      if(o.redirect){
           this.redirectURL = decodeURIComponent(o.redirect)      
      }
  },
  onUploadLic() {
    wx.chooseImage({
      success: res => {
        if (res.tempFilePaths.length > 0) {
          this.setData({
            licImgUrl: res.tempFilePaths[0]
          })
          setTimeout(() => {
             this.setData({
               licNo:'12343313',
                Name:'张三',
                genderIndex:1,
                birthDate:'1999-01-01'
             })
          }, 1000);
        }
      }
    })
  },
  onGenderChange(e:any){
       this.setData({
         genderIndex:e.detail.value,
       })
  },
  onBirthDateChange(e:any){
      this.setData({
         birthDate:e.detail.value,
      })
  },
  onSubmit(){
     this.setData({
         state:'VERIFYING',
     }),
     setTimeout(() => {
        this.onLicVerified()
     }, 3000);

  },
  onResSubmit(){
      this.setData({
          state:'UNSUBMIT',
          licImgUrl:'',
          checkImgUrl:'',
      })
  },
  onLicVerified(){
      this.setData({
        state:"VERIFIED",
        checkImgUrl:'/pages/image/check.png'
      })
      if(this.redirectURL){
        wx.redirectTo({
          url:this.redirectURL,
       }) 
      }
  }
})