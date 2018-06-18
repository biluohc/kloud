<template>
<Home>
  <template slot="home">
    <div v-if="isLogin">
      <div>
        <Row v-if="selected">
                <Col span="2"><bold class="count">{{ selectedNum }} </bold><Button type="primary" size="large" v-on:click="selectAll">全选</Button></Col>
                <Col span="2"><Button type="success" size="large" v-on:click="untrashSelected">恢复</Button></Col>      <Col span="2"><Button type="error" size="large" v-on:click="deleteSelected">彻底删除</Button></Col>
                <Col span="2" offset="16"><Button  type="error" size="large" v-on:click="removeAll">清空回收站</Button></Col>
        </Row>
        <Row v-else>
                <Col span="2"><Button type="primary" size="large" v-on:click="selectAll">全选</Button></Col>
                <Col span="2" offset="20"><Button type="error" size="large" v-on:click="removeAll">清空回收站</Button></Col>
        </Row>
      </div>
    <Row>
         <div class="count attr">
            <Col span="12" offset="1">名字</Col>
            <Col span="1" >大小</Col>
            <Col span="3" offset="2" >删除时间</Col>
        </div>
    </Row>
      <div v-for="(entry,index) in trashes" :key='""+selected+trashes.length+entry.id'>
      <TrashedEntry :entry="entry" :index="index" @update="fetchData" @check="entryCheck" />
      </div>
    </div>
  </template>
</Home>
</template>

<script>
import Home from "@/components/Home.vue";
import TrashedEntry from "@/components/TrashedEntry.vue";

export default {
  name: "HomeSettings",
  components: { Home, TrashedEntry },
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
      trashes: [],
      tmpCount: 0
    };
  },
  computed: {
    isLogin: function() {
      return this.$store.state.isLogin;
    },
    selected: function() {
      return this.selectedNum > 0;
    },
    selectedNum: function() {
      var count = 0;
      for (var i = 0; i < this.trashes.length; i++) {
        if (this.trashes[i].state === true) {
          count++;
        }
      }
      return count;
    }
  },
  methods: {
    fetchData: function() {
      var u = this.$store.state.userInfo;
      if (u !== undefined && u !== null && u.id !== "") {
        this.http
          .get(this.api.entryTrashes)
          .then(res => {
            console.log(res.status);
            if (res.status === 200) {
              this.updateTrashes(res.data.data);
            } else {
              this.$Message.error(res.data.msg);
            }
          })
          .catch(e => e);
      }
    },
    updateTrashes: function(val) {
      console.log(this.trashes.length, "->", val.length);
      for (var i = 0; i < val.length; i++) {
        var tmp = val[i];

        tmp.createTimeStr = this.api.dateTime(tmp.createTime);
        tmp.deleteTimeStr = this.api.dateTime(tmp.deleteTime);

        if (!tmp.deleteTimeStr) {
          tmp.deleteTimeStr = "-";
        }

        tmp.state = false;
        if (tmp.fileId === "") {
          tmp.sizeStr = "-";
        } else {
          tmp.sizeStr = this.api.size(tmp.file.size);
        }
        // 属性要先加入，加入后再改变Vue就日狗了。。

        this.trashes.splice(i, 1, tmp);
      }
      this.trashes.length = val.length;
      // 为0时只剩改长度，vue检测不到，垃圾js
      if (val.length == 0) {
        this.trashes.splice(0, this.trashes.length);
      }
    },
    entryCheck: function(i, val) {
      console.log(i, this.trashes[i].name, this.trashes[i].state, " -> ", val);
      this.trashes[i].state = val;
    },
    selectAll: function() {
      console.log("selectAll", this.selected, " -> ", !this.selected);
      var all = this.selected;
      for (var i = 0; i < this.trashes.length; i++) {
        this.entryCheck(i, !all);
      }
      console.log("selectedNum: ", this.selectedNum);
    },
    untrashSelected: function() {
      this.tmpCount = 0;

      for (var i = 0; i < this.trashes.length; i++) {
        if (this.trashes[i].state) {
          this.http
            .put(this.api.entryUnTrash + this.trashes[i].id)
            .then(res => {
              console.log(res.status);
              if (res.status === 200) {
                this.$Message.success("恢复成功");
                this.tmpCount++;
                if (this.tmpCount === this.selectedNum) {
                  this.fetchData();
                }
              } else {
                console.log(res.data.msg);
                this.$Message.error("恢复失败: " + res.data.msg);
                this.fetchData();
              }
            })
            .catch(e => e);
        }
      }
    },
    deleteSelected: function() {
      this.tmpCount = 0;

      for (var i = 0; i < this.trashes.length; i++) {
        if (this.trashes[i].state) {
          this.http
            .delete(this.api.entryDelete + this.trashes[i].id)
            .then(res => {
              console.log(res.status);
              if (res.status === 200) {
                this.$Message.success("删除成功");

                this.tmpCount++;
                if (this.tmpCount === this.selectedNum) {
                  this.fetchData();
                }
              } else {
                console.log(res.data.msg);
                this.$Message.error("删除失败: " + res.data.msg);
                this.fetchData();
              }
            })
            .catch(e => e);
        }
      }
    },
    removeAll: function() {
      this.tmpCount = 0;

      for (var i = 0; i < this.trashes.length; i++) {
        this.http
          .delete(this.api.entryDelete + this.trashes[i].id)
          .then(res => {
            console.log(res.status);
            if (res.status === 200) {
              this.$Message.success("删除成功");
              this.tmpCount++;
              if (this.tmpCount === this.trashes.length) {
                this.fetchData();
              }
            } else {
              console.log(res.data.msg);
              this.$Message.error("删除失败: " + res.data.msg);
              this.fetchData();
            }
          })
          .catch(e => e);
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