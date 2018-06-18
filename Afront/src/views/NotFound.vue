<template>
<Share>
  <template slot="share">
    <div class="_404">
    <h1 >{{  title }}</h1>

    <div class="bigb"></div>
    <p>
    <router-link to="/login">登录</router-link>
    或
    <router-link to="/signin">注册</router-link>    
    </p>
    <div class="bigb"></div>    
  </div>
  </template>
</Share>
</template>

<script>
import Share from "@/components/Share.vue";

export default {
  name: "404",
  components: {
    Share
  },
  data: function() {
    return {
      title: "404 页面不存在"
    };
  },
  created: function() {
    // 除非每个页面进入时都设置 title，否则如果不刷新就一直是404了
    // document.title = this.title;
  },
  methods: {}
};
</script>

<style lang="scss" scoped>
._404 {
  padding: 10%, 0, 10%, 0;
  border: 10%, 0, 10%, 0;
  text-align: center;
  font-size: large;
  font-weight: bold;
}
.bigb {
  width: 200px;
  height: 200px;
}
</style>
