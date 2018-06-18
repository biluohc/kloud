<template>
<Home>
  <template slot="home">
    <Row>
        <div class="count attr">
            <Col span="3" offset="1">用户</Col>
            <Col span="2" >邮箱</Col>
            <Col span="1" offset="1">角色</Col>
            <Col span="3" offset="3" >创建时间</Col>
            <Col span="2" offset="1" >容量</Col>
            <Col span="2" offset="3" ><router-link :to="'/admin/newUser'">新建用户</router-link></Col>
        </div>
    </Row>
    <div v-for="user in users" :key='user.id'>
      <User :user="user" @update="fetchData" />
    </div>
  </template>
</Home>
</template>

<script>
import Home from "@/admin/components/Home.vue";
import User from "@/admin/components/User.vue";

export default {
  name: "Users",
  components: { Home, User },
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
      users: []
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
          .get(this.api.users)
          .then(res => {
            console.log(res.status);
            if (res.status === 200) {
              this.updateUsers(res.data.data);
            } else {
              this.$Message.error(res.data.msg);
            }
          })
          .catch(e => e);
      }
    },
    updateUsers: function(val) {
      console.log(this.users.length, "->", val.length);
      for (var i = 0; i < val.length; i++) {
        var tmp = val[i];

        tmp.createTimeStr = this.api.dateTime(tmp.createTime);

        tmp.curSizeStr = this.api.size(tmp.curSize);
        tmp.maxSizeStr = this.api.size(tmp.maxSize);
        tmp.roleStr = this.api.role(tmp.role);
        // 属性要先加入，加入后再改变Vue就日狗了。。

        this.users.splice(i, 1, tmp);
      }
      this.users.length = val.length;
      // 为0时只剩改长度，vue检测不到，垃圾js
      if (val.length == 0) {
        this.users.splice(0, this.users.length);
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.count {
  font-size: large;
  font-weight: bold;
  border: 1%, 0, 1%, 0;
}

.attr {
  padding: 10px, 0, 10px, 0;
}
</style>
