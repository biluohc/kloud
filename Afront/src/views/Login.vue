<template>
<Share>
  <template slot="share">
    <section class="login">
        <div class="box">
            <div class="input last">
                <input 
                    type="text" 
                    maxlength="50" 
                    v-model.trim="email" 
                    placeholder="帐号" 
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
            <div class="get-access-token">
               <router-link to="/signin">注册</router-link>
            </div>
            <div class="submit user-select-none" @click="submit">登录</div>
        </div>
    </section>
  </template>
</Share>
</template>

<script>
import Share from "@/components/Share.vue";

export default {
  name: "Login",
  components: {
    Share
  },
  data() {
    return {
      email: "",
      password: ""
    };
  },
  methods: {
    valid: function() {
      try {
        this.email = this.email.trim();
        if (this.email === "") {
          throw "帐号不能为空";
        }
        this.password = this.password.trim();
        if (this.password === "") {
          throw "密码不能为空";
        }
        return true;
      } catch (e) {
        this.$Message.error(e);
      }
    },
    submit: function() {
      console.log(this.email, " --> ", this.password);
      if (!this.valid()) {
        return;
      }
      var passwordHash = this.api.sha2Password(this.password);

      // todu: Chose maxAge
      this.http
        .post(this.api.login, {
          email: this.email,
          password: passwordHash,
          maxAge: 86400 * 365
        })
        .then(res => {
          console.log(res.status);
          if (res.status === 200) {
            this.$store.dispatch("login", res.data);
            this.$router.push("/home/");
          } else {
            console.log(res.data.msg);
            this.$Message.error("登录失败: " + res.data.msg);
          }
        })
        .catch(e => e);
    }
  }
};
</script>
<style lang='scss' scoped>
</style>