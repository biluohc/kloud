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
            <Col span="4" >Referer正则</Col>
            <Col span="3" offset="1" >过期时间</Col>
            <Col span="3" offset="2" >公开时间</Col>
        </div>
    </Row>
      <div v-for="(entry,index) in publics" :key='""+ selected+ publics.length+ entry.id'>
      <PublicedEntry :entry="entry" :index="index" @update="fetchData" @check="entryCheck" />
      </div>
    </div>
  </template>
</Home>
</template>

<script>
import Home from "@/components/Home.vue";
import PublicedEntry from "@/components/PublicedEntry.vue";

export default {
  name: "Homepublics",
  components: { Home, PublicedEntry },
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
      publics: [],
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
      for (var i = 0; i < this.publics.length; i++) {
        if (this.publics[i].state === true) {
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
          .get(this.api.publics)
          .then(res => {
            console.log(res.status);
            if (res.status === 200) {
              this.updatePublics(res.data.data);
            } else {
              this.$Message.error(res.data.msg);
            }
          })
          .catch(e => e);
      }
    },
    updatePublics: function(val) {
      console.log(this.publics.length, "->", val.length);
      for (var i = 0; i < val.length; i++) {
        var tmp = val[i];

        tmp.createTimeStr = this.api.dateTime(tmp.createTime);
        tmp.expireTimeStr = this.api.dateTime(tmp.expireTime);

        if (!tmp.expireTimeStr) {
          tmp.expireTimeStr = "-";
        }

        tmp.state = false;
        // 属性要先加入，加入后再改变Vue就日狗了。。

        this.publics.splice(i, 1, tmp);
      }
      this.publics.length = val.length;
      // 为0时只剩改长度，vue检测不到，垃圾js
      if (val.length == 0) {
        this.publics.splice(0, this.publics.length);
      }
    },
    entryCheck: function(i, val) {
      console.log(i, this.publics[i].name, this.publics[i].state, " -> ", val);
      this.publics[i].state = val;
    },
    selectAll: function() {
      console.log("selectAll", this.selected, " -> ", !this.selected);
      var all = this.selected;
      for (var i = 0; i < this.publics.length; i++) {
        this.entryCheck(i, !all);
      }
      console.log("selectedNum: ", this.selectedNum);
    },
    deleteSelected: function() {
      this.tmpCount = 0;

      for (var i = 0; i < this.publics.length; i++) {
        if (this.publics[i].state) {
          this.http
            .delete(this.api.publicDelete + this.publics[i].id)
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