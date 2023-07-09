
function load_function(uid)
{
    $.get("/functions", function(data) {
        let f = _.find(data, d => {
            return d.uid == uid
        })
        console.log(f)

        const {
            host, hostname, href, origin, pathname, port, protocol, search
          } = window.location

        $("#name").html(f.name)
        $("#target").html(origin + "/invoke" + f.spec.endpoint)
    });
}

$(document).ready(function () {
    load_function($("#showfunction").val())

    $("#invokefu").on("click", function() {
        var target = $("#target").html()

        $("#spinner").css("display", "inline");
        $("#invokefu").prop("disabled", true);

        // reset response view
        $('#jsonContainer').html("");

        // invoke function
        $.get(target, function(response) {
            $('#jsonContainer').JSONView(response);

            $("#spinner").css("display", "none");
            $("#invokefu").prop("disabled", false);
        }).fail(function(xhr, status, error) {
            // Handle the error
            console.log(error);
            $("#spinner").css("display", "none");
            $("#invokefu").prop("disabled", false);
        });
    })
})

