/*
 * 多语言翻译，
 *
 * @param 管雷鸣
 * @address https://gitee.com/mail_osc/translate
 */
layui.define(['table', 'form', 'element'], function(exports) {

    var translate = {
        version: "2.1.8.20230110",
        useVersion: "v1",
        setUseVersion2: function() {
            translate.useVersion = "v2"
        },
        translate: null,
        includedLanguages: "zh-CN,zh-TW,en",
        resourcesUrl: "//res.zvo.cn/translate",
        selectLanguageTag: {
            show: !0,
            languages: "zh-CN,zh-TW,en",
            alreadyRender: !1,
            render: function() {
                if (!translate.selectLanguageTag.alreadyRender && (translate.selectLanguageTag
                    .alreadyRender = !0, translate.selectLanguageTag.show)) {
                    if (null == document.getElementById("translate")) {
                        var e = document.getElementsByTagName("body")[0],
                            t = document.createElement("div");
                        t.id = "translate", e.appendChild(t)
                    } else if (null != document.getElementById("translateSelectLanguage")) return;
                    translate.request.post("https://api.translate.zvo.cn/language.json", {}, function(
                        e) {
                        if (0 != e.result) {
                            var t = function(e) {
                                    var t = e.target.value;
                                    translate.changeLanguage(t)
                                },
                                n = document.createElement("select");
                            n.id = "translateSelectLanguage", n.className =
                                "translateSelectLanguage";
                            for (var a = 0; a < e.list.length; a++) {
                                var l = document.createElement("option");
                                l.setAttribute("value", e.list[a].id), null != translate.to &&
                                void 0 !== translate.to && translate.to.length > 0 ?
                                    translate.to == e.list[a].id && l.setAttribute("selected",
                                        "selected") : e.list[a].id == translate.language
                                    .local && l.setAttribute("selected", "selected"), l
                                    .appendChild(document.createTextNode(e.list[a].name)), n
                                    .appendChild(l)
                            }
                            window.addEventListener ? n.addEventListener("change", t, !1) : n
                                .attachEvent("onchange", t), document.getElementById(
                                "translate").appendChild(n)
                        } else console.log("load language list error : " + e.info)
                    })
                }
            }
        },
        localLanguage: "zh-CN",
        googleTranslateElementInit: function() {
            var e = "";
            null != document.getElementById("translate") && (e = "translate"), translate.translate =
                new google.translate.TranslateElement({
                    pageLanguage: "zh-CN",
                    includedLanguages: translate.selectLanguageTag.languages,
                    layout: 0
                }, e)
        },
        init: function() {
            var e = window.location.protocol;
            "file:" == window.location.protocol && (e = "http:"), -1 == this.resourcesUrl.indexOf(
                "://") && (this.resourcesUrl = e + this.resourcesUrl)
        },
        execute_v1: function() {
            if (null == document.getElementById("translate") && translate.selectLanguageTag.show) {
                var e = document.getElementsByTagName("body")[0],
                    t = document.createElement("div");
                t.id = "translate", e.appendChild(t)
            }
            "zh-CN,zh-TW,en" != translate.includedLanguages && (translate.selectLanguageTag.languages =
                translate.includedLanguages, console.log(
                "translate.js tip: translate.includedLanguages obsolete, please use the translate.selectLanguageTag.languages are set"
            ));
            var n = document.getElementsByTagName("head")[0],
                a = document.createElement("script");
            a.type = "text/javascript", a.src = this.resourcesUrl + "/js/element.js", n.appendChild(a)
        },
        setCookie: function(e, t) {
            var n = e + "=" + escape(t);
            document.cookie = n
        },
        getCookie: function(e) {
            for (var t = document.cookie.split("; "), n = 0; n < t.length; n++) {
                var a = t[n].split("=");
                if (a[0] == e) return unescape(a[1])
            }
            return ""
        },
        currentLanguage: function() {
            translate.check();
            var e = translate.getCookie("googtrans");
            return e.length > 0 ? e.substr(e.lastIndexOf("/") + 1, e.length - 1) : translate
                .localLanguage
        },
        changeLanguage: function(e) {
            if (",en,de,hi,lt,hr,lv,ht,hu,zh-CN,hy,uk,mg,id,ur,mk,ml,mn,af,mr,uz,ms,el,mt,is,it,my,es,et,eu,ar,pt-PT,ja,ne,az,fa,ro,nl,en-GB,no,be,fi,ru,bg,fr,bs,sd,se,si,sk,sl,ga,sn,so,gd,ca,sq,sr,kk,st,km,kn,sv,ko,sw,gl,zh-TW,pt-BR,co,ta,gu,ky,cs,pa,te,tg,th,la,cy,pl,da,tr,"
                .indexOf("," + e + ",") > -1) {
                translate.check();
                var t = "/" + translate.localLanguage + "/" + e,
                    n = document.location.host.split(".");
                if (n.length > 2) {
                    var a = n[n.length - 2] + "." + n[n.length - 1];
                    document.cookie = "googtrans=;expires=" + new Date(1) + ";domain=" + a + ";path=/",
                        document.cookie = "googtrans=" + t + ";domain=" + a + ";path=/"
                }
                return translate.setCookie("googtrans", "" + t), void location.reload()
            }
            translate.setUseVersion2(), translate.to = e, translate.storage.set("to", e), location
                .reload()
        },
        check: function() {
            "file:" == window.location.protocol && console.log(
                "\r\n---WARNING----\r\ntranslate.js 主动翻译组件自检异常，当前协议是file协议，翻译组件要在正常的线上http、https协议下才能正常使用翻译功能\r\n------------"
            )
        },
        to: "",
        autoDiscriminateLocalLanguage: !1,
        documents: [],
        ignore: {
            tag: ["style", "script", "img", "head", "link", "i", "pre", "code"],
            class: ["ignore", "translateSelectLanguage"]
        },
        setAutoDiscriminateLocalLanguage: function() {
            translate.autoDiscriminateLocalLanguage = !0
        },
        nodeQueue: {},
        setDocuments: function(e) {
            null != e && void 0 !== e && (void 0 === e.length ? translate.documents[0] = e : translate
                .documents = e, translate.nodeQueue = {}, console.log(
                "set documents , clear translate.nodeQueue"))
        },
        getDocuments: function() {
            return null != translate.documents && void 0 !== translate.documents && translate.documents
                .length > 0 ? translate.documents : document.all
        },
        listener: {
            isExecuteFinish: !1,
            isStart: !1,
            start: function() {
                translate.temp_linstenerStartInterval = setInterval(function() {
                    "complete" == document.readyState && (clearInterval(translate
                        .temp_linstenerStartInterval), translate.listener.addListener())
                }, 50)
            },
            addListener: function() {
                translate.listener.isStart = !0;
                const e = {
                        attributes: !0,
                        childList: !0,
                        subtree: !0
                    },
                    t = new MutationObserver(function(e, t) {
                        var n = [];
                        for (let t of e) "childList" === t.type && t.addedNodes.length > 0 && n.push
                            .apply(n, t.addedNodes);
                        n.length > 0 && translate.execute(n)
                    });
                for (var n = translate.getDocuments(), a = 0; a < n.length; a++) {
                    var l = n[a];
                    null != l && t.observe(l, e)
                }
            }
        },
        renderTask: class {
            constructor() {
                this.taskQueue = [], this.nodes = []
            }
            add(e, t, n) {
                var a = translate.util.hash(e.nodeValue);
                void 0 === this.nodes[a] && (this.nodes[a] = new Array), this.nodes[a].push(e);
                var l = this.taskQueue[a];
                null != l && void 0 !== l || (l = new Array);
                var r = new Array;
                r.originalText = t, r.resultText = n, l.push(r), this.taskQueue[a] = l
            }
            execute() {
                for (var e in this.taskQueue) {
                    (t = this.taskQueue[e]).sort(function(e, t) {
                        return t.originalText.length - e.originalText.length
                    }), this.taskQueue[e] = t
                }
                for (var e in this.nodes)
                    for (var t = this.taskQueue[e], n = 0; n < this.nodes[e].length; n++)
                        for (var a = 0; a < t.length; a++) {
                            var l = t[a];
                            this.nodes[e][a].nodeValue = this.nodes[e][a].nodeValue.replace(
                                new RegExp(l.originalText, "g"), l.resultText)
                        }
            }
        },
        execute: function(e) {
            if ("undefined" != typeof doc && (translate.useVersion = "v2"), "v1" != translate
                .useVersion) {
                var t = translate.util.uuid();
                if (translate.nodeQueue[t] = new Array, translate.nodeQueue[t].expireTime = Date.now() +
                    12e4, translate.nodeQueue[t].list = new Array, null == translate.to || "" ==
                translate.to) {
                    var n = translate.storage.get("to");
                    null != n && void 0 !== n && n.length > 0 && (translate.to = n)
                }
                try {
                    translate.selectLanguageTag.render()
                } catch (e) {
                    console.log(e)
                }
                if (null != translate.to && void 0 !== translate.to && 0 != translate.to.length) {
                    var a;
                    if (void 0 !== e) {
                        if (null == e) return void cnosole.log(
                            "translate.execute(...) 中传入的要翻译的目标区域不存在。");
                        void 0 === e.length ? (a = new Array)[0] = e : a = e
                    } else a = translate.getDocuments();
                    for (var l = 0; l < a.length & l < 20; l++) {
                        var r = a[l];
                        translate.whileNodes(t, r)
                    }
                    var s = {},
                        o = {};
                    for (var u in translate.nodeQueue[t].list) {
                        if (null == u || void 0 === u || 0 == u.length || "undefined" == u) continue;
                        s[u] = [], o[u] = [];
                        let e = new translate.renderTask;
                        for (var i in translate.nodeQueue[t].list[u]) {
                            var c = translate.nodeQueue[t].list[u][i].original,
                                d = translate.storage.get("hash_" + translate.to + "_" + i);
                            if (null != d && d.length > 0)
                                for (var g = 0; g < translate.nodeQueue[t].list[u][i].nodes.length; g++)
                                    e.add(translate.nodeQueue[t].list[u][i].nodes[g], c, d);
                            else s[u].push(c), o[u].push(i)
                        }
                        e.execute()
                    }
                    var h = [];
                    for (var u in translate.nodeQueue[t].list) s[u].length < 1 || h.push(u);
                    if (translate.listener.isExecuteFinish || (translate.temp_executeFinishNumber = 0,
                        translate.temp_executeFinishInterval = setInterval(function() {
                            translate.temp_executeFinishNumber == h.length && (translate
                                .listener.isExecuteFinish = !0, clearInterval(translate
                                .temp_executeFinishInterval))
                        }, 50)), 0 != h.length)
                        for (var f in h) {
                            u = h[f];
                            if (s[u].length < 1) return;
                            var v = {
                                from: u,
                                to: translate.to,
                                text: encodeURIComponent(JSON.stringify(s[u]))
                            };
                            translate.request.post("https://api.translate.zvo.cn/translate.json", v,
                                function(e) {
                                    if (0 == e.result) return console.log(
                                        "=======ERROR START======="), console.log(s[e
                                        .from]), console.log("response : " + e.info), console
                                        .log("=======ERROR END  ======="), void translate
                                        .temp_executeFinishNumber++;
                                    let n = new translate.renderTask;
                                    for (var a = 0; a < o[e.from].length; a++) {
                                        var l = e.text[a],
                                            r = o[e.from][a],
                                            u = e.from,
                                            i = "";
                                        try {
                                            i = translate.nodeQueue[t].list[u][r].original
                                        } catch (e) {
                                            console.log("uuid:" + t + ", originalWord:" + i +
                                                ", lang:" + u + ", hash:" + r + ", text:" + l +
                                                ", queue:" + translate.nodeQueue[t]), console
                                                .log(e);
                                            continue
                                        }
                                        for (var c = 0; c < translate.nodeQueue[t].list[u][r].nodes
                                            .length; c++) n.add(translate.nodeQueue[t].list[u][r]
                                            .nodes[c], i, l);
                                        translate.storage.set("hash_" + e.to + "_" + r, l)
                                    }
                                    n.execute(), translate.temp_executeFinishNumber++
                                })
                        }
                } else translate.autoDiscriminateLocalLanguage && translate.executeByLocalLanguage()
            } else translate.execute_v1()
        },
        whileNodes: function(e, t) {
            if (null != t && void 0 !== t) {
                var n = t.childNodes;
                if (n.length > 0)
                    for (var a = 0; a < n.length; a++) translate.whileNodes(e, n[a]);
                else translate.findNode(e, t)
            }
        },
        findNode: function(e, t) {
            if (null != t && void 0 !== t && null != t.parentNode) {
                var n = t.parentNode.nodeName;
                if (null != n && !(translate.ignore.tag.indexOf(n.toLowerCase()) > -1)) {
                    for (var a = !1, l = t.parentNode; t != l && null != l;) null != l.className &&
                    translate.ignore.class.indexOf(l.className) > -1 && (a = !0), l = l.parentNode;
                    if (!a)
                        if ("INPUT" == t.nodeName || "TEXTAREA" == t.nodeName) {
                            if (null == t.attributes || void 0 === t.attributes) return;
                            void 0 !== t.attributes.placeholder && translate.addNodeToQueue(e, t
                                .attributes.placeholder)
                        } else if (null != t.nodeValue && t.nodeValue.trim().length > 0) {
                            if (!(null != t.nodeValue && "string" == typeof t.nodeValue && t.nodeValue
                                .length > 0)) return;
                            translate.addNodeToQueue(e, t)
                        }
                }
            }
        },
        addNodeToQueue: function(e, t) {
            if (null != t.nodeValue && 0 != t.nodeValue.length) {
                translate.util.hash(t.nodeValue);
                if (!translate.util.findTag(t.nodeValue)) {
                    var n = translate.language.get(t.nodeValue);
                    for (var a in void 0 !== n[translate.to] && delete n[translate.to], n) {
                        null != translate.nodeQueue[e].list[a] && void 0 !== translate.nodeQueue[e]
                            .list[a] || (translate.nodeQueue[e].list[a] = new Array);
                        for (var l = 0; l < n[a].length; l++) {
                            var r = n[a][l],
                                s = translate.util.hash(r);
                            null != translate.nodeQueue[e].list[a][s] && void 0 !== translate.nodeQueue[
                                e].list[a][s] || (translate.nodeQueue[e].list[a][s] = new Array,
                                translate.nodeQueue[e].list[a][s].nodes = new Array, translate
                                .nodeQueue[e].list[a][s].original = r), translate.nodeQueue[e].list[
                                a][s].nodes[translate.nodeQueue[e].list[a][s].nodes.length] = t
                        }
                    }
                }
            }
        },
        language: {
            local: "chinese_simplified",
            setLocal: function(e) {
                translate.language.local = e
            },
            get: function(e) {
                for (var t = new Array, n = new Array, a = "", l = 0; l < e.length; l++) {
                    var r = e.charAt(l),
                        s = translate.language.getCharLanguage(r);
                    "" == s && (s = "unidentification"), n = translate.language.analyse(s, n, a, r), a =
                        s, t.push(s)
                }
                return void 0 !== n.unidentification && delete n.unidentification, void 0 !== n
                    .specialCharacter && delete n.specialCharacter, n
            },
            getCharLanguage: function(e) {
                return null == e || void 0 === e ? "" : this.english(e) ? "english" : this
                    .specialCharacter(e) ? "specialCharacter" : this.chinese_simplified(e) ?
                    "chinese_simplified" : this.japanese(e) ? "japanese" : this.korean(e) ? "korean" : (
                        console.log("not find is language , char : " + e + ", unicode: " + e.charCodeAt(
                            0).toString(16)), "")
            },
            analyse: function(e, t, n, a) {
                void 0 === t[e] && (t[e] = new Array);
                var l = 0;
                return "" == n || (l = n == e ? t[e].length - 1 : t[e].length), void 0 === t[e][l] && (
                    t[e][l] = ""), t[e][l] = t[e][l] + a, t
            },
            chinese_simplified: function(e) {
                return !!/.*[\u4e00-\u9fa5]+.*$/.test(e)
            },
            english: function(e) {
                return !!/.*[\u0041-\u005a]+.*$/.test(e) || !!/.*[\u0061-\u007a]+.*$/.test(e)
            },
            japanese: function(e) {
                return !!/.*[\u0800-\u4e00]+.*$/.test(e)
            },
            korean: function(e) {
                return !!/.*[\uAC00-\uD7AF]+.*$/.test(e)
            },
            specialCharacter: function(e) {
                return !!/.*[\u2460-\u24E9]+.*$/.test(e) || (!!/.*[\u2500-\u25FF]+.*$/.test(e) || (!!
                    /.*[\u3200-\u33FF]+.*$/.test(e) || (!!/.*[\uFF00-\uFF5E]+.*$/.test(e) || (!!
                    /.*[\u2000-\u22FF]+.*$/.test(e) || (!!/.*[\u3001-\u3036]+.*$/.test(
                    e) || (!!/.*[\u0020-\u007e]+.*$/.test(e) || (!!
                    /.*[\u0009\u000a\u0020\u00A0\u1680\u180E\u202F\u205F\u3000\uFEFF]+.*$/
                        .test(e) || (!!/.*[\u2000-\u200B]+.*$/.test(e) || (!!
                    /.*[\u00A1-\u0105]+.*$/.test(e) || !!
                    /.*[\u2C60-\u2C77]+.*$/.test(e))))))))))
            }
        },
        executeByLocalLanguage: function() {
            translate.request.post("https://api.translate.zvo.cn/ip.json", {}, function(e) {
                0 == e.result ? (console.log("==== ERROR 获取当前用户所在区域异常 ===="), console.log(e
                    .info), console.log("==== ERROR END ====")) : (translate
                    .setUseVersion2(), translate.storage.set("to", e.language), translate
                    .to = e.language, translate.selectLanguageTag, translate.execute())
            })
        },
        util: {
            uuid: function() {
                var e = (new Date).getTime();
                return window.performance && "function" == typeof window.performance.now && (e +=
                    performance.now()), "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx".replace(/[xy]/g,
                    function(t) {
                        var n = (e + 16 * Math.random()) % 16 | 0;
                        return e = Math.floor(e / 16), ("x" == t ? n : 3 & n | 8).toString(16)
                    })
            },
            findTag: function(e) {
                return /<[^>]+>/g.test(e)
            },
            arrayFindMaxNumber: function(e) {
                for (var t = {}, n = [], a = 0, l = 0, r = e.length; l < r; l++) t[e[l]] ? t[e[l]]++ :
                    t[e[l]] = 1, t[e[l]] > a && (a = t[e[l]]);
                for (var s in t) t[s] === a && n.push(s);
                return n
            },
            hash: function(e) {
                if (null == e || void 0 === e) return e;
                var t, n = 0;
                if (0 === e.length) return n;
                for (t = 0; t < e.length; t++) n = (n << 5) - n + e.charCodeAt(t), n |= 0;
                return n + ""
            },
            charReplace: function(e) {
                return null == e ? "" : e = (e = e.trim()).replace(/\t|\n|\v|\r|\f/g, "")
            }
        },
        request: {
            post: function(e, t, n) {
                this.send(e, t, n, "post", !0, {
                    "content-type": "application/x-www-form-urlencoded"
                }, null)
            },
            send: function(e, t, n, a, l, r, s) {
                var o = "";
                if (null != t)
                    for (var u in t) o.length > 0 && (o += "&"), o = o + u + "=" + t[u];
                var i = null;
                try {
                    i = new XMLHttpRequest
                } catch (e) {
                    i = new ActiveXObject("Microsoft.XMLHTTP")
                }
                if (i.open(a, e, l), null != r)
                    for (var u in r) i.setRequestHeader(u, r[u]);
                i.send(o), i.onreadystatechange = function() {
                    if (4 == i.readyState)
                        if (200 == i.status) {
                            var e = null;
                            try {
                                e = JSON.parse(i.responseText)
                            } catch (e) {
                                console.log(e)
                            }
                            n(null == e ? i.responseText : e)
                        } else null != s && s(i)
                }
            }
        },
        storage: {
            set: function(e, t) {
                localStorage.setItem(e, t)
            },
            get: function(e) {
                return localStorage.getItem(e)
            }
        }
    };
    try {
        translate.init()
    } catch (e) {
        console.log(e)
    }

    // 默认版本
    translate.setUseVersion2();

    // 外置实现
    translate.translate = function() {

        var template_temp_pearInterval = setInterval(function() {

            // 等待
            if (typeof(parent.window.pearTranslateConfig) == 'undefined') {
                return;
            }
            // 获取
            var translateConfig = parent.window.pearTranslateConfig;

            clearInterval(template_temp_pearInterval);

            // 初始化配置
            if (typeof(translateConfig.autoDiscriminateLocalLanguage) != 'undefined' && (
                translateConfig.autoDiscriminateLocalLanguage == true || translateConfig
                    .autoDiscriminateLocalLanguage == 'true')) {
                translate.setAutoDiscriminateLocalLanguage(); //设置用户第一次用时，自动识别其所在国家的语种进行切换
            }
            if (typeof(translateConfig.currentLanguage) != 'undefined' && translateConfig
                .currentLanguage.length > 0) {
                translate.language.setLocal(translateConfig.currentLanguage);
            }
            if (typeof(translateConfig.ignoreClass) != 'undefined' && translateConfig.ignoreClass
                .length > 0) {
                var classs = translateConfig.ignoreClass.split(',');
                for (var ci = 0; ci < classs.length; ci++) {
                    var className = classs[ci].trim();
                    if (className.length > 0) {
                        if (translate.ignore.class.indexOf(className.toLowerCase()) > -1) {
                            // 该 class 存在忽略标注，不进行翻译
                        } else {
                            // 翻译列表，将其加入
                            translate.ignore.class.push(className);
                        }
                    }
                }
            }
            if (typeof(translateConfig.ignoreTag) != 'undefined' && translateConfig.ignoreTag
                .length > 0) {
                var tags = translateConfig.ignoreTag.split(',');
                for (var ti = 0; ti < tags.length; ti++) {
                    var tagName = tags[ti].trim();
                    if (tagName.length > 0) {
                        if (translate.ignore.tag.indexOf(tagName.toLowerCase()) > -1) {
                            // 该 class 存在忽略标注，不进行翻译
                        } else {
                            // 翻译列表，将其加入
                            translate.ignore.tag.push(tagName);
                        }
                    }
                }
            }
            // 监听页面变化，实现局部翻译
            translate.listener.start();

            // 等待 document 加载完成，再执行翻译
            if (document.readyState == 'complete') {
                translate.execute();
            } else {
                window.onload = function() {
                    translate.execute();
                }
            }
            // 避免大型组件加载慢，无法立即得到翻译，通过一定程度的延迟来兼容
            setTimeout(translate.execute, 1500);

        }, 30);
    }

    window.translate = translate;

    exports('translate', translate);
});
