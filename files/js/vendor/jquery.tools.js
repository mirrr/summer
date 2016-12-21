/*
 * Методы управления всплывающими уведомлениями
 */

/* global define */

(function (factory) {
		if (typeof define === 'function' && define.amd) {
			define(['jquery'], factory);
		} else {
			factory(jQuery);
		}
	}

	(function ($) {
		var progressInterval = null;
		var progressTimeout = null;
		$.tools = function () {};
		$.progress = function () {};

		$.progress.start = function () {
			if (progressInterval) {
				tools.progressStop(0);
			}
			progressTimeout = setTimeout(function () {
				var $s = $('#summer-progress');
				var $sp = $('#summer-progress div');
				var i = $s.outerWidth() / 20;
				$sp.width(1);
				progressInterval = setInterval(function () {
					console.log($sp.outerWidth());
					$sp.width($sp.outerWidth() + i);
					i = (i / 1.05) + 1;
				}, 20);
			}, 100);
		};
		$.progress.stop = function () {
			clearTimeout(progressTimeout);
			if (typeof t === "undefined") {
				t = 500;
			}

			setTimeout(function (t) {
				if (progressInterval) {
					clearInterval(progressInterval);
					progressInterval = null;
					$('#summer-progress div').width(0);
				}
			}, t);
		};
		$.tools.addButton = function (obj) {
			var $button = $('<button/>', obj);
			$("#right-panel").append($button);
			return $button;
		};
		$.tools.addLink = function (obj) {
			var $link = $('<a/>', obj);
			$("#right-panel").append($link);
			return $link;
		};
		$.tools.addSearch = function (obj) {
			if ($('#right-panel .search').length === 0) {
				obj = obj || {};
				obj.class = 'search hidden';
				obj.type = 'text';
				var $search = $('<input/>', obj);
				$("#right-panel").append($search);

				var $allsearch = $('input[type=text].allsearch');
				if ($allsearch.length) {
					$allsearch.show();
					$allsearch.on('keyup', function (e) {
						$search.val($allsearch.val());
						var event = new Event('keyup');
						event.keyCode = e.originalEvent.keyCode;
						event.keyIdentifier = e.originalEvent.keyIdentifier;
						event.charCode = e.originalEvent.charCode;
						event.code = e.originalEvent.code;
						event.which = e.originalEvent.which;
						$search.get(0).dispatchEvent(event);
					});
					$allsearch.on('input', function () {
						$search.val($allsearch.val());
						var event = new Event('input');
						$search.get(0).dispatchEvent(event);
					});
				}
			}
			return $search;
		}
	}));
