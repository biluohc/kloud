<template>
<span>{{ name }} <select v-model="tmpAge">
    <option v-for="option in ageList"  :key="option.value" :value="option.value">
    {{ option.name }}
    </option>
 </select>
</span>
</template>

<script>
export default {
  name: "AgeSelector",
  props: {
    name: String,
    age: Number
  },
  data: function() {
    return {
      tmpAge: this.age || 0,
      ageList: [
        {
          name: "永久",
          value: 0
        },
        {
          name: "一天",
          value: 86400
        },
        {
          name: "一周",
          value: 86400 * 7
        },
        {
          name: "一月",
          value: 86400 * 30
        },
        {
          name: "半年",
          value: 86400 * 30 * 6
        },
        {
          name: "一年",
          value: 86400 * 365
        }
      ]
    };
  },
  created: function() {
    this.change();
  },
  watch: {
    tmpAge: function() {
      this.change();
    }
  },
  methods: {
    change: function() {
      this.$emit("AgeSelectorChange", parseInt(this.tmpAge));
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
