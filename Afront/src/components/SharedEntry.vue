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
    <Col span="9"><a  target="_blank" :href="url">{{ entry.name }}</a></Col>
    <Col span="3">{{ entry.password }}</Col>
    <Col span="1" v-if="entry.password">
      <Tooltip content="编辑密码" placement="left">
          <Button type="primary" shape="circle" icon="edit" v-on:click="editPassword"></Button>
      </Tooltip>
    </Col>
    <Col span="1" offset="3" v-if="!entry.password">
      <Tooltip content="编辑密码" placement="left">
          <Button type="primary" shape="circle" icon="edit" v-on:click="editPassword"></Button>
      </Tooltip>
    </Col>
    <Col span="3" offset="1">{{ entry.expireTimeStr }}</Col>
    <Col span="1">
      <Tooltip content="编辑过期时间" placement="left">
          <Button type="primary" shape="circle" icon="android-time" v-on:click="editMaxAge"></Button>
      </Tooltip>
    </Col>
    <Col span="3" offset="1">{{ entry.createTimeStr }}</Col>
    <Col span="1">
      <Tooltip content="删除此分享" placement="left">
          <Button type="error" shape="circle" icon="close" v-on:click="remove"></Button>
      </Tooltip>
    </Col>
</Row>
</div>
</template>

<script>
import AgeSelector from "@/components/AgeSelector.vue";

export default {
  name: "SharedEntry",
  components: { AgeSelector },
  props: {
    entry: Object,
    index: Number
  },
  data: function() {
    return {
      tmpState: this.entry.state,
      tmpValue: "",
      tmpAge: 0
    };
  },
  computed: {
    url: function() {
      return  "/#/share/" + this.entry.id + "/info";
    }
  },

  methods: {
    update: function() {
      this.$emit("update");
    },
    check: function() {
      this.$emit("check", this.index, this.tmpState);
    },
    editPassword: function() {
      console.log("editPassword", this.index, this.entry.name);
      this.tmpValue = this.entry.password;

      this.$Modal.confirm({
        render: h => {
          return h("Input", {
            props: {
              value: this.tmpValue,
              autofocus: true,
              placeholder: "空表示没有密码"
            },
            on: {
              input: val => {
                console.log(this.tmpValue, " -> ", val);
                this.tmpValue = val;
              }
            }
          });
        },
        title: "编辑密码：" + this.entry.name,
        closable: true,
        scrollable: true,
        onOk: () => {
          if (this.tmpValue === this.entry.password) {
            return;
          }
          this.http
            .put(this.api.sharePassword + this.entry.id, {
              password: this.tmpValue
            })
            .then(res => {
              console.log(res.status);
              if (res.status === 200) {
                this.$Message.success("编辑密码成功");

                this.update();
              } else {
                console.log(res.data.msg);
                this.$Message.error("编辑密码失败: " + res.data.msg);
              }
            })
            .catch(e => e);
        }
      });
    },
    editMaxAge: function() {
      console.log("editMaxAge", this.index, this.entry.name);
      this.tmpAge = 0;

      this.$Modal.confirm({
        render: h => {
          return h(AgeSelector, {
            props: {
              age: this.tmpAge,
              name: "过期时间"
            },
            on: {
              AgeSelectorChange: val => {
                console.log(this.tmpAge, " -> ", val);
                this.tmpAge = val;
              }
            }
          });
        },
        title: "编辑 过期时间：" + this.entry.name,
        closable: true,
        scrollable: true,
        onOk: () => {
          this.http
            .put(this.api.shareMaxAge + this.entry.id, {
              maxAge: this.tmpAge
            })
            .then(res => {
              console.log(res.status);
              if (res.status === 200) {
                this.$Message.success("编辑过期时间成功");

                this.update();
              } else {
                console.log(res.data.msg);
                this.$Message.error("编辑过期时间失败: " + res.data.msg);
              }
            })
            .catch(e => e);
        }
      });
    },
    remove: function() {
      this.http
        .delete(this.api.shareDelete + this.entry.id)
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

