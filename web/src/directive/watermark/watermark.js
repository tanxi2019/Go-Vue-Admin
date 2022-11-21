const watermark = {
  bind: function(el, { value }) {
    el.__handleClick__ = function() {
      let canvas = document.createElement('canvas')
      canvas.width = 250
      canvas.height = 200
      let ctx = canvas.getContext('2d')
      ctx.rotate((-20 * Math.PI) / 180) //旋转弧度
      ctx.font = font || '16px Microsoft JhengHei' //字体
      ctx.fillStyle = color || 'rgba(180, 180, 180, 0.3)' //字体填充颜色
      ctx.textAlign = 'center' //对齐方式
      ctx.textBaseline = 'Middle' //基线
      ctx.fillText(text, canvas.width / 3, canvas.height / 2) //被填充的文本
      ctx.fillText(time, canvas.width / 3, 130) //被填充的文本
      el.style.backgroundImage = `url(${canvas.toDataURL('image/png')})` //插入背景图
      el.style.zIndex = '999'


    }
    const { text, font, color, time } = value
    el.__handleClick__(text, font, color, time)
    //插入节点
  },
  inserted() {},
  componentUpdated: function(el) {},
  unbind: function(el) {}
}

export default watermark
