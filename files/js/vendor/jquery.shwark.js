(function ($) {
	$.fn.shwark = function (template, options) {
		var $dropdown = $('<div class="shwark-sandbox"></div>');
		$('body').append($dropdown);
		$dropdown.hide();

		options = $.extend({
			data: {},
			toRight: true,
			toDown: true,
			speed: 300,
			offset: 0
		}, options);

		// Обработчик
		$('html').on('click', this.selector, function (event) {
			event = event || window.event; // For IE
			event.preventDefault();

			var $this = $(event.target);
			var data = options.data;

			if ($.isFunction(options.data)) {
				data = options.data($this);
			}

			$dropdown.tpl(template, data);
			var $wrapper = $('<div class="shwark-wrapper"></div>');
			var $circle = $('<div class="shwark-reddot"></div>');
			var $circleIn = $('<div class="shwark-in-reddot"></div>');
			var $button = $this;
			var $target = options.target === 'body' ? $wrapper :
				($(options.target).length ? $($(options.target).get(0)) : $button);

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
				}

				if (winSize.height - 10 < newOfst.top + dropSize.height) {
					options.toDown = false;
				}

				var dropDirection = {};
				var dropStart = {};

				if (options.toDown === null || $button !== $target) {
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

				if (options.toRight === null || $button !== $target) {
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
				$circleIn.css(dropStart).animate(dropDirection, options.speed, function () {});
				$circle.css({
					'left': ofst.left + 'px',
					'top': ofst.top + 'px'
				}).animate({
					'width': '+=' + d + 'px',
					'height': '+=' + d + 'px',
					'left': newOfst.left - (d / 2) + 'px',
					'top': newOfst.top - (d / 2) + 'px',
					'backgroundColor': 'rgba(0,0,0,0)'
				}, options.speed, function () {
					$circle.css({
						'borderRadius': 0
					});
				});


				$circleIn.on('mouseup', function (event) {
					event = event || window.event; // For IE
					event.preventDefault();
					event.stopPropagation();
					return false;
				});

				$wrapper.on('mouseup', function (event) {
					event = event || window.event; // For IE
					event.preventDefault();
					$circleIn.animate(dropStart, options.speed / 2, function () {});
					$circle.css({
						'borderRadius': '50%'
					}).animate({
						'width': '-=' + d + 'px',
						'height': '-=' + d + 'px',
						'left': ofst.left + 'px',
						'top': ofst.top + 'px',
						'backgroundColor': 'rgba(0,0,0,0)'
					}, options.speed / 2, function () {
						$wrapper.remove();
					});


					return false;
				});
			}

			return false;
		});
	};
})(jQuery);


(function (factory) {
		if (typeof define === 'function' && define.amd) {
			define(['jquery'], factory);
		} else {
			factory(jQuery);
		}
	}

	(function ($) {
		$.shwarkCloseAll = function () {
			$('.shwark-wrapper').trigger('mouseup');
		};
	})
);
