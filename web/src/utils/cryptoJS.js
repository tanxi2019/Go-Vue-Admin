import CryptoJS from 'crypto-js'

/**
 * @Description: 签名
 * @author 风很大
 * @date 2022/1/20 0020
 */
const sign = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9'

/**
 * @Description: 加密文件
 * @author 风很大
 * @sign 签名
 * @encryptString 加密字符串
 * @date 2021/12/30 0030
 */
export const aesEncrypt = function(encryptString) {
  let key = CryptoJS.enc.Utf8.parse(sign)
  let encryptStr = CryptoJS.enc.Utf8.parse(encryptString)
  let encrypted = CryptoJS.AES.encrypt(encryptStr, key, {
    mode: CryptoJS.mode.ECB,
    padding: CryptoJS.pad.Pkcs7
  })
  return encrypted.toString()
}

/**
 * @Description: 解密文件
 * @author 风很大
 * @sign 签名
 * @encryptString 解密字符串
 * @date 2021/12/30 0030
 */
export const aesDecrypt = function(decryptString) {
  let key = CryptoJS.enc.Utf8.parse(sign)
  let decryptStr = CryptoJS.AES.decrypt(decryptString, key, {
    mode: CryptoJS.mode.ECB,
    padding: CryptoJS.pad.Pkcs7
  })
  return CryptoJS.enc.Utf8.stringify(decryptStr).toString()
}

