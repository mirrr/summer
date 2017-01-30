function shwOpen(template, options, $this, content) {
	if (!$this) {
		$this = $('body')
		options.toRight = false;
		options.toDown = false;
	} else {
		$this = $($this)
	}
	options = $.extend({
		wrapperBackground: 'transparent',
		wrapperClass: '',
		template: false,
		data: {},
		toRight: true,
		toDown: true,
		speed: 300,
		closeAll: true,
		beforeClose: function () {},
		afterOpen: function () {},
		beforeOpen: function (data) {
			return data
		},
		offset: 0
	}, options);

	var $dropdown = $('<div class="shwark-sandbox"></div>');
	var data = {};
	$('body').append($dropdown);
	$dropdown.hide();

	data = options.data;
	if ($.isFunction(options.data)) {
		data = options.data($this);
	}
	if ($.isFunction(options.beforeOpen)) {
		data = options.beforeOpen(data);
	}

	if (content) {
		$dropdown.html(content);
		var $cnt = $dropdown.find(".w-content-body")
		if ($cnt.length && template) {
			$cnt.tpl(template, data);
		}
	} else {
		$dropdown.tpl(template, data);
	}
	if ($.isFunction(options.afterOpen)) {
		options.afterOpen();
	}
	var $wrapper = $('<div class="shwark-wrapper"></div>');
	var $circle = $('<div class="shwark-reddot"></div>');
	var $circleIn = $('<div class="shwark-in-reddot"></div>');
	var $button = $this || $wrapper;
	var $target = options.target === 'body' ? $wrapper : (
		$(options.target).length ? $($(options.target).get(0)) : $button
	);

	if ($dropdown.length) {
		$dropdown = $($dropdown.get(0));
		$('body').append($wrapper);
		$wrapper.append($circle);
		$circleIn.append($dropdown);
		$circle.append($circleIn);
		$dropdown.show();
		$circleIn.show();

		var buttonSize = {
			height: $button.outerHeight(),
			width: $button.outerWidth()
		};
		var dropSize = {
			height: $circleIn.outerHeight(),
			width: $circleIn.outerWidth()
		};
		var targetSize = {
			height: $target.outerHeight(),
			width: $target.outerWidth()
		};
		var winSize = {
			height: $(window).height(),
			width: $(window).width()
		};

		var ofst = {
			left: Math.floor($button.offset().left - $('body').scrollLeft()),
			top: Math.floor($button.offset().top - $('body').scrollTop())
		};
		var newOfst = {
			left: Math.floor($target.offset().left - $('body').scrollLeft()),
			top: Math.floor($target.offset().top - $('body').scrollTop())
		};

		var d = Math.sqrt(Math.pow(dropSize.height, 2) + Math.pow(dropSize.width, 2)) * 2 + 20; // diameter

		if (winSize.width - 10 < newOfst.left + dropSize.width) {
			options.toRight = false;
			options.toLeft = true;
		}

		if (winSize.height - 10 < newOfst.top + dropSize.height) {
			options.toDown = false;
		}

		var dropDirection = {};
		var dropStart = {};

		if (!options.toDown || $button !== $target) {
			dropDirection.top = (d / 2 - dropSize.height / 2) + 'px';
			dropStart.top = -(dropSize.height / 2) + 'px';
			newOfst.top += targetSize.height / 2;
			ofst.top += buttonSize.height / 2;
		} else if (options.toDown) {
			dropDirection.top = (d / 2) + 'px';
			dropStart.top = '0px';
		} else {
			dropDirection.top = (d / 2 - dropSize.height) + 'px';
			dropStart.top = -dropSize.height + 'px';
			newOfst.top += targetSize.height;
			ofst.top += buttonSize.height;
		}

		if (!options.toRight && !options.toLeft || $button !== $target) {
			$wrapper.addClass('centered');
			dropDirection.left = (d / 2 - dropSize.width / 2) + 'px';
			dropStart.left = -(dropSize.width / 2) + 'px';
			newOfst.left += targetSize.width / 2;
			ofst.left += buttonSize.width / 2;
		} else if (options.toRight) {
			dropDirection.left = (d / 2) + 'px';
			dropStart.left = '0px';
		} else {
			dropDirection.left = (d / 2 - dropSize.width) + 'px';
			dropStart.left = -dropSize.width + 'px';
			newOfst.left += targetSize.width;
			ofst.left += buttonSize.width;
		}

		// animation
		$wrapper.css({
			'background': options.wrapperBackground
		}).addClass(options.wrapperClass);
		$circleIn.css(dropStart).animate(dropDirection, {
			duration: options.speed,
			queue: false
		}, function () {});
		$circle.css({
			'borderRadius': 0,
			'left': ofst.left + 'px',
			'top': ofst.top + 'px',
		}).animate({
			'width': '+=' + d + 'px',
			'height': '+=' + d + 'px',
			'left': newOfst.left - (d / 2) + 'px',
			'top': newOfst.top - (d / 2) + 'px',
			'backgroundColor': 'rgba(0,0,0,0)'
		}, {
			duration: options.speed,
			queue: false,
			done: function () {
				$circle.css({
					'borderRadius': 0
				});
			}
		});

		var close = function () {
			options.beforeClose();
			$circleIn.animate(dropStart, {
				duration: options.speed / 1.5,
				queue: false
			}, function () {});
			$circle.css({
				'borderRadius': '0'
			}).animate({
				'width': '-=' + d + 'px',
				'height': '-=' + d + 'px',
				'left': ofst.left + 'px',
				'top': ofst.top + 'px',
				'backgroundColor': 'rgba(0,0,0,0)'
			}, {
				duration: options.speed / 1.5,
				queue: false,
				done: function () {
					$wrapper.remove();
				}
			});
		}

		$circleIn.on('mouseup', function (event) {
			event = event || window.event; // For IE
			event.preventDefault();
			event.stopPropagation();
			return false;
		});


		$("body").on("all-shwark-close", close);
		$dropdown.find(".shwark-close").on('mouseup', close);
		$dropdown.find(".all-shwark-close").on('mouseup', function (event) {
			$("body").trigger("all-shwark-close");
		});

		$wrapper.on('mouseup', close);
		return close;
	}
	return function () {};
}


(function ($) {
	$.fn.shwark = function (template, options, content) {
		options = $.extend({
			afterOpen: function () {},
			beforeClose: function () {},
			beforeOpen: function (data) {
				return data
			}
		}, options);

		var ret = {
			afterOpen: options.afterOpen,
			beforeClose: options.beforeClose,
			beforeOpen: options.beforeOpen,
			close: function () {}
		};

		// Обработчик
		$('html').on('click', this.selector, function (event) {
			event = event || window.event; // For IE
			event.preventDefault();
			ret.close = shwOpen(template, options, $(event.target), content);
			return false;
		});

		return ret;
	};

	$.fn.shwin = function (title, text, options) {
		options = $.extend({
			target: 'body',
			template: false,
			beforeClose: function () {},
			afterOpen: function () {},
			beforeOpen: function (data) {
				return data
			}
		}, options);

		var content = '<div class="w-box"><div class="w-title unselectable">' + title + '</div><div class="w-content"><div class="w-content-body">' + text +
			'</div></div><div class="w-win-bottom"></div><div class="w-close unselectable fa fa-remove ' + (options.closeAll ? 'all-' : '') + 'shwark-close"></div></div>';
		var ret = {
			afterOpen: options.afterOpen,
			beforeClose: options.beforeClose,
			beforeOpen: options.beforeOpen,
			close: shwOpen(options.template, options, $(this), content)
		};

		return ret;
	};

})(jQuery);


(function (factory) {
	if (typeof define === 'function' && define.amd) {
		define(['jquery'], factory);
	} else {
		factory(jQuery);
	}
}(function ($) {
	$.shwarkCloseAll = function () {
		$('.shwark-wrapper').trigger('mouseup');
	};
}));
