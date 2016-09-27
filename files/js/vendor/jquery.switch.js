/**
 * switchCheckbox.js
 * Version: 1.0.0
 * Author: Ron Masas
 */

(function ($) {
	$.fn.extend({
		switchCheckbox: function () {
			$(this).each(function () {
				var $orgCheckbox = $(this);
				var $newCheckbox = $('<div>', {
					class: 'switch-ui-select'
				}).append($('<div>', {
					class: 'inner'
				}));
				$orgCheckbox.hide().after($newCheckbox);

				// If the original checkbox is checked, add checked class to the switch checkbox.
				if ($orgCheckbox.is(':checked')) {
					$newCheckbox.addClass('checked');
				}

				// If the original checkbox is disabled, add disabled class to the switch checkbox.
				if ($orgCheckbox.is(':disabled')) {
					$newCheckbox.addClass('disabled');
				} else {
					$newCheckbox.click(function () {
						$newCheckbox.toggleClass('checked');
						$orgCheckbox.trigger('click');

						if ($newCheckbox.hasClass('checked')) {
							$orgCheckbox.prop('checked', true);
						} else {
							$orgCheckbox.prop('checked', false);
						}
					});
				}
			});
		}
	});
})(jQuery);
