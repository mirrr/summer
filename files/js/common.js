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

	$('.timepicker').datetimepicker(timepk);
	$('.datepicker').datetimepicker({
		lang: 'ru',
		timepicker: false,
		format: 'd.m.Y',
		formatDate: 'd.m.Y',
		onChangeDateTime: function () {
			$(this).datetimepicker('hide');
		}
	});

	// Захват поля поиска и перенос его в шапку
	var $search = $('input[type=text].search');
	var $allsearch = $('input[type=text].allsearch');
	if ($search.length && $allsearch.length) {
		$search.hide();
		$allsearch.show();
		$allsearch.on('keyup', function (e) {
			$search.val($allsearch.val());
			var event = new Event('keyup');
			event.keyCode = e.originalEvent.keyCode;
			event.keyIdentifier = e.originalEvent.keyIdentifier;
			event.charCode = e.originalEvent.charCode;
			event.code = e.originalEvent.code;
			event.which = e.originalEvent.which;
			$search.get(0).dispatchEvent(event);
		});
		$allsearch.on('input', function () {
			$search.val($allsearch.val());
			var event = new Event('input');
			$search.get(0).dispatchEvent(event);
		});
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

});
