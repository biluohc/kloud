<template>
<div class="Ps">
<span>{{ valueName }}:  <input type="text" v-model.trim="tmpValue" @input="input" :placeholder="placeHolder"></input></span>
<p/>
 <AgeSelector :name="'过期时间'" :age="age"  @AgeSelectorChange="ageChange" />
</div>
</template>

<script>
import AgeSelector from "@/components/AgeSelector.vue";

export default {
  name: "PsEditor",
  components: {
    AgeSelector
  },
  props: {
    value: String,
    valueName: String,
    placeHolder: String
  },
  data: function() {
    return {
      age: 0,
      tmpValue: this.value
    };
  },
  created: function() {
    this.change();
  },
  methods: {
    change: function() {
      this.$emit("PsEditorChange", this.tmpValue, parseInt(this.age));
    },
    input: function() {
      this.change();
    },
    ageChange: function(val) {
      this.age = val;
      this.change();
    }
  }
};
</script>

<style lang="scss" scoped>
span {
  right: 20px;
  font-size: large;
  //   font-weight: bold;
}
</style>
