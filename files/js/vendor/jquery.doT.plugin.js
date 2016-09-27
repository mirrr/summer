/* global doT, $ */
/* exported tplRet, tplGlobRet, tplFormatNumber */

$.fn.tpl = function (tplId, data) {
	var tpl = doT.template($('#tpl_' + tplId).html());

	if (!$.isArray(data)) {
		data = [data];
	}

	return this.each(function () {
		var html = '';

		for (var itemIdx = 0; itemIdx < data.length; itemIdx++) {
			html += tpl(data[itemIdx]);
		}

		$(this).html(html);
	});
};

$.fn.tplReplace = function (tplId, data) {
	var tpl = doT.template($('#tpl_' + tplId).html());

	if (!$.isArray(data)) {
		data = [data];
	}

	return this.each(function () {
		var html = '';

		for (var itemIdx = 0; itemIdx < data.length; itemIdx++) {
			html += tpl(data[itemIdx]);
		}

		$(this).replaceWith(html);
	});
};

$.fn.tplAppend = function (tplId, data) {
	var tpl = doT.template($('#tpl_' + tplId).html());

	if (!$.isArray(data)) {
		data = [data];
	}

	return this.each(function () {
		var html = '';

		for (var itemIdx = 0; itemIdx < data.length; itemIdx++) {
			html += tpl(data[itemIdx]);
		}

		$(this).append(html);
	});
};

$.fn.tplPrepend = function (tplId, data) {
	var tpl = doT.template($('#tpl_' + tplId).html());

	if (!$.isArray(data)) {
		data = [data];
	}

	return this.each(function () {
		var html = '';

		for (var itemIdx = 0; itemIdx < data.length; itemIdx++) {
			html += tpl(data[itemIdx]);
		}

		$(this).prepend(html);
	});
};

$.fn.tplBefore = function (tplId, data) {
	var tpl = doT.template($('#tpl_' + tplId).html());

	if (!$.isArray(data)) {
		data = [data];
	}

	return this.each(function () {
		var html = '';

		for (var itemIdx = 0; itemIdx < data.length; itemIdx++) {
			html += tpl(data[itemIdx]);
		}

		$(this).before(html);
	});
};

$.fn.tplAfter = function (tplId, data) {
	var tpl = doT.template($('#tpl_' + tplId).html());

	if (!$.isArray(data)) {
		data = [data];
	}

	return this.each(function () {
		var html = '';

		for (var itemIdx = 0; itemIdx < data.length; itemIdx++) {
			html += tpl(data[itemIdx]);
		}

		$(this).after(html);
	});
};

function tplRet (tplId, data) {
	var tpl = doT.template($('#tpl_' + tplId).html());

	if (!$.isArray(data)) {
		data = [data];
	}

	var html = '';

	for (var itemIdx = 0; itemIdx < data.length; itemIdx++) {
		html += tpl(data[itemIdx]);
	}

	return html;
}

function tplGlobRet(tplId, data) {
	var tpl = doT.template($('#tpl_' + tplId).html());
	return tpl(data);
}

function tplFormatNumber(num, decimal) {
	num = Number(num) || 0;
	var separator = '&nbsp;';
	var decpoint = '.';

	var parts = num.toFixed(decimal).split('.');
	parts[0] = parts[0].replace(/(\d{1,3}(?=(\d{3})+(?:\.\d|\b)))/g, '\$1' + separator);
	return (parts[1] ? parts[0] + decpoint + parts[1] : parts[0]);
}
