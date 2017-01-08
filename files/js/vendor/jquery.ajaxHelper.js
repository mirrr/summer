(function ($) {
	/**
	 * Вешаем обработчик на отправку формы через ajax
	 * @param  {Object} options
	 * @example
	 * $('#form-ajax').ajaxFormSender ({
	 * 		url:     'http://...',              // адрес отправки (по умолчанию берется из формы или из options.action)
	 * 		method:  'POST',                    // метод отправки (по умолчанию берется из формы или 'POST')
	 * 		timeout: 15000,                     // 15 sec.
	 * 		check:   function () {....},        // колбек проверки (return true/false - успешность проверки)
	 * 		success: function (result) {....},  // колбек на успешное выполнение (return true - закрыть окно)
	 * 		error:   function (result) {....}   // колбек на ошибку в ответе сервера
	 * });
	 */
	$.fn.ajaxFormSender = function (options) {
		options = options || {};

		// Обработчик
		$('body').off('submit', this.selector);
		$('body').on('submit', this.selector, function (event) {
			event = event || window.event;
			event.preventDefault();

			var $this = $(event.target);
			var action = options.action || $this.attr('action') || '';
			var settings = $.extend({
				url: ((window.baseUrl && action.indexOf('/') === -1) ? (window.baseUrl + '/') : '') + action,
				method: $this.attr('method') || 'POST',
				timeout: 15000,
				check: (function () {
					return true;
				}),
				success: (function (result) {
					$.message.ok(result.message);
					return true;
				}),
				error: (function (result) {
					$.message.ajaxWarn(result);
				}),
				after: (function () {})
			}, options);

			var data = $this.serialize();

			var obj = {};
			$this.serializeArray().forEach(function (el) {
				if (typeof obj[el.name] === 'undefined' && el.name.indexOf('[]') === -1) {
					obj[el.name] = el.value;
				} else {
					var name = el.name.replace('[]', '');

					if (typeof obj[name] === 'undefined') {
						obj[name] = [];
					} else if (!Array.isArray(obj[name])) {
						obj[name] = [obj[name]];
					}

					obj[name].push(el.value);
				}
			});

			settings.ok = function () {
				$.progress.start();
				$.ajax({
					url: settings.url,
					type: settings.method,
					data: data,
					timeout: settings.timeout,
					success: function (result) {
						$.progress.stop();
						if (settings.success(result)) {
							$.wbox.close();
						}

						settings.after(result);
					},
					error: function (result) {
						$.progress.stop();
						settings.error(result);
						settings.after(result);
					}
				});
			};

			if (settings.check(obj, $this, settings)) {
				settings.ok();
			}

			return false;
		});
	};


	/**
	 * Вешаем обработчик на кнопки действий в таблице
	 * (INFO: Вместо объекта опций можно передать ф-ю, отдающую этот объект)
	 *
	 * @param  {Object/Function} options
	 * @example
	 * $('table .status, table .remove, table .trash').ajaxActionSender ({
	 * 		url:     'http://...',              // адрес отправки (options.url или window.baseUrl + '/' + options.target)
	 * 		method:  'POST',                    // метод отправки (по умолчанию 'GET')
	 * 		action:  'edit',                    // ajax-действие  (по умолчанию берется из data-action)
	 * 		timeout: 15000,                     // 15 sec.
	 * 		check:   function () {....},        // колбек проверки (return true/false - успешность проверки)
	 * 		success: function (result) {....},  // колбек на успешное выполнение (return true - закрыть окно)
	 * 		error:   function (result) {....}   // колбек на ошибку в ответе сервера
	 * });
	 */
	$.fn.ajaxActionSender = function (options) {

		// Обработчик
		$(this.selector).forceClick(function (event) {
			var opt = {};

			if (typeof options === 'function') {
				opt = options.call(this, event);
			} else {
				opt = typeof options === 'object' ? options : {};
			}

			var $this = $(event.target);
			// для вложенных span и т.д.
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
		});
	};


	/**
	 * "Умный" обработчик для линков
	 *
	 * @param  {Function} func
	 */
	$.fn.forceClick = function (func, target) {
		if (typeof func !== 'function') {
			throw new Error('Not specified handler function');
		}
		if (!target) {
			target = $('body');
		} else {
			target = $(target);
		}

		target.off('click', this.selector);
		target.on('click', this.selector, function (event) {
			event = event || window.event;
			event.preventDefault();
			if ($(this).hasClass('need-confirm')) {
				$.tools.confirm('Are you sure?', 'Are you sure that you want to perform this action?', function () {
					func.call(this, event);
				}, $(this));
			} else {
				func.call(this, event);
			}
			return false;
		});
		return this;
	};


	/**
	 * Простейшая загрузка элементов таблицы или списка (with doT.js)
	 *
	 * @param  {Object} options
	 *
	 * @example
	 * $('table>body').listLoad ({
	 * 		url:     'http://...',              // адрес отправки (options.url или window.baseUrl + '/' + options.target)
	 * 		// или :
	 * 		target:  'edit',                    // ajax-контроллер (для кабмина)
	 * 		itemTpl: 'item',                    // шаблон doT.js
	 * 		noitemsTpl: 'noitems',              // шаблон doT.js
	 * 		timeout: 15000,                     // 15 sec.
	 * 		success: function (result) {....},  // колбек на успешное выполнение (return true - закрыть окно)
	 * 		error:   function (result) {....}   // колбек на ошибку в ответе сервера
	 * });
	 *
	 */
	$.fn.listLoad = function (options) {
		var $this = $(this.selector);
		options = $.extend({
			target: $this.data('target') || '',
			itemTpl: 'item',
			data: {},
			method: 'GET',
			timeout: 15000,
			success: (function () {}),
			emptylist: (function () {}),
			error: (function (result) {
				$.message.ajaxWarn(result);
			}),
			after: (function () {})
		}, options);

		if (!options.url && options.target) {
			options.url = (window.baseUrl ? window.baseUrl + '/' : '') + options.target;
		}

		$.progress.start();
		$.ajax({
			url: options.url,
			type: options.method,
			data: options.data,
			timeout: options.timeout,
			success: function (result) {
				$.progress.stop();
				empty = false;
				if (Array.isArray(result.data) && result.data.length) {
					$this.tpl(options.itemTpl, result.data);
				} else {
					if (options.noitemsTpl) {
						$this.tpl(options.noitemsTpl);
					} else {
						$this.tpl(options.itemTpl, []);
					}
					empty = true;
				}

				if (typeof options.success === 'function') {
					options.success(result);
				}

				if (empty) {
					options.emptylist($this);
				}

				options.after(result);
			},
			error: function (result) {
				$.progress.stop();
				if (options.noitemsTpl) {
					$this.tpl(options.noitemsTpl);
				}

				if (typeof options.error === 'function') {
					options.error(result);
				}

				options.after(result);
			}
		});
	};

})(jQuery);
