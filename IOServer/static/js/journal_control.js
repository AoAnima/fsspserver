// function ПросмотрФайла(ИдГаллереи, ИдФайла){
//     console.log(ИдГаллереи)
//     ОбёрткаГаллереи = document.getElementById(ИдГаллереи)
//     МетаОбъектыФайлов = ОбёрткаГаллереи.querySelectorAll('.file')
//     console.log(МетаОбъектыФайлов)
//     ОбёрткаГаллереи.classList.add('modal_carusel')
//     for (let номер in МетаОбъектыФайлов){
//
//         let МетаОбъектФайла = МетаОбъектыФайлов[номер]
//         console.log(МетаОбъектФайла)
//
//         ИдDOMОбъектаФайла = МетаОбъектФайла.dataset.file_object
//
//         DOMОбъектФайла = МетаОбъектФайла.querySelector(`#${ИдDOMОбъектаФайла}`)
//         DOMОбъектФайла.src = МетаОбъектФайла.dataset.src
//         if (МетаОбъектФайла.id === `file_${ИдФайла}`){
//             МетаОбъектФайла.classList.add('active')
//         }
//     }
// }

function ПросмотрФайла(ИдФайла, src, НажатаяКнопка){


    let ПолныйПросмотр = document.getElementById('file_preview')
    let Кнопки = document.querySelectorAll(`.btn_preview`)

    let  Документ;
    if (src.includes(".pdf")){
         Документ = `<embed class="file_object" id="full_${ИдФайла}" src="${src}" type="application/pdf" width="100%" height="100%">`
        ПолныйПросмотр.innerHTML = Документ
    } else if (src.includes(".doc") || src.includes(".docx") || src.includes(".odt")) {

        let  Документ = document.createElement('a');

        Документ.setAttribute('href', src);
        Документ.setAttribute('download',"");
        Документ.setAttribute('target','_blank');
        Документ.style.display = 'none';
        console.log(Документ)
        ПолныйПросмотр.appendChild(Документ);
        Документ.click();
        ПолныйПросмотр.removeChild(Документ);

    } else {
         Документ =` <img class="file_object" id="full_${ИдФайла}" src="${src}" alt="">`
        ПолныйПросмотр.innerHTML = Документ
    }


    for (let Кнопка of Кнопки){
        Кнопка.classList.remove('active')
    }
    console.log(НажатаяКнопка)
    НажатаяКнопка.classList.add('active')
}
