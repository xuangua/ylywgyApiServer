<template>
  <div class="admin-update">
    <div class="update-info">
      <h4>修改个人信息</h4>
      <input type="text" v-model="user.nickname">
      <input type="text" v-model="user.motto">
      <input type="text" v-model="user.email">
      <button @click="updateInfo">确认修改</button>
    </div>
    <div class="update-password">
      <h4>修改密码</h4>
      <input type="password" v-model="oldPassword" placeholder="输入旧密码">
      <input type="password" v-model="newPassword" placeholder="输入新密码">
      <input type="password" v-model="vertifyPassword" placeholder="再次输入新密码">
      <button @click="updatePassword">确认修改</button>
    </div>
    <top-tip ref="tip"/>
  </div>
</template>
<script>
export default {
  middleware: 'auth',
  data () {
    return {
      oldPassword: '',
      newPassword: '',
      vertifyPassword: '',
      user: {}
    }
  },
  mounted () {
    this.$store.dispatch('ADMIN_INFO').then((data) => {
      this.user = data.data
    })
  },
  head () {
    return {
      title: '修改信息 - ' + this.$store.state.user.nickname
    }
  },
  methods: {
    updateInfo () {
      this.$store.dispatch('UPDATE_ADMIN', this.user).then((data) => {
        if (data.success) {
          this.$refs.tip.openTip('信息修改完成')
        }
      })
    },
    updatePassword () {
      if (!this.oldPassword || !this.newPassword || !this.vertifyPassword) {
        return false
      }
      if (this.newPassword !== this.vertifyPassword) {
        this.$refs.tip.openTip('两次密码不一致！')
        return false
      }
      this.$store.dispatch('UPDATE_ADMIN', { oldPassword: this.oldPassword, newPassword: this.newPassword }).then((data) => {
        if (data.success) {
          this.$refs.tip.openTip('密码重置完成')
          // clear token
          this.$store.dispatch('LOGOUT').then(ret => {
            if (ret.success) {
              this.$store.state.token = ''
              this.$router.push('/login')
            }
          })
        }
      })
    }
  }
}
</script>
