<template>
<div class="path-tree">
    <Tree :data="data"  :load-data="fetchSubs" @on-select-change="pathTreeChange"></Tree>
</div>
</template>

<script>
export default {
  name: "PathTree",
  data: function() {
    return {
      data: [
        {
          title: "/",
          id: "",
          loading: false,
          expand: true,
          children: []
        }
      ]
    };
  },
  methods: {
    fetchSubs: function(item, callback) {
      this.http
        .get(this.api.entrySubs + item.id)
        .then(res => {
          console.log(res.status);
          if (res.status === 200) {
            var childrens = res.data.data.filter(function(item) {
              item.title = item.name;
              item.expand = false;
              item.children = [];
              item.loading = false;
              return item.fileId === "";
            });

            callback(childrens);
          }
        })
        .catch(e => e);
    },
    // 见鬼的 /的 val是 undefined
    pathTreeChange: function(val) {
      console.log("pathTreeChange:", val[0].id, val[0].name);
      this.$emit("pathTreeChange", val[0].id);
    }
  }
};
</script>

<style lang="scss" scoped>
.path {
  position: relative;
  padding: 0% 1% 1% 0%;
  background: #fff;
  font-size: large;
  font-weight: bold;

  .span {
    font-size: large;
    font-weight: bold;
  }

  a {
    position: relative;
    color: #ccc;
    margin-left: 20px;

    &:hover {
      color: #fff;
    }
  }
}
</style>
