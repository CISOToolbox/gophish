/*
    scanner.js
    Loads and saves the global scanner-detection settings (admin only).
*/

function load() {
    api.scanner.get()
        .success(function (s) {
            $("#window_seconds").val(s.window_seconds)
            $("#min_chrome").val(s.min_chrome)
            $("#min_firefox").val(s.min_firefox)
            $("#min_safari").val(s.min_safari)
            $("#user_agents").val(s.user_agents || "")
        })
        .error(function () {
            errorFlash("Error fetching scanner settings")
        })
}

function save() {
    var settings = {
        window_seconds: parseInt($("#window_seconds").val(), 10) || 0,
        min_chrome: parseInt($("#min_chrome").val(), 10) || 0,
        min_firefox: parseInt($("#min_firefox").val(), 10) || 0,
        min_safari: parseInt($("#min_safari").val(), 10) || 0,
        user_agents: $("#user_agents").val()
    }
    api.scanner.put(settings)
        .success(function () {
            successFlash("Scanner settings saved")
        })
        .error(function (data) {
            errorFlash((data.responseJSON && data.responseJSON.message) || "Error saving scanner settings")
        })
}

$(document).ready(function () {
    load()
})
