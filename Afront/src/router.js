import Vue from 'vue'
import Router from 'vue-router'

import Test from '@/views/Test.vue'
import Login from '@/views/Login.vue'
import Signin from '@/views/Signin.vue'
import NotFound from '@/views/NotFound.vue'
import HomeSettings from '@/views/HomeSettings.vue'

import HomeFiles from '@/views/HomeFiles.vue'
import HomeTrashes from '@/views/HomeTrashes.vue'
import HomeShares from '@/views/HomeShares.vue'
import ShareInit from '@/views/ShareInit.vue'
import ShareInfo from '@/views/ShareInfo.vue'
import HomePublics from '@/views/HomePublics.vue'

const Stat = () => import(/* webpackChunkName: "admin" */ '@/admin/views/Stat.vue')
const Users = () => import(/* webpackChunkName: "admin" */ '@/admin/views/Users.vue')
const NewUser = () => import(/* webpackChunkName: "admin" */ '@/admin/views/NewUser.vue')
const User = () => import(/* webpackChunkName: "admin" */ '@/admin/views/User.vue')
const Files = () => import(/* webpackChunkName: "admin" */ '@/admin/views/Files.vue')

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/home',
      component: HomeFiles
    },
    { path: '', redirect: '/home' },
    { path: '/', redirect: '/home' },
    {
      path: '/home/:id',
      component: HomeFiles
    },
    {
      path: '/trashes',
      component: HomeTrashes
    },
    {
      path: '/shares',
      component: HomeShares
    },
    {
      path: '/share/:id/info',
      component: ShareInfo
    },
    {
      path: '/share/:id/init',
      component: ShareInit
    },
    {
      path: '/publics',
      component: HomePublics
    },
    {
      path: '/settings',
      component: HomeSettings
    },
    {
      path: '/login',
      component: Login
    },
    {
      path: '/signin',
      component: Signin
    },
    // /admin 
    {
      path: '/admin',
      component: Stat,
    },
    {
      path: '/admin/users',
      component: Users
    },
    {
      path: '/admin/user/:id',
      component: User
    },
    {
      path: '/admin/newUser',
      component: NewUser
    },
    {
      path: '/admin/files',
      component: Files
    },
    {
      path: '/test',
      component: Test
    },
    { path: '*', component: NotFound }
  ]
})


