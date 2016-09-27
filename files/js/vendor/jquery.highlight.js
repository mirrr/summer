(function ($) {
	$.fn.highlight = function (time, color) {
		var $that = $(this);
		time = time || 300;
		color = color || '#ffff88';

		if (!$that.data('bc')) {
			$that.data('bc', $that.css('background-color'));
		}

		$that.css({'background-color': color});
		setTimeout(function () {
			$that.addClass('highlight-kd');
			setTimeout(function () {
				$that.css({'background-color': $that.data('bc')});
				setTimeout(function () {
					$that.removeClass('highlight-kd');
				}, 500);
			}, 10);
		}, time);
	};
})(jQuery);
