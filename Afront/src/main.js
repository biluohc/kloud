// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import Vuex from 'vuex'

import iView from 'iview';
import 'iview/dist/styles/iview.css';

import http from "axios";
http.defaults.withCredentials = true //cookie


// https://blog.csdn.net/u012702547/article/details/79066107
// axios 请求超时。
http.interceptors.request.use(config => {
  return config;
}, err => {
  console.log(err)
  iView.Message.error({ message: '请求超时!' });
  return Promise.resolve(err);
})
// axios 服务端错误处理。
http.interceptors.response.use(response => {
  return response
}, err => {
  console.log(err.response.data)
  return Promise.resolve(err.response);
})

Vue.config.productionTip = false

import App from './App'
import router from './router'
import store from './store'
import api from './api'

Vue.use(iView);
Vue.prototype.http = http
Vue.prototype.api = api

/* eslint-disable no-new */
new Vue({
  el: '#app',
  store,
  router,
  components: { App },
  template: '<App/>'
})
