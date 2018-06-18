// import Base64 from 'crypto-js/enc-base64';
import sha256 from 'crypto-js/sha256';
import Clipboard from 'clipboard';  //剪切板

const host = "http://localhost:8080"
// const host = "http://biluohc.me:8080"

// 日狗。。不就是想匹配\..
const entryName = RegExp("[\\\\/]")

// 放在 api 方便统一管理和 grep api.
const API = {
    github: "https://github.com/biluohc/kloud",
    signin: host + "/apiv1/signin",
    // post/get/delete
    login: host + "/apiv1/login",
    // 后面接 / 的表示要接 id 作为完整的 apiurl
    // post/put 除了文件上传之外表单都是 JSON
    // get /:id

    avatarURL: host + "/apiv1/user/avatar",
    password: host + "/apiv1/user/password",

    // 从 / 到当前id的 info 列表
    // 空的指 / 的
    entryPath: host + "/apiv1/entry/path/",
    // subs, 空的指Root的
    entrySubs: host + "/apiv1/entry/subs/",

    // 空的指 root的 info    
    entryInfo: host + "/apiv1/entry/info/",
    entryFile: host + "/apiv1/entry/file/",

    entryTrashes: host + "/apiv1/entry/trashes",
    // post /:id
    entryDirPost: host + "/apiv1/entry/dir/",
    entryFilePost: host + "/apiv1/entry/file/",
    entrySha1Post: host + "/apiv1/entry/sha1/",
    // put /:id
    entryRename: host + "/apiv1/entry/rename/",
    entryMoveto: host + "/apiv1/entry/moveto/",
    entryTrash: host + "/apiv1/entry/trash/",
    entryUnTrash: host + "/apiv1/entry/untrash/",
    // delete /:id
    entryDelete: host + "/apiv1/entry/",

    entryShare: host + "/apiv1/entry/share/",
    shares: host + "/apiv1/shares",
    shareFile: host + "/apiv1/share/",
    // shareFile: host + "/apiv1/share/:id/:name",
    // 有密码时的验证地址
    shareInit: host + "/apiv1/share/:id/:name/init",
    sharePassword: host + "/apiv1/share/password/",
    shareInfo: host + "/apiv1/share/info/",
    shareInit: host + "/apiv1/share/init/",
    shareMaxAge: host + "/apiv1/share/maxage/",
    shareDelete: host + "/apiv1/share/",

    entryPublic: host + "/apiv1/entry/public/",
    publics: host + "/apiv1/publics",
    publicFile: host + "/apiv1/public/",
    // publicFile: host + "/apiv1/public/:id/:name",
    publicReferer: host + "/apiv1/public/referer/",
    publicMaxAge: host + "/apiv1/public/maxage/",
    publicDelete: host + "/apiv1/public/",

    stat: host + "/Apiv1/stat",
    logFile: host + "/Apiv1/log.txt",

    users: host + "/Apiv1/users",
    user: host + "/Apiv1/user/",
    newUser: host + "/Apiv1/user",
    userDisable: host + "/Apiv1/user/disable/",
    userUnDisable: host + "/Apiv1/user/undisable/",

    files: host + "/Apiv1/files",
    file: host + "/Apiv1/file/",
    fileDisable: host + "/Apiv1/file/disable/",
    fileUnDisable: host + "/Apiv1/file/undisable/",

    //  2次 hash，以后两次都加盐 @kloud
    sha2Password: function (password) {
        var passwordOne = '' + sha256(password+'kloud');
        var passwordTwo = '' + sha256(passwordOne+'kloud');
        console.log(password, passwordOne, passwordTwo);
        return passwordTwo
    },
    // 注意不可逆，别提交给服务器
    dateTime: function (str) {
        var tmp = new Date(str)
        // console.log(tmp, Number(tmp) === Number(new Date('0001-01-01T00:00:00Z')) )
        // Go的莫名其妙的空时间
        if (Number(tmp) === Number(new Date('0001-01-01T00:00:00Z'))) {
            return "永久"
        }
        return tmp.toLocaleString()
    },
    size: function (size) {
        // B，KB，MB，GB，TB，PB，EB，ZB，YB，BB
        var units = ["", "K", "M", "G", "T", "P", "E", "Z", "Y", "B"];
        var s = size;
        var count = 0;
        while (s / 1024 > 1) {
            s = s / 1024;
            count += 1;
        }
        return s.toFixed(2) + units[count];
    },
    role: function (role) {
        var tmp = "游客"
        if (role === 0) {
            tmp = "管理员";
        } else if (role === 1) {
            tmp = "普通用户";
        }
        return tmp;
    },
    isValidEntryName: function (str) {
        return (!entryName.test(str)) && str !== "." && str !== ".."
    },
    copyLink: function () {
        var clipboard = new Clipboard(".btn");
        clipboard.on("success", e => {
            console.log("复制成功");
            // 释放内存
            clipboard.destroy();
        });
        clipboard.on("error", e => {
            console.log("该浏览器不支持自动复制");
            clipboard.destroy();
        });
    }
}

export default API