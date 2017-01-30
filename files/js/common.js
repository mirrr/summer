function createWebSocket(path) {
	return new WebSocket((location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + path);
}

Number.prototype.formatMoney = function (c, d, t) {
	var n = this,
		c = isNaN(c = Math.abs(c)) ? 2 : c,
		d = d == undefined ? "." : d,
		t = t == undefined ? "&nbsp;" : t,
		s = n < 0 ? "-" : "",
		i = parseInt(n = Math.abs(+n || 0).toFixed(c)) + "",
		j = (j = i.length) > 3 ? j % 3 : 0;
	return s + (j ? i.substr(0, j) + t : "") + i.substr(j).replace(/(\d{3})(?=\d)/g, "$1" + t) + (c ? d + Math.abs(n - i).toFixed(c).slice(2) : "");
};

$(function () {
	var pixelRatio = (window && window.devicePixelRatio || 1);

	if (pixelRatio > 1.5) {
		$('#content img.hd').each(function (index, el) {
			var $that = $(el);

			function x2() {
				var src = $that.prop("src").split(".");
				if (src.length > 1) {
					var w = $that.innerWidth();
					var h = $that.innerHeight();
					src[src.length - 2] += "_x2";
					$that.removeClass('hd').addClass('hd_x2').prop('src', src.join(".")).css({
						width: w + 'px',
						height: h + 'px'
					});
				}
			}
			if ($that.innerWidth() > 0) {
				x2();
			} else {
				$(el).one("load", x2);
			}
		});
	}

	$('.navbar-default').find('li.dropdown').each(function (index, el) {
		if (!$(el).children('ul').length || !$(el).children('ul').children('li').length) {
			$(el).hide();
		}
	});
	$('.navbar-default').find('li').each(function (index, el) {
		var $a = $(el).children('a')
		if ($a && $a.length === 1 && ($a.attr('href') === panelPath + moduleName || $a.attr('href') === panelPath + moduleName + "/")) {
			var $li = $(el);
			while ($li.length > 0) {
				$li.addClass('active');
				$li = $li.parent().closest('li');
			}
		}
	});


	// Select2 starter
	function selectFn() {
		var $select = $(this);
		$select.select2({
			language: "ru",
			placeholder: $select.attr("placeholder")
		});
	}
	$('select:not(.custom)').each(selectFn);

	/* Инициализация библиотеки всплывающего окна */
	$.wbox.init({
		parent: 'body',
		blures: '#all',
		afterOpen: function () {
			// Кастомные чекбоксы в окне
			$('.w-box input.switch:checkbox').switchCheckbox();
			// Кастомный выпадающий список
			$('.w-box select:not(.custom)').each(selectFn);
			// Redactor
			$('.w-box textarea.htmlText').redactor();
		},
		beforeClose: function () {
			$('.w-box select.select2-hidden-accessible').select2('close');
		}
	});


	function updatePage() {
		if ($(document).scrollTop() < 10) {
			$('.back-to-top').fadeOut();
		} else {
			$('.back-to-top').fadeIn();
		}

		var contentSize = $('#inside-content').outerHeight() - 5;
		var contentSpace = $(window).height() - $('footer').outerHeight() - $('#content').offset().top - 5;
		$('#content').innerHeight(Math.max(contentSpace, contentSize));
	}

	if ($('#content').length) {
		$(window).scroll(updatePage);
		$(window).resize(updatePage);
		setInterval(updatePage, 300);
		updatePage();
	}

	$('.timepicker').each(function (index, el) {
		$(el).datetimepicker(timepk);
	});
	$('.datepicker').each(function (index, el) {
		$(el).datetimepicker({
			lang: 'ru',
			timepicker: false,
			format: $(el).data("format") || 'd.m.Y',
			formatDate: $(el).data("format") || 'd.m.Y',
			onChangeDateTime: function () {
				$(this).datetimepicker('hide');
			}
		});
	});


	// выделение текста при клике
	$('body').on('click', '.clickselect', function (event) {
		var e = this;
		if (window.getSelection) {
			var s = window.getSelection();
			if (s.setBaseAndExtent) {
				s.setBaseAndExtent(e, 0, e, e.innerText.length - 1);
			} else {
				var r = document.createRange();
				r.selectNodeContents(e);
				s.removeAllRanges();
				s.addRange(r);
			}
		} else if (document.getSelection) {
			var s = document.getSelection();
			var r = document.createRange();
			r.selectNodeContents(e);
			s.removeAllRanges();
			s.addRange(r);
		} else if (document.selection) {
			var r = document.body.createTextRange();
			r.moveToElementText(e);
			r.select();
		}
	});

	$(".li-search").click(function () {
		$(this).children(".allsearch").focus();
	});
	$(".allsearch").focus(function () {
		$(this).parent(".li-search").addClass('focus')
	}).blur(function () {
		$(this).parent(".li-search").removeClass('focus')
	});
});
