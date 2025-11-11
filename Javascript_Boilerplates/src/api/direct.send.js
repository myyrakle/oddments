const URL = `https://directsend.co.kr/index.php/api_v2/sms_change_word`;
const axios = require("axios");

async function send(phone, text) {
    const obj = {
        title: "title",
        message: null,
        sender: "01044455555",
        username: "아이디",
        receiver: [{ name: "unknown", mobile: null }],
        key: "API 키값...",
    };

    obj.receiver[0].mobile = phone.replace("-", "");
    obj.message = text;

    try {
        const res = await axios({
            url: URL,
            method: "post",
            headers: {
                "cache-control": "no-cache",
                "content-type": "application/json",
                charset: "utf-8",
            },
            data: obj,
        });
        console.log(res.data);
    } catch (error) {
        console.error(error);
    }
}

module.exports = send;
