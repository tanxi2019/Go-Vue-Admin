<template>
<div class="container">
  <div class="wow pic animate__animated animate__backInLeft"></div>
  <div class="flex">
      <div
        class="wow card animate__animated animate__fadeInUp"
        v-for="(i, index) in arr"
        :key="index"></div>
  </div>

  <div
    class="item wow animate__animated animate__fadeInRight"
    v-for="(j, index) in 100"
    :key="index">
    <div class="square"></div>
    <div class="content">
      <div></div>
      <div></div>
    </div>
  </div>
</div>
</template>

<script>
export default {
  name: 'index',
  data(){
    return {
      arr:[1,2,3,4]
    }
  },
  methods: {},
  mounted(){
    new this.$wow.WOW().init()
  }
}
</script>

<style scoped>
.pic {
  width: 100%;
  height: 120px;
  margin: 20px 0;
  border-radius: 8px;
  background: #CAE5E8;
}
.item {
  margin: 5px 0;
  overflow: hidden;
}
.item div {
  display: inline;
  float: left;
}
.flex {
  display: flex;
  justify-content: space-between
}
.card {
  width: 65px;
  height: 65px;
  margin: 15px 0;
  margin-top: 0;
  border-radius: 8px;
  background: #CAE5E8;
}
.item::after {
  content: "";
  clear: both;
  display: block;
}
.square {
  width: 20%;
  height: 75px;
  background: #CAE5E8;
  border-radius: 8px;
}
.content {
  width: 79%;
}
.content div {
  height: 30px;
  margin: 15px;
  margin-top: 0;
  padding: 0;
  width: 100%;
  border-radius: 8px;
  line-height: 30px;
  background: #99D1D3;
}
.content div:last-child {
  width: 65%;
  background: #CAE5E8;
}
@keyframes one-in {
  from {
    padding-left: 100%;
    height: 0;
  }
  to {
    padding-left: 0;
    height: 100px;
  }
}
@keyframes two-in {
  from {
    padding-right: 100%;
    padding-top: 100px;
    height: 0;
  }
  to {
    padding-right: 0;
    padding-top: 0;
    height: 100px;
  }
}

</style>
