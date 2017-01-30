(function ($) {
	$.fn.highlight = function (time, color) {
		var $that = $(this);
		$that.addClass('n-yellow');
		setTimeout(function () {
			$that.addClass('nn-yellow');
		}, 10);
		setTimeout(function () {
			$that.removeClass('nn-yellow').removeClass('n-yellow');
		}, 1100);
	};
})(jQuery);
