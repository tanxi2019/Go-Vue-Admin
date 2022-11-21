import watermark from './watermark'

const install = function(Vue) {
  Vue.directive('watermark', watermark)
}

if (window.Vue) {
  window.throttle = watermark
  Vue.use(install)
}

watermark.install = install
export default watermark
