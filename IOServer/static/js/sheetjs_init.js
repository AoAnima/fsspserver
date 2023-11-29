function uploadXlsx (e) {
    var files = e.target.files, f = files[0];
    var reader = new FileReader();
    reader.onload = function(e) {
        var data = new Uint8Array(e.target.result);
        var workbook = XLSX.read(data, {type: 'array'});
        console.log(workbook);
        /* DO SOMETHING WITH workbook HERE */
    };
    reader.readAsArrayBuffer(f);
     console.log(reader);

}
