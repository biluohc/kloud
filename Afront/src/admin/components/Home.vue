<template>
    <div class="layout">
        <Layout>
            <Header>
                <Menu mode="horizontal" theme="dark" active-name="1">
                    <div class="layout-logo"><router-link to="/admin"><img src= "/static/favicon.ico"></router-link> </div>
                    <div class="layout-nav">
                        <MenuItem name="0">
                          <router-link to="/home">                        
                           <div class="user">
                            <Icon type="android-home"></Icon>
                            Home
                           </div>
                          </router-link>                 
                        </MenuItem>
                        <MenuItem name="1">
                        <router-link to="/settings">                        
                          <Avatar v-if="avatarURL" :src="avatarURL" />
                          <Avatar v-if="!avatarURL" icon="person" />
                         <span class="user"> {{ userName }} </span>
                         </router-link>
                        </MenuItem>
                        <MenuItem/>
                        <MenuItem name="2">
                           <div @click="logout" class="user">
                            <Icon type="log-out"></Icon>
                            登出
                           </div>                          
                        </MenuItem>
                    </div>
                </Menu>
            </Header>
            <Layout :style="{background: '#fff'}">
                <Sider hide-trigger :style="{background: '#fff'}">
                    <Menu theme="dark" active-name="1"  width="auto" id="menu-box">
                            <MenuItem name="1" class="nav-item" >
                               <router-link to="/admin">
                               <Icon type="stats-bars"></Icon>
                                系统概况
                               </router-link>
                            </MenuItem>
                            <MenuItem></MenuItem>
                            <MenuItem name="1" class="nav-item" >
                               <router-link to="/admin/files">                            
                                <Icon type="filing"></Icon>
                                文件列表
                               </router-link>
                            </MenuItem>
                            <MenuItem></MenuItem>
                            <MenuItem name="2" class="nav-item" >
                               <router-link to="/admin/users">                            
                                <Icon type="android-people"></Icon>
                                用户列表
                               </router-link>
                            </MenuItem>
                    <!-- 这里考虑加上用户容量情况 -->
                    <!-- 加上版权等信息 -->
                    </Menu>
                </Sider>
<Content :style="{background: '#fff'}">
 <slot name="home"></slot>
</Content>
</Layout>
</Layout>
</div>
</template>

<script>
export default {
  name: "Home",
  components: {},
  created: function() {
    this.$store.dispatch("checkAdmin");
  },
  data: function() {
    return {};
  },
  computed: {
    avatarURL: function() {
      if (this.$store.state.isLogin) {
        return this.$store.state.userInfo.avatarURL;
      }
      return "";
    },
    userName: function() {
      if (this.$store.state.isLogin) {
        return this.$store.state.userInfo.name;
      }
      return "";
    }
  },
  methods: {
    logout() {
      this.$store.dispatch("logout");
      this.$Message.success("已登出！");
      this.$router.replace("/login/");
    }
  }
};
</script>

<style lang='scss' scoped>
.layout {
  border: 1px solid #d7dde4;
  background: #f5f7f9;
  position: relative;
  border-radius: 4px;
  overflow: hidden;
}
.layout-logo {
  width: 100px;
  height: 30px;
  /* background: #5b6270; */
  border-radius: 3px;
  float: left;
  position: relative;
  top: 15px;
  left: 20px;
}
.layout-nav {
  width: 420px;
  margin: 0 auto;
  margin-right: 20px;
  font-size: large;
}

#menu-box {
  background: #444;
  /* border-radius: 3px; */
  position: fixed;
  opacity: 0.88;
  height: 100%;
  min-height: 100%;
}

.nav-item {
  font-size: x-large;
  font-weight: bold;
}

.user {
  font-size: large;
  font-weight: bold;

  .logout {
    left: 20%;
  }
}
</style>