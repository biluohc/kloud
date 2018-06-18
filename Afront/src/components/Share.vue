<style lang="scss">
.layout {
  //   border: 1px solid #d7dde4;
  background: #f5f7f9;
  position: relative;
  border-radius: 4px;
  overflow: hidden;
  opacity: 0.9;
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
}
.layout-footer-center {
  position: relative;
  height: 100px;
  text-align: center;
}
.login {
  .box {
    width: 300px;
    margin: 80px auto;
    padding: 70px 20px;
    background: #444;
    border-radius: 3px;
  }

  .get-access-token {
    margin-top: 10px;

    text-align: center;
    a {
      color: DodgerBlue;

      &:hover {
        text-decoration: underline;
      }
    }
  }

  .input {
    height: 45px;

    input {
      width: 100%;
      height: 100%;
      padding: 0 0 0 8px;
      box-sizing: border-box;
      font-size: 16px;
    }
  }

  .last input {
    border-radius: 0 0 5px 5px;
  }

  .submit {
    position: relative;
    margin-top: 15px;
    padding: 10px 0;
    background: #f9671e;
    color: #fff;
    text-align: center;
    border-radius: 5px;
    font-size: 16px;
    font-weight: 500;
    cursor: pointer;

    &:after {
      content: "";
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: #000;
      border-radius: 5px;
      opacity: 0;
    }

    &:active:after {
      opacity: 0.2;
    }
  }
}
</style>
<template>
    <div class="layout">
        <Layout>
            <Header>
                <Menu mode="horizontal" theme="dark" active-name="1">
                    <div class="layout-logo">
                        <router-link to="/home"><img src= "/static/favicon.ico"></router-link>
                    </div>
                    <div class="layout-nav">
                    </div>
                </Menu>
            </Header>
            <Content :style="{padding: '0 50px', background: 'rgba(255,255,255,0.2)'}">
              <div class="main"><slot name="share"></slot></div>
            </Content>
            <Footer class="layout-footer-center">
                2017-2018 &copy;
                <a :href="this.api.github" target="_blank">kloud <Icon type="social-github"></Icon></a>
            </Footer>
        </Layout>
    </div>
</template>

<script>
export default {
  name: "Share",
  components: {},
  props: {},
  data: function() {
    return {};
  },
  methods: {}
};
</script>