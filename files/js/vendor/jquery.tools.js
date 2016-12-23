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
			var $s = $('#summer-progress');
			var $sp = $('#summer-progress div');
			var i = $s.outerWidth() / 20;
			$sp.width(1);
			progressInterval = setInterval(function () {
				$sp.width($sp.outerWidth() + i);
				i = (i / 1.08) + 1;
			}, 10);
		};
		$.progress.stop = function () {
			clearTimeout(progressTimeout);
			if (typeof t === "undefined") {
				t = 400;
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
			var onclick = obj.onClick;
			if (typeof obj.onClick !== "undefined") {
				delete obj.onClick;
			}
			var $button = $('<button/>', obj);
			$("#right-panel").append($button);
			if (onclick) {
				$button.on("click", function (event) {
					event = event || window.event;
					event.preventDefault();
					onclick(this, event);
					return false;
				});
			}
			return $button;
		};
		$.tools.addLink = function (obj) {
			var $link = $('<a/>', obj);
			$link.attr('target', '_blank');
			$("#right-panel").append($link);
			return $link;
		};

		var oldText = "";
		var timerId = null;
		$.tools.searcher = function (onChange) {
			var $search = $('input[type=text].allsearch');
			if ($search.length && !$search.parent(".li-search").is(":visible")) {
				$search.parent(".li-search").show();
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
