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
    <Col span="3">{{ entry.referer }}</Col>
    <Col span="1" v-if="entry.referer">
      <Tooltip content="编辑 Referer正则" placement="left">
          <Button type="primary" shape="circle" icon="edit" v-on:click="editReferer"></Button>
      </Tooltip>
    </Col>
    <Col span="1" offset="3" v-if="!entry.referer">
      <Tooltip content="编辑 Referer正则" placement="left">
          <Button type="primary" shape="circle" icon="edit" v-on:click="editReferer"></Button>
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
      <Tooltip content="删除此公开" placement="left">
          <Button type="error" shape="circle" icon="close" v-on:click="remove"></Button>
      </Tooltip>
    </Col>
</Row>
</div>
</template>

<script>
import AgeSelector from "@/components/AgeSelector.vue";

export default {
  name: "PublicedEntry",
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
      return this.api.publicFile + this.entry.id + "/" + this.entry.name;
    }
  },

  methods: {
    update: function() {
      this.$emit("update");
    },
    check: function() {
      this.$emit("check", this.index, this.tmpState);
    },
    editReferer: function() {
      console.log("editReferer", this.index, this.entry.name);
      this.tmpValue = this.entry.referer;

      this.$Modal.confirm({
        render: h => {
          return h("Input", {
            props: {
              value: this.tmpValue,
              autofocus: true,
              placeholder: "空表示允许任意 Referer"
            },
            on: {
              input: val => {
                console.log(this.tmpValue, " -> ", val);
                this.tmpValue = val;
              }
            }
          });
        },
        title: "编辑 Referer正则：" + this.entry.name,
        closable: true,
        scrollable: true,
        onOk: () => {
          if (this.tmpValue === this.entry.referer) {
            return;
          }
          this.http
            .put(this.api.publicReferer + this.entry.id, {
              referer: this.tmpValue
            })
            .then(res => {
              console.log(res.status);
              if (res.status === 200) {
                this.$Message.success("编辑 Referer正则 成功");

                this.update();
              } else {
                console.log(res.data.msg);
                this.$Message.error("编辑 Referer正则 失败: " + res.data.msg);
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
            .put(this.api.publicMaxAge + this.entry.id, {
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
        .delete(this.api.publicDelete + this.entry.id)
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

