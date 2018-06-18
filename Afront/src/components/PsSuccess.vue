<template>
<div>
<h1>{{ title }}</h1>
<span><a :href="url" target="_blank" ref="url">{{ tmpUrl }}</a>
<Tooltip content="点击复制" placement="top">
<Button type="primary" shape="circle" icon="clipboard" v-on:click="copyLink" class="btn" :data-clipboard-text="tmpUrl">
</Button></Tooltip></span>
</div>
</template>

<script>
export default {
  name: "PsSuccess",
  props: {
    title: String,
    url: String
  },
  data: function() {
    return {
      tmpUrl: this.url
    };
  },
  mounted: function() {
    // Vue里要用 ref 来访问元素/组件，直接原生js容易碰上undefind
    // 这是因为 share生成的url只有路径部分，没有协议和host，所以要补上。
    var tmpUrl = this.$refs.url.href;
    console.log(tmpUrl);
    this.tmpUrl = tmpUrl;
  },
  methods: {
    copyLink: function() {
      this.api.copyLink();
    }
  },
  computed: {}
};
</script>

<style lang="scss" scoped>
.span {
  font-size: large;
  font-weight: bold;
}

.a {
  right: 20px;
}
</style>
