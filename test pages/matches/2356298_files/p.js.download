(function () {
	var retargetUrl = window.location.href;
	if (window.parent && window.location !== window.parent.location) {
		retargetUrl = window.parent.document.referrer;
	}
	const ticker = window.setInterval(function() {
		if (!window.document.body) {
			return;
		}
		clearInterval(ticker);

		var pixel = document.createElement('img');
		pixel.setAttribute(
			'style',
			'position: absolute; width: 1px; height: 1px; left: 0px; bottom: 0px; opacity: 0;',
		)
		pixel.src = 'https://my.rtmark.net/img.gif?f=sync&partner=68c84d97734bf7f44eec38b5777a048387b640c517aa53ad20b91247a58f7296&ttl=&rurl=' + encodeURIComponent(retargetUrl);
		window.document.body.appendChild(pixel);
	}, 500);
})();
