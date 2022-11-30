<template>
  <div class="container">
    <h4>本地文件上传</h4>
    <el-upload
      action="/api/upload/file"
      name="file"
      :auto-upload="true"
      list-type="picture-card"
      :on-preview="handlePictureCardPreview"
      :on-success="uploadSuccess"
    >
      <i class="el-icon-plus"></i>
    </el-upload>
    <el-dialog :visible.sync="dialogVisible">
      <img width="100%" :src="dialogImageUrl" alt="">
    </el-dialog>

    <h4>七牛文件上传</h4>
    <el-upload
      action="/api/upload/qiniu"
      name="file"
      :auto-upload="true"
      list-type="picture-card"
      :on-preview="handlePictureCardPreview"
      :on-success="uploadSuccess"
    >
      <i class="el-icon-plus"></i>
    </el-upload>
    <el-dialog :visible.sync="dialogVisible">
      <img width="100%" :src="dialogImageUrl" alt="">
    </el-dialog>

  </div>
</template>

<script>
export default {
  name: 'index',
  data() {
    return {
      dialogVisible: false,
      dialogImageUrl: '',
      image:''
    }
  },
  methods: {
    handlePictureCardPreview(file) {
      this.dialogImageUrl = file.url
      this.dialogVisible = true
    },
    uploadSuccess(res, file) {
      this.image = res.data.url;
    }
  }
}
</script>


