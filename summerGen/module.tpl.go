package main

import (
	"text/template"
)

var moduleTpl = template.Must(template.New("module.go").Parse(`package {{if .Vendor}}{{.FullName}}{{else}}main{{end}}

import (
	"github.com/kennygrant/sanitize"
	"github.com/night-codes/mgo-ai"
	"gopkg.in/gin-gonic/gin.v1"{{if .AddSearch}}
	"gopkg.in/mgo.v2/bson"{{end}}
	"gopkg.in/night-codes/summer.v1"
	"gopkg.in/night-codes/types.v1"
	"time"
)

type (
	{{.Name}}Struct struct {
		ID          uint64    ` + "`" + `form:"id"  json:"id"  bson:"_id"` + "`" + `
		Name        string    ` + "`" + `form:"name" json:"name" bson:"name" valid:"required"` + "`" + `
		Description string    ` + "`" + `form:"description" json:"description" bson:"description"` + "`" + `
		Created     time.Time ` + "`" + `form:"-" json:"created" bson:"created"` + "`" + `
		Updated     time.Time ` + "`" + `form:"-" json:"updated" bson:"updated"` + "`" + `
		Deleted     bool      ` + "`" + `form:"-" json:"deleted" bson:"deleted"` + "`" + `
	}
	{{.Name}}Module struct {
		summer.Module
	}
	{{if .Vendor}}obj map[string]interface{}
	arr []interface{}{{end}}
)

{{if .Vendor}}
func New(panel *summer.Panel, groupTo ...summer.Simple) summer.Simple {
	if len(groupTo) == 0 {
		groupTo = []summer.Simple{nil}
	}
{{else}}
var (
{{end}}
	{{.name}} {{if .Vendor}}:{{end}}= panel.AddModule(
		&summer.ModuleSettings{
			Name:           "{{.name}}",{{if .Collection}}
			CollectionName: "{{.Collection}}",{{end}}
			Title:          "{{.Title}}",{{if .Menu}}
			MenuOrder:      0,
			MenuTitle:      "{{.Title}}",
			{{if .SubDir}}TemplateName:   "{{.name}}/{{.SubDir}}",{{end}}
			Menu:           panel.{{.Menu}},{{end}}{{if .GroupTo}}
			GroupTo:        {{if .Vendor}}groupTo[0]{{else}}{{.GroupTo}}{{end}},
			GroupTitle:     "{{.Title}}",{{end}}
		},
		&{{.Name}}Module{},
	)
{{if .Vendor}}
return {{.name}}
}
{{else}}
)
{{end}}

// Add new record
func (m *{{.Name}}Module) Add(c *gin.Context) {
	var result {{.Name}}Struct
	if !summer.PostBind(c, &result) {
		return
	}
	result.ID = ai.Next("{{.name}}")
	result.Created = time.Now()
	result.Updated = time.Now()
	result.Name = sanitize.HTML(result.Name)
	result.Description = sanitize.HTML(result.Description)

	if err := m.Collection.Insert(result); err != nil {
		c.String(400, "DB error")
		return
	}
	c.JSON(200, obj{"data": result})
}

// Edit record
func (m *{{.Name}}Module) Edit(c *gin.Context) {
	id := types.Uint64(c.PostForm("id"))
	var result {{.Name}}Struct
	var newValue {{.Name}}Struct
	if !summer.PostBind(c, &newValue) {
		return
	}
	if err := m.Collection.FindId(id).One(&result); err == nil {
		result.Name = sanitize.HTML(newValue.Name)
		result.Description = sanitize.HTML(newValue.Description)
		result.Updated = time.Now()
		if err := m.Collection.UpdateId(newValue.ID, obj{"$set": result}); err != nil {
			c.String(400, "DB error")
			return
		}
	}
	c.JSON(200, obj{"data": result})
}

// Get record from DB
func (m *{{.Name}}Module) Get(c *gin.Context) {
	id := types.Uint64(c.PostForm("id"))
	result := {{.Name}}Struct{}
	if err := m.Collection.FindId(id).One(&result); err != nil {
		c.String(404, "Not found")
	}
	c.JSON(200, obj{"data": result})
}

// GetAll records
func (m *{{.Name}}Module) GetAll(c *gin.Context) { {{if or .AddSort .AddSearch .AddPages .AddTabs}}
	filter := struct { {{if .AddSort}}
		Sort    string ` + "`" + `form:"sort"  json:"sort"` + "`" + `{{end}}{{if .AddSearch}}
		Search  string ` + "`" + `form:"search"  json:"search"` + "`" + `{{end}}{{if .AddPages}}
		Page    int    ` + "`" + `form:"page"  json:"page"` + "`" + `{{end}}{{if .AddTabs}}
		Deleted bool   ` + "`" + `form:"deleted"  json:"deleted"` + "`" + `{{end}}
	}{}
	summer.PostBind(c, &filter){{end}}
	results := []{{.Name}}Struct{}
	request := obj{"deleted": {{if .AddTabs}}filter.Deleted{{else}}false{{end}} }
{{if .AddSort}}
	sort := "-_id"

	// sort engine
	if len(filter.Sort) > 0 {
		sort = filter.Sort
	}
{{end}}{{if .AddSearch}}
	// search engine
	if len(filter.Search) > 0 {
		regex := bson.RegEx{Pattern: filter.Search, Options: "i"}
		request["$or"] = arr{
			obj{"name": regex},
			obj{"description": regex},
		}
	}
{{end}}{{if .AddPages}}
	// records pagination
	count, _ := m.Collection.Find(request).Count()
	limit := 0
	skip := 0
	if filter.Page > 0 {
		limit = 50
		skip = limit * (filter.Page - 1)
	}
{{end}}
	// request to DB
	if err := m.Collection.Find(request).Sort({{if .AddSort}}sort{{else}}"-_id"{{end}}){{if .AddPages}}.Limit(limit).Skip(skip){{end}}.All(&results); err != nil {
		c.String(404, "Not found")
		return
	}

	c.JSON(200, obj{"data": results{{if .AddPages}}, "page": filter.Page, "count": count, "limit": limit{{end}} })
}

// {{if .AddTabs}}Action - remove/restore{{else}}Delete - remove{{end}} record
func (m *{{.Name}}Module) {{if .AddTabs}}Action{{else}}Delete{{end}}(c *gin.Context) {
	id := types.Uint64(c.PostForm("id"))

	if err := m.Collection.UpdateId(id, obj{"$set": obj{"deleted": {{if .AddTabs}}c.PostForm("action") == "remove"{{else}}true{{end}} }}); err != nil {
		c.String(404, "Not found")
		return
	}
	c.JSON(200, obj{"data": obj{"id": id}})
}


`))

func init() {
	moduleTpl = template.Must(moduleTpl.Parse(`
{{define "module.html"}}
{{if .AddTabs}}
<div class="tabs">
	<a href="#" data-id="active" class="active">Active</a>
	<a href="#" data-id="deleted">Deleted</a>
</div>{{else}}<p>{{.Name}} list</p>{{end}}
<div class="tablerunner">
	<table id="maintable">
		<thead>
			<tr>
				<th class="td-short">#</th>
				<th{{if .AddSort}} data-sorter="name"{{end}}>Name</th>
				<th{{if .AddSort}} data-sorter="description"{{end}}>Description</th>
				<th>Created</th>
				<th>Updated</th>
				<th class="td-short">Actions</th>
			</tr>
		</thead>
		<tbody>
		</tbody>
	</table>
</div>{{if .AddPages}}
<div class="pagination">
	Page <span id="page-current">1</span> from <span id="page-count">1</span> &nbsp;
	<button id="page-before" disabled="true"> < </button>&nbsp; <button id="page-next" disabled="true"> > </button>
</div>{{end}}
{{end}}
`))

	moduleTpl = template.Must(moduleTpl.Parse(`
{{define "script.js"}}
$(function () {
	var $tbl = $('#maintable>tbody');

	// Filter - sent each time when list of users is loaded
	var filter = { {{if .AddSearch}}
		'search': '',{{end}}{{if .AddPages}}
		'page': 1,{{end}}{{if .AddTabs}}
		'deleted': false,{{end}}
	}

	// load list of users
	function update() {
		$tbl.listLoad({
			url: ajaxUrl + 'getAll',
			method: 'POST',
			noitemsTpl: '{{if .SubDir}}{{.SubDir}}-{{end}}noitems',
			itemTpl: '{{if .SubDir}}{{.SubDir}}-{{end}}item',
			data: filter,{{if .AddPages}}
			success: updatePages{{end}}
		});
	}
	update();{{if .AddSort}}

	// Sort
	$.tools.addSorterFn(function (name, direction) {
		filter.sort = (direction === -1 ? '-' : '') + name;
		update();
	});{{end}}{{if .AddSearch}}

	// Search
	$.tools.addSearchFn(function (value) {
		filter.search = value;
		update();
	});{{end}}{{if .AddTabs}}

	// Ajax tabs
	$('div.tabs').on('change', function () {
		filter.deleted = $('div.tabs').data('active')[0] === "deleted";
		update();
	});{{end}}

	// "New item" button
	var $newAdmin = $.tools.addButton({
		html: '<span class="fa fa-plus"></span> Add record',
		onClick: function () {
			$.wbox.open('Add new record', window.tplRet('{{if .SubDir}}{{.SubDir}}-{{end}}form-add'));
		}
	});

	// Submit form "New item"
	$('#add-form').ajaxFormSender({
		url: ajaxUrl + 'add',
		success: function (result) {
			$tbl.tplPrepend('{{if .SubDir}}{{.SubDir}}-{{end}}item', result.data);
			$('#noitems').hide().remove();
			$tbl.children('tr[data-id=' + result.data.id + ']').children().highlight(500);
			return true;
		}
	});

	// "Edit" button pressed
	$('.edit').ajaxActionSender({
		url: ajaxUrl + 'get',
		method: 'POST',
		success: function (result) {
			$.wbox.open('Change record', window.tplRet('{{if .SubDir}}{{.SubDir}}-{{end}}form-edit', result.data));
		}
	});

	// Submit form "Edit"
	$('#edit-form').ajaxFormSender({
		url: ajaxUrl + 'edit',
		success: function (result) {
			result.data.teasersCount = '...';
			$tbl.children('tr[data-id=' + result.data.id + ']').tplReplace('{{if .SubDir}}{{.SubDir}}-{{end}}item', result.data);
			$tbl.children('tr[data-id=' + result.data.id + ']').children().highlight(500);
			return true;
		}
	});

	// "Remove{{if .AddTabs}}/Restore{{end}}" button pressed
	$('.remove{{if .AddTabs}}, .restore{{end}}').ajaxActionSender({
		url: ajaxUrl + '{{if .AddTabs}}action{{else}}delete{{end}}',
		method: 'POST',
		remove: true // remove from list if success
	});{{if .AddPages}}

	// Pagination
	function updatePages(data) {
		var pages = data.count / data.limit;
		$('#page-count').text(Math.ceil(pages) || 1);
		$('#page-current').text(Math.ceil(data.page) || 1);
		$('#page-next').attr('disabled', 'disabled');
		$('#page-before').attr('disabled', 'disabled');
		if (data.page < pages) {
			$('#page-next').removeAttr('disabled');
		}
		if (data.page > 1) {
			$('#page-before').removeAttr('disabled');
		}
	}
	$('#page-next').forceClick(function () {
		filter.page++;
		update(filter);
	});
	$('#page-before').forceClick(function () {
		filter.page--;
		update(filter);
	});
{{end}}
});
{{end}}
`))

	moduleTpl = template.Must(moduleTpl.Parse(`
{{define "item.html"}}
<tr data-id="{{"{{"}}= it.id {{"}}"}}">
	<td class="td-short">{{"{{"}}= it.id {{"}}"}}</td>
	<td>{{"{{"}}= it.name {{"}}"}}</td>
	<td>{{"{{"}}= it.description {{"}}"}}</td>
	<td>{{"{{"}}= moment(it.created).format("YYYY-MM-DD") {{"}}"}}</td>
	<td>{{"{{"}}= moment(it.updated).format("YYYY-MM-DD") {{"}}"}}</td>
	<td class="td-short">
		<span class="fa fa-pencil edit" title="Edit"></span>{{if .AddTabs}}
	{{"{{"}}? it.deleted {{"}}"}}
		<span class="fa fa-mail-reply-all restore need-confirm"  data-action="restore"></span>
	{{"{{"}}??{{"}}"}}{{end}}
		<span class="fa fa-trash remove need-confirm" title="Remove" data-action="remove"></span>{{if .AddTabs}}
	{{"{{"}}?{{"}}"}}
{{end}}
	</td>
</tr>
{{end}}
`))

	moduleTpl = template.Must(moduleTpl.Parse(`
{{define "noitems.html"}}
<tr id="noitems">
	<td colspan="100">
		<p>{{.Name}} not found</p>
	</td>
</tr>
{{end}}
`))
	moduleTpl = template.Must(moduleTpl.Parse(`
{{define "form-add.html"}}
<form id="add-form">
	<div>
		<label for="name" class="required">Name</label>
		<input type="text" name="name" id="name" />
	</div>
	<div>
		<label for="description">Description</label>
		<textarea name="description" id="description"></textarea>
	</div>
	<div class="form-footer required">
		<button type="submit" name="submit">Create</button>
	</div>
</form>
{{end}}
`))
	moduleTpl = template.Must(moduleTpl.Parse(`
{{define "form-edit.html"}}
<form id="edit-form">
	<div>
		<label for="name" class="required">Name</label>
		<input type="text" name="name" id="name" value="{{"{{"}}= it.name{{"}}"}}" />
	</div>
	<div>
		<label for="description">Description</label>
		<textarea name="description" id="description">{{"{{"}}= it.description{{"}}"}}</textarea>
	</div>
	<div class="form-footer required">
		<input type="hidden" name="id" value="{{"{{"}}= it.id{{"}}"}}" />
		<button type="submit" name="submit">Save</button>
	</div>
</form>

{{end}}
`))

}
