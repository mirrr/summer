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
		$.tools = function () {};
		$.tools.addButton = function (obj) {
			var onclick = obj.onClick;
			if (typeof obj.onClick !== "undefined") {
				delete obj.onClick;
			}
			var $button = $('<button/>', obj);
			$("#right-panel>div").append($button);
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
			$link = $.tools.createBoxLink(obj);
			$('#right-panel>div').append($link);
			return $link;
		};
		$.tools.createBoxLink = function (obj) {
			var $link = $('<a/>', obj);
			$link.attr({
				'target': '_blank',
				'title': $link.text()
			}).html($('<span/>', {
				'class': 'text',
				'html': $link.html(),
			})).addClass('button');
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
