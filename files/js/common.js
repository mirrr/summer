$(function () {
	/* Инициализация библиотеки всплывающего окна */
	$('select').select2();
	$.wbox.init({
		parent: 'body',
		blures: '#all',
		afterOpen: function () {
			// Кастомные чекбоксы в окне
			$('.w-box input.switch:checkbox').switchCheckbox();
			// Кастомный выпадающий список
			$('.w-box select:not(.custom)').select2({
				language: "ru"
			});
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
	}

	$(window).scroll(updatePage);
	$(window).resize(updatePage);

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

});
