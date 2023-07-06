console.log("loading functions")

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
			functionentry += "<td><button class='button is-primary'>Invoke</button></td>";
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