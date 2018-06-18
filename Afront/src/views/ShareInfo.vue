<template>
<Share>
  <template slot="share">
    <div class="shareInfo">
    <div v-if="share">
    <h1>下载：{{  share.name }}</h1>
    <p/>
    <h2>大小：{{  share.sizeStr }}</h2>
    <div class="bigb"></div>
    <a :href="url" download><Button type="primary" icon="android-download" size="large"></Button></a>    
    </div>

    <div class="shareInfo" v-if="notf">
      <h1>找不到该分享</h1>
      <div class="bigb"></div>      
    </div>
    <div class="bigb"></div>    
      <p>
        <router-link to="/login">登录</router-link>
        或
        <router-link to="/signin">注册</router-link>    
      </p>
    </div>    
  </template>
</Share>
</template>

<script>
import Share from "@/components/Share.vue";

// 关于点击下载与浏览的日狗属性
// http://www.zhangxinxu.com/wordpress/2016/04/know-about-html-download-attribute/
export default {
  name: "ShareInfo",
  components: {
    Share
  },
  data: function() {
    return {
      share: null,
      notf: false
    };
  },
  created: function() {
    var id = this.$route.params.id;
    this.http
      .get(this.api.shareInfo + id)
      .then(res => {
        console.log(res.status);
        if (res.status === 200) {
          // 以后还要处理失效
          res.data.data.sizeStr = this.api.size(res.data.data.entry.file.size)
          this.share = res.data.data;
        } else {
          this.$Message.error(res.data.msg);
          this.notf = true;
        }
      })
      .catch(e => e);
  },
  computed: {
    url: function() {
      return this.api.shareFile + this.share.id + "/" + this.share.name;
    }
  },
  methods: {}
};
</script>

<style lang="scss" scoped>
.shareInfo {
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
