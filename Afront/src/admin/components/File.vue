<style lang="scss" scoped>
.file {
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
.rc {
  text-indent: 40%;
}
</style>

<template>
<div v-if="file" class="file">
 <Row>
    <!-- <Col span="9"><router-link :to="url" >{{ file.name }}</router-link></Col> -->
    <Col span="9"><a :href="url" target="_blank">{{ file.name }}</a></Col>
    <Col span="3">{{ file.createTimeStr }}</Col>
    <Col span="1" offset="1">{{ file.sizeStr }}</Col>    
    <Col span="1" offset="1"><div class="rc">{{ file.rc }}</div></Col>
    <Col span="1" offset="1" v-if="file.disabled">
      <Tooltip content="解除禁用" placement="left">
          <Button type="success" shape="circle" icon="ios-undo" v-on:click="undisable"></Button>
      </Tooltip>
    </Col>
    <Col span="1" offset="1" v-if="!file.disabled">
      <Tooltip content="禁用" placement="left">
          <Button type="error" shape="circle" icon="close" v-on:click="disable"></Button>
      </Tooltip>
    </Col>
</Row>
</div>
</template>

<script>
export default {
  name: "File",
  components: {},
  props: {
    file: Object
  },
  data: function() {
    return {};
  },
  computed: {
    url: function() {
      // return "/admin/file/" + this.file.id;
      return this.api.file + this.file.id + "/" + this.file.name;
    }
  },

  methods: {
    update: function() {
      this.$emit("update");
    },
    undisable: function() {
      this.http
        .put(this.api.fileUnDisable + this.file.id)
        .then(res => {
          if (res.status === 200) {
            this.$Message.success("解除禁用成功");
            this.update();
          } else {
            console.log(res.data.msg);
            this.$Message.error("解除禁用失败: " + res.data.msg);
          }
        })
        .catch(e => e);
    },
    disable: function() {
      this.http
        .put(this.api.fileDisable + this.file.id)
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

