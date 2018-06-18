<template>
<div class="path" >
  <span ><router-link :to='"/home/"+paths[0].id' :class="1===paths.length&&'focus'">全部文件</router-link></span>
  <span v-for="(path,index) in paths" :key="path.id" v-if="path.name!=='/'">
    <router-link :to='"/home/"+ path.id' :class="index+1===paths.length&&'focus'">{{ path.name }}/</router-link>
  </span>
</div>
</template>

<script>
//  Path貌似被某个占用了，一点反应都没有

export default {
  name: "Paths",
  props: {
    paths: Array
  },
  data: function() {
    return {};
  },
  methods: {},
  computed: {}
};
</script>

<style lang="scss" scoped>
.path {
  position: relative;
  padding: 1% 0% 1% 0%;
  background: #fff;
  font-size: large;
  font-weight: bold;
  color: #2d8cf0;

  .span {
    font-size: large;
    font-weight: bold;
  }
  .focus {
    color: #ff00a9;
  }
  a {
    position: relative;
    margin-left: 16px;
  }
}
</style>
