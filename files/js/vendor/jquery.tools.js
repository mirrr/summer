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
		$.tools.addButton = function (obj, onclick) {
			if (typeof obj.onClick !== "undefined") {
				onclick = obj.onClick;
				delete obj.onClick;
			}
			var $button = $('<button/>', obj);
			$("#right-panel>div").append($button);
			if ($.isFunction(onclick)) {
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
		$.tools.addSorterFn = function (fn) {
			if ($.isFunction(fn)) {
				$(window).on('table-sorter', function (event, name, direction) {
					fn(name, direction)
				});
			}
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
		$.tools.addSearchFn = $.tools.searcher = function (onChange) {
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
	})
);

$(function () {
	function clearSorterD() {
		$('#content th[data-sorter]').data('sort-direction', 0)
			.find('.sort-ind').removeClass('fa-caret-down').removeClass('fa-caret-up').addClass('fa-unsorted');
	}
	$("#content th[data-sorter]").each(function (index, el) {
		$(el).css({
			'font-weight': 'bold',
			'cursor': 'pointer'
		});
		if (!$(el).find('.sort-ind').length) {
			$(el).append($('<span/>', {
				'class': 'fa fa-unsorted sort-ind'
			}));
			$(el).data('sort-direction', 0);
		}
		$(el).on('mousedown', function (event) {
			event = event || window.event;
			event.preventDefault();

			var $sortIn = $(el).find('.sort-ind');
			if ($(el).data('sort-direction') === 1) {
				clearSorterD();
				$sortIn.removeClass('fa-unsorted').addClass('fa-caret-up');
				$(el).data('sort-direction', -1)
			} else {
				clearSorterD();
				$sortIn.removeClass('fa-unsorted').addClass('fa-caret-down');
				$(el).data('sort-direction', 1)
			}
			$(window).trigger('table-sorter', [$(el).data('sorter'), $(el).data('sort-direction')]);
			return false;
		});
	});
});
