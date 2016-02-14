var api = "/api/v1/values/search/eq"
var findItemsBuildURL = function(mode, keys) {
	return "/api/v1/values/search/"+encodeURIComponent(mode)+"/"+encodeURIComponent(keys.join(","));
};

var getItemBuildURL = function(value_id, fields) {
	return "/api/v1/values/"+encodeURIComponent(value_id)+
		(fields && fields.length > 0? 
			"?"+m.route.buildQueryString({fields: fields}):
			"");
}

var keysToString = function(keys) {
	return _.sortBy(keys).join(":")
}

var NewItem = function(item){
	return {
		value_id: m.prop(item && item.value_id || ""),
		value: m.prop(item && item.value || ""),
		keys: m.prop(item && item.keys || []),
		props: m.prop(item && item.props || {ts: 0}),
		flags: m.prop(item && item.flags || []),
		created_at: m.prop(item && item.created_at || "0"),
		updated_at: m.prop(item && item.updated_at || "0"),
		is_enabled: m.prop(item && item.is_enabled || true),
		is_removed: m.prop(item && item.is_removed || false),
	}
}

// var NewItem = function(item){
// 	return m.prop({
// 		value_id: item && item.value_id || "",
// 		value: item && item.value || "",
// 		keys: item && item.keys || [],
// 		props: item && item.props || {ts: 0},
// 		flags: item && item.flags || [],
// 		created_at: item && item.created_at || "0",
// 		updated_at: item && item.updated_at || "0",
// 		is_enabled: item && item.is_enabled || true,
// 		is_removed: item && item.is_removed || false,
// 	})
// }

var deleteItem = function(value_id) {
	return m.request({url: getItemBuildURL(value_id, []), "method": "GET"})
}

var createItem = function(data, fields) {
	return m.request({url: getItemBuildURL("", fields), "method": "POST", data: data})
}

var updateItem = function(value_id, data, fields) {
	return m.request({url: getItemBuildURL(value_id, fields), "method": "PUT", data: data})
}

var app = {
	items: m.prop([]),
	_addItem: function(item) {
		var _item = NewItem(item);
		_item = {value_id: _item.value_id(), keys: _item.keys(), isNew: m.prop(false), isPending: m.prop(false), data: _item};

		this.items().push(_item);

		return _item;
	},

	isPendingItem: function(value_id) {
		return _.find(this.items(), {value_id: value_id}).isPending();
	},

	isNewItem: function(value_id) {
		return _.find(this.items(), {value_id: value_id}).isNew();
	},

	createItem: function(item, fields) {
		item.value_id = "new:ts:"+_.now();
		if (!item.hasOwnProperty("props")) {
			item.props = {}	;
		}
		item.props["ts"] = _.now();
		var newItem = this._addItem(item);

		m.endComputation();

		newItem.isNew(true);
		newItem.isPending(true);

		return createItem(NewItem(item), fields)
			.then(function(res){
				if (res.status_code == 200 && res.data && res.data.value_id && res.data.value_id.length != 0) {
					newItem.value_id = res.data.value_id;
					newItem.data.value_id(res.data.value_id);
					newItem.data.created_at(res.data.created_at);
					newItem.data.updated_at(res.data.updated_at);

					newItem.isNew(false);
					newItem.isPending(false);

					return newItem;
				} else {
					var message = ["Ошибка создания"];
					message.push("item: "+JSON.stringify(item))
					message.push("fields: "+JSON.stringify(fields))
					message.push("response: "+JSON.stringify(res))
					alert(message.join("\n"));
				}

				return res;
			}.bind(this)); 	
	},

	updateItem: function(value_id, item, fields) {
		var _item = _.find(this.items(), {value_id: value_id});

		if (!_item) {
			var res = {message: "Ошибка поиска"};
			return {
				then: function(cbSuccess) {
					cbSuccess(res);
				}
			}
		}

		var _itemData = {
			value_id: _item.data.value_id(),
			value: _item.data.value(),
			keys: _item.data.keys(),
			props: _item.data.props(),
			flags: _item.data.flags(),
			created_at: _item.data.created_at(),
			updated_at: _item.data.updated_at(),
			is_enabled: _item.data.is_enabled(),
			is_removed: _item.data.is_removed(),
		}

		_itemData = _.assign({}, _itemData, item);
		console.log(_itemData.flags)
		_item.data = NewItem(_itemData)

		_item.isPending(true);
		m.endComputation();

		return updateItem(value_id, _itemData, fields).then(function(res){
			if (res.status_code == 200 && res.data && res.data.value_id && res.data.value_id.length != 0) {
					_item.data.created_at(res.data.created_at);
					_item.data.updated_at(res.data.updated_at);
					_item.isPending(false);
					return _item;
				} else {
					var message = ["Ошибка обновления"];
					message.push("value_id: "+value_id)
					message.push("item: "+JSON.stringify(item))
					message.push("fields: "+JSON.stringify(fields))
					message.push("response: "+JSON.stringify(res))
					alert(message.join("\n"));
				}

				return res;
		})
	},

	init: function(keys) {
		// Загрузка всех позиций в которых ключи пересекаются с keys
		return m.request({url: findItemsBuildURL("contains", keys)})
			.then(function(res){
				if (res.status_code == 200 && res.data) {
					_.each(res.data, this._addItem.bind(this));
				}
			}.bind(this));
	},

	// Выборка всех позиций с точным совпадением ключей
	getItemsByKeys: function(keys) {

		return _(this.items())
			.filter(function(item){
				return keysToString(item.keys) == keysToString(keys);
			})
			.orderBy(function(item){return item.data.props().ts}, "desc")
			.value();
	},

	getItemsByIds: function(ids) {
		return _(this.items())
			.filter(function(item){
				return _.intersection(ids, [item.value_id]).length == 1;
			})
			.orderBy(function(item){return item.data.props().ts}, "desc")
			.value();
	}
}

var ItemList = {
	controller: function(c) {
		var api = {
			element: "ul.uk-list",
			item_config: {
				element: "li",
				element_config: {},
			},
			item_view: {
				controller: function(c) {
					return c;
				},
				view: function(c) {

					return m(c.item_config.element, c.item_config.element_config, c.item.value());
				}
			}
		};

		api = _.merge({}, api, c);

		return api;
	},
	view: function(c) {
		var list = _.map(app.getItemsByKeys(c.keys), function(item){
			var config = _.merge({}, c.item_config, {item: item.data, item_config: {element_config: {key: item.data.value_id()}}}, c)
			
			return m.component(c.item_view, config);
		});

		return m(c.element, list);
	}
}

// Helpfull

var ItemListCreater = {
	controller: function(c) {
		var api = {
			newItem: NewItem({keys: c.keys}),
			onChangeValue: function() {
				return function(value) {
					this.newItem.value(value);
					m.redraw.strategy("none");
				}.bind(this);
			},
			onCreateItem: function() {
				return function(e) {
					e.preventDefault();

					var newItem = {
						keys: this.newItem.keys(),
						value: this.newItem.value(),
					}

					this.newItem.value("");

					app.createItem(newItem, ["value", "props", "is_enabled"])
						.then(function(res){}.bind(this)); 	

					return false;
				}.bind(this)
			}
		};

		api = _.merge({}, api, c);
		return api;
	}, 
	view: function(c) {

		var form = m("form.uk-form", {onsubmit: c.onCreateItem()}, [
			m("input[type='text'].uk-width-1-1", {
				placeholder: c.placeholder || "",
				oninput: m.withAttr("value", c.onChangeValue()), value: c.newItem.value()
				}),
			]);
		return m("div", form);
	}
}