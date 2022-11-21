import request from '@/utils/request'

// 创建
export function createExample(data) {
  return request({
    url: '/api/example/create',
    method: 'post',
    data
  })
}
// 单条数据
export function GetExample(params) {
  return request({
    url: '/api/example/id',
    method: 'get',
    params
  })
}

// 列表
export function listExample(params) {
  return request({
    url: '/api/example/list',
    method: 'get',
    params
  })
}

// 更新
export function putExample(data) {
  return request({
    url: '/api/example/put',
    method: 'put',
    data
  })
}

// 删除
export function deleteExample(data) {
  return request({
    url: '/api/example/delete',
    method: 'delete',
    data
  })
}

// 批量删除
export function removeExample(data) {
  return request({
    url: '/api/example/remove',
    method: 'delete',
    data
  })
}

// 列表
export function rank(params) {
  return request({
    url: '/api/example/rank',
    method: 'get',
    params
  })
}
