layui.define(['jquery', 'layer'], function (exports) {

    let MOD_NAME = 'sail',
        $ = layui.jquery,
        layer = layui.layer;

    var sail = new function () {

        this.ajaxError = function (jqXHR, textStatus) {
            if (jqXHR.status === 401) {
                return
            }

            let sjson = jqXHR.responseJSON;
            if (sjson.hasOwnProperty("message")) {
                layer.msg(sjson.message, {
                    icon: 2,
                    time: 1000
                });
            } else {
                layer.msg(textStatus.toString(), {
                    icon: 2,
                    time: 1000
                });
            }
        }

        this.checkSuccess = function (result, callback) {
            if (result.code === 0) {
                callback(result);
            } else {
                layer.msg(result.message, {
                    icon: 2,
                    time: 1000
                });
            }
        }

        this.setAuth = function (token) {
            return {
                "Authorization": "Bearer " + token,
            }
        }

        let setAuthor = function (token) {
            return {
                "Authorization": "Bearer " + token,
            }
        }

        let delRefreshAndRelogin = function () {
            localStorage.removeItem("accessToken");
            localStorage.removeItem("refreshToken");
            layer.msg("登录过期，请重新登录", {
                icon: 3,
                time: 1000
            });
            top.location.href = "/ui/login";
        }

        this.prefilterAjax = function () {
            $.ajaxPrefilter(function (options, originalOptions, jqXHR) {
                originalOptions._error = originalOptions.error;
                // overwrite error handler for current request
                options.error = function (_jqXHR, _textStatus, _errorThrown) {
                    if (jqXHR.status !== 401) {
                        if (originalOptions._error) originalOptions._error(_jqXHR, _textStatus, _errorThrown);
                        return;
                    }
                    let refreshToken = localStorage.getItem("refreshToken");
                    $.ajax({
                        type: "POST",
                        url: "/sail/v1/login/refresh",
                        data: JSON.stringify({old_refresh_token: refreshToken}),
                        contentType: "application/json;charset=utf-8",
                        cache: false,
                        dataType: "JSON",
                        success: function (result) {
                            if (result.code !== 0) {
                                delRefreshAndRelogin();
                                return
                            }
                            localStorage.setItem("accessToken", result.data.access_token);
                            localStorage.setItem("refreshToken", result.data.refresh_token);
                            originalOptions.headers = setAuthor(result.data.access_token);
                            $.ajax(originalOptions);
                        },
                        error: function (jqXHR, textStatus) {
                            delRefreshAndRelogin();
                        }
                    })
                };
            });
        }
    }
    exports(MOD_NAME, sail);
});
