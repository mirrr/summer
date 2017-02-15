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
			if (typeof obj.onClick !== 'undefined') {
				onclick = obj.onClick;
				delete obj.onClick;
			}
			var $button = $('<button/>', obj);
			$('#right-panel>div').append($button);
			if ($.isFunction(onclick)) {
				$button.on('click', function (event) {
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

		$.tools.addGroupper = function (tableSelector, data, onclick) {
			var $table = $(tableSelector);
			if (!$table.data('chb') && !$table.find('.chb').length) {
				$table.data('chb', true);
				$table.find('tbody>tr[data-id]').each(function (index, el) {
					var $el = $(el);
					$el.prepend('<td class="chb"><input type="checkbox" data-id="' + $el.data('id') + '" id="check_' + $el.data('id') + '"><label for="check_' + $el.data('id') + '"></label></td>')
				});
				$table.find('thead>tr').first().prepend('<th rowspan="100" style="width: 2rem" class="chb"><input type="checkbox" id="check_all"><label for="check_all"></label></th>')
			}
			if (!$('#float-grouper').length) {
				$('#all').append($('<div/>', {
					id: 'float-grouper',
					class: 'unselectable',
				}));
				$('#float-grouper').append($('<div/>', {
					id: 'float-grouper-in'
				}));
			}
			if (data) {
				if (!Array.isArray(data)) data = [data];
				data.forEach(function (el, i) {
					var $el = $('<span/>', el);
					$el.addClass = 'float-grouper-item';
					$('#float-grouper-in').append($el);

					if ($.isFunction(onclick)) {
						$el.on('click', function (event) {
							event = event || window.event;
							event.preventDefault();
							onclick(this, event);
							return false;
						});
					}
				});
			}
		};
		$.tools.updateGroupper = function (tableSelector) {
			var $table = $(tableSelector);
			if ($table.data('chb')) {
				$table.find('tbody>tr[data-id]').each(function (index, el) {
					var $el = $(el);
					if (!$el.find('.chb').length) {
						$el.prepend('<td class="chb"><input type="checkbox" data-id="' + $el.data('id') + '" id="check_' + $el.data('id') + '"><label for="check_' + $el.data('id') + '"></label></td>')
					}
				});
			}
			if ($table.find('td.chb input[type=checkbox]:checked').length) {
				$('#float-grouper').fadeIn('300')
				$table.find('th.chb input[type=checkbox]').prop('checked', true);
			} else {
				$('#float-grouper').fadeOut('300')
				$table.find('th.chb input[type=checkbox]').prop('checked', false);
			}
		}

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

		$.tools.confirm = function (title, text, callback, target) {
			title = title || window.toolsConfirmTitle || 'Are you sure?';
			text = text || window.toolsConfirmMessage || 'Are you sure that you want to perform this action?';
			var cancel = window.toolsConfirmCancel || 'Cancel';
			var ok = window.toolsConfirmOk || 'OK';
			if (typeof callback !== 'function') {
				throw new Error('Not specified handler function');
			}
			$confirm = $('<div />', {
				'class': 'summer-confirm'
			});
			$confirm.html('<div class="w-title unselectable">' + title +
				'</div>' + (text ? '<div class="w-content-body">' + text + '</div>' : '') +
				'<div class="form-footer"><button type="cancel" class="shwark-close">' + cancel + '</button>' +
				'<button type="submit">' + ok + '</button></div>')

			close = shwOpen(null, {
					wrapperBackground: 'rgba(0,0,0,0.3)',
					wrapperClass: 'summer-confirm-wrapper'
				}, target,
				$confirm
			)
			$confirm.find('button[type="submit"]').click(function () {
				setTimeout(close, 1);
				callback();
			});
		};


		/**
		 * "Smart" handler for links
		 *
		 * @param  {String} selector
		 * @param  {Function} func
		 * @param  {String} target
		 * @param  {Function} confirm
		 */
		$.tools.forceClick = function (selector, func, target, confirm) {

			if (typeof func !== 'function') {
				throw new Error('Not specified handler function');
			}
			if (!target) {
				target = $('body');
			} else {
				target = $(target);
			}
			if (typeof confirm !== 'function') {
				confirm = function (element, event, callback) {
					if ($(element).hasClass('need-confirm')) {
						$.tools.confirm(null, null, function () {
							setTimeout(callback, 0);
						}, $(element));
					} else {
						setTimeout(callback, 0);
					}
				}
			}

			target.off('click', selector);
			target.on('click', selector, function (event) {
				var that = this;
				event = event || window.event;
				event.preventDefault();
				confirm(that, event, function () {
					func.call(that, event);
				});
				return false;
			});
			return this;
		};


		/**
		 * Handler for action buttons
		 * (INFO: Instead the options object you can put function that return the options object
		 *
		 * @param  {Object/Function} options
		 * @example
		 * $.tools.ajaxActionSender('table .status, table .remove, table .trash', {
		 * 		url:     'http://...',              // адрес отправки (options.url или window.baseUrl + '/' + options.target)
		 * 		method:  'POST',                    // метод отправки (по умолчанию 'GET')
		 * 		action:  'edit',                    // ajax-действие  (по умолчанию берется из data-action)
		 * 		timeout: 15000,                     // 15 sec.
		 * 		check:   function () {....},        // колбек проверки (return true/false - успешность проверки)
		 * 		success: function (result) {....},  // колбек на успешное выполнение (return true - закрыть окно)
		 * 		error:   function (result) {....}   // колбек на ошибку в ответе сервера,
		 * 		confirm: function (element, event, callback) {callback();} // функция подтверждения
		 * 		                                                           // (вместо стандартной, при подтверждении вызвать callback)
		 * });
		 */
		$.tools.ajaxActionSender = function (selector, options) {

			// Handler
			$.tools.forceClick(selector, function (event) {
				var opt = {};

				if (typeof options === 'function') {
					opt = options.call(this, event);
				} else {
					opt = typeof options === 'object' ? options : {};
				}

				var $this = $(event.target);
				// for embedded 'span' and simillar.
				if (!$this.data('id') && $this.parent().data('id')) {
					$this = $this.parent();
				}

				var settings = $.extend({
					action: $this.data('action') || 'edit',
					data: {},
					dataKeys: [],
					dataId: 'id',
					id: null,
					method: 'GET',
					remove: false,
					timeout: 15000,
					url: opt.target ? ((window.baseUrl ? window.baseUrl + '/' : '') + opt.target) : '',
					check: (function () {
						return true;
					}),
					success: (function (result) {
						$.message.ok(result.message);
					}),
					error: (function (result) {
						$.message.ajaxWarn(result);
					}),
					after: (function () {})
				}, opt);

				if ($this.data('id')) {
					if (typeof settings.selector === 'string') {
						settings.selector = settings.selector.replace('$id', $this.data('id'));
					}

					if (typeof settings.findSelector === 'string') {
						settings.findSelector = settings.findSelector.replace('$id', $this.data('id'));
					}

					if (typeof settings.bodySelector === 'string') {
						settings.bodySelector = settings.bodySelector.replace('$id', $this.data('id'));
					}
				}

				var sel = settings.selector ? $this.closest(settings.selector) : (
					settings.findSelector ? $this.find(settings.findSelector) : (
						settings.bodySelector ? $('body').find(settings.bodySelector) : (
							settings.thisSelector ? $this : $this.closest('tr')
						)
					)
				);
				var id = settings.id || sel.data(settings.dataId);

				settings.dataKeys.forEach(function (el) {
					if (sel.data(el)) {
						settings.data[el] = sel.data(el);
					}
				});

				var data = $.extend({
					action: settings.action,
					id: id
				}, settings.data);
				settings.ok = function () {
					$.progress.start();
					$.ajax({
						url: settings.url,
						type: settings.method,
						data: data,
						timeout: settings.timeout,
						success: function (result) {
							$.progress.stop();
							if (settings.success(result, $this, settings, sel) || settings.remove) {
								sel.hide().remove();
							}

							settings.after(result);
						},
						error: function (result) {
							$.progress.stop();
							settings.error(result, $this, settings);
							settings.after(result);
						}
					});
				};

				if (settings.url && settings.check(data, $this, settings)) {
					settings.ok();
				}
			}, null, options.confirm);
		};


		var oldText = '';
		var timerId = null;
		$.tools.addSearchFn = $.tools.searcher = function (onChange) {
			var $search = $('input[type=text].allsearch');
			if ($search.length && !$search.parent('.li-search').is(':visible')) {
				$search.parent('.li-search').show();
				$search.on('keyup', function (e) {
					if (oldText !== $search.val()) {
						if (timerId) {
							clearTimeout(timerId);
							timerId = null;
						}
						timerId = setTimeout(function () {
							onChange($search.val());
						}, 350);
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
	$('#content th[data-sorter]').each(function (index, el) {
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
