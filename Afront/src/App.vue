<template>
	<div id="app">
		  <div class="router-view">
        <router-view />
      </div>
	</div>
</template>

<script>
export default {
  name: "App",
  components: {},
  data() {
    return {};
  },
  created() {},
  methods: {}
};
</script>

<style lang="scss">
html body {
  background: url(/static/ffsend_bg.svg);
}
</style>
