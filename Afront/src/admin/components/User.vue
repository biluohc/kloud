<style lang="scss" scoped>
.user {
  position: relative;
  padding: 0% 1% 1% 0%;
  background: #fff;

  border: 0;
  border-top: 1px solid #000;
  margin-top: -1px;
  border-bottom: 1px solid #000;
  margin-bottom: -1px;
  font-size: large;
}
</style>

<template>
<div v-if="user" class="user">
 <Row>
    <Col span="3"><router-link :to="url" >
        <Avatar v-if="user.avatarURL" :src="user.avatarURL" />
        <Avatar v-if="!user.avatarURL" icon="person" />
        {{ user.name }}
    </router-link><span v-if="user.id===userInfo.id">(你)</span></Col>
    <Col span="3">{{ user.email }}</Col>
    <Col span="3" offset="1">{{ user.roleStr }}</Col>
    <Col span="3" offset="1">{{ user.createTimeStr }}</Col>
    <Col span="2" offset="1">{{ user.curSizeStr }} / {{ user.maxSizeStr }}</Col>
    <Col span="1" offset="1" v-if="user.disabled">
      <Tooltip content="启用" placement="left">
          <Button type="success" shape="circle" icon="ios-undo" v-on:click="undisable"></Button>
      </Tooltip>
    </Col>
    <Col span="1" offset="1" v-if="!user.disabled">
      <Tooltip content="禁用" placement="left">
          <Button type="error" shape="circle" icon="close" v-on:click="disable"></Button>
      </Tooltip>
    </Col>
</Row>
</div>
</template>

<script>
export default {
  name: "User",
  components: {},
  props: {
    user: Object
  },
  data: function() {
    return {};
  },
  computed: {
    userInfo: function() {
      return this.$store.state.userInfo;
    },
    url: function() {
      return "/admin/user/" + this.user.id;
    }
  },

  methods: {
    update: function() {
      this.$emit("update");
    },
    undisable: function() {
      this.http
        .put(this.api.userUnDisable + this.user.id)
        .then(res => {
          if (res.status === 200) {
            this.$Message.success("启用成功");
            this.update();
          } else {
            console.log(res.data.msg);
            this.$Message.error("启用失败: " + res.data.msg);
          }
        })
        .catch(e => e);
    },
    disable: function() {
      this.http
        .put(this.api.userDisable + this.user.id)
        .then(res => {
          if (res.status === 200) {
            this.$Message.success("禁用成功");
            this.update();
          } else {
            console.log(res.data.msg);
            this.$Message.error("禁用失败: " + res.data.msg);
          }
        })
        .catch(e => e);
    }
  }
};
</script>

