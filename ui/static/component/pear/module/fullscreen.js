layui.define(['message', 'table', 'jquery', 'element', 'yaml', 'form', 'tab', 'menu', 'frame', 'theme', 'convert'],
    function(exports) {
        "use strict";

        var $ = layui.jquery;
        var fullScreen = new function() {
            this.fullScreen = function(dom){
                return new Promise(function(res, rej) {
                    var docElm = dom && document.querySelector(dom) || document.documentElement;
                    if (docElm.requestFullscreen) {
                        docElm.requestFullscreen();
                    } else if (docElm.mozRequestFullScreen) {
                        docElm.mozRequestFullScreen();
                    } else if (docElm.webkitRequestFullScreen) {
                        docElm.webkitRequestFullScreen();
                    } else if (docElm.msRequestFullscreen) {
                        docElm.msRequestFullscreen();
                    }else{
                        rej("");
                    }
                    res("返回值");
                });
            }
            this.fullClose = function(){
                if (document.exitFullscreen) {
                    document.exitFullscreen();
                } else if (document.mozCancelFullScreen) {
                    document.mozCancelFullScreen();
                } else if (document.webkitCancelFullScreen) {
                    document.webkitCancelFullScreen();
                } else if (document.msExitFullscreen) {
                    document.msExitFullscreen();
                }
                return new Promise(function(res, rej) {
                    res("返回值");
                });
            }
            this.isFullscreen = function(){
                return document.fullscreenElement ||
                    document.msFullscreenElement ||
                    document.mozFullScreenElement ||
                    document.webkitFullscreenElement || false;
            }
        };
        exports('fullscreen', fullScreen);
    })
