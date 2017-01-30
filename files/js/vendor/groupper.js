$(function () {
	$('body').on('mouseup', 'td.chb>label', function (event) {
		var $chb = $(this).parent();
		var $checkbox = $chb.find('input[type=checkbox]');
		var $tr = $chb.closest('tr');
		var $table = $tr.closest('table');
		var top = $tr.offset().top;
		setTimeout(function () {
			if ($checkbox.prop('checked')) {
				$tr.addClass('checked');
			} else {
				$tr.removeClass('checked');
				var $trnext = $tr.nextAll('.checked')
				if ($trnext.length) {
					top = $trnext.offset().top;
				} else {
					$trnext = $tr.prevAll('.checked')
					if ($trnext.length) {
						top = $trnext.offset().top;
					}
				}
			}
			if ($table.find('td.chb input[type=checkbox]:checked').length) {
				setTimeout(function () {
					$table.find('th.chb input[type=checkbox]').prop('checked', true);
					$('#float-grouper').fadeIn('300')
				}, 1);
			} else {
				setTimeout(function () {
					$table.find('th.chb input[type=checkbox]').prop('checked', false);
					$('#float-grouper').fadeOut('300')
				}, 1);
			}
			$('#float-grouper').css({
				'top': top + $tr.outerHeight() / 2 + 'px',
				'left': $table.offset().left + 5 + 'px',
			});
		}, 40);
	});
	$('body').on('mouseup', 'th.chb>label', function (event) {
		var $chb = $(this).parent();
		var $checkbox = $chb.find('input[type=checkbox]');
		var $tr = $chb.closest('tr');
		var $table = $chb.closest('table');
		var $checkboxes = $table.find('td.chb input[type=checkbox]');
		setTimeout(function () {
			if ($checkbox.prop('checked')) {
				$checkboxes.prop('checked', true);
				$checkboxes.closest('tr').addClass('checked');
				$('#float-grouper').fadeIn('300');
				$('#float-grouper').css({
					'top': $tr.offset().top + $tr.outerHeight() / 2 + 'px',
					'left': $table.offset().left + 5 + 'px',
				});
			} else {
				$checkboxes.prop('checked', false);
				$checkboxes.closest('tr').removeClass('checked');
				$('#float-grouper').fadeOut('300');
			}
		}, 40);
	});
});
