const {Workbook} = require('exceljs');

function makeExcel(list)
{
  const workbook = new Workbook();
  const sheet = workbook.addWorksheet('download');
  sheet.columns = [
    {header:'이름', key:'name', width:6},
    {header:'번호', key:'phone', width:12},
    {header:'이메일주소', key:'email', width:40},
  ];
  sheet.addRows(list);

  //첫째줄 스타일 가공
  sheet.eachRow((row, number)=>{
    row.alignment = {horizontal:'center'};

    if(number==1){
      row.eachCell(cell=>{
        cell.fill = {
          type: 'pattern',
          pattern: 'solid',
          fgColor: { argb:'D9D9D9' }
        };
        cell.font = { bold: true };
      })
    }
  });
  
  return workbook;
}

/*
    res.setHeader('Content-Type', 'application/vnd.openxmlformats');
    res.setHeader('Content-Disposition', 'attachment; filename=download.xlsx');
    await workbook.xlsx.write(res);
    res.end();

*/
