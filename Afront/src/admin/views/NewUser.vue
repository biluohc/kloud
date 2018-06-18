<template>
<Home>
  <template slot="home">
    <div v-if="isLogin" class="count">
    <Form :label-width="80">
            <Row>
            <Col span="2" offset="1">
                <Avatar icon="person" size="large"/>
            </Col>
           </Row>
            <Row>
              <Col span="2">名字：</Col>
              <Col span="12"><input type="text" v-model="name" size="80%"></input></Col>
            </Row>
            <Row>
              <Col span="2">邮箱：</Col>
              <Col span="12"><input type="text" v-model="email" size="80%"></input></Col>
            </Row>
            <Row>
              <Col span="2">角色：</Col>
              <Col span="12">
                <select v-model="role" class="count">
                    <option value="0">管理员</option>
                    <option value="1">普通用户</option>                
                </select>
              </Col>
            </Row>
            <Row>
              <Col span="2">总容量：</Col>
              <Col span="12">{{maxSizeStr}}</Col>
            </Row>
            <Row>
              <Col span="2">编辑总容量：</Col>
              <Col span="4"><span><input type="number" v-model="maxSize"></input> 比特</span></Col>
            </Row>
            <Row>
              <Col span="2">编辑头像：</Col>
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
import Home from "@/admin/components/Home.vue";

var defaultMaxSize = 4294967296;

export default {
  name: "NewUser",
  components: { Home },
  data: function() {
    return {
      name: "",
      email: "",
      role: 1,
      // 4G
      maxSize: defaultMaxSize,
      password: "",
      passwordCheck: "",
      avatarURL: ""
    };
  },
  computed: {
    isLogin: function() {
      return this.$store.state.isLogin;
    },
    maxSizeStr: function() {
      return this.api.size(this.maxSize);
    }
  },
  methods: {
    submit: function() {
      if (this.name === "" || this.email === "") {
        this.$Message.error("用户名和邮箱不能为空");
      }
      if (this.maxSize < 0) {
        this.$Message.error("总容量不能小于0");
      }
      if (this.password !== this.passwordCheck) {
        this.$Message.error("密码确认失败： 不相等");
        return;
      }
      if (this.password.length < 6) {
        this.$Message.error("密码长度不能小于6");
        return;
      }
      this.role = parseInt(this.role);
      this.maxSize = parseInt(this.maxSize);
      var passwordHash = this.api.sha2Password(this.password);

      console.log(
        JSON.stringify({
          name: this.name,
          email: this.email,
          role: this.role,
          avatarURL: this.avatarURL,
          maxSize: this.maxSize,
          password: passwordHash
        })
      );
      this.http
        .post(this.api.newUser, {
          name: this.name,
          email: this.email,
          role: this.role,
          avatarURL: this.avatarURL,
          maxSize: this.maxSize,
          password: passwordHash
        })
        .then(res => {
          console.log(res.status);
          if (res.status === 200) {
            this.$Message.success("新建用户成功");
            this.cancel();
            this.$router.replace("/admin/user/" + res.data.data.id);
          } else {
            console.log(res.data.msg);
            this.$Message.error("新建用户失败: " + res.data.msg);
          }
        })
        .catch(e => e);
    },
    cancel: function() {
      this.name = "";
      this.email = "";
      this.avatarURL = "";
      this.role = 1;
      this.maxSize = defaultMaxSize;
      this.password = "";
      this.passwordCheck = "";
    }
  }
};
</script>

<style lang="scss" scoped>
.count {
  font-size: large;
  font-weight: bold;
}
</style>
