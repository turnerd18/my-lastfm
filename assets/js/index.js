var curMonth = "";
var monthRegex = /((?:Jan(?:uary)?|Feb(?:ruary)?|Mar(?:ch)?|Apr(?:il)?|May|Jun(?:e)?|Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|Sept|Oct(?:ober)?|Nov(?:ember)?|Dec(?:ember)?)),?\s[0-9]+/;
var months = {
	"Jan": "January",
	"Feb": "February",
	"Mar": "March",
	"Apr": "April",
	"May": "May",
	"Jun": "June",
	"Jul": "July",
	"Aug": "August",
	"Sep": "September",
	"Oct": "October",
	"Nov": "November",
	"Dec": "December"
};
var page = 1;

function make_first_month() {
	var str = $($("ul#grid").children()[0]).text();
	var matches = monthRegex.exec(str);
	var firstMonth = matches[1];
	$("ul#grid").before("<div class=\"row\"><div class=\"month columns\"><h3>" + months[firstMonth] + "</h3></div></div>");
	curMonth = firstMonth;
}

function make_months() {
	var children = $("ul#grid").children();
	for( var i = 0, size = children.length; i < size; i++ ) {
		var matches = monthRegex.exec($(children[i]).text())
		var thisMonth = matches[1];
		if( thisMonth != curMonth ) {
			$($(children[i]).parent())
				.before("<ul class=\"small-block-grid-2 medium-block-grid-3 large-block-grid-5\"/>");
			var beforeList = $(children.slice(0,i).clone());
			$("ul#grid").prev().html(beforeList);
			$("ul#grid").before("<div class=\"row\"><div class=\"month columns\"><h3>" + months[thisMonth] + "</h3></div></div>");
			$(children.slice(0,i)).remove();
			curMonth = thisMonth;
		}
	}
}

$(document).ready(function() {
	make_first_month();
	make_months();
	$("a#more").click(function() {
		// Variables
		var url = "/indexmore";
		var data = "page=" + page;
		var request = new XMLHttpRequest();

		// AJAX request
		request.open("POST", url, true);
		request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
		// request response actions
		request.onreadystatechange = function() {
			if(request.readyState == 4 && request.status == 200) {
				$("ul#grid").append(request.responseText);
				page++;
				make_months();
			}
		}
		// send request to server with data
		request.send(data);
	});
});
