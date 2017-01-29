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
	 * Do not use this!
	 * The function is deprecated! Use $.tools.ajaxActionSender instead current function!
	 */
	$.fn.ajaxActionSender = function (options) {
		$.tools.ajaxActionSender(this.selector, options)
		console.error("Function $(selector).ajaxActionSender(options) is deprecated and will be removed in the next release! Please use $.tools.ajaxActionSender(selector, options)");
	};


	/**
	 * Do not use this!
	 * The function is deprecated! Use $.tools.forceClick instead current function!
	 */
	$.fn.forceClick = function (func, target, confirm) {
		$.tools.forceClick(this.selector, func, target, confirm)
		console.error("Function $(selector).forceClick(func) is deprecated and will be removed in the next release! Please use $.tools.forceClick(selector, func)");
	};


	/**
	 * The simplest elements loading table or list (with doT.js)
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
			after: (function () {}),
			before: (function () {}),
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
				options.before(result);
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

				if ($this.parent('table').length || $this.parent().parent('table').length || $this.parent().children('table').length) {
					$.tools.updateGroupper($this.closest('table'));
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
				options.before(result);
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
