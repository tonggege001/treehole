//index.js
//获取应用实例
const app = getApp()

Page({
  data: {
    popup: false,
    commentstate: {256: false, 321: false, 148: false},
    up:{256: 3, 148: 1},
    down:{256: 0, 321: 2, 148: 1},
    commentlist:{
      256: [{nickname: "小红",  content: "加油，勇敢一点！"}, {nickname: "兰香", content: "人生很美好的。"}],
      321:[{nickname:"", content: "你一定可以的。"}],
    },

    postlist:[{
      id: 256,
      content: "历史总是那么的相似。\n以为交到比较合得来的异性朋友。聊着聊着又发现原来他喜欢我好友。我觉得好友对他也蛮有好感。但我不想当这助攻。行吗？\n我没喜欢这男的。但还是很失落。",
      timestamp: 1604490667
    },
    {
      id: 321,
      content: "俗人儿，生日快乐呀！虽然偶尔还是会抱怨为什么要来到这个世界受苦，但是还是要感谢妈妈超级不容易地把你带到这个不容易的世界来，为了不辜负她，你也要好好的过好你的人生啊，或许。。。今天就是为了提醒你，你的人生又少一年了，你还在原地踏步吗？？？是的呀，我还在原地踏步。。。所以从明天开始，我可以在事业上步步高升，在爱情上幸福美满，在生活上平安健康吗？我希望是可以的！",
      timestamp: 1604470667
    },

    {
      id: 148,
      content: "今天又想离婚\n原因是他不肯扫落叶，左右邻居都扫了，就我们一家，太突兀，不知道他啥脑回路\n把娃哄着后我自己扫了，胡乱扫了扫，装了三袋\n结果他说今天不收树叶\n发现每次吵架，都是因为他不肯按照我的意愿去干一件事情",
      timestamp: 1604370667
    }],


    nptextarea: "",
    npinput: "",
       
  },

  ttime:  function(time = +new Date()) {
    var date = new Date(time + 8 * 3600 * 1000); // 增加8小时
    return date.toJSON().substr(0, 19).replace('T', ' ');
   },


  updatedata: function(){
    var that = this;
    wx.request({
      url: 'https://tonggege.work/getpost', 
      method: "POST",
      data:{
        "page_num": 0,  
        "page_count": 0  
      }, 
      header:{
        'content-type': 'application/json' // 默认值
      },
      success(res){
        if (res.data.code == 0){console.log(res.data);
          console.log(res.data); 
          var ppostlist = []
          res.data.postlist.forEach(element => {
            ppostlist.push(JSON.parse(element))
          });
          for(var i=0;i<ppostlist.length;i++){
            ppostlist[i]["time"] = that.ttime(ppostlist[i]["timestamp"]*1000);
          }

          var ccommentlist = res.data.commentlist;
          for(var key in ccommentlist){ 
            var tmpcmtlist = [];
            ccommentlist[key].forEach(element => {
              tmpcmtlist.push(JSON.parse(element))
            });
            ccommentlist[key] = tmpcmtlist;
          }

          that.setData({
            postlist: ppostlist,
            up: res.data.up,
            down: res.data.down,
            commentlist: ccommentlist
          })

        }
        console.log(res.data)
      },
      fail(res){
        console.log("error", res)
      }
    })
  },

  //事件处理函数
  commentClick: function(e){
    var c = this.data.commentstate;
    var id = this.data.postlist[e.currentTarget.dataset.idx].id;
    c[id] = !c[id];
    this.setData({
      commentstate: c,
    })
  },

  commentFormSubmit: function(e){
    console.log("commentsubmit:", e)
    var id = this.data.postlist[e.currentTarget.dataset.idx].id;
    console.log("id=",id)
    var content = e.detail.value.comment;
    var nickname = e.detail.value.nickname==""?"匿名":e.detail.value.nickname;
    var that = this;
    wx.request({
      url: 'https://tonggege.work/newcomment',
      method: "POST",
      data:{
        "id": id,
        // "timestamp": Date.parse(new Date())/1000,
        "content": content,
        "nickname": nickname
      }, 
      success(res){
        that.updatedata();
      },
      fail(res){
        console.log("error", res)
      }
    })

  },


  upclick:function(e){
    var id = this.data.postlist[e.currentTarget.dataset.idx].id;
    var that = this;
    wx.request({
      url: 'https://tonggege.work/upgood',
      method: "POST",
      data:{
        "id": id,
      }, 
      success(res){
        that.updatedata();
      },
      fail(res){
        console.log("error", res)
      }
    })
  },

  downclick: function(e){
    var id = this.data.postlist[e.currentTarget.dataset.idx].id;
    var _up = this.data.up;
    var that = this;
    wx.request({
      url: 'https://tonggege.work/downbad',
      method: "POST",
      data:{
        "id": id,
      }, 
      success(res){
        that.updatedata();
      },
      fail(res){
        console.log("error", res)
      }
    })
  },

  onShow: function () {
    this.updatedata();
  },
 
  newpost:function(){
    this.setData({
      popup: true
    })
  },
 
  nptextareainput:function(e){
    this.setData({ 
      nptextarea: e.detail.value 
    }); 
  },

  npinput:function(e){
    this.setData({ 
      npinput: e.detail.value 
    });
  },

  newpostSubmit:function(e){
    var that = this;
    wx.request({
      url: 'https://tonggege.work/newpost',
      method: "POST",
      data:{
        "timestamp": Date.parse(new Date())/1000,
        "content": that.data.nptextarea,
        "nickname": that.data.npinput
      }, 
      success(res){
        that.updatedata();
      }, 
      fail(res){
        console.log("error", res)
      }
    })

    this.setData({
      popup: false,
      nptextarea: "",
      npinput: ""
    })
  },

  newpostCancel: function(e){
    this.setData({
      popup: false,
      nptextarea: "",
      npinput: ""
    });
  }
})


