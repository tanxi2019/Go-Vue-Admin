import axios from 'axios'
import { Message, MessageBox } from 'element-ui'
import store from '@/store'
import { getToken } from '@/utils/auth'
import router from '@/router'

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // interface 的 base_url
  withCredentials: true, // send cookies when cross-domain requests
  timeout: 5000 // request timeout
})
// console.log(service.defaults.baseURL)
// request interceptor
service.interceptors.request.use(
  config => {
    // do something before request is sent
    if (store.getters.token) {
      config.headers['Authorization'] = 'Bearer ' + getToken()
      // config.headers['Content-Type'] = 'application/json'
    }
    return config
  },
  error => {
    // do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// response interceptor
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
  */

  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  response => {
    return response.data
  },
  error => {
    if (error.response.status === 401) {
      if (error.response.data.msg.indexOf('JWT认证失败') !== -1) {
        MessageBox.confirm(
          '登录超时, 重新登录或继续停留在当前页？',
          '登录状态已失效',
          {
            confirmButtonText: '重新登录',
            cancelButtonText: '继续停留',
            type: 'warning'
          }
        ).then(() => {
          // store.dispatch('user/resetToken').then(() => {
          //   location.reload()
          // })
          store.dispatch('user/logout').then(() => {
            location.reload()
          })
        })
      } else {
        Message({
          showClose: true,
          message: error.response.data.msg,
          type: 'error',
          duration: 5 * 1000
        })
        return Promise.reject(error)
      }
    } else if (error.response.status === 403) {
      router.push({ path: '/401' })
    } else {
      Message({
        showClose: true,
        message: error.response.data.msg || error.msg,
        type: 'error',
        duration: 5 * 1000
      })
      return Promise.reject(error)
    }
  }
)

export default service
