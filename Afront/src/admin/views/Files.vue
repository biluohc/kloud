<template>
<Home>
  <template slot="home">
    <Row>
        <div class="count attr">
            <Col span="9">文件名</Col>
            <Col span="3">创建时间</Col>            
            <Col span="1" offset="1">大小</Col>
            <Col span="1" offset="1">计数</Col>
        </div>
    </Row>
    <div v-for="file in files" :key='file.id'>
      <File :file="file" @update="fetchData" />
    </div>
  </template>
</Home>
</template>

<script>
import Home from "@/admin/components/Home.vue";
import File from "@/admin/components/File.vue";

export default {
  name: "Files",
  components: { Home, File },
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
    files: []
    };
  },
  computed: {
    isLogin: function() {
      return this.$store.state.isLogin;
    }
  },
  methods: {
    fetchData() {
      if (this.isLogin) {
        this.http
          .get(this.api.files)
          .then(res => {
            console.log(res.status);
            if (res.status === 200) {
              this.updateFiles(res.data.data);
            } else {
              this.$Message.error(res.data.msg);
            }
          })
          .catch(e => e);
      }
    },
    updateFiles: function(val) {
      console.log(this.files.length, "->", val.length);
      for (var i = 0; i < val.length; i++) {
        var tmp = val[i];

        tmp.createTimeStr = this.api.dateTime(tmp.createTime);
        tmp.sizeStr = this.api.size(tmp.size);
        // 属性要先加入，加入后再改变Vue就日狗了。。

        this.files.splice(i, 1, tmp);
      }
      this.files.length = val.length;
      // 为0时只剩改长度，vue检测不到，垃圾js
      if (val.length == 0) {
        this.files.splice(0, this.files.length);
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.count {
  font-size: large;
  font-weight: bold;
  border: 1%, 0, 1%, 0;
}

.attr {
  padding: 10px, 0, 10px, 0;
}
</style>
