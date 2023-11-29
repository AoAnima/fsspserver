

var ОткрытыеРедакторыMonaco={};
require.config({
    'vs/nls' : {
        availableLanguages: {
            '*': 'ru'
        }
    }
});


function ПодогнатьРазмерMonaco(MonacoId){

    document.getElementById(MonacoId).layout()



    // ОткрытыеРедакторыMonaco[`monaco_${ИсточникКода}`]["IStandaloneCodeEditor"].layout()
}
var SynthWave = {
    base: 'vs',
    inherit: true,
    name: "SynthWave 84",
    type: "dark",
    colors: {
        "focusBorder": "#1f212b",
        'foreground' : "#ffffff",
        'editor.foreground' : "#ffffff",
        "widget.shadow": "#2a2139",
        "selection.background": "#ffffff36",
        "errorForeground": "#fe4450",
        "textLink.activeForeground": "#ff7edb",
        "textLink.foreground": "#f97e72",
        "button.background": "#2a2139",
        "dropdown.background": "#232530",
        "dropdown.listBackground": "#2a2139",
        "input.background": "#2a2139",
        "inputOption.activeBorder": "#ff7edb99",
        "inputValidation.errorBackground": "#fe445080",
        "inputValidation.errorBorder": "#fe445000",
        "scrollbar.shadow": "#2a2139",
        "scrollbarSlider.activeBackground": "#9d8bca20",
        "scrollbarSlider.background": "#9d8bca30",
        "scrollbarSlider.hoverBackground": "#9d8bca50",
        "badge.foreground": "#ffffff",
        "badge.background": "#2a2139",
        "progressBar.background": "#f97e72",
        "list.activeSelectionBackground": "#ffffff36",
        "list.activeSelectionForeground": "#ffffff",
        "list.dropBackground": "#34294f66",
        "list.focusBackground": "#2a213999",
        "list.focusForeground": "#ffffff",
        "list.highlightForeground": "#f97e72",
        "list.hoverBackground": "#2a213999",
        "list.hoverForeground": "#ffffff",
        "list.inactiveSelectionBackground": "#34294f66",
        "list.inactiveSelectionForeground": "#ffffff",
        "list.inactiveFocusBackground": "#2a213999",
        "list.errorForeground": "#fe4450E6",
        "list.warningForeground": "#72f1b8bb",
        "activityBar.background": "#171520",
        "activityBar.dropBackground": "#34294f66",
        "activityBar.foreground": "#ffffffCC",
        "activityBarBadge.background": "#f97e72",
        "activityBarBadge.foreground": "#2a2139",
        "sideBar.background": "#241b2f",
        "sideBar.foreground": "#ffffff99",
        "sideBar.dropBackground": "#34294f4c",
        "sideBarSectionHeader.background": "#241b2f",
        "sideBarSectionHeader.foreground": "#ffffffca",
        "menu.background": "#463465",
        "editorGroup.border": "#495495",
        "editorGroup.dropBackground": "#34294f4a",
        "editorGroupHeader.tabsBackground": "#241b2f",
        "tab.border": "#241b2f00",
        "tab.activeBorder": "#880088",
        "tab.inactiveBackground": "#262335",
        "editor.background": "#1d0331",
        "editorLineNumber.foreground": "#ffffff73",
        "editorLineNumber.activeForeground": "#ffffffcc",
        "editorCursor.background": "#241b2f",
        "editorCursor.foreground": "#f97e72",
        "editor.selectionBackground": "#3A5070",
        "editor.selectionHighlightBackground": "#ffffff36",
        "editor.wordHighlightBackground": "#34294f88",
        "editor.wordHighlightStrongBackground": "#34294f88",
        "editor.findMatchBackground": "#D18616bb",
        "editor.findMatchHighlightBackground": "#D1861655",
        "editor.findRangeHighlightBackground": "#34294f1a",
        "editor.hoverHighlightBackground": "#463465",
        "editor.lineHighlightBackground": "#34294f66",
        "editor.rangeHighlightBackground": "#49549539",
        "editorIndentGuide.background": "#49549539",
        "editorIndentGuide.activeBackground": "#2a2139",
        "editorRuler.foreground": "#34294f33",
        "editorCodeLens.foreground": "#262335",
        "editorCodeLens.background": "#262335",
        "editorBracketMatch.background": "#34294f66",
        "editorBracketMatch.border": "#495495",
        "editorOverviewRuler.border": "#34294fb3",
        "editorOverviewRuler.findMatchForeground": "#D1861699",
        "editorOverviewRuler.modifiedForeground": "#b893ce99",
        "editorOverviewRuler.addedForeground": "#09f7a099",
        "editorOverviewRuler.deletedForeground": "#fe445099",
        "editorOverviewRuler.errorForeground": "#fe4450dd",
        "editorOverviewRuler.warningForeground": "#72f1b8cc",
        "editorError.foreground": "#fe4450",
        "editorWarning.foreground": "#72f1b8cc",
        "editorGutter.modifiedBackground": "#b893ce8f",
        "editorGutter.addedBackground": "#206d4bd6",
        "editorGutter.deletedBackground": "#fa2e46a4",
        "diffEditor.insertedTextBackground": "#0beb9916",
        "diffEditor.removedTextBackground": "#fe445016",
        "editorWidget.background": "#372d4b",
        "editorWidget.border": "#ffffff22",
        "editorWidget.resizeBorder": "#ffffff44",
        "editorSuggestWidget.highlightForeground": "#f97e72",
        "editorSuggestWidget.selectedBackground": "#ffffff36",
        "peekView.border": "#495495",
        "peekViewEditor.background": "#232530",
        "peekViewEditor.matchHighlightBackground": "#D18616bb",
        "peekViewResult.background": "#232530",
        "peekViewResult.matchHighlightBackground": "#D1861655",
        "peekViewResult.selectionBackground": "#2a213980",
        "peekViewTitle.background": "#232530",
        "panelTitle.activeBorder": "#f97e72",
        "statusBar.background": "#241b2f",
        "statusBar.foreground": "#ffffff80",
        "statusBar.debuggingBackground": "#f97e72",
        "statusBar.debuggingForeground": "#08080f",
        "statusBar.noFolderBackground": "#241b2f",
        "statusBarItem.prominentBackground": "#2a2139",
        "statusBarItem.prominentHoverBackground": "#34294f",
        "titleBar.activeBackground": "#241b2f",
        "titleBar.inactiveBackground": "#241b2f",
        "extensionButton.prominentBackground": "#f97e72",
        "extensionButton.prominentHoverBackground": "#ff7edb",
        "pickerGroup.foreground": "#f97e72ea",
        "terminal.foreground": "#ffffff",
        "terminal.ansiBlue": "#03edf9",
        "terminal.ansiBrightBlue": "#03edf9",
        "terminal.ansiBrightCyan": "#03edf9",
        "terminal.ansiBrightGreen": "#72f1b8",
        "terminal.ansiBrightMagenta": "#ff7edb",
        "terminal.ansiBrightRed": "#fe4450",
        "terminal.ansiBrightYellow": "#fede5d",
        "terminal.ansiCyan": "#03edf9",
        "terminal.ansiGreen": "#72f1b8",
        "terminal.ansiMagenta": "#ff7edb",
        "terminal.ansiRed": "#fe4450",
        "terminal.ansiYellow": "#f97e72",
        "terminal.selectionBackground": "#34294f4d",
        "terminalCursor.background": "#ffffff",
        "terminalCursor.foreground": "#03edf9",
        "debugToolBar.background": "#241b2f",
        "walkThrough.embeddedEditorBackground": "#232530",
        "gitDecoration.modifiedResourceForeground": "#b893ceee",
        "gitDecoration.deletedResourceForeground": "#fe4450",
        "gitDecoration.addedResourceForeground": "#72f1b8cc",
        "gitDecoration.untrackedResourceForeground": "#72f1b8",
        "gitDecoration.ignoredResourceForeground": "#ffffff59",
        // "minimapGutter.addedBackground": "#09f7a099",
        // "minimapGutter.modifiedBackground": "#b893ce",
        // "minimapGutter.deletedBackground": "#fe4450",
        "breadcrumbPicker.background": "#232530"
    },
    rules: [
        { token: '', foreground: '#d7d7d7'},
        { token: '', background: '#1d0331'},
        { token: 'delimiter', foreground: '#f88126' },
        { token: 'delimiter.html', foreground: '#fdeac5' },
        { token: 'tag', foreground: '#ffb443' },
        { token: 'attribute.name', foreground: '#50FA7B' },
        { token: 'metatag.content.html', foreground: '#13ff00' },
        {token: 'attribute.value.html', foreground: '#ffde00' },
        { token: 'metatag.html', foreground: '#13ff00' },
        { token: 'tag.html', foreground: '#36f9f6' },
        { token: 'string.html', foreground: '#ffffff' },

    ]
}

/*
InitMonaco
* КлассыЭлементовИсходников  по умолчанию ["source_code"] можно передавать массив классов
* MonacoWrapper - id контейнера куда поместить редактор  по умолчанию monaco_wrapper
* */

function InitMonaco(КлассыЭлементовИсходников=["source_code"], MonacoWrapper="monaco_wrapper") {

     let КонтейнерИсходников = document.getElementById(КлассыЭлементовИсходников);
     let ОсновнойКонтейнерРедактор = document.getElementById(MonacoWrapper);

     let ПанельВкладок = document.createElement('div')
         ПанельВкладок.id="monaco_tabs"



     let КонтейнерРедакторов = document.createElement('div')
         КонтейнерРедакторов.id="monacos"

     for (let Исходник of КонтейнерИсходников.children) {
          console.log(Исходник);
         let Вкладка = document.createElement('button')
             Вкладка.type=`button`
             Вкладка.id=`tab_${Исходник.id}`
             Вкладка.classList.add(`monaco_tab`)
             Вкладка.innerText=Исходник.id
             // Вкладка.setAttribute('onclick',ПоказатьРедактор(`monaco_${Исходник.id}`, "monacos"))


             // Вкладки.append(Вкладка)

         let NewMonaco = document.createElement('div')
             NewMonaco.id=`monaco_${Исходник.id}`
             NewMonaco.classList.add("monaco")
             NewMonaco.style.width="100%"
             NewMonaco.style.height="100vh"
             КонтейнерРедакторов.append(NewMonaco)

         let editor
            require(['vs/editor/editor.main'], function () {
             monaco.editor.defineTheme('SynthWave', SynthWave)
              editor = monaco.editor.create(NewMonaco, {
                 value: Исходник.value,
                 language: язык,
                 // theme: "SynthWave",
                 renderWhitespace: "all",
                 wordWrap: 'bounded',
                 wordWrapColumn: 100,
                 wordWrapMinified: true,
                 fontSize: 15,
                 fontFamily:'Fira Code',
                 // try "same", "indent" or "none"
                 wrappingIndent: "indent",
                 wrappingStrategy: "advanced"
             });

             // console.log(editor.getLayoutInfo());
             // editor.layout()
             let KeyMap = editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KEY_S, function () {
                 console.log("Получаем значение редактора для сохранения", editor.getValue());
             });
             // ОткрытыеРедакторыMonaco[NewMonaco.id]={"monaco_window":NewMonaco,
             //             "IStandaloneCodeEditor":editor
             // }

             });
         /*Подгоним размер редактора и активируем вкладку*/
         Вкладка.addEventListener('click', () => {
             ПоказатьРедактор(`monaco_${Исходник.id}`, "monacos");
             editor.layout()
             // for (let Кнопка of Вкладки.children){
             //     Кнопка.classList.remove(`active`);
             // }
             Вкладка.classList.add(`active`)
         })
         /*активируем вкладку и кнопку*/
         if(!!Исходник.dataset.main_tpl && Исходник.dataset.main_tpl!== "" && Исходник.dataset.main_tpl === Исходник.id){
             Вкладка.classList.add(`active`)
             NewMonaco.classList.add("active")
         }
     }
    // ОсновнойКонтейнерРедактор.append(Вкладки)
    ОсновнойКонтейнерРедактор.append(КонтейнерРедакторов)
}

function ПоказатьРедактор(ИдРедактора, КонтейнерРедакторов){
    let котейнерРедактора = document.getElementById(КонтейнерРедакторов)
    for (let редактор of котейнерРедактора.children){
        редактор.classList.remove('active')
    }
    document.getElementById(ИдРедактора).classList.add('active')
}



