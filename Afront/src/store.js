
import Vue from 'vue'
import Vuex from 'vuex'
import http from "axios"
import api from './api'

Vue.use(Vuex);
import Router from './router'

var store = new Vuex.Store({
  state: {
    // App 仅在不登录时才检查 Login 状态
    // isLogin: true,
    isLogin: false,
    rootEntryId: "",
    userInfo: {
      id: '',
      name: '',
      email: '',
      role: 2,
      avatarURL: '',
      createTime: '',
      disabled: false,
      maxSize: 0,
      curSize: 0
    }
  },
  mutations: {
    // 退出登录
    unlogin(state) {
      state.isLogin = false;
    },
    // 更新登录状态
    updateLoginStatus(state, data) {
      state.userInfo = data
      state.isLogin = true
      console.log("isLogin:", state.isLogin)
    },
    updateRootEntryId(state, data) {
      var tmp = state.rootEntryId
      state.rootEntryId = data
      console.log("rootEntryId: ", tmp, " -> ", state.rootEntryId)
    }
  },
  actions: {
    // 登录
    login({ state, commit }, data) {
      commit('updateLoginStatus', data);
    },
    logout({ state, dispatch, commit }) {

      commit("unlogin")
      // 删除服务端当前 session
      http.delete(api.login)
        .then(res => {
          console.log("logoutSucess: ", res)
        })
        .catch(e => {
          console.log("logoutError: ", e)
        });
    },
    checkLogin({ state, dispatch, commit }) {
      // if (!state.isLogin) {
      http.get(api.login)
        .then(res => {
          if (res.status == 200) {
            console.log("checkLoginSucess: ", res.data.data)
            commit('updateLoginStatus', res.data.data);
          } else {
            // 检查登录失败就扔到/login,
            Router.replace("/login")
          }
        })
        .catch(e => {
          console.log("checkLoginError: ", e)
        });
      // }
    },
    checkAdmin({ state, dispatch, commit }) {
      console.log("checkAdmin")
      if (!state.isLogin) {
        http.get(api.login)
          .then(res => {
            if (res.status == 200) {
              console.log("checkLoginSucess: ", res.data.data)
              commit('updateLoginStatus', res.data.data);
              // 检查admin失败就扔到/home
              if (res.data.data.role !== 0) {
                Router.replace("/home")
              }
            } else {
              // 检查登录失败就扔到/login,
              Router.replace("/login")
            }
          })
          .catch(e => {
            console.log("checkLoginError: ", e)
          });
      } else {
        if (state.userInfo.role !== 0) {
          Router.replace("/home")
        }
      }
    }
  }
});

export default store;