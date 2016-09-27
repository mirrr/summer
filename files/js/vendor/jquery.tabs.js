$(function () {
	// Вкладки и табы
	$('div.tabs>*').forceClick(function (event) {
		var $tabs = $(this).parent();

		if ($tabs.attr('multiple')) {
			if (!event.ctrlKey) {
				$tabs.children().removeClass('active');
				$(this).addClass('active');
			} else {
				$(this).toggleClass('active');
			}

			var arr = [];
			$tabs.children().each(function () {
				if ($(this).hasClass('active')) {
					arr.push($(this).data('id'));
				}
			});
			$tabs.data('active', arr);
		} else {
			$tabs.children().removeClass('active');
			$(this).addClass('active');
			$tabs.data('active', [$(this).data('id')]);
		}

		$tabs.trigger('change');
	});
});
