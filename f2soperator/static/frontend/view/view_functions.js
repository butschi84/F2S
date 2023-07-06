
//Button: Anclicken eines Mitglieds
$("#functionstable").delegate( "#invoke", "click", function() {
	var id = $(this).attr("title");
	$("#showfunction").val(id);
	$("#main").load("frontend/view/view_invoke.html", function() {
		$.getScript("frontend/view/view_invoke.js")
	});

});


function load_functions()
{
    $.get("/functions", function(data) {
		//Inhalt der Tabelle zuruecksetzen
		$("#functionstable").html("");

		//Alle Resultate durchgehen
		for(i=0;i<data.length;i++)
		{
			//Eintrag vorbereiten
			var functionentry = "<tr>";
			functionentry += "<td>"+data[i].uid+"</td>";
			functionentry += "<td>"+data[i].name+"</td>";
			functionentry += "<td>"+data[i].spec.endpoint+"</td>";
			functionentry += "<td>"+data[i].target.containerImage+"</td>";
			functionentry += "<td><button class='button is-primary' id=invoke title="+data[i].uid+">Invoke</button></td>";
			functionentry += "</tr>";

			//Eintrag ausgeben
			$("#functionstable").append(functionentry);
		}

 		$('#tfunctions').DataTable( {
			info: false,
			lengthChange: false,
			columnDefs: [ {
				targets: [ 1 ],
				orderData: [ 1, 0 ]
			}]
		});
    });
}

load_functions()