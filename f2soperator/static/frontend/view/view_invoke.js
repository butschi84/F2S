
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
        $.get(target, function(response) {
            $("#result").html(response)
            console.log(response)
        }).fail(function(xhr, status, error) {
            // Handle the error
            console.log(error);
        });
    })
})

