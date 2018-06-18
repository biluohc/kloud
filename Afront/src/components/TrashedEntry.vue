<style lang="scss" scoped>
.entry {
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
<div v-if="entry" class="entry">
 <Row>
    <Col span="1"><Checkbox v-model="tmpState" @on-change="check"></Checkbox></Col>
    <Col span="12" v-if="entry.fileId"><a  target="_blank" :href="url">{{ entry.name }}</a></Col>
    <Col span="12" v-if="!entry.fileId"><router-link :to="url" >{{ entry.name }}/</router-link></Col>
    <Col span="1">{{ entry.sizeStr }}</Col>
    <Col span="3" offset="1">{{ entry.deleteTimeStr }}</Col>
    <Col span="1">
      <Tooltip content="恢复" placement="left">
          <Button type="success" shape="circle" icon="ios-undo" v-on:click="untrash"></Button>
      </Tooltip>
    </Col>
    <Col span="1">
      <Tooltip content="彻底删除" placement="left">
          <Button type="error" shape="circle" icon="close" v-on:click="remove"></Button>
      </Tooltip>
    </Col>
</Row>
</div>
</template>

<script>
export default {
  name: "TrashedEntry",
  components: {},
  props: {
    entry: Object,
    index: Number
  },
  data: function() {
    return {
      tmpState: this.entry.state
    };
  },
  computed: {
    url: function() {
      if (this.entry.fileId) {
        return this.api.entryFile + this.entry.id + "/" + this.entry.name;
      }
      return "/home/" + this.entry.id;
    }
  },

  methods: {
    update: function() {
      this.$emit("update");
    },
    check: function() {
      this.$emit("check", this.index, this.tmpState);
    },
    untrash: function() {
      this.http
        .put(this.api.entryUnTrash + this.entry.id)
        .then(res => {
          if (res.status === 200) {
            this.$Message.success("恢复成功");
          } else {
            console.log(res.data.msg);
            this.$Message.error("恢复失败: " + res.data.msg);
          }
          this.update();
        })
        .catch(e => e);
    },
    remove: function() {
      this.http
        .delete(this.api.entryDelete + this.entry.id)
        .then(res => {
          if (res.status === 200) {
            this.$Message.success("删除成功");
          } else {
            console.log(res.data.msg);
            this.$Message.error("删除失败: " + res.data.msg);
          }
          this.update();
        })
        .catch(e => e);
    }
  }
};
</script>

