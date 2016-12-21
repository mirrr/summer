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

		var oldText = "";
		var timerId = null;
		$.tools.searcher = function (onChange) {
			var $search = $('input[type=text].allsearch');
			if ($search.length && !$search.is(":visible")) {
				$search.show();
				$search.on('keyup', function (e) {
					if (oldText !== $search.val()) {
						if (timerId) {
							clearTimeout(timerId);
							timerId = null;
						}
						timerId = setTimeout(function () {
							onChange($search.val());
						}, 300);
						oldText = $search.val();
					}
				});
			}
		}
	}));
