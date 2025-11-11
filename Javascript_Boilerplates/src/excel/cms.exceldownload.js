const {Workbook} = require('exceljs');
// const { Workbook } = require("exceljs-lightweight");

// 엑셀 파일을 생성해 전송합니다.
//
// response = express 리스폰스 객체
// list = DB에서 읽어온 값 그대로
// columns = 생성 형태. [{header:'맨 위에 열 이름', key:'DB select값의 이름', width:'가로폭(평균 10 정도 적당)'}]
async function sendCMSExcel(response, list, columns)
{
  const workbook = new Workbook();
  const sheet = workbook.addWorksheet('cms.download');
  sheet.columns = columns;
  sheet.addRows(list);

  //첫째줄 스타일 가공
  sheet.eachRow((row, number)=>{
    row.alignment = {horizontal:'center'};

    if(number==1){
      row.eachCell(cell=>{
        cell.fill = {
          type: 'pattern',
          pattern: 'solid',
          fgColor: { argb:'D9D9D9' } //회색
        };
        cell.font = { bold: true }; //볼트체
        cell.border = { //테두리
          top: { style: "thin" },
          left: { style: "thin" },
          right: { style: "thin" },
          bottom: { style: "thin" },
        };
      })
    }
  });
  
  response.setHeader('Content-Type', 'application/vnd.openxmlformats');
  response.setHeader('Content-Disposition', 'attachment; filename=cms.download.xlsx');
  await workbook.xlsx.write(response);
  response.end();
}

/* Example
  
  app.get('/download', async (req, res)=>{
    //...
    
    const selectResult = await req.sequalize.query(`select no, name, age from tb_person`);
    await sendCMSExcel(res, selectResult, [{header:'번호', key:'no', width:5}, {header:'이름', key:'name', width:10}, {header:'나이', key:'age', width:10}]);
    // done
  });
*/
