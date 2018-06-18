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
    <Col span="12" v-if="entry.fileId"><a target="_blank" :href="url">{{ entry.name }}</a></Col>
    <Col span="12" v-if="!entry.fileId"><router-link :to="url" >{{ entry.name }}/</router-link></Col>
    <Col span="1">{{ entry.sizeStr }}</Col>
    <Col span="3" offset="1">{{ entry.modifyTimeStr }}</Col>

    <Col span="1">
      <Tooltip content="重命名" placement="left">
          <Button type="primary" shape="circle" icon="edit" v-on:click="rename"></Button>
      </Tooltip>
    </Col>
    <Col span="1">
      <Tooltip content="移动" placement="left">
          <Button type="primary"  shape="circle" icon="arrow-move" v-on:click="moveto"></Button>
      </Tooltip>
    </Col>
    <Col span="1">
      <Tooltip content="删除，移动到回收站" placement="left">
          <Button type="error" shape="circle" icon="trash-a" v-on:click="trash"></Button>
      </Tooltip>
    </Col>
    <Col span="1">
      <Tooltip content="分享该文件" placement="left" v-if="entry.fileId">
          <Button type="primary" shape="circle" icon="android-share-alt" v-on:click="share"></Button>
      </Tooltip>
    </Col>
    <Col span="1">
      <Tooltip content="公开该文件" placement="left" v-if="entry.fileId">
          <Button type="primary" shape="circle" icon="earth" v-on:click="publik"></Button>
      </Tooltip>
    </Col>
</Row>
</div>
</template>

<script>
import PathTree from "@/components/PathTree.vue";
import PsEditor from "@/components/PsEditor.vue";
import PsSuccess from "@/components/PsSuccess.vue";

export default {
  name: "Entry",
  components: {
    PsEditor,
    PsSuccess
  },
  props: {
    curEntry: Object,
    entry: Object,
    index: Number,
    subs: Array
  },
  data: function() {
    return {
      tmpState: this.entry.state,
      tmpValue: "",
      tmpValueName: "",
      tmpAge: 0,
      tmpPlaceHolder: "",
      tmpMoveToId: ""
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
    rename: function() {
      console.log("reName", this.index, this.entry.name);
      this.tmpValue = this.entry.name;

      this.$Modal.confirm({
        render: h => {
          return h("Input", {
            props: {
              value: this.tmpValue,
              autofocus: true
            },
            on: {
              input: val => {
                console.log(this.tmpValue, " -> ", val);
                this.tmpValue = val;
                console.log(this.api.isValidEntryName(this.tmpDir));
              }
            }
          });
        },
        title: "重命名：" + this.entry.name,
        closable: true,
        scrollable: true,
        onOk: () => {
          if (!this.api.isValidEntryName(this.tmpDir)) {
            return this.$Message.error("名字不能含/或\\且不能是.或者..");
          }
          if (this.tmpValue === this.entry.name) {
            return;
          }
          for (var i2 = 0; i2 < this.subs.length; i2++) {
            if (this.tmpValue === this.subs[i2].name) {
              return this.$Message.error(
                "'" + this.tmpValue + "' 在当前目录已经存在"
              );
            }
          }
          this.http
            .put(this.api.entryRename + this.entry.id, {
              name: this.tmpValue
            })
            .then(res => {
              console.log(res.status);
              if (res.status === 200) {
                this.$Message.success("重命名成功");

                this.update();
              } else {
                console.log(res.data.msg);
                this.$Message.error("重命名失败: " + res.data.msg);
              }
            })
            .catch(e => e);
        }
      });
    },
    moveto: function() {
      this.tmpMoveToId = this.$store.state.rootEntryId;
      console.log("moveTo:", this.tmpMoveToId, this.index, this.entry.name);

      this.$Modal.confirm({
        render: h => {
          return h(PathTree, {
            on: {
              pathTreeChange: val => {
                if (val) {
                  this.tmpMoveToId = val;
                } else {
                  this.tmpMoveToId = this.$store.state.rootEntryId;
                }
              }
            }
          });
        },
        title: "移动到",
        closable: true,
        scrollable: true,
        onOk: () => {
          if (this.tmpMoveToId === this.curEntry.id) {
            return;
          }

          if (this.tmpMoveToId === this.entry.id) {
            this.$Message.error("不能移动到自身或其子目录下");
            return;
          }
          this.http
            .put(this.api.entryMoveto + this.tmpMoveToId, {
              id: this.entry.id
            })
            .then(res => {
              if (res.status === 200) {
                this.$Message.success("移动成功");
              } else {
                console.log(res.data.msg);
                this.$Message.error("移动失败: " + res.data.msg);
              }
              this.update();
            })
            .catch(e => e);
        }
      });
    },
    trash: function() {
      this.http
        .put(this.api.entryTrash + this.entry.id)
        .then(res => {
          if (res.status === 200) {
            this.$Message.success("移动到回收站成功");
          } else {
            console.log(res.data.msg);
            this.$Message.error("移动到回收站失败: " + res.data.msg);
          }
          this.update();
        })
        .catch(e => e);
    },
    share: function() {
      this.tmpValueName = "密码";
      this.tmpValue = "";
      this.tmpAge = 0;
      this.tmpPlaceHolder = "空表示没有密码";

      console.log("share:", this.index, this.entry.name);

      this.$Modal.confirm({
        render: h => {
          return h(PsEditor, {
            props: {
              value: this.tmpValue,
              valueName: this.tmpValueName,
              placeHolder: this.tmpPlaceHolder
            },
            on: {
              PsEditorChange: (val, age) => {
                this.tmpValue = val;
                this.tmpAge = age;
                console.log(this.tmpValue, this.tmpAge);
              }
            }
          });
        },
        title: "分享：" + this.entry.name,
        closable: true,
        scrollable: true,
        onOk: () => {
          this.http
            .post(this.api.entryShare + this.entry.id, {
              password: this.tmpValue,
              maxAge: this.tmpAge
            })
            .then(res => {
              if (res.status === 200) {
                this.tmpValue = "/#/share/" + res.data.data.id + "/info";

                this.$Message.success({
                  render: h => {
                    return (
                      <PsSuccess title={"分享文件成功"} url={this.tmpValue} />
                    );
                  },
                  duration: 8,
                  closable: true
                });
              } else {
                console.log(res.data.msg);
                this.$Message.error("分享文件失败: " + res.data.msg);
              }
            })
            .catch(e => e);
        }
      });
    },
    publik: function() {
      this.tmpValueName = "Referer正则";
      this.tmpValue = "";
      this.tmpAge = 0;
      this.tmpPlaceHolder = "空表示没有限制 Referer";

      console.log("public:", this.index, this.entry.name);

      this.$Modal.confirm({
        render: h => {
          return h(PsEditor, {
            props: {
              value: this.tmpValue,
              valueName: this.tmpValueName,
              placeHolder: this.tmpPlaceHolder
            },
            on: {
              PsEditorChange: (val, age) => {
                this.tmpValue = val;
                this.tmpAge = age;
                console.log(this.tmpValue, this.tmpAge);
              }
            }
          });
        },
        title: "公开：" + this.entry.name,
        closable: true,
        scrollable: true,
        onOk: () => {
          this.http
            .post(this.api.entryPublic + this.entry.id, {
              referer: this.tmpValue,
              maxAge: this.tmpAge
            })
            .then(res => {
              if (res.status === 200) {
                this.tmpValue =
                  this.api.publicFile +
                  res.data.data.id +
                  "/" +
                  res.data.data.name;

                this.$Message.success({
                  render: h => {
                    return (
                      <PsSuccess title={"公开文件成功"} url={this.tmpValue} />
                    );
                  },
                  duration: 8,
                  closable: true
                });
              } else {
                console.log(res.data.msg);
                this.$Message.error("公开文件失败: " + res.data.msg);
              }
            })
            .catch(e => e);
        }
      });
    }
  }
};
</script>

