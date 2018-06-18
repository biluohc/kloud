<template>
<Home>
  <template slot="home">
    <div class="count">
    <Form ref="user"  v-if="user" :label-width="80">
            <Row>
            <Col span="2" offset="1">
                <Avatar v-if="user.avatarURL" :src="user.avatarURL" size="large" />
                <Avatar v-if="!user.avatarURL" icon="person" size="large"/>
            </Col>
           </Row>
            <Row>
              <Col span="2">名字：</Col>
              <Col span="12">{{user.name}}</Col>
            </Row>
            <Row>
              <Col span="2">邮箱：</Col>
              <Col span="12">{{user.email}}</Col>
            </Row>
            <Row>
              <Col span="2">角色：</Col>
              <Col span="12">{{user.roleStr}}</Col>
            </Row>
            <Row>
              <Col span="2">创建时间：</Col>
              <Col span="12">{{user.createTimeStr}}</Col>
            </Row>
            <Row>
              <Col span="2">总容量：</Col>
              <Col span="12">{{user.maxSizeStr}}</Col>
            </Row>
            <Row>
              <Col span="2">已用容量：</Col>
              <Col span="12">{{user.curSizeStr}}</Col>
            </Row>
            <Row>
              <Col span="2">头像：</Col>
              <Col span="12"><input type="text" v-model="avatarURL" size="80%"></input></Col>
            </Row>
            <Row>
              <Col span="2">密码：</Col>
              <Col span="12"><input type="password" v-model.trim="password" size="80%"></input></Col>
            </Row>
            <Row>
              <Col span="2">确认密码：</Col>
              <Col span="12"><input type="password" v-model.trim="passwordCheck" size="80%"></input></Col>
            </Row>
            <Row>
              <FormItem>
                <Col span="3" offset="3">
                <Button type="primary" @click="submit">提交</Button>
                </Col>
                <Col span="3">
                <Button type="ghost" @click="cancel" style="margin-left: 8px">取消</Button>
                </Col>
              </FormItem>
            </Row>
    </Form>
    </div>
  </template>
</Home>
</template>

<script>
import Home from "@/components/Home.vue";

export default {
  name: "HomeSettings",
  components: { Home },
  created: function() {},
  watch: {
    userInfo: function() {
      this.user = this.fetchData();
      this.avatarURL = this.user.avatarURL;
    }
  },
  data: function() {
    return {
      user: this.fetchData(),
      avatarURL: this.fetchData().avatarURL,
      password: "",
      passwordCheck: ""
    };
  },
  computed: {
    isLogin: function() {
      return this.$store.state.isLogin;
    },
    userInfo: function() {
      return this.$store.state.userInfo;
    }
  },
  methods: {
    fetchData: function() {
      var tmp = this.$store.state.userInfo;
      if (tmp) {
        tmp.createTimeStr = this.api.dateTime(tmp.createTime);
        tmp.curSizeStr = this.api.size(tmp.curSize);
        tmp.maxSizeStr = this.api.size(tmp.maxSize);
        tmp.roleStr = this.api.role(tmp.role);
        return tmp;
      }
    },
    submit: function() {
      console.log(
        "submit :",
        this.user.avatarURL,
        "\n",
        this.password,
        " --> ",
        this.passwordCheck
      );

      if (this.password !== this.passwordCheck) {
        this.$Message.error("密码确认失败： 不相等");
        return;
      }
      if (this.password !== "" && this.password.length < 6) {
        this.$Message.error("密码长度不能小于6");
        return;
      }

      if (this.user.avatarURL !== this.avatarURL) {
        this.http
          .post(this.api.avatarURL, {
            avatarURL: this.avatarURL
          })
          .then(res => {
            console.log(res.status);
            if (res.status === 200) {
              this.$Message.success("修改头像成功");

              this.cancel();
            } else {
              console.log(res.data.msg);
              this.$Message.error("修改头像失败: " + res.data.msg);
            }
          })
          .catch(e => e);
      }

      if (this.password !== "") {
        var passwordHash = this.api.sha2Password(this.password);

        this.http
          .post(this.api.password, {
            password: passwordHash
          })
          .then(res => {
            console.log(res.status);
            if (res.status === 200) {
              this.$Message.success("修改密码成功");

              console.log("logout");
              this.$store.dispatch("logout");
              this.$Message.success("已登出！");
              this.$router.replace("/login/");
            } else {
              console.log(res.data.msg);
              this.$Message.error("修改密码失败: " + res.data.msg);
            }
          })
          .catch(e => e);
      }
    },
    cancel: function() {
      this.$store.dispatch("checkLogin");

      this.password = "";
      this.passwordCheck = "";
      this.avatarURL = this.user.avatarURL;
    }
  }
};
</script>

<style scoped>
.count {
  font-size: large;
  font-weight: bold;
}
</style>