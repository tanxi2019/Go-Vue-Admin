import throttle from './throttle'

const install = function(Vue) {
  Vue.directive('throttle', throttle)
}

if (window.Vue) {
  window.throttle = throttle
  Vue.use(install)
}

throttle.install = install
export default throttle
