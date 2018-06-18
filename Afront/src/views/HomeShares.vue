<template>
<Home>
  <template slot="home">
    <div v-if="isLogin">
      <div>
        <Row v-if="selected">
            <Col span="2"><bold class="count">{{ selectedNum }} </bold><Button type="primary" size="large" v-on:click="selectAll">全选</Button></Col>
            <Col span="2"><Button type="error" size="large" v-on:click="deleteSelected">删除</Button></Col>
        </Row>
        <Row v-else>
            <Col span="2"><Button type="primary" size="large" v-on:click="selectAll">全选</Button></Col>
            <Col span="2"><Button type="error" size="large" v-on:click="deleteSelected">删除</Button></Col>
        </Row>
      </div>
    <Row>
         <div class="count attr">
            <Col span="9" offset="1">名字</Col>
            <Col span="4" >密码</Col>
            <Col span="3" offset="1" >过期时间</Col>
            <Col span="3" offset="3" >分享时间</Col>
        </div>
    </Row>
      <div v-for="(entry,index) in shares" :key='""+ selected+ shares.length+ entry.id'>
      <SharedEntry :entry="entry" :index="index" @update="fetchData" @check="entryCheck" />
      </div>
    </div>
  </template>
</Home>
</template>

<script>
import Home from "@/components/Home.vue";
import SharedEntry from "@/components/SharedEntry.vue";

export default {
  name: "HomeShares",
  components: { Home, SharedEntry },
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
      shares: [],
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
      for (var i = 0; i < this.shares.length; i++) {
        if (this.shares[i].state === true) {
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
          .get(this.api.shares)
          .then(res => {
            console.log(res.status);
            if (res.status === 200) {
              this.updateShares(res.data.data);
            } else {
              this.$Message.error(res.data.msg);
            }
          })
          .catch(e => e);
      }
    },
    updateShares: function(val) {
      console.log(this.shares.length, "->", val.length);
      for (var i = 0; i < val.length; i++) {
        var tmp = val[i];

        tmp.createTimeStr = this.api.dateTime(tmp.createTime);
        tmp.expireTimeStr = this.api.dateTime(tmp.expireTime);

        if (!tmp.expireTimeStr) {
          tmp.expireTimeStr = "-";
        }

        tmp.state = false;
        // 属性要先加入，加入后再改变Vue就日狗了。。

        this.shares.splice(i, 1, tmp);
      }
      this.shares.length = val.length;
      // 为0时只剩改长度，vue检测不到，垃圾js
      if (val.length == 0) {
        this.shares.splice(0, this.shares.length);
      }
    },
    entryCheck: function(i, val) {
      console.log(i, this.shares[i].name, this.shares[i].state, " -> ", val);
      this.shares[i].state = val;
    },
    selectAll: function() {
      console.log("selectAll", this.selected, " -> ", !this.selected);
      var all = this.selected;
      for (var i = 0; i < this.shares.length; i++) {
        this.entryCheck(i, !all);
      }
      console.log("selectedNum: ", this.selectedNum);
    },
    deleteSelected: function() {
      this.tmpCount = 0;

      for (var i = 0; i < this.shares.length; i++) {
        if (this.shares[i].state) {
          this.http
            .delete(this.api.shareDelete + this.shares[i].id)
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