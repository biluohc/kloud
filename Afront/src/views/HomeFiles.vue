<template>
<Home>
  <template slot="home">
  <div v-if="curEntry">
      <div>
        <Row v-if="selected">
                <Col span="2"><span class="count">{{ selectedNum }}</span><Button type="primary" size="large" v-on:click="selectAll">全选</Button></Col>       
                <Col span="2"><Button type="error" size="large" v-on:click="deleteSelected">删除</Button></Col>
                <Col span="2"><Button type="primary" size="large" v-on:click="moveSelectedTo">移动到</Button></Col>
                <Col span="2"><Button type="primary" size="large" v-on:click="makedir">+ 创建文件夹</Button></Col>
                <Col span="2"><Upload action="/apiv1/entry/file">
            <Button type="ghost" icon="ios-cloud-upload-outline">上传文件</Button></Upload></Col>
        </Row>
        <Row v-else>
                <Col span="2"><Button type="primary" size="large" v-on:click="selectAll">全选</Button></Col>
                <Col span="2"><Button type="primary" size="large" v-on:click="makedir">+ 创建文件夹</Button></Col>
                <Col span="2">
                <Upload :action="fileApi" name="file"
                :with-credentials="true" :show-upload-list="false" :on-success="uploadSuccess" :on-error="uploadFailed"
                >
            <Button type="ghost" icon="ios-cloud-upload-outline">上传文件</Button></Upload></Col>
        </Row>
      </div>
  <Paths :paths="paths"/>
    <Row>
         <div class="count attr">
            <Col span="12" offset="1">名字</Col>
            <Col span="1" >大小</Col>
            <Col span="3" offset="2" >修改时间</Col>
        </div>
    </Row>
      <div v-for="(entry,index) in subs" :key='""+selected+subs.length+entry.id'>
      <Entry :entry="entry" :index="index" :curEntry="curEntry" :subs="subs" @update="fetchSubs" @check="entryCheck" />
      </div>
    </div>
  </template>
</Home>
</template>

<script>
import Home from "@/components/Home.vue";
import Paths from "@/components/Paths.vue";
import PathTree from "@/components/PathTree.vue";
import Entry from "@/components/Entry.vue";

export default {
  name: "HomeFiles",
  components: { Home, Paths, PathTree, Entry },
  created: function() {
    var id = (this.$root && this.$route.params.id) || "";
    this.fetchData(id);
  },
  data: function() {
    return {
      loading: false,
      //   从 / 到 当前的结构数组。
      paths: [],
      //   当前的子条目
      subs: [],

      // 新建文件夹的缓冲值
      tmpDir: "新建文件夹",
      // 删除/移动选中的一批文件个数成功/失败后 fetchSubs
      tmpMoveToId: "",
      tmpSelectedNum: 0,
      tmpCount: 0
    };
  },
  computed: {
    userInfo: function() {
      return this.$store.state.userInfo;
    },
    curEntry: function() {
      if (this.paths.length > 0) {
        return this.paths[this.paths.length - 1];
      }
      return null;
    },
    // 上传文件的api
    fileApi: function() {
      return this.api.entryFilePost + this.curEntry.id;
    },
    selected: function() {
      return this.selectedNum > 0;
    },
    selectedNum: function() {
      var count = 0;
      for (var i = 0; i < this.subs.length; i++) {
        if (this.subs[i].state === true) {
          count++;
        }
      }
      return count;
    }
  },
  watch: {
    $route(to) {
      console.log("watch$routeTo: '", to.params.id, "'");
      this.fetchData(to.params.id);
    }
  },
  methods: {
    fetchData: function(paramsId) {
      var id = "";
      // 当发现 this.$route.params.id 是字符串 "undefined" 时，真是醉了。。
      if (paramsId && paramsId !== "undefined") {
        id = paramsId;
        // id不空时才有必要获取目录
        this.fetchSubs(id);
      }
      console.log(id);
      this.http
        .get(this.api.entryPath + id)
        .then(res => {
          console.log(res.status, this.paths.length);
          var lastCurEntryId = "";
          if (this.paths.length > 0) {
            lastCurEntryId = this.paths[this.paths.length - 1];
          }
          this.updatePaths(res.data.data);

          var curEntry = this.paths[this.paths.length - 1];

          if (this.$store.state.rootEntryId !== this.paths[0].id) {
            this.$store.commit("updateRootEntryId", this.paths[0].id);
          }

          if (!this.$route.params.id || this.$route.params.id !== curEntry.id) {
            console.log("this.$router.replaceTo: /home/" + curEntry.id);
            this.$router.replace("/home/" + curEntry.id);
          }
        })
        .catch(e => e);
    },
    fetchSubs: function(eid) {
      var id = eid || this.$route.params.id;
      this.http
        .get(this.api.entrySubs + id)
        .then(res => {
          console.log(res.status);
          if (res.status === 200) {
            this.updateSubs(res.data.data);
          }
        })
        .catch(e => e);
    },
    updatePaths: function(val) {
      console.log(this.paths.length, "->", val.length);

      this.paths = val;
    },
    updateSubs: function(val) {
      console.log(this.subs.length, "->", val.length);
      for (var i = 0; i < val.length; i++) {
        var tmp = val[i];

        tmp.createTimeStr = this.api.dateTime(tmp.createTime);
        tmp.modifyTimeStr = this.api.dateTime(tmp.modifyTime);
        tmp.state = false;
        if (tmp.fileId === "") {
          tmp.sizeStr = "-";
        } else {
          tmp.sizeStr = this.api.size(tmp.file.size);
        }
        // 属性要先加入，加入后再改变Vue就日狗了。。

        this.subs.splice(i, 1, tmp);
      }
      this.subs.length = val.length;
      // 为0时只剩改长度，vue检测不到，垃圾js
      if (val.length == 0) {
        this.subs.splice(0, this.subs.length);
      }
    },
    uploadSuccess: function() {
      console.log("上传成功");
      this.$Message.success("上传成功");
      this.fetchSubs();
    },
    uploadFailed: function() {
      console.log("上传失败");
      this.$Message.error("上传失败");
    },
    entryCheck: function(i, val) {
      console.log(i, this.subs[i].name, this.subs[i].state, " -> ", val);
      this.subs[i].state = val;
    },
    selectAll: function() {
      console.log("selectAll", this.selected, " -> ", !this.selected);
      var all = this.selected;
      for (var i = 0; i < this.subs.length; i++) {
        this.entryCheck(i, !all);
      }
      console.log("selectedNum: ", this.selectedNum);
    },
    makedir: function() {
      console.log("makedir");
      this.$Modal.confirm({
        render: h => {
          return h("Input", {
            props: {
              value: this.tmpDir,
              autofocus: true
            },
            on: {
              input: val => {
                console.log(this.tmpDir, " -> ", val);
                this.tmpDir = val;
                console.log(this.api.isValidEntryName(this.tmpDir));
              }
            }
          });
        },
        title: "新建文件夹",
        closable: true,
        scrollable: true,
        onCancel: () => {},
        onOk: () => {
          if (!this.api.isValidEntryName(this.tmpDir)) {
            return this.$Message.error("名字不能含/或\\且不能是.或者..");
          }
          for (var i2 = 0; i2 < this.subs.length; i2++) {
            if (this.tmpDir === this.subs[i2].name) {
              return this.$Message.error(
                "'" + this.tmpDir + "' 在当前目录已经存在"
              );
            }
          }
          this.http
            .post(this.api.entryDirPost + this.curEntry.id, {
              name: this.tmpDir
            })
            .then(res => {
              console.log(res.status);
              if (res.status === 200) {
                this.$Message.success("新建文件夹成功");
                this.fetchSubs();
              } else {
                console.log(res.data.msg);
                this.$Message.error("新建文件夹失败: " + res.data.msg);
              }
            })
            .catch(e => e);
        }
      });
    },
    moveSelectedTo: function() {
      this.tmpSelectedNum = this.selectedNum;
      this.tmpCount = 0;

      // 初始化
      this.tmpMoveToId = this.$store.state.rootEntryId;
      console.log("moveSelectedTo:", this.tmpCount, "/", this.tmpSelectedNum);

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
        onCancel: () => {},
        onOk: () => {
          for (var i = 0; i < this.subs.length; i++) {
            if (this.subs[i].state) {
              if (this.tmpMoveToId === this.curEntry.id) {
                return;
              }

              if (this.tmpMoveToId === this.subs[i].id) {
                this.$Message.error("不能移动到自身或其子目录下");
                return;
              }
              this.http
                .put(this.api.entryMoveto + this.tmpMoveToId, {
                  id: this.subs[i].id
                })
                .then(res => {
                  if (res.status === 200) {
                    this.$Message.success("移动成功");
                    this.tmpCount++;
                    console.log(this.tmpCount, "/", this.tmpSelectedNum);
                    if (this.tmpCount == this.tmpSelectedNum) {
                      this.fetchSubs();
                    }
                  } else {
                    console.log(res.data.msg);
                    this.$Message.error("移动失败: " + res.data.msg);
                    this.fetchSubs();
                  }
                })
                .catch(e => e);
            }
          }
        }
      });
    },
    deleteSelected: function() {
      console.log("deleteSelected", this.selectedNum);
      this.tmpSelectedNum = this.selectedNum;
      this.tmpCount = 0;

      console.log("Trash:", this.tmpCount, "/", this.tmpSelectedNum);

      for (var i = 0; i < this.subs.length; i++) {
        if (this.subs[i].state) {
          this.http
            .put(this.api.entryTrash + this.subs[i].id)
            .then(res => {
              if (res.status === 200) {
                this.$Message.success("移动到回收站成功");
                this.tmpCount++;
                console.log(this.tmpCount, "/", this.tmpSelectedNum);
                if (this.tmpCount == this.tmpSelectedNum) {
                  this.fetchSubs();
                }
              } else {
                console.log(res.data.msg);
                this.$Message.error("移动到回收站失败: " + res.data.msg);
                this.fetchSubs();
              }
            })
            .catch(e => e);
        }
      }
    }
  }
};
</script>

<style scoped>
.count {
  font-size: large;
  font-weight: bold;
  margin-right: 20px;
}

.attr {
  padding: 10px, 0, 10px, 0;
}
</style>