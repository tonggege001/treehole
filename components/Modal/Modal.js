// components/Modal/Modal.js
Component({
  /**
   * 组件的属性列表
   * show: 控制Modal显示与隐藏
   * height：modal的高度
   * bindcancel:点击取消按钮的回调函数
   * bindconfirm：点击确定按钮的回调函数
   */
  properties: {
    show:{
      type:Boolean,
      value:false
    },
    height:{
      type:String,
      value:'70%'
    }
  },

  /**
   * 组件的初始数据
   */
  data: {

  },

  /**
   * 组件的方法列表
   */
  methods: {
    clickMask(){
      // this.setData({show:false})
    },
    
    stopmove(){
      return true;
    }
  }
})
