<template>
<Share>
  <template slot="share">
    <section class="login">
        <div class="box">
            <div class="input last">
                <input 
                    type="text" 
                    maxlength="50" 
                    v-model.trim="password" 
                    placeholder="请输入分享的密钥" 
                   @keyup.enter="submit"
                    />
            </div>
            <div class="get-access-token">
                <router-link to="/login">登录</router-link>
                或
                <router-link to="/signin">注册</router-link>    
            </div>
            <div class="submit user-select-none" @click="submit">提交</div>
        </div>
    </section>
  </template>
</Share>
</template>

<script>
import Share from "@/components/Share.vue";

export default {
  name: "ShareInit",
  components: {
    Share
  },
  data() {
    return {
      password: ""
    };
  },
  methods: {
    submit: function() {
      console.log("submit:", this.$route.params.id, this.password);
      this.http
        .post(this.api.shareInit + this.$route.params.id, {
          password: this.password
        })
        .then(res => {
          console.log(res.status);
          if (res.status === 200) {
            this.$Message.success("验证成功，马上跳转");
            var url = "/#/share/" + res.data.data.id + "/info";
            // Referer: http://127.0.0.1:8000/
            // 好像只能手动加上自己服务器的 地址 了，否则自己有自己的 referer 也打不开。。
            window.location.href = url;
          } else {
            console.log(res.data.msg);
            this.$Message.error("验证失败: " + res.data.msg);
          }
        })
        .catch(e => e);
    }
  }
};
</script>
<style lang='scss' scoped>
</style>