
var InputMask = function ( opts ) {

    if ( opts && opts.masked ) {
        // Make it easy to wrap this plugin and pass elements instead of a selector
        opts.masked = typeof opts.masked === 'string' ? document.querySelectorAll( opts.masked ) : opts.masked;
    }

    if ( opts ) {
        this.opts = {
            masked: opts.masked || document.querySelectorAll( this.d.masked ),
            mNum: opts.mNum || this.d.mNum,
            mChar: opts.mChar || this.d.mChar,
            error: opts.onError || this.d.onError
        }
    } else {
        this.opts = this.d;
        this.opts.masked = document.querySelectorAll( this.opts.masked );
    }

    this.refresh( true );
};
var inputMask = {

    // Default Values
    d: {
        masked : '.masked',
        mNum : 'XdDmMyY9',
        mChar : '_',
        onError: function(){}
    },

    refresh: function(init) {
        var t, parentClass;

        if ( !init ) {
            this.opts.masked = document.querySelectorAll( this.opts.masked );
        }

        for(i = 0; i < this.opts.masked.length; i++) {
            t = this.opts.masked[i]
            parentClass = t.parentNode.getAttribute('class');

            if ( !parentClass || ( parentClass && parentClass.indexOf('shell') === -1 ) ) {
                this.createShell(t);
                this.activateMasking(t);
            }
        }
    },

    // replaces each masked t with a shall containing the t and it's mask.
    createShell : function (t) {
        var wrap = document.createElement('span'),
            mask = document.createElement('span'),
            emphasis = document.createElement('i'),
            tClass = t.getAttribute('class'),
            pTxt = t.getAttribute('placeholder'),
            placeholder = document.createTextNode(pTxt);

        t.setAttribute('maxlength', placeholder.length);
        t.setAttribute('data-placeholder', pTxt);
        t.removeAttribute('placeholder');


        if ( !tClass || ( tClass && tClass.indexOf('masked') === -1 ) ) {
            t.setAttribute( 'class', tClass + ' masked');
        }

        mask.setAttribute('aria-hidden', 'true');
        mask.setAttribute('id', t.getAttribute('id') + 'Mask');
        mask.appendChild(emphasis);
        mask.appendChild(placeholder);

        wrap.setAttribute('class', 'shell');
        wrap.appendChild(mask);
        t.parentNode.insertBefore( wrap, t );
        wrap.appendChild(t);
    },

    setValueOfMask : function (e) {
        var value = e.target.value,
            placeholder = e.target.getAttribute('data-placeholder');

        return "<i>" + value + "</i>" + placeholder.substr(value.length);
    },

    // add event listeners
    activateMasking : function (t) {
        var that = this;
        if (t.addEventListener) { // remove "if" after death of IE 8
            t.addEventListener('keyup', function(e) {
                that.handleValueChange.call(that,e);
            }, false);
        } else if (t.attachEvent) { // For IE 8
            t.attachEvent('onkeyup', function(e) {
                e.target = e.srcElement;
                that.handleValueChange.call(that, e);
            });
        }
    },

    handleValueChange : function (e) {
        var id = e.target.getAttribute('id');

        if(e.target.value == document.querySelector('#' + id + 'Mask i').innerHTML) {
            return; // Continue only if value hasn't changed
        }

        document.getElementById(id).value = this.handleCurrentValue(e);
        document.getElementById(id + 'Mask').innerHTML = this.setValueOfMask(e);

    },

    handleCurrentValue : function (e) {
        var isCharsetPresent = e.target.getAttribute('data-charset'),
            placeholder = isCharsetPresent || e.target.getAttribute('data-placeholder'),
            value = e.target.value, l = placeholder.length, newValue = '',
            i, j, isInt, isLetter, strippedValue;

        // strip special characters
        strippedValue = isCharsetPresent ? value.replace(/\W/g, "") : value.replace(/\D/g, "");

        for (i = 0, j = 0; i < l; i++) {
            isInt = !isNaN(parseInt(strippedValue[j]));
            isLetter = strippedValue[j] ? strippedValue[j].match(/[A-Z]/i) : false;
            matchesNumber = this.opts.mNum.indexOf(placeholder[i]) >= 0;
            matchesLetter = this.opts.mChar.indexOf(placeholder[i]) >= 0;
            if ((matchesNumber && isInt) || (isCharsetPresent && matchesLetter && isLetter)) {
                newValue += strippedValue[j++];
            } else if ((!isCharsetPresent && !isInt && matchesNumber) || (isCharsetPresent && ((matchesLetter && !isLetter) || (matchesNumber && !isInt)))) {
                //this.opts.onError( e ); // write your own error handling function
                return newValue;
            } else {
                newValue += placeholder[i];
            }
            // break if no characters left and the pattern is non-special character
            if (strippedValue[j] == undefined) {
                break;
            }
        }
        if (e.target.getAttribute('data-valid-example')) {
            return this.validateProgress(e, newValue);
        }
        return newValue;
    },

    validateProgress : function (e, value) {
        var validExample = e.target.getAttribute('data-valid-example'),
            pattern = new RegExp(e.target.getAttribute('pattern')),
            placeholder = e.target.getAttribute('data-placeholder'),
            l = value.length, testValue = '';

        //convert to months
        if (l == 1 && placeholder.toUpperCase().substr(0,2) == 'MM') {
            if(value > 1 && value < 10) {
                value = '0' + value;
            }
            return value;
        }
        // test the value, removing the last character, until what you have is a submatch
        for ( i = l; i >= 0; i--) {
            testValue = value + validExample.substr(value.length);
            if (pattern.test(testValue)) {
                return value;
            } else {
                value = value.substr(0, value.length-1);
            }
        }

        return value;
    }
};

for ( var property in inputMask ) {
    if (inputMask.hasOwnProperty(property)) {
        InputMask.prototype[ property ] = inputMask[ property ];
    }
}

//  Declaritive initalization
(function(){
    console.log("!!!!!!!!!!!!!!!!!InputMask")
    // var scripts = document.getElementsByTagName('script'),
    //     script = scripts[ scripts.length - 1 ];
    // if ( script.getAttribute('data-autoinit') ) {
    new InputMask();
    // }
})();
// let  plugin_command = {
//     // @Required @Unique
//     // plugin name
//     name: 'customCommand',
//     // @Required
//     // data display
//     display: 'command',
//
//     // @Options
//     title: 'Add range tag',
//     buttonClass: '',
//     innerHTML: '<i class="fas fa-carrot"></i>',
//
//     // @Required
//     // add function - It is called only once when the plugin is first run.
//     // This function generates HTML to append and register the event.
//     // arguments - (core : core object, targetElement : clicked button element)
//     add: function (core, targetElement) {
//         console.log(core, targetElement);
//     },
//
//     // @Overriding
//     // Plugins with active methods load immediately when the editor loads.
//     // Called each time the selection is moved.
//     active: function (element) {
//         console.log(element);
//         return false;
//     },
//
//     // @Required
//     // The behavior of the "command plugin" must be defined in the "action" method.
//     action: function () {
//         console.log("action");
//     }
// };

// function initRichEditor(block){
//
//     let Editor = SUNEDITOR.create(block, {
//         stickyToolbar:false,
//         height : 'auto',
//         width : '100%',
//         minHeight:'200px',
//         maxHeight:'60vh',
//         // plugins: [],
//         buttonList: [
//             ['undo', 'redo'],
//             ['font', 'fontSize', 'formatBlock'],
//             ['paragraphStyle'],
//             ['bold', 'underline', 'italic', 'strike', 'subscript', 'superscript'],
//             ['fontColor', 'hiliteColor', 'textStyle'],
//             ['removeFormat'],
//             ['outdent', 'indent'],
//             ['align', 'horizontalRule', 'list', 'lineHeight'],
//             ['table', 'link', 'image', 'video'],
//             ['fullScreen', 'showBlocks', 'codeView'],
//             ['preview', 'print'],
//             ['customCommand', 'template']
//         ],
//         "lang": SUNEDITOR_LANG.ru,
//
//         // callBackSave: СохранитьСтатью(this),
//     });
//
//     Editor.onChange = function (contents) {
//         СохранитьСтатью (Editor)
//     }
// }

















/*частицы */
function ParticleInit()
{
    particlesJS("particles-js", {
        "particles": {
            "number": {
                "value": 30,
                "density": {
                    "enable": true,
                    "value_area": 300
                }
            },
            "color": {
                "value": "#ffffff"
            },
            "shape": {
                "type": "circle",
                "stroke": {
                    "width": 0,
                    "color": "#000000"
                },
                "polygon": {
                    "nb_sides": 5
                },
                "image": {
                    "src": "img/github.svg",
                    "width": 100,
                    "height": 100
                }
            },
            "opacity": {
                "value": 0.5,
                "random": false,
                "anim": {
                    "enable": false,
                    "speed": 1,
                    "opacity_min": 0.1,
                    "sync": false
                }
            },
            "size": {
                "value": 3,
                "random": true,
                "anim": {
                    "enable": false,
                    "speed": 40,
                    "size_min": 0.1,
                    "sync": false
                }
            },
            "line_linked": {
                "enable": true,
                "distance": 150,
                "color": "#ffffff",
                "opacity": 0.4,
                "width": 1
            },
            "move": {
                "enable": true,
                "speed": 6,
                "direction": "none",
                "random": false,
                "straight": false,
                "out_mode": "out",
                "bounce": false,
                "attract": {
                    "enable": false,
                    "rotateX": 600,
                    "rotateY": 1200
                }
            }
        },
        "interactivity": {
            "detect_on": "canvas",
            "events": {
                "onhover": {
                    "enable": true,
                    "mode": "grab"
                },
                "onclick": {
                    "enable": true,
                    "mode": "push"
                },
                "resize": true
            },
            "modes": {
                "grab": {
                    "distance": 140,
                    "line_linked": {
                        "opacity": 1
                    }
                },
                "bubble": {
                    "distance": 400,
                    "size": 40,
                    "duration": 2,
                    "opacity": 8,
                    "speed": 3
                },
                "repulse": {
                    "distance": 200,
                    "duration": 0.4
                },
                "push": {
                    "particles_nb": 4
                },
                "remove": {
                    "particles_nb": 2
                }
            }
        },
        "retina_detect": true
    });
}
// postgres.covid.данные_осп.всего_каранти_в_мед
// postgres.covid.данные_осп.всего_каранти_дома
// postgres.covid.данные_осп.всего_прошли_иccледование
// postgres.covid.данные_осп.спи_каранти_в_мед
// postgres.covid.данные_осп.спи_каранти_дома
// postgres.covid.данные_осп.спи_прошли_иccледование
// postgres.covid.данные_осп.дозн_каранти_в_мед
// postgres.covid.данные_осп.дозн_каранти_дома
// postgres.covid.данные_осп.дозн_прошли_иccледование
// postgres.covid.данные_осп.оупдс_каранти_в_мед
// postgres.covid.данные_осп.оупдс_каранти_дома
// postgres.covid.данные_осп.оупдс_прошли_иccледование


