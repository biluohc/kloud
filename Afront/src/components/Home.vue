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

.state {
  left: 20%;
  color: #cccc33;
  text-align: center;
  font-weight: bold;
  font-size: large;
}
</style>
<template>
    <div class="layout">
        <Layout>
            <Header>
                <Menu mode="horizontal" theme="dark" active-name="1">
                    <div class="layout-logo"><router-link to="/home"><img src= "/static/favicon.ico"></router-link> </div>
                    <div class="layout-nav">
                        <MenuItem name="0" v-if="user.role===0">
                          <router-link to="/admin">                        
                           <div class="user">
                            <Icon type="stats-bars"></Icon>
                            系统管理
                           </div>
                          </router-link>                 
                        </MenuItem>
                        <MenuItem name="1">
                        <router-link to="/settings">                        
                          <Avatar v-if="user.avatarURL" :src="user.avatarURL" />
                          <Avatar v-if="!user.avatarURL" icon="person" />
                         <span class="user"> {{ user.name }} </span>
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
                               <router-link to="/home">
                               <Icon type="folder"></Icon>
                                全部文件
                               </router-link>
                            </MenuItem>
                            <MenuItem></MenuItem>
                            <MenuItem name="1" class="nav-item" >
                               <router-link to="/trashes">                            
                                <Icon type="ios-trash-outline"></Icon>
                                回收站
                               </router-link>
                            </MenuItem>
                            <MenuItem></MenuItem>
                            <MenuItem name="2" class="nav-item" >
                               <router-link to="/shares">
                            <Icon type="android-share-alt"></Icon>
                                我的分享
                               </router-link>
                            </MenuItem>
                            <MenuItem></MenuItem>            
                            <MenuItem name="3" class="nav-item" >
                               <router-link to="/publics">                            
                                <Icon type="earth"></Icon>
                                公共文件
                               </router-link>                                
                            </MenuItem>
                            <MenuItem></MenuItem>
                            <MenuItem name="4" class="nav-item" > 
                            </MenuItem> 
                            <div class="state">
                            <p>{{ spaceStr }}</p>                               
                            <p><Progress :percent="spacePercent"></Progress></p>
                            <p/>
                            <p/>
                            <p/>
                            <p>Powered by</p>
                              <a :href="this.api.github" target="_blank">                            
                                <Icon type="social-github"></Icon>
                                 kloud
                              </a>
                            </div>         
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
    this.$store.dispatch("checkLogin");
  },
  data: function() {
    return {};
  },
  computed: {
    user: function() {
      if (this.$store.state.isLogin) {
        return this.$store.state.userInfo;
      }
      return "";
    },
    spaceStr: function() {
      if (this.user) {
        return (
          this.api.size(this.user.curSize) +
          "/" +
          this.api.size(this.user.maxSize)
        );
      }
      return "";
    },
    spacePercent: function() {
      if (this.user) {
        console.log(this.user.curSize, this.user.maxSize);
        console.log(this.user.curSize / this.user.maxSize);
        console.log((this.user.curSize / this.user.maxSize).toFixed(0));
        return (this.user.curSize / this.user.maxSize).toFixed(0) * 100;
      }
      return 0;
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
