// const { Workbook } = require("exceljs");
const { Workbook } = require("exceljs-lightweight");

const S3_BUCKET_NAME = "버킷명";

// 엑셀 파일을 생성해 전송합니다.
//
// response = express 리스폰스 객체
// list = DB에서 읽어온 값 그대로
// columns = 생성 형태. [{header:'맨 위에 열 이름', key:'DB select값의 이름', width:'가로폭(평균 10 정도 적당)'}]
async function sendCMSExcel(req, list, columns) {
    const workbook = new Workbook();
    const sheet = workbook.addWorksheet("cms.download");
    sheet.columns = columns;
    sheet.addRows(list);

    //첫째줄 스타일 가공
    sheet.eachRow((row, number) => {
        row.alignment = { horizontal: "center" };

        if (number == 1) {
            row.eachCell((cell) => {
                cell.fill = {
                    type: "pattern",
                    pattern: "solid",
                    fgColor: { argb: "D9D9D9" }, //회색
                };
                cell.font = { bold: true }; //볼트체

                const border = { style: "thin" };
                cell.border = {
                    //테두리
                    top: border,
                    left: border,
                    right: border,
                    bottom: border,
                };
            });
        }
    });

    const now = new Date();
    const _ = require("lodash");

    const dirname = `${now.getFullYear()}-${("00" + (now.getMonth() + 1)).slice(
        -2
    )}-${("00" + now.getDate()).slice(-2)}`;
    const filename = `${("00" + (now.getHours() + 1)).slice(-2)}-${(
        "00" +
        (now.getMinutes() + 1)
    ).slice(-2)}-${("00" + now.getSeconds()).slice(-2)}-${(
        "00" + now.getMilliseconds()
    ).slice(-2)}-${_.random(1000)}`;

    const full_path = `private/${dirname}/${filename}.cms.excel.xlsx`;

    await req.s3.upload({
        key: full_path,
        bucket: S3_BUCKET_NAME,
        data: await workbook.xlsx.writeBuffer(),
        ContentEncoding: "base64",
        ContentType:
            "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        option: {
            Expires: 300,
            ACL: "private",
        },
    });

    const url = await req.s3.getSignedUrl({
        key: full_path,
        bucket: S3_BUCKET_NAME,
        expires: 60,
    });

    return url;
}
