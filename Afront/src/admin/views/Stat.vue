<template>
<Home>
  <template slot="home">
    <div class="count" >
    <Form v-if="stat" :label-width="80">
            <Row>
              <Col span="2"><router-link to="/admin/users">用户</router-link>数量：</Col>
              <Col span="12">{{ stat.usersNumber }}</Col>
            </Row>
            <Row>
            </Row>          
            <Row>
              <Col span="2"><router-link to="/admin/users">文件</router-link>数量：</Col>
              <Col span="12">{{ stat.filesNumber }}</Col>
            </Row>
            <Row>
              <Col span="2"><a :href="this.api.logFile">日志</a>大小：</Col>
              <Col span="2">{{ stat.logSizeStr }}</Col>
              <Col span="2">
                <Tooltip content="清空日志" placement="left">
                <Button type="error" shape="circle" icon="close" v-on:click="clear">
                </Button>
                </Tooltip>
            </Col>
            </Row>
    </Form>
    </div>
  </template>
</Home>
</template>

<script>
import Home from "@/admin/components/Home.vue";

export default {
  name: "Stat",
  components: { Home },
  data: function() {
    return {
      stat: null
    };
  },
  created: function() {
    this.fetchData();
  },
  computed: {
    isLogin: function() {
      return this.$store.state.isLogin;
    }
  },
  watch: {
    isLogin: function() {
      this.fetchData();
    }
  },
  methods: {
    fetchData: function() {
      if (this.isLogin) {
        this.http
          .get(this.api.stat)
          .then(res => {
            console.log(res.status);
            if (res.status === 200) {
              res.data.data.logSizeStr = this.api.size(res.data.data.logSize);
              this.stat = res.data.data;
            } else {
              console.log(res.data.msg);
              this.$Message.error(res.data.msg);
            }
          })
          .catch(e => e);
      }
    },
    clear: function() {
      this.http
        .delete(this.api.logFile)
        .then(res => {
          console.log(res.status);
          if (res.status === 200) {
            this.$Message.success("清空日志成功");
            this.fetchData();
          } else {
            console.log(res.data.msg);
            this.$Message.error("清空日志失败：" + res.data.msg);
          }
        })
        .catch(e => e);
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
