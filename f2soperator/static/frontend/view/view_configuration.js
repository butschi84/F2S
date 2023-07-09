$.get("/config", function(data) {
    $("#http_timeout").val(data.Config.F2S.Timeouts.HttpTimeout)
    $("#request_timeout").val(data.Config.F2S.Timeouts.RequestTimeout)
    $("#scaling_timeout").val(data.Config.F2S.Timeouts.ScalingTimeout)
});