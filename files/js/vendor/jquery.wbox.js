/*
 * Методы управления окном
 */
(function (factory) {
	if (typeof define === 'function' && define.amd) {
		define(['jquery'], factory);
	} else {
		factory(jQuery);
	}
}(function ($) {

	var Wbox = function () {
		var that = this;
		var settings = {
			id: '',
			class: '',
			parent: 'body',
			blures: '.w-blures',
			afterOpen: function () {},
			beforeClose: function () {}
		};

		this.wrap = $('<div/>', {
			'class': 'w-wrap'
		});
		this.win = $('<div/>', {
			'class': 'w-box'
		}).appendTo(this.wrap);
		this.title = $('<div/>', {
			'class': 'w-title unselectable'
		}).appendTo(this.win);
		this.content = $('<div/>', {
			'class': 'w-content'
		}).appendTo(this.win);
		this.bottom = $('<div/>', {
			'class': 'w-win-bottom shadowed'
		}).appendTo(this.win);
		this.close = $('<div/>', {
			'class': 'w-close unselectable fa fa-remove'
		}).appendTo(this.win);
		this.contentDiv = $('<div/>', {
			'class': 'w-content-body'
		}).appendTo(this.content);


		this.init = function (options) {
			settings = $.extend(settings, options);
			that.win.addClass(settings.class);

			if (settings.id.length) {
				that.win.attr('id', settings.id);
			}

			$(settings.parent).append(that.wrap);
			$(settings.parent).on('mousedown', '.w-wrap, .w-close, .w-box [type="cancel"]', function (event) {
				event = event || window.event;

				if ($(event.target).hasClass('w-wrap') || $(event.target).hasClass('w-close') || $(event.target).attr('type') === 'cancel') {
					event.preventDefault();
					that.close();
					return false;
				}
			});
			$(window).resize(that.updatePosition);
			that.content.scroll(that.updateShadow);
		};

		this.updatePosition = function () {
			var maxHeight = Number(($(window).height() * 0.92).toFixed(0));

			that.content.css({
				'max-height': (maxHeight - that.title.outerHeight()) + 'px'
			});
			var margin = ($(window).height() - that.win.outerHeight()) / 2;
			that.win.css({
				'margin-top': margin + 'px'
			});

			if (typeof that.content.perfectScrollbar === 'function') {
				that.content.perfectScrollbar('update');
			}

			that.updateShadow();
		};

		this.updateShadow = function () {
			if (that.content.scrollTop() < 10) {
				that.title.removeClass('shadowed');
			} else {
				that.title.addClass('shadowed');
			}

			if (that.content.scrollTop() >= (that.contentDiv.outerHeight() - that.content.outerHeight() - 40)) {
				that.bottom.removeClass('shadowed');
			} else {
				that.bottom.addClass('shadowed');
			}
		};

		// open window
		this.open = function (title, text, width) {
			that.title.text(title);
			that.contentDiv.html(text);


			if (width) {
				that.win.css({
					'width': width + 'px',
					'min-width': width + 'px'
				});
			}

			if (typeof settings.afterOpen === 'function') {
				settings.afterOpen();
			}

			that.wrap.css({
				opacity: 0
			}).show();

			if (typeof that.content.perfectScrollbar === 'function') {
				that.content.css({
					'position': 'relative'
				});
				that.content.perfectScrollbar({
					suppressScrollX: true,
					includePadding: false
				});
			}
			$('body').on("touchmove", function (e) {
				e.preventDefault();
			});

			$(settings.blures).addClass('w-blured');

			that.wrap.css({
				opacity: 1
			});
			that.content.scrollTop(0);
			that.updatePosition();
		};


		// change content
		this.set = function (title, text, width) {
			if (width) {
				that.win.css({
					'width': width + 'px',
					'min-width': width + 'px'
				});
			}

			that.title.text(title);
			that.contentDiv.html(text);
			that.content.scrollTop(0);
			that.updatePosition();
		};

		// close window
		this.close = function () {
			if (typeof settings.beforeClose === 'function') {
				settings.beforeClose();
			}

			that.wrap.css({
				opacity: 0
			});
			$(settings.blures).removeClass('w-blured');
			setTimeout(function () {
				that.wrap.hide();
				that.contentDiv.html('');
				that.content.scrollTop(0);

				if (typeof that.content.perfectScrollbar === 'function') {
					that.content.perfectScrollbar('destroy');
				}
				$('body').off("touchmove");
			}, 400);
		};

		// Удалить окно
		this.destroy = function () {
			that.close();
			that.title.remove();
			that.content.remove();
			that.close.remove();
			that.win.remove();
			that.wrap.remove();
		};
	};

	$.wbox = new Wbox('#page-top');
}));
$(function () {
	window.wboxwait = window.wboxWait = function () {
		$.wbox.open('Loading...', '<div class="preloader" data-name="preloader"><div class="preloader-wrapper active"><div class="spinner" ><div class="circle-clipper left"><div class="circle"></div></div><div class="gap-patch"><div class="circle"></div></div><div class="circle-clipper right"><div class="circle"></div></div></div></div></div>');
	};
});
