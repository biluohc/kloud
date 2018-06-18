<template>
<Share>
  <template slot="share">
    <section class="login">
        <div class="box">
            <div class="input last">
                <input 
                    type="text" 
                    maxlength="50" 
                    v-model.trim="name" 
                    placeholder="名字" 
                    @keyup.enter="submit"
                    />
            </div>
            <div class="input last">
                <input 
                    type="text" 
                    maxlength="50" 
                    v-model.trim="email" 
                    placeholder="邮箱" 
                    @keyup.enter="submit"
                    />
            </div>
            <div class="input last">
                <input 
                    type="password" 
                    maxlength="50" 
                    v-model.trim="password" 
                    placeholder="密码" 
                   @keyup.enter="submit"
                    />
            </div>
            <div class="input last">
                <input 
                    type="password" 
                    maxlength="50" 
                    v-model.trim="passwordCheck" 
                    placeholder="确认密码" 
                   @keyup.enter="submit"
                    />
            </div>
            <div class="get-access-token">
                <a href="/#/login">登录</a>
            </div>
            <div class="submit user-select-none" @click="submit">注册</div>
        </div>
    </section>
  </template>
</Share>
</template>


<script>
import Share from "@/components/Share.vue";

export default {
  name: "Signin",
  components: {
    Share
  },
  data() {
    return {
      name: "",
      email: "",
      password: "",
      passwordCheck: ""
    };
  },
  methods: {
    valid: function() {
      try {
        if (this.name === "") {
          throw "名字不能为空";
        }
        if (this.email === "") {
          throw "邮箱不能为空";
        }
        if (this.password === "") {
          throw "密码不能为空";
        }
        if (this.password.length < 6) {
          throw "密码不能小于6位";
        }
        if (this.password !== this.passwordCheck) {
          throw "确认密码失败：不相等";
        }
        return true;
      } catch (e) {
        this.$Message.error(e);
      }
    },
    submit: function() {
      if (!this.valid()) {
        return;
      }
      var passwordHash = this.api.sha2Password(this.password);
      this.http
        .post(this.api.signin, {
          name: this.name,
          email: this.email,
          password: passwordHash
        })
        .then(res => {
          if (res.status === 200) {
            this.$Message.success("注册成功，转到登录页面");
            this.$router.replace("/login");
          } else {
            console.log(res.data);
            this.$Message.error(res.status + ": " + res.data.msg);
          }
        })
        .catch(e => {
          this.$Message.error("网络连接错误 ", e);
        });
    }
  }
};
</script>
<style lang='scss' scoped>
</style>