<template>
  <div class="container">
    <h2>websocket</h2>
    <div class="msg">{{ new Date().toLocaleTimeString() }} ：{{ msg }}</div>
    <el-button @click="handClickTest">Test</el-button>
  </div>
</template>

<script>
export default {
  name: 'index',
  data() {
    return {
      msg: '',
      lockReconnect: false, // 是否真正建立连接
      timeout: 5 * 1000, // 30秒一次心跳
      timeoutObj: null, // 心跳心跳倒计时
      serverTimeoutObj: null, // 心跳倒计时
      timeoutNum: null,// 断开 重连倒计时
      reconnectTime: 5 * 1000 // 断开 重连倒计时
    }
  },
  created() {
    // 页面刚进入时开启长连接
    this.initWebSocket()
  },
  destroyed() {
    // 页面销毁时关闭长连接，必须是挂载方法。直接调用 this.websocketclose() 不生效
    this.$websocket.close()
    // 清除时间
    clearTimeout(this.timeoutObj)
    clearTimeout(this.serverTimeoutObj)
    console.log("222")
  },

  methods: {
    //初始化websocket
    initWebSocket() {
      this.$websocket = new WebSocket(process.env.VUE_APP_WEB_SOCKET)
      this.$websocket.onopen = this.webSocketOnopen
      this.$websocket.onerror = this.webSocketOnerror
      this.$websocket.onmessage = this.webSocketOnmessage
      this.$websocket.onclose = this.webSocketClose
    },
    // 连接成功
    webSocketOnopen() {
      if (this.$websocket.readyState === 1) {
        console.log('连接成功', '状态：' + this.$websocket.readyState)
        // 开启心跳
        this.start()
      }else {
        this.reconnect()
      }

    },
    // 链接错误
    webSocketOnerror(e) {
      console.log('连接发生错误', e.type, '状态：' + this.$websocket.readyState)
      // 重连
      this.reconnect()
    },
    // 接收消息
    webSocketOnmessage(e) {
      let obj = e.data
      console.log(obj)
      this.msg = obj
      this.reset()

      // let obj = JSON.parse(e.data)[0]
      // console.log(obj)
      // switch (obj.type) {
      //   case 'heartbeat':
      //     //收到服务器信息，心跳重置
      //     console.log('接收到的服务器消息：', obj.msg)
      //     this.reset()
      //     break
      //   case 'good':
      //     this.msg = obj
      //     console.log('接收到的服务器消息：', obj.msg)
      // }
    },
    // 发送消息
    handClickTest:function() {
      if (this.$websocket.readyState === 1) {
        this.$websocket.send('123')
      }
    },
    // 关闭
    webSocketClose(e) {
      if (e.code === 10000 && this.$websocket.readyState === 3) {
        console.log('手动关闭连接，不再重连', e.code, '状态：' + this.$websocket.readyState)
      } else {
        console.log('断线重连', e.code, '状态：' + this.$websocket.readyState)
        // 断线重连
        this.reconnect()
      }
    },
    //重新连接
    reconnect() {
      if (this.lockReconnect) {
        return
      }
      this.lockReconnect = true
      //没连接上会一直重连，设置延迟避免请求过多
      console.log('每5秒发起一次重连', new Date().toLocaleTimeString())
      this.timeoutNum && clearTimeout(this.timeoutNum)
      this.timeoutNum = setTimeout(() => {
        this.initWebSocket()//新连接
        this.lockReconnect = false
      }, this.reconnectTime)
    },
    // 重置心跳
    reset() {
      clearTimeout(this.timeoutObj)//清除时间
      clearTimeout(this.serverTimeoutObj)//清除时间
      console.log('每5秒发起一次心跳', new Date().toLocaleTimeString())
      this.start() //重启心跳
    },
    // 开启心跳
    start() {
      this.timeoutObj && clearTimeout(this.timeoutObj)
      this.serverTimeoutObj && clearTimeout(this.serverTimeoutObj)
      this.timeoutObj = setTimeout(() => {
        // 这里发送一个心跳，后端收到后，返回一个心跳消息，
        if (this.$websocket.readyState === 1) {
          // 如果连接正常
          this.$websocket.send('heartBeat')
        } else {
          // 断线重连
          this.reconnect()
        }
        this.serverTimeoutObj = setTimeout(() => {
          //超时关闭
          this.$websocket.close()
        }, this.timeout)

      }, this.timeout)
    }

  }
}
</script>

<style lang="scss" scoped>
.msg {
  padding: 20px;
  box-sizing: border-box;
}
</style>
