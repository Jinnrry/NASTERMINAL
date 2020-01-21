<template>

  <div class="container fullscreen">
    <form class="form-signin">
      <h1 class="h3 mb-3 font-weight-normal">需要登录</h1>

      <label for="inputPassword" class="sr-only">密码</label>
      <input v-model="pwd" type="password" id="inputPassword" class="form-control" placeholder="Password" required>

      <button @click="btn_login" class="btn btn-lg btn-primary btn-block" type="submit">登录</button>
    </form>

  </div>


</template>

<script>
  import CryptoJS from "crypto-js";
  import md5 from 'js-md5'

  export default {
    name: "login",
    data() {
      return {
        pwd: ""
      }
    },
    methods: {
      btn_login() {

        var key = this.CalcuMD5("12345678")

        console.log(key)
        var mi = this.Encrypt("Word data", key);
        console.log(mi)
        console.log(this.Decrypt(mi, key))
        // $.get("https://raw.githubusercontent.com/Jinnrry/NASTERMINAL/master/localhost/hostinfo").then((res) => {
        //   cosole.log(this.decrypt(res, this.CalcuMD5(this.pwd)))
        // })
      },
      // 解密
      Decrypt(word, key) {

        let iv = "";

        key = CryptoJS.enc.Utf8.parse(key);
        iv = CryptoJS.enc.Utf8.parse(iv);

        let base64 = CryptoJS.enc.Base64.parse(word);

        let src = CryptoJS.enc.Base64.stringify(base64);

        // 解密模式为CBC，补码方式为PKCS5Padding（也就是PKCS7）
        let decrypt = CryptoJS.AES.decrypt(src, key, {
          iv: iv,
          mode: CryptoJS.mode.CBC,
          padding: CryptoJS.pad.Pkcs7
        });

        let decryptedStr = decrypt.toString(CryptoJS.enc.Utf8);
        return decryptedStr.toString();

      },

      // 加密
      Encrypt(word, key) {
        let iv = "";

        key = CryptoJS.enc.Utf8.parse(key);
        iv = CryptoJS.enc.Utf8.parse(iv);

        let srcs = CryptoJS.enc.Utf8.parse(word);
        // 加密模式为CBC，补码方式为PKCS5Padding（也就是PKCS7）
        let encrypted = CryptoJS.AES.encrypt(srcs, key, {
          iv: iv,
          mode: CryptoJS.mode.CBC,
          padding: CryptoJS.pad.Pkcs7
        });

        //返回base64
        return CryptoJS.enc.Base64.stringify(encrypted.ciphertext);
      },
      CalcuMD5(pwd) {
        pwd = pwd.toUpperCase();
        pwd = md5(pwd);
        return pwd;
      }
    }
  }
</script>

<style scoped>
  .form-signin {
    width: 100%;
    max-width: 330px;
    padding: 0 35% 10% 35%;
    margin: auto;
    display: table-cell;
    vertical-align: middle;
  }

  .form-signin .checkbox {
    font-weight: 400;
  }

  .form-signin .form-control {
    position: relative;
    box-sizing: border-box;
    height: auto;
    padding: 10px;
    font-size: 16px;
  }

  .form-signin .form-control:focus {
    z-index: 2;
  }

  .form-signin input[type="email"] {
    margin-bottom: -1px;
    border-bottom-right-radius: 0;
    border-bottom-left-radius: 0;
  }

  .form-signin input[type="password"] {
    margin-bottom: 10px;
    border-top-left-radius: 0;
    border-top-right-radius: 0;
  }

  .fullscreen {
    height: 100%;
    display: table;
  }
</style>
