<template>
<Home>
  <template slot="home">
    <div class="count">
    <Form v-if="user" :label-width="80">
            <Row>
            <Col span="2" offset="1">
                <Avatar v-if="user.avatarURL" :src="user.avatarURL" size="large" />
                <Avatar v-if="!user.avatarURL" icon="person" size="large"/>
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
              <Col span="12">{{user.roleStr}}</Col>
            </Row>
            <Row>
              <Col span="2">编辑角色：</Col>
              <Col span="12">
                <select v-model="role" class="count">
                    <option value="0">管理员</option>
                    <option value="1">普通用户</option>                
                </select>
              </Col>
            </Row>
            <Row>
              <Col span="2">创建时间：</Col>
              <Col span="12">{{user.createTimeStr}}</Col>
            </Row>
            <Row>
              <Col span="2">已用容量：</Col>
              <Col span="12">{{user.curSizeStr}}</Col>
            </Row>
            <Row>
              <Col span="2">总容量：</Col>
              <Col span="12">{{user.maxSizeStr}}</Col>
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
              <Col span="2">修改密码：</Col>
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

export default {
  name: "User",
  components: { Home },
  created: function() {
    this.fetchData();
  },
  watch: {
    isLogin: function() {
      this.fetchData();
    }
  },
  data: function() {
    return {
      user: null,
      name: "",
      email: "",
      role: 1,
      maxSize: 0,
      password: "",
      passwordCheck: "",
      avatarURL: ""
    };
  },
  computed: {
    isLogin: function() {
      return this.$store.state.isLogin;
    }
  },
  methods: {
    fetchData() {
      if (this.isLogin) {
        this.http
          .get(this.api.user + this.$route.params.id)
          .then(res => {
            console.log(res.status);
            if (res.status === 200) {
              this.updateUser(res.data.data);
            } else {
              this.$Message.error(res.data.msg);
            }
          })
          .catch(e => e);
      }
    },
    updateUser: function(val) {
      console.log(this.user, "->", val);
      var tmp = val;

      tmp.createTimeStr = this.api.dateTime(tmp.createTime);
      tmp.curSizeStr = this.api.size(tmp.curSize);
      tmp.maxSizeStr = this.api.size(tmp.maxSize);
      tmp.roleStr = this.api.role(tmp.role);

      this.user = tmp;

      this.name = this.user.name;
      this.email = this.user.email;
      this.avatarURL = this.user.avatarURL;
      this.role = this.user.role;
      this.maxSize = this.user.maxSize;
    },
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
      if (this.password !== "" && this.password.length < 6) {
        this.$Message.error("密码长度不能小于6");
        return;
      }
      this.role = parseInt(this.role);
      this.maxSize =parseInt(this.maxSize)
      
      var passwordHash = "";
      if (this.password.length > 0) {
        passwordHash = this.api.sha2Password(this.password);
      }
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
        .put(this.api.user + this.user.id, {
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
            this.$Message.success("编辑用户信息成功");

            this.fetchData();
          } else {
            console.log(res.data.msg);
            this.$Message.error("编辑用户信息失败: " + res.data.msg);
          }
        })
        .catch(e => e);
    },
    cancel: function() {
      this.fetchData();

      this.name = this.user.name;
      this.email = this.user.email;
      this.avatarURL = this.user.avatarURL;
      this.role = this.user.role;
      this.maxSize = this.user.maxSize;
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
