<template>
  <!--
  @author:风很大
  @description: 分页组件
  @pagination: Object， pagination: {  page: number,  size: number, total: number },
  @page:分页点击自定义组件
  @time: 2022/1/20 0020
  -->
  <div class="pagination">
    <el-pagination
      @current-change="handleCurrentChange"
      @size-change="handSizeChange"
      :current-page="pagination.page"
      :page-size="pagination.size"
      background
      layout="total, prev, pager, next, sizes"
      :total="pagination.total">
    </el-pagination>
  </div>
</template>

<script>
  export default {
    name: 'index',
    props:{
      pagination:{
        type:Object
      }
    },
    data(){
      return{}
    },
    methods:{
      /**
       * @Description: 分页点击自定义事件
       * @author 风很大
       * @page 点击事件
       * @date 2022/1/17 0017
      */
      handleCurrentChange:function(currentPage) {
        this.$emit('page',currentPage)
      },
      /**
       * @Description: 每页条数
       * @author 风很大
       * @page 点击事件
       * @date 2022/1/17 0017
       */
      handSizeChange:function(val) {
        this.$emit('size',val)
      }
    }
  }
</script>

<style scoped>
  .pagination {
    margin: 10px 0;
    float: right;
  }
</style>
